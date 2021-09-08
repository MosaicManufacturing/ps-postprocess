package palette

import (
    "errors"
    "fmt"
    "math"
    "sort"
    "strconv"
)

type Splice struct {
    Drive int
    Length float32
}

type Ping struct {
    Length float32
    Extrusion float32
}

type Algorithm struct {
    Ingoing int
    Outgoing int
    HeatFactor float32
    CompressionFactor float32
    CoolingFactor float32
    Reverse bool
}

type MSF struct {
    Palette Palette
    DrivesUsed []bool
    SpliceList []Splice
    PingList []Ping
}

func NewMSF(paletteData Palette) MSF {
    return MSF{
        Palette:    paletteData,
        DrivesUsed: make([]bool, paletteData.GetInputCount()),
        SpliceList: make([]Splice, 0),
        PingList:   make([]Ping, 0),
    }
}

func (msf *MSF) addSplice(splice Splice) error {
    // splice length validation first
    if len(msf.SpliceList) == 0 {
        // first splice
        minLength := msf.Palette.GetFirstSpliceMinLength()
        if splice.Length < minLength - 5 {
            message := "First Piece Too Short\n"
            message += fmt.Sprintf("The first piece created by Palette would be %.2f mm long, but must be at least %.2f mm.", splice.Length, minLength)
            return errors.New(message)
        }
    } else {
        // all other splices
        spliceDelta := splice.Length - msf.SpliceList[len(msf.SpliceList)-1].Length
        if spliceDelta < MinSpliceLength - 5 {
            message := "Piece Too Short\n"
            message += fmt.Sprintf("Canvas attempted to create a splice that was %.2f mm long, but Palette's minimum splice length is %.2f mm.", splice.Length, MinSpliceLength)
            return errors.New(message)
        }
    }
    msf.SpliceList = append(msf.SpliceList, splice)
    msf.DrivesUsed[splice.Drive] = true
    return nil
}

func (msf *MSF) AddSplice(drive int, length float32) error {
    return msf.addSplice(Splice{
        Drive:  drive,
        Length: length,
    })
}

func (msf *MSF) AddLastSplice(drive int, finalLength float32) error {
    prevSpliceLength := float32(0)
    requiredLength := msf.Palette.GetFirstSpliceMinLength()
    if len(msf.SpliceList) > 0 {
        prevSplice := msf.SpliceList[len(msf.SpliceList)-1]
        prevSpliceLength = prevSplice.Length
        requiredLength = MinSpliceLength
    }
    extraLength := (finalLength - prevSpliceLength) * 0.04
    if (finalLength - prevSpliceLength) < requiredLength {
        extraLength += requiredLength - (finalLength - prevSpliceLength)
    }
    splice := Splice{
        Drive:  drive,
        Length: finalLength + extraLength + msf.Palette.BowdenTubeLength,
    }
    effectiveLoadingOffset := msf.Palette.GetEffectiveLoadingOffset()
    if effectiveLoadingOffset > 2000 {
        if splice.Length < 2000 {
            splice.Length = 2000
        }
    } else {
        if splice.Length < (effectiveLoadingOffset * 1.02) {
            splice.Length = effectiveLoadingOffset * 1.02
        }
    }
    return msf.addSplice(splice)
}

func (msf *MSF) addPing(ping Ping) {
    msf.PingList = append(msf.PingList, ping)
}

func (msf *MSF) AddPing(length float32) {
    msf.addPing(Ping{
        Length: length,
        Extrusion: 0,
    })
}

func (msf *MSF) AddPingWithExtrusion(length, extrusion float32) {
    msf.addPing(Ping{
        Length: length,
        Extrusion: extrusion,
    })
}

func (msf *MSF) GetFilamentLengthsByDrive() []float32 {
    lengths := make([]float32, msf.Palette.GetInputCount())
    if len(msf.SpliceList) == 0 {
        return lengths
    }
    cumulativeLength := float32(0)
    for _, splice := range msf.SpliceList {
        lengths[splice.Drive] += splice.Length - cumulativeLength
        cumulativeLength = splice.Length
    }
    return lengths
}

func (msf *MSF) GetTotalFilamentLength() float32 {
    if len(msf.SpliceList) == 0 {
        return 0
    }
    return msf.SpliceList[len(msf.SpliceList)-1].Length
}

func (msf *MSF) GetOutputAlgorithmsList() []Algorithm {
    numInputs := msf.Palette.GetInputCount()
    algIsPresent := make([][]bool, 0, numInputs)
    for i := 0; i < numInputs; i++ {
        algIsPresent = append(algIsPresent, make([]bool, numInputs))
    }
    algs := make([]Algorithm, 0)

    firstSplice := true
    outgoingExt := 0
    var ingoingExt int

    // "combination" algorithms (splicing two different drives, when transitioning)
    for _, splice := range msf.SpliceList {
        ingoingExt = splice.Drive
        if !firstSplice {
            ingoingIndex := msf.Palette.MaterialMeta[ingoingExt].Index
            outgoingIndex := msf.Palette.MaterialMeta[outgoingExt].Index
            ingoingId := strconv.Itoa(ingoingIndex)
            outgoingId := strconv.Itoa(outgoingIndex)
            if !algIsPresent[ingoingIndex - 1][outgoingIndex - 1] {
                for _, spliceSettings := range msf.Palette.SpliceSettings {
                    if spliceSettings.IngoingID == ingoingId &&
                        spliceSettings.OutgoingID == outgoingId {
                        alg := Algorithm{
                            Ingoing:           ingoingIndex,
                            Outgoing:          outgoingIndex,
                            HeatFactor:        spliceSettings.HeatFactor,
                            CompressionFactor: spliceSettings.CompressionFactor,
                            CoolingFactor:     spliceSettings.CoolingFactor,
                            Reverse:           spliceSettings.Reverse,
                        }
                        algs = append(algs, alg)
                        break
                    }
                }
                algIsPresent[ingoingIndex - 1][outgoingIndex - 1] = true
            }
        }
        outgoingExt = ingoingExt
        firstSplice = false
    }

    // "self-splicing" algorithms (splicing a drive with itself, for run-out detection or hot-swapping)
    // (included manually as MSF will not contain any self-splices)
    for drive := 0; drive < numInputs; drive++ {
        if msf.DrivesUsed[drive] {
            materialIndex := msf.Palette.MaterialMeta[drive].Index
            materialId := strconv.Itoa(materialIndex)
            if !algIsPresent[materialIndex - 1][materialIndex - 1] {
                for _, spliceSettings := range msf.Palette.SpliceSettings {
                    if spliceSettings.IngoingID == materialId &&
                        spliceSettings.OutgoingID == materialId {
                        alg := Algorithm{
                            Ingoing:           materialIndex,
                            Outgoing:          materialIndex,
                            HeatFactor:        spliceSettings.HeatFactor,
                            CompressionFactor: spliceSettings.CompressionFactor,
                            CoolingFactor:     spliceSettings.CoolingFactor,
                            Reverse:           spliceSettings.Reverse,
                        }
                        algs = append(algs, alg)
                        break
                    }
                }
                algIsPresent[materialIndex - 1][materialIndex - 1] = true
            }
        }
    }
    sort.Slice(algs, func(i, j int) bool {
        a := algs[i]
        b := algs[j]
        if a.Ingoing != b.Ingoing {
            return a.Ingoing < b.Ingoing
        }
        return a.Outgoing < b.Outgoing
    })
    return algs
}

func (msf *MSF) GetMSF2Header(filename string) string {
    header := msf.createMSF2()
    // start multicolor mode
    msfFilename := replaceSpaces(truncate(filename, charLimitMSF2))
    printLength := msf.GetTotalFilamentLength()
    intLength := uint(math.Ceil(float64(printLength)))
    header += "O1 D" + msfFilename + " D" + intToHexString(intLength, 8) + EOL
    header += "M0" + EOL
    return header
}

func (msf *MSF) getMSF2PingLine(ping Ping) string {
    line := "O31 D" + floatToHexString(ping.Length)
    if ping.Extrusion > 0 {
        line += " D" + floatToHexString(ping.Extrusion)
    }
    line += EOL
    return line
}

func (msf *MSF) getMSF3PingLine(ping Ping) string {
    if ping.Extrusion > 0 {
        return fmt.Sprintf("O31 L%.2f E%.2f%s", ping.Length, ping.Extrusion, EOL)
    }
    return fmt.Sprintf("O31 L%.2f%s", ping.Length, EOL)
}

func (msf *MSF) GetConnectedPingLine() string {
    pingCount := len(msf.PingList)
    if pingCount == 0 || !msf.Palette.ConnectedMode {
        return ""
    }
    if msf.Palette.Type == TypeP2 {
        return msf.getMSF2PingLine(msf.PingList[pingCount-1])
    }
    return msf.getMSF3PingLine(msf.PingList[pingCount-1])
}

func (msf *MSF) createMSF1() string {
    // todo
    return ""
}

func (msf *MSF) createMSF2() string {
    // todo
    return ""
}

func (msf *MSF) createMSF3() string {
    // todo
    return ""
}

func (msf *MSF) CreateMSF() string {
    if msf.Palette.Type == TypeP1 {
        return msf.createMSF1()
    }
    if msf.Palette.Type == TypeP2 {
        return msf.createMSF2()
    }
    return msf.createMSF3()
}

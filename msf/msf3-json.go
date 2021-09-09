package msf

import "encoding/json"

type palette3Algorithm struct {
    IngoingID int `json:"ingoingId"`
    OutgoingID int `json:"outgoingId"`
    Heat float32 `json:"heat"`
    Compression float32 `json:"compression"`
    Cooling float32 `json:"cooling"`
}

type palette3Splice struct {
    ID int `json:"id"`
    Length float32 `json:"length"`
}

type palette3Ping struct {
    Length float32 `json:"length"`
    Extrusion float32 `json:"extrusion,omitempty"` // exclude if 0 to save space
}

type palette3Json struct {
    Version string `json:"version"`
    Drives []int `json:"drives"`
    Splices []palette3Splice `json:"splices"`
    Pings []palette3Ping `json:"pingCount"`
    Algorithms []palette3Algorithm `json:"algorithms"`
}

type palette3JsonConnected struct {
    Version string `json:"version"`
    Drives []int `json:"drives"`
    Splices []palette3Splice `json:"splices"`
    PingCount int `json:"pingCount"`
    Algorithms []palette3Algorithm `json:"algorithms"`
}

func (p *palette3Json) marshal(connected bool) (string, error) {
    if connected {
        // in connected mode, simply include the ping count
        // rather than the actual list of pings
        pc := palette3JsonConnected{
            Version:    p.Version,
            Drives:     p.Drives,
            Splices:    p.Splices,
            PingCount:  len(p.Pings),
            Algorithms: p.Algorithms,
        }
        bytes, err := json.Marshal(pc)
        if err != nil {
            return "", err
        }
        return string(bytes), nil
    }
    // in accessory mode, include the actual ping data
    bytes, err := json.Marshal(p)
    if err != nil {
        return "", err
    }
    return string(bytes), nil
}

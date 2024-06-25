package ptp

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"sort"
	"strconv"
)

func lerp(minVal, maxVal, t float32) float32 {
	boundedT := float32(math.Max(0, math.Min(1, float64(t))))
	return ((1 - boundedT) * minVal) + (t * maxVal)
}

func roundZ(z float32) float32 {
	return float32(math.Round(float64(z)*10000) / 10000)
}

func getSegmentLengths(ax, ay, az, bx, by, bz, cx, cy, cz float32) (lenAB, lenAC, lenBC float64) {
	ABx := float64(bx - ax)
	ABy := float64(by - ay)
	ABz := float64(bz - az)
	ACx := float64(cx - ax)
	ACy := float64(cy - ay)
	ACz := float64(cz - az)
	BCx := float64(cx - bx)
	BCy := float64(cy - by)
	BCz := float64(cz - bz)
	lenAB = math.Sqrt((ABx * ABx) + (ABy * ABy) + (ABz * ABz))
	lenAC = math.Sqrt((ACx * ACx) + (ACy * ACy) + (ACz * ACz))
	lenBC = math.Sqrt((BCx * BCx) + (BCy * BCy) + (BCz * BCz))
	return
}

func collinear(ax, ay, az, bx, by, bz, cx, cy, cz float32) bool {
	// points A, B, and C are collinear IFF the largest of the lengths of AB, AC, and BC
	// is equal to the sum of the other two, as a corollary to triangle inequality
	lenAB, lenAC, lenBC := getSegmentLengths(ax, ay, az, bx, by, bz, cx, cy, cz)
	return math.Abs(lenAB+lenAC-lenBC) < collinearityEpsilon || // AB + AC == BC, or
		math.Abs(lenAB+lenBC-lenAC) < collinearityEpsilon || // AB + BC == AC, or
		math.Abs(lenAC+lenBC-lenAB) < collinearityEpsilon // AC + BC == AB
}

func directionallyCollinear(ax, ay, az, bx, by, bz, cx, cy, cz float32) bool {
	// same logic as collinear(), except the condition narrows to AB + BC == AC
	lenAB, lenAC, lenBC := getSegmentLengths(ax, ay, az, bx, by, bz, cx, cy, cz)
	return math.Abs(lenAB+lenBC-lenAC) < collinearityEpsilon
}

func floatsToHex(r, g, b float32) string {
	rInt := uint8(r * 255)
	gInt := uint8(g * 255)
	bInt := uint8(b * 255)
	return fmt.Sprintf("#%02x%02x%02x", rInt, gInt, bInt)
}

func sortFloat32Slice(floats []float32) {
	sort.Slice(floats, func(i, j int) bool {
		return floats[i] < floats[j]
	})
}

func openForWrite(w *Writer, name string) error {
	path, ok := w.paths[name]
	if !ok {
		return errors.New(fmt.Sprintf("attempt to open invalid file '%s'", name))
	}
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	writer := bufio.NewWriter(file)
	w.files[name] = file
	w.writers[name] = writer
	return nil
}

func flushAndClose(w *Writer, name string) error {
	if _, ok := w.paths[name]; !ok {
		return errors.New(fmt.Sprintf("attempt to flush and close invalid file '%s'", name))
	}
	writer := w.writers[name]
	if err := writer.Flush(); err != nil {
		return err
	}
	file := w.files[name]
	return file.Close()
}

func concatOntoWriter(w *Writer, writername, filename string) error {
	if _, ok := w.paths[writername]; !ok {
		return errors.New(fmt.Sprintf("attempt to write to invalid file '%s'", writername))
	}
	if _, ok := w.paths[filename]; !ok {
		return errors.New(fmt.Sprintf("attempt to read from invalid file '%s'", filename))
	}
	filepath := w.paths[filename]
	writer := w.writers[writername]
	file, openErr := os.Open(filepath)
	if openErr != nil {
		return openErr
	}
	_, err := io.Copy(writer, file)
	if err != nil {
		return err
	}
	return file.Close()
}

func writeUint32LE(writer *bufio.Writer, val uint32) error {
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, val)
	_, err := writer.Write(buf)
	return err
}

func writeFloat32LE(writer *bufio.Writer, val float32) error {
	return writeUint32LE(writer, math.Float32bits(val))
}

func writeUint8(writer *bufio.Writer, val uint8) error {
	return writer.WriteByte(val)
}

func prepareFloatForJSON(val float32, maxDecimals int) string {
	roundingFactor := math.Pow(10, float64(maxDecimals))
	val64 := math.Round(float64(val)*roundingFactor) / roundingFactor
	return strconv.FormatFloat(val64, 'f', -1, 64)
}

func setToSlice[T comparable](set map[T]bool, sortFn func([]T)) []T {
	values := make([]T, 0, len(set))
	for value := range set {
		values = append(values, value)
	}
	sortFn(values)
	return values
}

func MinFloat32(a, b float32) float32 {
	return float32(math.Min(float64(a), float64(b)))
}

func MaxFloat32(a, b float32) float32 {
	return float32(math.Max(float64(a), float64(b)))
}

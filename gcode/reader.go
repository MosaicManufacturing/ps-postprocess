package gcode

import (
	"bufio"
	"errors"
	"io"
	"os"
)

var ErrEarlyExit = errors.New("Early Exit")
type LineCallback func(Command, int) error

// if callback returns an error, reading will stop before EOF

func ReadByLine(path string, callback LineCallback) (err error) {
	infile, openErr := os.Open(path)
	if openErr != nil {
		err = openErr
		return
	}
	defer func() {
		if closeErr := infile.Close(); closeErr != nil {
			err = closeErr
		}
	}()
	reader := bufio.NewReader(infile)
	lineNumber := 0
	for {
		line, isPrefix, readErr := reader.ReadLine()
		if readErr == io.EOF {
			break
		}
		if readErr != nil {
			err = readErr
			return
		}
		if isPrefix {
			var fragment []byte
			for isPrefix {
				fragment, isPrefix, readErr = reader.ReadLine()
				if readErr == io.EOF {
					break
				}
				if readErr != nil {
					err = readErr
					return
				}
				line = append(line, fragment...)
			}
		}
		gcode := ParseLine(string(line))
		cbErr := callback(gcode, lineNumber)
		if(cbErr == ErrEarlyExit ){
			break; 
		} else if cbErr != nil {
			err = cbErr
			return
		}
		lineNumber++
	}
	return nil
}

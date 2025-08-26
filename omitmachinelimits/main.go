package omitmachinelimits

import (
	"bufio"
	"log"
	"os"
	"strings"

	"mosaicmfg.com/ps-postprocess/gcode"
)

func RemoveDefaultMotionParameterCommands(argv []string) {
	if len(argv) != 2 {
		log.Fatalln("expected 2 command-line arguments")
	}

	const EOL = "\r\n"

	inPath := argv[0]
	outPath := argv[1]

	outfile, createErr := os.Create(outPath)
	if createErr != nil {
		log.Fatalln(createErr)
	}
	defer func() {
		if closeErr := outfile.Close(); closeErr != nil {
			log.Fatalln(closeErr)
		}
	}()

	writer := bufio.NewWriter(outfile)
	defer func() {
		if flushErr := writer.Flush(); flushErr != nil {
			log.Fatalln(flushErr)
		}
	}()

	begunStartSequence := false
	readErr := gcode.ReadByLine(inPath, func(line gcode.Command, lineNum int) error {
		// check if we've reached TYPE:Custom
		if !begunStartSequence && strings.HasPrefix(line.Comment, "TYPE:Custom") {
			begunStartSequence = true
		}

		// filter out default motion parameter commands before TYPE:Custom
		if !begunStartSequence {
			switch line.Command {
			case "M201", "M203", "M204", "M205":
				return nil // skip this line
			}
		}

		// write all other lines
		if _, writeErr := writer.WriteString(line.String() + EOL); writeErr != nil {
			return writeErr
		}
		return nil
	})

	if readErr != nil {
		log.Fatalln(readErr)
	}
}

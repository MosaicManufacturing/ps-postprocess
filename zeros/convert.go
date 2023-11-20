package zeros

import (
	"bufio"
	"log"
	"mosaicmfg.com/ps-postprocess/gcode"
	"os"
	"regexp"
)

var reParamWithNoLeadingZero = regexp.MustCompile("( [XYZEF]-?)\\.([0-9]+)")

const EOL = "\r\n"

func RestoreLeadingZeros(argv []string) {
	argc := len(argv)

	if argc != 2 {
		log.Fatalln("expected 2 command-line arguments")
	}
	inpath := argv[0]
	outpath := argv[1]

	outfile, createErr := os.Create(outpath)
	if createErr != nil {
		log.Fatalln(createErr)
	}
	writer := bufio.NewWriter(outfile)

	err := gcode.ReadByLine(inpath, func(command gcode.Command, _ int) error {
		// ignore non-command lines
		if len(command.Raw) == 0 || len(command.Command) == 0 {
			_, err := writer.WriteString(command.Raw + EOL)
			return err
		}

		// ignore commands that include arbitary text content
		if command.Command == "M70" || command.Command == "M118" {
			_, err := writer.WriteString(command.Raw + EOL)
			return err
		}

		// transform raw command output where needed
		result := replaceAllStringSubmatchFunc(
			reParamWithNoLeadingZero,
			command.Raw,
			func(groups []string) string {
				return groups[1] + "0." + groups[2]
			},
		)
		_, err := writer.WriteString(result + EOL)
		return err
	})
	if err != nil {
		log.Fatalln(err)
	}

	if err = writer.Flush(); err != nil {
		log.Fatalln(err)
	}
	if err = outfile.Close(); err != nil {
		log.Fatalln(err)
	}
}

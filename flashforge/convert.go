package flashforge

import (
	"bufio"
	"fmt"
	"log"
	"mosaicmfg.com/ps-postprocess/gcode"
	"os"
	"strconv"
)

const EOL = "\r\n"

func convert(inpath, outpath string, printExtruder int) error {
	outfile, createErr := os.Create(outpath)
	if createErr != nil {
		return createErr
	}
	writer := bufio.NewWriter(outfile)

	err := gcode.ReadByLine(inpath, func(line gcode.Command, _ int) error {
		// perform the following conversions:
		//   - stabilize print temperature: M109 S<temp> T<ext> -> M6 T<ext>
		//   - stabilize bed temperature: M190 S<temp> -> M7
		//   - toolchange: T<ext> -> M108 T<ext>
		outputLine := line.Raw
		if line.Command == "M109" {
			tool := printExtruder
			if t, ok := line.Params["t"]; ok {
				tool = int(t)
			}
			outputLine = fmt.Sprintf("M6 T%d", tool)
		} else if line.Command == "M190" {
			outputLine = "M7"
		} else if len(line.Command) > 1 && line.Command[0] == 'T' {
			tool, err := strconv.ParseInt(line.Command[1:], 10, 32)
			if err != nil {
				return err
			}
			outputLine = fmt.Sprintf("M108 T%d", tool)
		}
		if _, err := writer.WriteString(outputLine + EOL); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	if err := writer.Flush(); err != nil {
		return err
	}
	if err := outfile.Close(); err != nil {
		return err
	}
	return nil
}

func ConvertCommands(argv []string) {
	argc := len(argv)

	if argc != 3 {
		log.Fatalln("expected 3 command-line arguments")
	}
	inpath := argv[0]
	outpath := argv[1]
	printExtruder, err := strconv.ParseInt(argv[2], 10, 32)
	if err != nil {
		log.Fatalln(err)
	}

	if err := convert(inpath, outpath, int(printExtruder)); err != nil {
		log.Fatalln(err)
	}
}

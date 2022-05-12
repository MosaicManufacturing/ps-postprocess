package comments

import (
	"bufio"
	"log"
	"os"
	"strings"
)

const EOL = "\r\n"

func Strip(argv []string) {
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

	err := ReadByLine(inpath, func(line string, _ int) error {
		if len(line) == 0 {
			return nil
		}
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			return nil
		}
		if commentStart := strings.IndexByte(line, ';'); commentStart >= 0 {
			line = strings.TrimRight(line[:commentStart], " \t")
		}
		if len(line) == 0 {
			return nil
		}
		_, err := writer.WriteString(line + EOL)
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

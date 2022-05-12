package main

import (
	"fmt"
	"log"
	"mosaicmfg.com/ps-postprocess/comments"
	"mosaicmfg.com/ps-postprocess/flashforge"
	"mosaicmfg.com/ps-postprocess/msf"
	"mosaicmfg.com/ps-postprocess/ptp"
	"mosaicmfg.com/ps-postprocess/sequences"
	"mosaicmfg.com/ps-postprocess/ultimaker"
	"os"
)

func main() {
	argv := os.Args[1:]
	if len(argv) == 0 {
		log.Fatalln("expected command as first argument")
	}
	switch argv[0] {
	case "msf":
		msf.ConvertForPalette(argv[1:])
	case "ptp":
		ptp.GenerateToolpath(argv[1:])
	case "comments":
		comments.Strip(argv[1:])
	case "ultimaker":
		ultimaker.AddHeader(argv[1:])
	case "flashforge":
		flashforge.ConvertCommands(argv[1:])
	case "printerscript":
		sequences.ConvertSequences(argv[1:])
	default:
		log.Fatalln(fmt.Sprintf("unknown command '%s'", argv[0]))
	}
}

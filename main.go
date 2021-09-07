package main

import (
    "fmt"
    "log"
    "os"
)

func main() {
    argv := os.Args[1:]
    if len(argv) == 0 {
        log.Fatalln("expected command as first argument")
    }
    switch argv[0] {
    case "ptp":
        generateToolpath(argv[1:])
    default:
        log.Fatalln(fmt.Sprintf("unknown command '%s'", argv[0]))
    }
}

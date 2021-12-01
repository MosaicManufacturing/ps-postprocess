package main

import (
    "fmt"
    "log"
)

func main() {
    licenses, err := getAllRepoModules()
    if err != nil {
        log.Fatalln(err)
    }
    fmt.Println(licenses)
}

package main

import (
    "encoding/json"
    "io/ioutil"
    "log"
    "os"
)

const licenseFile = "./licenses.json"

func getLicenseJSON() (string, error) {
    licenses, err := getAllRepoModules()
    if err != nil {
        return "", err
    }
    jsonBytes, err := json.MarshalIndent(licenses, "", "  ")
    if err != nil {
        return "", err
    }
    return string(jsonBytes), nil
}

func update() {
    licenses, err := getLicenseJSON()
    if err != nil {
        log.Fatalln(err)
    }
    if err := ioutil.WriteFile(licenseFile, []byte(licenses), 0644); err != nil {
        log.Fatalln(err)
    }
}

func check() {
    licenses, err := getLicenseJSON()
    if err != nil {
        log.Fatalln(err)
    }
    fromDiskBytes, err := ioutil.ReadFile(licenseFile)
    fromDisk := string(fromDiskBytes)
    if licenses != fromDisk {
        log.Fatalln("License bundle is out of date.\nRun `./licensebot-update.sh` and commit the changes.")
    }
}

func main() {
    if len(os.Args) == 1 {
        log.Fatalln("expected command as argument")
    }
    command := os.Args[1]
    switch command {
    case "update":
        update()
    case "check":
        check()
    default:
        log.Fatalf("unexpected command '%s' as argument", command)
    }
}

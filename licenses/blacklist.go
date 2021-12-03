package main

import (
    "fmt"
    "io/ioutil"
    "log"
    "strings"
)

var blacklist []string

func init() {
    // load blacklist contents
    blacklistBytes, err := ioutil.ReadFile("./blacklist.txt")
    if err != nil {
        log.Fatalln(err)
    }
    blacklistLines := splitOnNewlines(string(blacklistBytes))
    for _, line := range blacklistLines {
        line = strings.TrimSpace(line)
        if len(line) == 0 {
            continue
        }
        if line[0] == '#' {
            continue
        }
        key := strings.ToLower(line)
        blacklist = append(blacklist, key)
    }
}

func checkBlacklist(licenseId string) error {
    lowercased := strings.ToLower(licenseId)
    for _, key := range blacklist {
        if strings.Contains(lowercased, key) {
            return fmt.Errorf("license %s is blacklisted", licenseId)
        }
    }
    return nil
}

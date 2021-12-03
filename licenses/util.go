package main

import (
    "bytes"
    "os/exec"
    "strings"
)

func runCommand(name string, arg ...string) (string, error) {
    cmd := exec.Command(name, arg...)
    var out bytes.Buffer
    cmd.Stdout = &out
    err := cmd.Run()
    return out.String(), err
}

func splitOnNewlines(str string) []string {
    return strings.Split(strings.ReplaceAll(str, "\r\n", "\n"), "\n")
}

func getGitHubRawUrl(normalUrl string) string {
    rawUrl := strings.Replace(normalUrl, "https://github.com", "https://raw.githubusercontent.com", 1)
    rawUrl = strings.Replace(rawUrl, "/blob", "", 1)
    return rawUrl
}

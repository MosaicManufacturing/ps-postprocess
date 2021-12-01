package main

import (
    "encoding/csv"
    "fmt"
    "io"
    "log"
    "net/http"
    "sort"
    "strings"
)

type License struct {
    Name string
    LicenseId string
    LicenseText string
}

type LicenseWithoutText struct {
    Name string
    LicenseId string
    LicenseUrl string
}

func init() {
    // install go-licenses
    if _, err := runCommand("go", "get", "-v", "github.com/google/go-licenses"); err != nil {
        log.Fatalln(err)
    }
    if _, err := runCommand("go", "build", "github.com/google/go-licenses"); err != nil {
        log.Fatalln(err)
    }
}

func getModuleLicenses(relativePath string) ([]LicenseWithoutText, error) {
    stdout, err := runCommand("./go-licenses", "csv", relativePath)
    if err != nil {
        return nil, err
    }

    r := csv.NewReader(strings.NewReader(stdout))
    var entries []LicenseWithoutText
    for {
        record, err := r.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            return nil, err
        }
        dependency := record[0]
        if strings.HasPrefix(dependency, "mosaicmfg.com") {
            // skip our own subpackages
            continue
        }
        licenseUrl := record[1]
        licenseId := record[2]
        entries = append(entries, LicenseWithoutText{
            Name:       dependency,
            LicenseId:  licenseId,
            LicenseUrl: licenseUrl,
        })
    }
    return entries, nil
}

func getAllRepoModules() ([]License, error) {
    // TODO: walk repository and determine all modules, submodules, etc.
    directories := []string{
        "..",
    }

    licensesMap := make(map[string]LicenseWithoutText)

    for _, directory := range directories {
        fmt.Printf("Getting dependencies of %s\n", directory)
        licenses, err := getModuleLicenses(directory)
        if err != nil {
            return nil, err
        }
        // add licenses to map where dependency is not already present
        for _, license := range licenses {
            if _, exists := licensesMap[license.Name]; !exists {
                licensesMap[license.Name] = license
            }
        }
    }

    // convert licenses map to licenses slice including text content
    licenses := make([]License, 0, len(licensesMap))
    for _, license := range licensesMap {
        licenseUrl := license.LicenseUrl
        if strings.HasPrefix(licenseUrl, "https://github.com") {
            licenseUrl = getGitHubRawUrl(licenseUrl)
        }

        // retrieve license text from URL
        resp, err := http.Get(licenseUrl)
        if err != nil {
            return nil, err
        }
        bodyBytes, err := io.ReadAll(resp.Body)
        if err != nil {
            return nil, err
        }
        licenseText := string(bodyBytes)
        if err := resp.Body.Close(); err != nil {
            return nil, err
        }

        licenses = append(licenses, License{
            Name:        license.Name,
            LicenseId:   license.LicenseId,
            LicenseText: licenseText,
        })
    }
    // sort the slice alphabetically
    sort.Slice(licenses, func(i, j int) bool {
        return licenses[i].Name < licenses[j].Name
    })

    return licenses, nil
}
package sequences

import (
    "../gcode"
    "strings"
)

type sequencesPreflight struct {
    totalLayers int
    totalTime int
}

func preflight(inpath string) (sequencesPreflight, error) {
   results := sequencesPreflight{}

   err := gcode.ReadByLine(inpath, func(line gcode.Command, lineNumber int) error {
       if line.Command != "" {
           return nil
       }
       if line.Comment == "LAYER_CHANGE" {
           results.totalLayers++
       } else if strings.HasPrefix(line.Comment, "estimated printing time (normal mode) = ") {
           timeEstimate, err := gcode.ParseTimeString(line.Comment)
           if err != nil {
               return err
           }
           results.totalTime = int(timeEstimate)
       }
       return nil
   })
   return results, err
}
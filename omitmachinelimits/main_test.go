package omitmachinelimits

import (
	"os"
	"testing"
)

func TestRemoveDefaultMotionParameterCommands(t *testing.T) {
	inpath := "test-files/in.gcode"
	expectedPath := "test-files/expected.gcode"

	outfile, err := os.CreateTemp(os.TempDir(), "*.gcode")
	if err != nil {
		t.Fatal(err)
	}
	if err = outfile.Close(); err != nil {
		t.Fatal(err)
	}
	outpath := outfile.Name()

	RemoveDefaultMotionParameterCommands([]string{
		inpath,
		outpath,
	})

	expectedContent, err := os.ReadFile(expectedPath)
	if err != nil {
		t.Fatal(err)
	}
	outpathContent, err := os.ReadFile(outpath)
	if err != nil {
		t.Fatal(err)
	}

	if string(expectedContent) != string(outpathContent) {
		t.Errorf("filtered G-code does not match expected")
	}

	if err = os.Remove(outpath); err != nil {
		t.Fatal(err)
	}
}

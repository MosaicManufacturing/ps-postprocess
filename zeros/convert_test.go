package zeros

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func testWithInputFiles(t *testing.T, dirname string) {
	inpath := path.Join("test-files", dirname, "in.gcode")
	expectedPath := path.Join("test-files", dirname, "expected.gcode")
	outfile, err := ioutil.TempFile(os.TempDir(), "*.gcode")
	fmt.Println(outfile.Name())
	if err != nil {
		t.Fatal(err)
	}
	if err = outfile.Close(); err != nil {
		t.Fatal(err)
	}
	outpath := outfile.Name()

	RestoreLeadingZeros([]string{
		inpath,
		outpath,
	})

	expectedContent, err := ioutil.ReadFile(expectedPath)
	if err != nil {
		t.Fatal(err)
	}
	outpathContent, err := ioutil.ReadFile(outpath)
	if err != nil {
		t.Fatal(err)
	}

	if string(expectedContent) != string(outpathContent) {
		t.Fatal("text content does not match expected")
	}

	if err = os.Remove(outpath); err != nil {
		t.Fatal(err)
	}
}

func TestRestoreLeadingZeros(t *testing.T) {
	testWithInputFiles(t, "1")
}

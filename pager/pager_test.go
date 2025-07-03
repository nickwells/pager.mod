package pager

import (
	"bytes"
	"os"
	"testing"

	"github.com/nickwells/testhelper.mod/v2/testhelper"
)

func TestGetTrialPagers(t *testing.T) {
	const pagerCmd = "somePager -x -y"

	var ec testhelper.EnvCache

	envWithoutPager := testhelper.EnvEntry{Key: "PAGER"}
	if err := (&ec).Setenv(envWithoutPager); err != nil {
		t.Fatal("couldn't set the 'PAGER' environment variable:", err)
		return
	}

	pagers := getTrialPagers()
	testhelper.DiffStringSlice(t, "with empty PAGER", "trial pagers",
		pagers, dfltPagers)

	(&ec).ResetEnv()

	envWithPager := testhelper.EnvEntry{Key: "PAGER", Value: pagerCmd}
	if err := (&ec).Setenv(envWithPager); err != nil {
		t.Fatal("couldn't set the 'PAGER' environment variable:", err)
		return
	}

	pagers = getTrialPagers()
	expPagers := []string{pagerCmd}
	expPagers = append(expPagers, dfltPagers...)
	testhelper.DiffStringSlice(t, "with PAGER='"+pagerCmd+"'", "trial pagers",
		pagers, expPagers)
}

func TestIsWriterATerminal(t *testing.T) {
	const filename = "testdata/plainFile"

	var wb bytes.Buffer

	if isWriterATerminal(&wb) {
		t.Error("a bytes.Buffer is not a terminal")
	}

	wf, err := os.OpenFile(filename, os.O_WRONLY, 0)
	if err != nil {
		t.Fatal("couldn't open the file to test isWriterATerminal: ", filename)
		return
	}

	if isWriterATerminal(wf) {
		t.Error("a plain file is not a terminal")
	}
}

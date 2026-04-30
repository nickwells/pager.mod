package pager

import (
	"bytes"
	"os"
	"testing"

	"github.com/nickwells/testhelper.mod/v2/testhelper"
)

func TestGetTrialPagers(t *testing.T) {
	const (
		pagerCmd1 = "somePager -x -y"
		pagerCmd2 = "otherPager -x -y"
	)

	testCases := []struct {
		testhelper.ID
		envVarNames []string
		envVals     []testhelper.EnvEntry
		expPagers   []string
	}{
		{
			ID: testhelper.MkID("no-env"),
		},
		{
			ID: testhelper.MkID("PAGER"),
			envVals: []testhelper.EnvEntry{
				{Key: "PAGER", Value: pagerCmd1},
			},
			expPagers: []string{pagerCmd1},
		},
		{
			ID:          testhelper.MkID("Other"),
			envVarNames: []string{"Other"},
			envVals: []testhelper.EnvEntry{
				{Key: "Other", Value: pagerCmd2},
			},
			expPagers: []string{pagerCmd2},
		},
		{
			ID:          testhelper.MkID("both"),
			envVarNames: []string{"Other"},
			envVals: []testhelper.EnvEntry{
				{Key: "PAGER", Value: pagerCmd1},
				{Key: "Other", Value: pagerCmd2},
			},
			expPagers: []string{pagerCmd2, pagerCmd1},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			var ec testhelper.EnvCache

			ec.Clearenv()

			if err := (&ec).Setenv(tc.envVals...); err != nil {
				t.Fatal("couldn't set the environment variables:", err)
				return
			}

			pagers := getTrialPagers(tc.envVarNames)
			expPagers := tc.expPagers
			expPagers = append(expPagers, dfltPagers...)
			testhelper.DiffStringSlice(t,
				tc.IDStr(), "trial pagers list",
				pagers, expPagers)

			(&ec).ResetEnv()
		})
	}
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

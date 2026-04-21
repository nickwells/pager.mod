package pager

import (
	"io"
	"os"
	"os/exec"
	"regexp"

	"golang.org/x/term"
)

var cmdSplit = regexp.MustCompile(`\s+`)

// Pager represents the input to be paged and the command to be used to
// perform the paging.
type Pager struct {
	pagerIn io.WriteCloser
	cmd     *exec.Cmd
}

var dfltPagers = []string{"less", "more"}

// getTrialPagers returns a list of Pagers to be searched for. The first
// entry will be the value of the PAGER environment variable and the rest
// will be the default pagers: "less" and "more".
func getTrialPagers(envVarNames []string) []string {
	trialPagers := []string{}

	envVarNames = append(envVarNames, "PAGER")

	for _, envVar := range envVarNames {
		envPager := os.Getenv(envVar)
		if envPager != "" {
			trialPagers = append(trialPagers, envPager)
		}
	}

	return append(trialPagers, dfltPagers...)
}

// getPagerCmd returns the pager command or nil if no executable command can
// be found.
func getPagerCmd(envVarNames []string) *exec.Cmd {
	for _, tp := range getTrialPagers(envVarNames) {
		parts := cmdSplit.Split(tp, -1)

		path, err := exec.LookPath(parts[0])
		if err == nil {
			return exec.Command(path, parts[1:]...) //nolint:gosec
		}
	}

	return nil
}

// isWriterATerminal returns true if the io.Writer is a Terminal
func isWriterATerminal(w io.Writer) bool {
	if outFile, ok := w.(*os.File); ok {
		if term.IsTerminal(int(outFile.Fd())) { //nolint:gosec
			return true
		}
	}

	return false
}

// Start checks that at least one of the Writers is a terminal and if so it
// will start a pager command and connect the terminal writers to the input
// of the pager command. It returns the pager which should have Done() called
// on it after any output is complete.
//
// The pager command is chosen from the contents of the PAGER environment
// variable or the commands 'less' or 'more' if that cannot be used.
func Start(sw SetW) *Pager {
	return StartWithEnvVars(sw)
}

// StartWithEnvVars behaves like [Start] but the supplied environment
// variable names are checked first for an executable pager in addition to
// those checked by Start.
func StartWithEnvVars(sw SetW, envVarNames ...string) *Pager {
	stdoutIsTty := isWriterATerminal(sw.StdW())
	stderrIsTty := isWriterATerminal(sw.ErrW())

	if !stdoutIsTty && !stderrIsTty {
		return nil
	}

	cmd := getPagerCmd(envVarNames)
	if cmd == nil {
		return nil
	}

	pagerIn, err := cmd.StdinPipe()
	if err != nil {
		return nil
	}

	cmd.Stdout = sw.StdW()
	cmd.Stderr = sw.ErrW()

	err = cmd.Start()
	if err != nil {
		return nil
	}

	if stdoutIsTty {
		sw.SetStdW(pagerIn)
	}

	if stderrIsTty {
		sw.SetErrW(pagerIn)
	}

	return &Pager{
		pagerIn: pagerIn,
		cmd:     cmd,
	}
}

// Done will wait for the pager to complete. Note that it is safe to call
// with a nil pointer.
func (p *Pager) Done() {
	if p == nil {
		return
	}

	p.pagerIn.Close()

	_ = p.cmd.Wait()
}

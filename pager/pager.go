package pager

import (
	"io"
	"os"
	"os/exec"
	"regexp"

	"golang.org/x/term"
)

var pagerCmdDflts = []string{"less", "more"} //nolint: gofumpt
var cmdSplit = regexp.MustCompile(`\s+`)

type Pager struct {
	pagerIn io.WriteCloser
	cmd     *exec.Cmd
}

// getPagerCmd returns the pager command or nil if no executable command can
// be found. It will first try the contents of the PAGER environment variable
// and then the list of pager command defaults.
func getPagerCmd() *exec.Cmd {
	trialPagers := []string{os.Getenv("PAGER")}
	trialPagers = append(trialPagers, pagerCmdDflts...)

	for _, tp := range trialPagers {
		parts := cmdSplit.Split(tp, -1)
		path, err := exec.LookPath(parts[0])
		if err == nil {
			return exec.Command(path, parts[1:]...)
		}
	}

	return nil
}

// isWriterATerminal returns true if the io.Writer is a Terminal
func isWriterATerminal(w io.Writer) bool {
	if outFile, ok := w.(*os.File); ok {
		if term.IsTerminal(int(outFile.Fd())) {
			return true
		}
	}
	return false
}

// Start checks that at least one of the Writers is a terminal and if so it
// will start a pager command and connect the terminal writers to the input
// of the pager command. It returns the pager which should have Done() called
// on it after any output is complete.
func Start(sw SetW) *Pager {
	outIsTty := isWriterATerminal(sw.StdW())
	errIsTty := isWriterATerminal(sw.ErrW())

	if !outIsTty && !errIsTty {
		return nil
	}
	cmd := getPagerCmd()
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

	if outIsTty {
		sw.SetStdW(pagerIn)
	}
	if errIsTty {
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

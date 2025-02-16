package git

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

type DiffOpt func(*DiffRange)

type DiffRange struct {
	StartCommit string
	EndCommit   string
}

func WithStart(start string) DiffOpt {
	return func(d *DiffRange) {
		d.StartCommit = start
	}
}

func WithEnd(end string) DiffOpt {
	return func(d *DiffRange) {
		d.EndCommit = end
	}
}

func NewDiffRange(opts ...DiffOpt) (*DiffRange, error) {
	cmd := exec.Command("git", "merge-base", "main", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("error getting merge-base: %w", err)
	}

	d := &DiffRange{
		StartCommit: strings.TrimSpace(string(output)),
		EndCommit:   "HEAD",
	}

	for _, opt := range opts {
		opt(d)
	}

	return d, nil
}

func (d *DiffRange) GetChangedFiles() ([][]byte, error) {

	cmd := exec.Command("git", "diff", d.StartCommit, d.EndCommit, "--name-only")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("error getting diff: %w", err)
	}

	files := bytes.Split(output, []byte("\n"))

	return files, nil

}

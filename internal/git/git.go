package git

import (
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
	cmd := exec.Command("git", "merge-base", "origin/main", "HEAD")
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

func (d *DiffRange) GetModifiedFilepaths() ([]string, error) {

	cmd := exec.Command("git", "diff", d.StartCommit, d.EndCommit, "--name-only")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("error getting diff: %w", err)
	}

	out_str := strings.TrimSpace(string(output))

	files := strings.Split(out_str, "\n")

	return files, nil

}

func (d *DiffRange) GetFileDiff(file string) (*string, error) {

	cmd := exec.Command("git", "diff", d.StartCommit, d.EndCommit, file)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("error getting file diff: %w", err)
	}

	out_str := strings.TrimSpace(string(output))

	return &out_str, nil
}

func GetStartCommitFromCount(endCommit string, count int) (string, error) {
	cmd := exec.Command("git", "rev-parse", fmt.Sprintf("%s~%d", endCommit, count))
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("error calculating start commit: %w", err)
	}
	return strings.TrimSpace(string(output)), nil
}

func NewDiffRangeWithCount(endCommit string, count int) (*DiffRange, error) {
	if endCommit == "" {
		endCommit = "HEAD"
	}

	startCommit, err := GetStartCommitFromCount(endCommit, count)
	if err != nil {
		return nil, err
	}

	return &DiffRange{
		StartCommit: startCommit,
		EndCommit:   endCommit,
	}, nil
}

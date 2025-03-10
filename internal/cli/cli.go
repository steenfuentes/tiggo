package cli

import (
	"flag"
	"fmt"

	"github.com/steenfuentes/tiggo/internal/git"
)

type Config struct {
	StartCommit string
	EndCommit   string
	PrevCount   int
}

func ParseFlags() (*Config, error) {
	var config Config

	prevCount := flag.Int("p", 0, "Number of previous commits to include")

	flag.Usage = func() {
		fmt.Printf("Usage of tiggo:\n")
		fmt.Printf("  tiggo [flags] [start-commit] [end-commit]\n")
		fmt.Printf("  tiggo [flags] [end-commit]\n")
		fmt.Printf("  tiggo [flags]\n\n")
		fmt.Printf("Flags:\n")
		flag.PrintDefaults()
		fmt.Printf("\nExamples:\n")
		fmt.Printf("  tiggo abc123 def456     # Analyze changes between two commits\n")
		fmt.Printf("  tiggo def456 -p 5       # Analyze 5 commits before def456\n")
		fmt.Printf("  tiggo -p 3              # Analyze last 3 commits from HEAD\n")
		fmt.Printf("  tiggo                   # Analyze changes from merge-base with main to HEAD\n")
	}

	flag.Parse()

	args := flag.Args()

	switch len(args) {
	case 0:
		if *prevCount > 0 {
			config.EndCommit = "HEAD"
			config.PrevCount = *prevCount
		}
	case 1:
		if *prevCount > 0 {
			config.EndCommit = args[0]
			config.PrevCount = *prevCount
		} else {
			return nil, fmt.Errorf("when providing single commit, -p flag is required")
		}
	case 2:
		if *prevCount > 0 {
			return nil, fmt.Errorf("cannot use -p flag with two commits")
		}
		config.StartCommit = args[0]
		config.EndCommit = args[1]
	default:
		return nil, fmt.Errorf("too many arguments")
	}

	return &config, nil
}

func BuildDiffRange(config *Config) (*git.DiffRange, error) {
	if config.StartCommit != "" && config.EndCommit != "" {
		return git.NewDiffRange(
			git.WithStart(config.StartCommit),
			git.WithEnd(config.EndCommit),
		)
	}

	if config.PrevCount > 0 {
		return git.NewDiffRangeWithCount(config.EndCommit, config.PrevCount)
	}

	return git.NewDiffRange()
}

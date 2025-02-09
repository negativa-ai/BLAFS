package main

import (
	"fmt"
	"os"

	"github.com/alexflint/go-arg"
)

type ShadowCmd struct {
	Images string `arg:"-i,--images" help:"Images to shadow, separated by comma"`
}
type DebloatCmd struct {
	Images string `arg:"-i,--images" help:"Images to debloat separated by comma"`
	Top    int    `arg:"-t,--top" help:"Top N layers to debloat" default:"-1"`
}

var args struct {
	Shadow  *ShadowCmd  `arg:"subcommand:shadow" help:"Shadow images"`
	Debloat *DebloatCmd `arg:"subcommand:debloat" help:"Debloat images"`
}

func main() {
	p := arg.MustParse(&args)
	if p.Subcommand() == nil {
		p.WriteHelp(os.Stdout)
		os.Exit(1)
	}

	switch {
	case args.Shadow != nil:
		fmt.Printf("checkout requested for branch %s\n", args.Shadow.Images)
	case args.Debloat != nil:
		fmt.Printf("commit requested with message \"%s\"\n", args.Debloat.Images)
	}
}

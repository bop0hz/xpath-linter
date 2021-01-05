package report // import "github.com/bop0hz/xpath-linter/report"

import (
	"fmt"
	"os"

	"github.com/gookit/color"
	"github.com/mattn/go-isatty"
)

type Reporter interface {
	Compile(name string, nodeIndex int, reason string, target string) string
}

func Red(t string) string {
	if isatty.IsTerminal(os.Stdout.Fd()) {
		return color.FgRed.Render(t)
	}
	return t
}

func Green(t string) string {
	if isatty.IsTerminal(os.Stdout.Fd()) {
		return color.FgGreen.Render(t)
	}
	return t
}

type ColorReport struct{}

func (r *ColorReport) Compile(name string, i int, reason string, target string) string {
	return fmt.Sprintf("[%s][%d] %s: %s", Red(name), i, reason, Green(target))
}

type Report struct{}

func (r *Report) Compile(name string, i int, reason string, target string) string {
	return fmt.Sprintf("[%s][%d] %s: %s", name, i, reason, target)
}

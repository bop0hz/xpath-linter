package version // import "github.com/bop0hz/xpath-linter/pkg/version"

import (
	"fmt"
)

var (
	builtAt = ""
	build   = ""
	version = ""
	hash    = ""
	branch  = ""
)

type AppInfo struct {
	BuiltAt string `json:"builtAt,omitempty"`
	Build   string `json:"build,omitempty"`
	Version string `json:"version,omitempty"`
	Hash    string `json:"hash,omitempty"`
	Branch  string `json:"branch,omitempty"`
}

func Version() *AppInfo {
	return &AppInfo{
		Build:   build,
		BuiltAt: builtAt,
		Version: version,
		Hash:    hash,
		Branch:  branch,
	}
}

func (app *AppInfo) Print() string {
	return fmt.Sprintf("Application was built at = %s, version = %s, build = %s\n", app.BuiltAt, app.Version, app.Build)
}

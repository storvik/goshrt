package version

import (
	"fmt"
	"runtime"
)

var (
	GitVersion     = "not injected"
	GitCommit      = "not injected"
	BuildTime      = "not injected"
	BuildGoVersion = "not injected"
)

func String() string {
	return fmt.Sprintf(`Goshrt - URL shortener written in Go
Version:     %s
Commit hash: %s
Build time:  %s
Build with:  %s
Run arch:    %s/%s`,
		GitVersion, GitCommit, BuildTime, BuildGoVersion, runtime.GOOS, runtime.GOARCH)
}

func Print() {
	fmt.Printf(`Goshrt - URL shortener written in Go
Version:     %s
Commit hash: %s
Build time:  %s
Build with:  %s
Run arch:    %s/%s`,
		GitVersion, GitCommit, BuildTime, BuildGoVersion, runtime.GOOS, runtime.GOARCH)
}

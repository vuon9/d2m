//go:build windows
// +build windows

package app

import (
	"os/exec"
)

var (
	windowsOpenCmd      string
	windowsPossibleCmds = []string{
		"Start-Process",
	}
)
var windowsPossibleExecs = []string{
	"powershell.exe",
	"explorer.exe",
}
var windowsOpenExec string

func init() {
	for _, exe := range windowsPossibleExecs {
		execPath, err := exec.LookPath(exe)
		if err != nil {
			continue
		}

		windowsOpenExec = execPath
		break
	}

	for _, cmd := range windowsPossibleCmds {
		err := exec.Command(windowsOpenExec, "Get-Command", cmd).Run()
		if err != nil {
			continue
		}

		windowsOpenCmd = cmd
		break
	}
}

// OpenURL open url
func OpenURL(url string) error {
	if windowsOpenCmd != "" {
		return exec.Command(windowsOpenExec, windowsOpenCmd, url).Run()
	}
	return nil
}

// CommandExists returns true if an 'open' command exists
func CommandExists() bool {
	return windowsOpenCmd != ""
}

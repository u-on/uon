package runcmd

import "testing"

func TestRunCommand(t *testing.T) {
	RunCommand("cmd.exe", "/c", "dir")
}

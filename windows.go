package main

import (
	"os/exec"
)

//Architecture specific ping command
func pinger(address string) *exec.Cmd {
	return exec.Command("ping", address, "-n", "1")
}

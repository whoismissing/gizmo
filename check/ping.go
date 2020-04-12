package check

import (
	"os/exec"
	"runtime"
)

// Ping() executes the windows / linux ping command on a user-specified ip string
func Ping(ip string) bool {
	cmd := exec.Command("")

	if runtime.GOOS == "windows" {
		cmd = exec.Command("ping", ip, "-n", "1")
	} else {
		cmd = exec.Command("ping", ip, "-c", "1")
	}

	err := cmd.Run()

	var status bool
	if err != nil {
		status = false
	} else {
		status = true
	}

	return status
}

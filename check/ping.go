package check

import (
    "os/exec"

    "runtime"
)

func Ping(ip string) bool {
    cmd := exec.Command("")

    if runtime.GOOS == "windows" {
        cmd = exec.Command("ping", ip, "-n", "1")
    } else {
        cmd = exec.Command("ping", ip, "-c", "1")
    }

    err := cmd.Run()

    return (err == nil)
}

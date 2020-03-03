package check

import (
    "os/exec"

    "log"
    "syscall"
)

/*
The External() function will execute an outside user-defined program 
that accepts a target IP as the first argument to implement a custom 
service check. The program should return 0 in the case of a succesful
check and non-zero in the case of failure.

Essentially, we are just doing the below commands.

./external_check 192.168.1.1
echo $?

*/
func External(ip string, filepath string) bool {
    cmd := exec.Command(filepath, ip)

    if err := cmd.Start(); err != nil {
        log.Printf("cmd.Start: %v", err)
        return false
    }

    if err := cmd.Wait(); err != nil {
        if exiterr, ok := err.(*exec.ExitError); ok {
            /* The program has exited with an exit code != 0
             * This should work on both Unix and Windows */
            if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
                _ = status
                //log.Printf("Exit Status: %d", status.ExitStatus())
                return false
            }
        } else {
            log.Fatalf("cmd.Wait: %v", err)
        }
    }

    // Program returned 0 so success!
    return true
}

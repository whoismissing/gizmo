package check

import (
	"golang.org/x/crypto/ssh"
	"time"
)

// Ssh() authenticates to an ssh server and executes a command
// on port 22 given a username, password, and command strings
func Ssh(ip string, user string, pass string, command string) bool {
	// TODO: make timeout adjustable
	sshConfig := &ssh.ClientConfig{
		User:    user,
		Auth:    []ssh.AuthMethod{ssh.Password(pass)},
		Timeout: time.Second,
	}

	sshConfig.HostKeyCallback = ssh.InsecureIgnoreHostKey()

	// TODO: make port adjustable
	client, err := ssh.Dial("tcp", ip+":22", sshConfig)
	if err != nil {
		return false
	}

	session, err := client.NewSession()
	if err != nil {
		return false
	}

	_, err = session.CombinedOutput(command)
	if err != nil {
		return false
	}

	client.Close()
	return true
}

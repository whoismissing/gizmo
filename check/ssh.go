package check

import (
	"golang.org/x/crypto/ssh"
	"time"
)

func Ssh(ip string, user string, pass string) bool {
	sshConfig := &ssh.ClientConfig{
		User:    user,
		Auth:    []ssh.AuthMethod{ssh.Password(pass)},
		Timeout: time.Second,
	}
	sshConfig.HostKeyCallback = ssh.InsecureIgnoreHostKey()
	client, err := ssh.Dial("tcp", ip+":22", sshConfig)
	if err != nil {
		return false
	}
	session, err := client.NewSession()
	if err != nil {
		return false
	}
	_, err = session.CombinedOutput("dir")
	if err != nil {
		return false
	}
	client.Close()
	return true
}

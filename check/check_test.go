package check

import (
	"fmt"
	"testing"
)

func TestWebCheck(t *testing.T) {
	t.Logf("TestWeb")

	up := Web("http://localhost")

	if up {
		fmt.Println("web is up")
	} else {
		fmt.Println("web is down")
	}

}

func TestDnsCheck(t *testing.T) {
	t.Logf("TestDns")

	down := Dns("127.0.0.1", "doesntexist.com")

	fmt.Println("dns is", down)

}

func TestSshCheck(t *testing.T) {
	t.Logf("TestSsh")

	up := Ssh("127.0.0.1", "username", "password", "ls > /tmp/i")
	fmt.Println("ssh is", up)

}

func TestFtpCheck(t *testing.T) {
	t.Logf("TestFtp")

	up := Ftp("127.0.0.1", "anonymous", "", "test.txt")

	fmt.Println("ftp is", up)
}

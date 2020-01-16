package check

import (
	"time"

	"github.com/jlaffaye/ftp"
)

func Ftp(ip string, user string, pass string, filename string) bool {
	client, err := ftp.Dial(ip+":21", ftp.DialWithTimeout(time.Second))
	if err != nil {
		return false
	}
	err = client.Login(user, pass)
	if err != nil {
		return false
	}
	_, err = client.Retr(filename)
	if err != nil {
		return false
	}
	client.Quit()
	return true
}

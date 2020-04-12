package check

import (
	"github.com/jlaffaye/ftp"

	"time"
)

// Ftp() authenticates and retrieves a file on ftp
// port 21 given a username, password, and filename.
func Ftp(ip string, user string, pass string, filename string) bool {

    // TODO: make port and timeout adjustable
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

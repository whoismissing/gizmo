package check

import (
	"testing"
    "fmt"
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

package sshutils

import (
	"github.com/fulldump/biff"
	"os"
	"testing"
)

func TestConnect(t *testing.T) {

	if os.Getenv("INTEGRATION") == "" {
		t.SkipNow()
	}

	keyBytes, err := os.ReadFile("../certs/id_rsa")
	biff.AssertNil(err)

	sc, err := Connect(keyBytes, "root", "127.0.0.1:2222")
	biff.AssertNil(err)
	biff.AssertNotNil(sc)

}

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

	s := SshConfig{
		Addr:    "127.0.0.1:2222",
		User:    "root",
		Private: keyBytes,
	}

	sc, err := s.Connect()
	biff.AssertNil(err)
	biff.AssertNotNil(sc)

}

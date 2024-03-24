package sshutils

import (
	"github.com/fulldump/biff"
	"os"
	"testing"
)

func TestConnect(t *testing.T) {

	keyBytes, err := os.ReadFile("../certs/id_rsa")
	biff.AssertNil(err)

	s := SshConfig{
		Addr:    "127.0.0.1:2222",
		User:    "root",
		Private: keyBytes,
	}

	err = s.Connect()
	biff.AssertNil(err)

}

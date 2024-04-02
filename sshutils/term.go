package sshutils

import (
	"golang.org/x/crypto/ssh"
	"os"

	"github.com/pkg/sftp"
)

type Term struct {
	client  *ssh.Client
	session *ssh.Session

	sftp *sftp.Client
}

func NewTerm(private []byte, user, addr string) (*Term, error) {

	t, err := Connect(private, user, addr)
	if err != nil {
		return nil, err
	}

	s, err := t.NewSession()
	if err != nil {
		return nil, err
	}

	return &Term{client: t, session: s}, nil
}

func (t *Term) Run(f string) (string, error) {

	b, err := t.session.CombinedOutput(f)

	return string(b), err
}

func (t *Term) Upload(path, remotePath string) (int64, error) {

	f, err := os.Open(path)
	if err != nil {
		return 0, err
	}

	return UploadFile(t.client, f, remotePath)
}

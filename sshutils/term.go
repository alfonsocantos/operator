package sshutils

import (
	"fmt"
	"golang.org/x/crypto/ssh"
)

type Term struct {
	client  *ssh.Client
	session *ssh.Session
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

func (t *Term) Run(f string, params ...any) (string, error) {

	b, err := t.session.CombinedOutput(fmt.Sprintf(f, params...))

	return string(b), err
}

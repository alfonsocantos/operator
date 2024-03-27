package sshutils

import "golang.org/x/crypto/ssh"

type Term struct {
	session *ssh.Session
}

func NewTerm(cl *ssh.Client) (*Term, error) {
	t, err := cl.NewSession()
	if err != nil {
		return nil, err
	}
	return &Term{session: t}, nil
}

func (t *Term) Run(f string, params ...string) {

	// todo run things

}

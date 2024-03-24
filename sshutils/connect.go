package sshutils

import (
	"golang.org/x/crypto/ssh"
	"net"
)

type SshConfig struct {
	Addr    string
	User    string
	Private []byte
}

func (s *SshConfig) Connect() error {

	signer, err := ssh.ParsePrivateKey(s.Private)
	if err != nil {
		return err
	}

	config := &ssh.ClientConfig{
		User: s.User,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	conn, err := ssh.Dial("tcp", s.Addr, config)
	if err != nil {
		return err
	}

	session, err := conn.NewSession()
	defer session.Close()

	return nil

}

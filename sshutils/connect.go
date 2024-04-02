package sshutils

import (
	"golang.org/x/crypto/ssh"
)

func Connect(private []byte, user, addr string) (*ssh.Client, error) {

	signer, err := ssh.ParsePrivateKey(private)
	if err != nil {
		return nil, err
	}

	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // todo do it configurable
	}

	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return nil, err
	}

	return client, nil
}

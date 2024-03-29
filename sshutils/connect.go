package sshutils

import (
	"golang.org/x/crypto/ssh"
	"net"
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
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return nil, err
	}

	return client, nil
}

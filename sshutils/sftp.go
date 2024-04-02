package sshutils

import (
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io"
)

func UploadFile(cl *ssh.Client, origin io.ReadCloser, remotePath string) (int64, error) {

	sc, err := sftp.NewClient(cl)
	if err != nil {
		return 0, err
	}
	defer sc.Close()

	dst, err := sc.Create(remotePath)
	if err != nil {
		return 0, err
	}
	defer dst.Close()

	size := int64(32 * 1024)
	copied := int64(0)
	for {
		bc, err := io.CopyN(dst, origin, size)
		if err == io.EOF {
			return copied + bc, nil
		}
		if err != nil {
			return 0, err
		}
	}
}

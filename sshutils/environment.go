package sshutils

import (
	"context"
	"github.com/fulldump/biff"
	"os"
)

func environment(f ...func(ctx context.Context)) {

	keyBytes, err := os.ReadFile("../certs/id_rsa")
	biff.AssertNil(err)

	sc, _ := NewTerm(keyBytes, "root", "127.0.0.1:2222")

	ctx := context.WithValue(context.Background(), "terminal", sc)

	for _, f := range f {
		f(ctx)
	}

}

package sshutils

import (
	"context"
	"github.com/fulldump/biff"
	"os"
	"testing"
)

func TestTerm_Run(t *testing.T) {

	if os.Getenv("INTEGRATION") == "" {
		t.SkipNow()
	}

	environment(func(ctx context.Context) {
		t, ok := ctx.Value("terminal").(*Term)
		biff.AssertTrue(ok)

		out, err := t.Run("whoami")
		biff.AssertNil(err)
		biff.AssertEqual(out, "root\n")

	})
}

func TestTerm_UploadFile(t *testing.T) {

	if os.Getenv("INTEGRATION") == "" {
		t.SkipNow()
	}

	f, _ := os.CreateTemp(".", "*")
	defer func() {
		os.Remove(f.Name())
	}()
	n, _ := f.WriteString("Hello world!")
	biff.AssertEqual(n, 12)

	environment(func(ctx context.Context) {

		t, ok := ctx.Value("terminal").(*Term)
		biff.AssertTrue(ok)

		s, err := t.Upload(f.Name(), "file.txt")
		biff.AssertNil(err)
		biff.AssertEqual(s, int64(12))
	})

}

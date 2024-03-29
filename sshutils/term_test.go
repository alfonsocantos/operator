package sshutils

import (
	"context"
	"github.com/fulldump/biff"
	"os"
	"testing"
)

func TestRunCommand(t *testing.T) {

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

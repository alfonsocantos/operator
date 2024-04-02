package sshutils

import (
	"context"
	"github.com/fulldump/biff"
	"io"
	"strings"
	"testing"
)

func TestUploadFile(t *testing.T) {

	path := "file.txt"
	rd := io.NopCloser(strings.NewReader("Hello World!"))

	environment(func(ctx context.Context) {

		t, ok := ctx.Value("terminal").(*Term)
		biff.AssertTrue(ok)

		s, err := UploadFile(t.client, rd, path)
		biff.AssertNil(err)
		biff.AssertEqual(s, int64(12))
	})

}

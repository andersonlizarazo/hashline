package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"

	"code.andersonlizarazo.com/hashline/pkg/database"
	"code.andersonlizarazo.com/hashline/pkg/session"
)

func cmdReadHash(ctx context.Context, db database.Database, stdout io.Writer, path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	data, _ = bytes.CutSuffix(data, []byte("\n"))
	lines := bytes.Split(data, []byte("\n"))

	tokens, err := session.LineTokens(ctx, db, path, 0, int64(len(lines)))
	if err != nil {
		return err
	}

	for idx := range lines {
		fmt.Fprintf(stdout, "%s\t%s\n", tokens[idx], lines[idx])
	}

	return nil
}

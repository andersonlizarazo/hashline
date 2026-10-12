package session

import (
	"bytes"
	"context"
	_ "embed"
	"errors"
	"hash/fnv"
	"os"

	"code.andersonlizarazo.com/hashline/pkg/database"
)

func LineTokens(ctx context.Context, db database.Database, path string, pos, span int64) ([]string, error) {
	// TODO: Detect simple shifts, and line changes to avoid regenerating tokens for unchanged lines.

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	hash := fnv.New64()
	_, err = hash.Write(data)
	if err != nil {
		return nil, err
	}

	fnvhash := hash.Sum64()

	filehash, tokens, err := db.ReadTokens(ctx, path, pos, span)
	if err != nil && !errors.Is(err, database.ErrNoTokens) {
		return nil, err
	}
	if tokens != nil && filehash == fnvhash {
		return tokens, nil
	}

	data, _ = bytes.CutSuffix(data, []byte("\n"))
	amount := len(bytes.Split(data, []byte("\n")))
	tokens, err = db.MintTokens(ctx, int64(amount))
	if err != nil {
		return nil, err
	}

	// TODO: Save tokens.

	return tokens[pos : pos+span], nil
}

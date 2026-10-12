package database

import (
	"context"
	"database/sql"
	_ "embed"
	"encoding/binary"
	"errors"
	"time"

	_ "github.com/ncruces/go-sqlite3/driver"
)

var (
	ErrTooManyTokens = errors.New("too many tokens")
	ErrSpanTooLarge  = errors.New("span too large")
	ErrNoTokens      = errors.New("no tokens")
)

//go:embed schema.sql
var schema string

const InMemory = ":memory:"

type Database struct {
	db   *sql.DB
	toks []string
}

func New(ctx context.Context, dataSource string, tokens []string) (Database, error) {
	if len(tokens) == 0 {
		panic("invalid [tokens] argument")
	}

	db, err := sql.Open("sqlite3", dataSource)
	if err != nil {
		return Database{}, err
	}

	pingCtx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	if err := db.PingContext(pingCtx); err != nil {
		db.Close()
		return Database{}, err
	}

	_, err = db.ExecContext(ctx, schema)
	if err != nil {
		db.Close()
		return Database{}, err
	}

	return Database{db: db, toks: tokens}, nil
}

func (db Database) Close() {
	db.db.Close()
}

func (db Database) MintTokens(ctx context.Context, amount int64) ([]string, error) {
	if amount < 1 {
		panic("invalid amount")
	}

	if amount > int64(len(db.toks)) {
		return nil, ErrTooManyTokens
	}

	// TODO: Use transactions.

	var ringCursor uint
	if err := db.db.QueryRowContext(ctx, "SELECT ring_cursor FROM token_pool").Scan(&ringCursor); err != nil {
		return nil, err
	}

	newRingCursor := (ringCursor + uint(amount)) % uint(len(db.toks))
	if _, err := db.db.ExecContext(ctx, "UPDATE token_pool SET emitted_at = unixepoch(), ring_cursor = ?", newRingCursor); err != nil {
		return nil, err
	}

	tokens := make([]string, amount)
	for idx := range amount {
		tokens[idx] = db.toks[(ringCursor+uint(idx))%uint(len(db.toks))]
	}

	return tokens, nil
}

func (db Database) WriteTokens(ctx context.Context, filePath string, fileHash uint64, tokens []string, pos int64) error {
	// TODO: Use transactions.

	hashBlob := make([]byte, 8)
	binary.NativeEndian.PutUint64(hashBlob, fileHash)

	if _, err := db.db.ExecContext(ctx, "INSERT INTO metadata (path, fnv64_hash, updated_at) VALUES (:path, :hash, unixepoch()) ON CONFLICT (path) DO UPDATE SET fnv64_hash = excluded.fnv64_hash, updated_at = excluded.updated_at", sql.Named("path", filePath), sql.Named("hash", hashBlob)); err != nil {
		return err
	}

	stmt, err := db.db.PrepareContext(ctx, "INSERT INTO line (path, token, pos) VALUES (:path, :token, :pos) ON CONFLICT (path, token) DO UPDATE SET pos = excluded.pos")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, token := range tokens {
		if _, err := stmt.ExecContext(ctx, sql.Named("path", filePath), sql.Named("token", token), sql.Named("pos", pos)); err != nil {
			return err
		}

		pos += 1
	}

	return nil
}

func (db Database) ReadTokens(ctx context.Context, path string, pos, span int64) (uint64, []string, error) {
	hashBlob := make([]byte, 8)
	err := db.db.QueryRowContext(ctx, "SELECT fnv64_hash FROM metadata WHERE path = :path", sql.Named("path", path)).Scan(&hashBlob)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, nil, ErrNoTokens
	}
	if err != nil {
		return 0, nil, err
	}

	hash := binary.NativeEndian.Uint64(hashBlob)

	rows, err := db.db.QueryContext(ctx, "SELECT token FROM line WHERE path = :path AND pos BETWEEN :start AND :end ORDER BY pos", sql.Named("path", path), sql.Named("start", pos), sql.Named("end", pos+span))
	if err != nil {
		return 0, nil, err
	}
	defer rows.Close()

	tokens := make([]string, 0, span)
	for rows.Next() {
		var token string
		if err := rows.Scan(&token); err != nil {
			return 0, nil, err
		}

		tokens = append(tokens, token)
	}
	if err := rows.Err(); err != nil {
		return 0, nil, err
	}

	if int64(len(tokens)) != span {
		return 0, nil, ErrSpanTooLarge
	}

	return hash, tokens, nil
}

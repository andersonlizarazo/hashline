package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"

	"code.andersonlizarazo.com/hashline/pkg/database"
)

func main() {
	ctx := context.Background()

	db, err := database.New(ctx, database.InMemory, []string{"1", "2", "3", "4"})
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}

	if err := do(ctx, db, os.Args, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(2)
	}
}

func do(ctx context.Context, db database.Database, args []string, stdout io.Writer) error {
	flagSet := flag.NewFlagSet(args[0], flag.ContinueOnError)

	if err := flagSet.Parse(args[1:]); err != nil {
		return err
	}

	file := flagSet.Arg(0)
	switch narg := flagSet.NArg(); narg {
	case 0:
		flagSet.Usage()
		return nil
	case 1:
		return cmdReadHash(ctx, db, stdout, file)
	case 3:
		return cmdReplaceHash(ctx, db, file, flagSet.Arg(1), flagSet.Arg(2))
	default:
		return fmt.Errorf("invalid number of arguments: %d", narg)
	}
}

func cmdReplaceHash(ctx context.Context, db database.Database, path string, hash string, newText string) error {
	panic("not implemented")
}

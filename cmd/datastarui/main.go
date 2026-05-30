package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	copycmd "github.com/coreycole/datastarui/cmd/datastarui/internal/copy"
)

func main() {
	if err := run(context.Background(), os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(ctx context.Context, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: datastarui add|update|diff|doctor [components...] --source <dir> --target <dir> --module <module>")
	}
	cmd, rest := args[0], args[1:]
	opts, err := parseOptions(rest)
	if err != nil {
		return err
	}

	var result copycmd.Result
	switch cmd {
	case "init", "add":
		result, err = copycmd.Add(ctx, opts)
	case "update":
		result, err = copycmd.Update(ctx, opts)
	case "diff":
		result, err = copycmd.Diff(ctx, opts)
	case "doctor":
		err = copycmd.Doctor(ctx, opts)
	default:
		return fmt.Errorf("unknown command %q", cmd)
	}
	if err != nil {
		return err
	}
	if len(result.CopiedFiles) > 0 {
		fmt.Printf("copied %d files\n", len(result.CopiedFiles))
	}
	if len(result.Drift) > 0 {
		for _, drift := range result.Drift {
			fmt.Printf("%s %s\n", drift.Status, drift.Path)
		}
		return fmt.Errorf("drift found: %s", strings.TrimSpace(result.LockPath))
	}
	if result.LockPath != "" {
		fmt.Println(result.LockPath)
	}
	return nil
}

func parseOptions(args []string) (copycmd.Options, error) {
	opts := copycmd.Options{SourceRoot: ".", TargetRoot: "./pkg/datastarui"}
	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch arg {
		case "--source":
			value, next, err := flagValue(args, i, arg)
			if err != nil {
				return copycmd.Options{}, err
			}
			opts.SourceRoot = value
			i = next
		case "--target":
			value, next, err := flagValue(args, i, arg)
			if err != nil {
				return copycmd.Options{}, err
			}
			opts.TargetRoot = value
			i = next
		case "--module":
			value, next, err := flagValue(args, i, arg)
			if err != nil {
				return copycmd.Options{}, err
			}
			opts.TargetModule = value
			i = next
		case "--force":
			opts.Force = true
		default:
			if strings.HasPrefix(arg, "--source=") {
				opts.SourceRoot = strings.TrimPrefix(arg, "--source=")
				continue
			}
			if strings.HasPrefix(arg, "--target=") {
				opts.TargetRoot = strings.TrimPrefix(arg, "--target=")
				continue
			}
			if strings.HasPrefix(arg, "--module=") {
				opts.TargetModule = strings.TrimPrefix(arg, "--module=")
				continue
			}
			if strings.HasPrefix(arg, "-") {
				return copycmd.Options{}, fmt.Errorf("unknown flag %q", arg)
			}
			opts.Components = append(opts.Components, arg)
		}
	}
	return opts, nil
}

func flagValue(args []string, index int, name string) (string, int, error) {
	next := index + 1
	if next >= len(args) || strings.HasPrefix(args[next], "--") {
		return "", index, fmt.Errorf("%s requires a value", name)
	}
	return args[next], next, nil
}

package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	copycmd "github.com/coreycole/datastarui/cmd/datastarui/internal/copy"
	"github.com/spf13/cobra"
)

func main() {
	if err := newRootCommand(context.Background(), os.Stdout, os.Stderr).Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func newRootCommand(ctx context.Context, stdout, stderr io.Writer) *cobra.Command {
	opts := &copycmd.Options{SourceRoot: ".", TargetRoot: "./pkg/datastarui"}

	cmd := &cobra.Command{
		Use:           "datastarui",
		Short:         "Copy DatastarUI source into consumer apps",
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	cmd.SetOut(stdout)
	cmd.SetErr(stderr)

	cmd.PersistentFlags().StringVar(&opts.SourceRoot, "source", opts.SourceRoot, "DatastarUI source checkout")
	cmd.PersistentFlags().StringVar(&opts.TargetRoot, "target", opts.TargetRoot, "consumer copied-source target")
	cmd.PersistentFlags().StringVar(&opts.TargetModule, "module", opts.TargetModule, "consumer Go module path")
	cmd.PersistentFlags().BoolVar(&opts.Force, "force", false, "overwrite modified managed files; reserved for explicit emergency use")

	cmd.AddCommand(
		newCopyCommand(ctx, "add", "Copy selected components into a consumer", opts, copycmd.Add),
		newCopyCommand(ctx, "init", "Initialize copied DatastarUI source in a consumer", opts, copycmd.Add),
		newCopyCommand(ctx, "update", "Refresh copied DatastarUI source", opts, copycmd.Update),
		newDiffCommand(ctx, opts),
		newDoctorCommand(ctx, opts),
	)
	return cmd
}

func newCopyCommand(ctx context.Context, use, short string, opts *copycmd.Options, run func(context.Context, copycmd.Options) (copycmd.Result, error)) *cobra.Command {
	return &cobra.Command{
		Use:   use + " [components...]",
		Short: short,
		Args:  cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Components = args
			result, err := run(ctx, *opts)
			if err != nil {
				return err
			}
			return printResult(cmd.OutOrStdout(), result)
		},
	}
}

func newDiffCommand(ctx context.Context, opts *copycmd.Options) *cobra.Command {
	return &cobra.Command{
		Use:   "diff",
		Short: "Report copied-source drift from datastarui.lock.json",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			result, err := copycmd.Diff(ctx, *opts)
			if err != nil {
				return err
			}
			if err := printResult(cmd.OutOrStdout(), result); err != nil {
				return err
			}
			if len(result.Drift) > 0 {
				return fmt.Errorf("drift found: %s", strings.TrimSpace(result.LockPath))
			}
			return nil
		},
	}
}

func newDoctorCommand(ctx context.Context, opts *copycmd.Options) *cobra.Command {
	return &cobra.Command{
		Use:   "doctor",
		Short: "Validate copied-source target health",
		Args:  cobra.NoArgs,
		RunE: func(*cobra.Command, []string) error {
			return copycmd.Doctor(ctx, *opts)
		},
	}
}

func printResult(w io.Writer, result copycmd.Result) error {
	if len(result.CopiedFiles) > 0 {
		if _, err := fmt.Fprintf(w, "copied %d files\n", len(result.CopiedFiles)); err != nil {
			return err
		}
	}
	for _, drift := range result.Drift {
		if _, err := fmt.Fprintf(w, "%s %s\n", drift.Status, drift.Path); err != nil {
			return err
		}
	}
	if result.LockPath != "" {
		_, err := fmt.Fprintln(w, result.LockPath)
		return err
	}
	return nil
}

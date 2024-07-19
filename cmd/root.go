package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/bamsammich/chop/internal/config"
	"github.com/bamsammich/chop/pkg/log"
	"github.com/spf13/cobra"
)

var (
	formatter = log.NewFormatter()
)

func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "chop [path]",
		Short: "Write structured logs in a human-readable way.",
		Args:  cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return FromStdin()
			}
			return FromFile(args[0])
		},
	}

	cmd.PersistentFlags().StringSliceVarP(&config.Include, "include", "i", config.Include, "fields to print (excludes others)")
	cmd.PersistentFlags().StringSliceVarP(&config.Exclude, "exclude", "e", config.Exclude, "fields to exclude")
	cmd.PersistentFlags().StringVarP(&config.DefaultField, "default-field", "d", config.DefaultField, "default field for unstructured logs")

	return cmd
}

func FromStdin() error {
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) != 0 {
		return fmt.Errorf("nothing passed to chop")
	}
	return printLogs(os.Stdin)
}

func printLogs(r io.Reader) error {
	var (
		scanner = bufio.NewScanner(r)
		count   = 0
	)

	for scanner.Scan() {
		log, err := formatter.FromString(count, scanner.Text())
		if err != nil {
			return err
		}
		fmt.Println(log)
		count++
	}
	return nil
}

func FromFile(path string) error {
	fi, err := os.Stat(path)
	if err != nil {
		return err
	}
	if fi.IsDir() {
		return fmt.Errorf("must not be a directory")
	}
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("unable to open log: %w", err)
	}
	defer file.Close()

	return printLogs(file)
}

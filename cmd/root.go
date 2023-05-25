package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/bamsammich/chop/internal/config"
	"github.com/bamsammich/chop/pkg/log"
	"github.com/spf13/cobra"
)

var (
	formatTuples []string
)

func NewRootCmd() *cobra.Command {
	format := log.NewFormatter()
	cmd := &cobra.Command{
		Use:   "chop [path]",
		Short: "Write structured logs in a human-readable way.",
		Args:  cobra.RangeArgs(0, 1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			for _, tuple := range formatTuples {
				if !regexp.MustCompile(`.*=\d+`).Match([]byte(tuple)) {
					return fmt.Errorf("format %q is invalid: must match <field_name>=<number>", tuple)
				}
				parts := strings.Split(tuple, "=")
				width, err := strconv.Atoi(parts[1])
				if err != nil {
					return err
				}
				field := parts[0]
				if len(field) > width {
					width = len(field)
				}
				format.Add(field, width)
			}
			if config.PrintExtras {
				format.Add(config.ExtraFieldsName, 50)
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return FromStdin(format)
			}
			return FromFile(args[0], format)
		},
	}
	cmd.PersistentFlags().StringSliceVarP(&formatTuples, "format", "f", formatTuples, "tuples of field names to print and column width")
	cmd.PersistentFlags().BoolVarP(&config.PrintExtras, "print-all", "a", false, "print all fields; fields without format defined will be printed as JSON")
	cmd.PersistentFlags().StringVarP(&config.DefaultField, "default-field", "d", config.DefaultField, "default field for unstructured logs")

	return cmd
}

func FromStdin(formatter *log.Formatter) error {
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) != 0 {
		return fmt.Errorf("nothing passed to chop")
	}
	return printLogs(os.Stdin, formatter)
}

func printLogs(r io.Reader, formatter *log.Formatter) error {
	var (
		scanner = bufio.NewScanner(os.Stdin)
		count   = 0
	)
	formatter.Header().Print()
	for scanner.Scan() {
		log, err := formatter.FromString(count, scanner.Text())
		if err != nil {
			return err
		}
		log.Print()
		count++
	}
	return nil
}

func FromFile(path string, formatter *log.Formatter) error {
	fi, err := os.Stat(path)
	if err != nil {
		return err
	}
	if fi.IsDir() {
		return fmt.Errorf("must not be a directory")
	}
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		return err
	}
	return printLogs(file, formatter)
}

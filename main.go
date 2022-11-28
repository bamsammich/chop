package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&messageField, "message-field", "m", defaultMessageField, "field containing log message")
	rootCmd.PersistentFlags().StringSliceVarP(&columns, "columns", "c", []string{defaultMessageField}, "field names to extract to columns")
	rootCmd.PersistentFlags().IntVar(&maxColumnWidth, "max-column-width", 60, "set maximum column width")
	rootCmd.PersistentFlags().IntVarP(&beforeLine, "before", "B", -1, "print lines before this count [-1 unsets this]")
	rootCmd.PersistentFlags().IntVarP(&afterLine, "after", "A", 0, "print lines after this count")
	rootCmd.PersistentFlags().BoolVarP(&excludeExtraFields, "exclude-extra-fields", "x", false, "exclude extra fields not defined by `columns`")
}

var (
	beforeLine int
	afterLine  int

	maxColumnWidth int
	columns        []string

	countHeader         = "   #"
	fieldsHeader        = "fields"
	defaultMessageField = "message"
	messageField        string
	excludeExtraFields  bool
)

var rootCmd = &cobra.Command{
	Use:   "chop [path]",
	Short: "Write structured log in a human-readable way.",
	Args:  cobra.RangeArgs(0, 1),
	PreRun: func(cmd *cobra.Command, args []string) {
		if messageField != defaultMessageField {
			columns = replace(columns, defaultMessageField, messageField)
		}
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var (
			w       = tabwriter.NewWriter(os.Stdout, 5, 1, 2, ' ', tabwriter.TabIndent)
			headers = append([]string{countHeader}, columns...)
		)
		if !excludeExtraFields {
			headers = append(headers, fieldsHeader)
		}
		if len(args) == 0 {
			err = FromStdin(w, headers...)
		} else {
			err = FromFile(w, args[0], headers...)
		}
		w.Flush()
		return
	},
}

func FromStdin(w io.Writer, headers ...string) error {
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) != 0 {
		return fmt.Errorf("nothing passed to chop")
	}

	var (
		scanner = bufio.NewScanner(os.Stdin)
		count   = 0
	)
	printHeader(w, headers...)
	for scanner.Scan() {
		if err := printLine(w, count, scanner.Text(), columns...); err != nil {
			return err
		}
		count++
	}

	return nil
}

func FromFile(w io.Writer, path string, headers ...string) error {
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

	var (
		scanner = bufio.NewScanner(file)
		count   = 0
	)
	printHeader(w, headers...)
	for scanner.Scan() {
		if err := printLine(w, count, scanner.Text(), columns...); err != nil {
			return err
		}
		count++
	}
	return nil
}

func printLine(w io.Writer, count int, line string, columns ...string) error {
	if beforeLine >= 0 && count > beforeLine {
		return nil
	}
	if count < afterLine {
		return nil
	}
	fieldMap := make(map[string]any)
	if !json.Valid([]byte(line)) {
		line = fmt.Sprintf("{%q:%q}", messageField, line)
	}
	if err := json.Unmarshal([]byte(line), &fieldMap); err != nil {
		return err
	}
	var message = fmt.Sprintf("%4d\t", count)
	for _, col := range columns {
		var (
			val string
			ok  bool
		)
		if val, ok = fieldMap[col].(string); !ok {
			message = fmt.Sprintf("%s --- \t", message)
			continue
		}
		if len(val) > maxColumnWidth && maxColumnWidth > 0 {
			val = fmt.Sprintf("%s...", string(val[:maxColumnWidth]))
		}
		message = fmt.Sprintf("%s%s\t", message, val)
		delete(fieldMap, col)
	}
	if !excludeExtraFields {
		if len(fieldMap) > 0 {
			message = fmt.Sprintf("%s%+v\t", message, fieldMap)
		} else {
			message = fmt.Sprintf("%s%+v\t", message, fieldMap)
		}
	}
	fmt.Fprintln(w, message)
	return nil
}

func printHeader(w io.Writer, columns ...string) {
	var line string
	for _, col := range columns {
		if line == "" {
			line = fmt.Sprintf("%s\t", col)
		} else {
			line = fmt.Sprintf("%s%s\t", line, col)
		}
	}
	fmt.Fprintf(w, "%s\n", line)
}

func replace[T comparable](l []T, old, new T) []T {
	for i, other := range l {
		if other == old {
			l[i] = new
		}
	}
	return l
}

func main() {
	rootCmd.Execute()
}

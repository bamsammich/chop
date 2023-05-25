package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

const extraFieldsName = "_fields"

var (
	formatTuples = []string{"message=60"}
	fieldOrder   []string
	format       = make(map[string]int)
	printExtras  bool
)

func newRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "chop [path]",
		Short: "Write structured logs in a human-readable way.",
		Args:  cobra.RangeArgs(0, 1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			for _, tuple := range formatTuples {
				parts := strings.Split(tuple, "=")
				width, err := strconv.Atoi(parts[1])
				if err != nil {
					return err
				}
				field := parts[0]
				fieldOrder = append(fieldOrder, field)
				if len(field) > width {
					width = len(field)
				}
				format[field] = width
			}
			if printExtras && !contains(fieldOrder, extraFieldsName) {
				fieldOrder = append(fieldOrder, extraFieldsName)
				format[extraFieldsName] = 50
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			printHeader(format)
			if len(args) == 0 {
				err = FromStdin()
			} else {
				err = FromFile(args[0])
			}
			return
		},
	}
	cmd.PersistentFlags().StringSliceVarP(&formatTuples, "format", "f", formatTuples, "tuples of field names to print and column width")
	cmd.PersistentFlags().BoolVarP(&printExtras, "print-all", "a", false, "print all fields; fields without format defined will be printed as JSON")

	return cmd
}

func FromStdin() error {
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) != 0 {
		return fmt.Errorf("nothing passed to chop")
	}

	var (
		scanner = bufio.NewScanner(os.Stdin)
		count   = 0
	)

	for scanner.Scan() {
		log, err := LogFromString(count, scanner.Text())
		if err != nil {
			return err
		}
		log.Print()
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
	defer file.Close()
	if err != nil {
		return err
	}

	var (
		scanner = bufio.NewScanner(file)
		count   = 0
	)
	for scanner.Scan() {
		log, err := LogFromString(count, scanner.Text())
		if err != nil {
			return err
		}
		log.Print()
		count++
	}
	return nil
}

type Log struct {
	Fields map[string]logField
	Number int
}

type logField struct {
	Value any
	Width int
}

func (lf logField) Format() string {
	return fmt.Sprintf("%*v", lf.Width, lf.Value)
}

func LogFromString(count int, line string) (*Log, error) {
	log := &Log{Number: count, Fields: map[string]logField{}}
	fields := make(map[string]any)
	extras := make(map[string]any)
	if !json.Valid([]byte(line)) {
		log.Fields["message"] = logField{Value: line, Width: format["message"]}
		return log, nil
	}
	if err := json.Unmarshal([]byte(line), &fields); err != nil {
		return nil, err
	}
	for k, v := range fields {
		if _, ok := format[k]; !ok {
			extras[k] = v
			continue
		}
		log.Fields[k] = logField{Value: v, Width: format[k]}
	}
	for k, v := range format {
		if _, ok := log.Fields[k]; !ok {
			log.Fields[k] = logField{Value: "", Width: v}
		}
	}
	if printExtras {
		b, err := json.Marshal(extras)
		if err != nil {
			return nil, err
		}
		log.Fields[extraFieldsName] = logField{Value: string(b), Width: format[extraFieldsName]}
	}
	return log, nil
}

func (l *Log) Print() {
	line := fmt.Sprintf("%6d", l.Number)
	for _, name := range fieldOrder {
		if _, ok := l.Fields[name]; ok {
			line = fmt.Sprintf("%s %s", line, l.Fields[name].Format())
			continue
		}
		line = fmt.Sprintf("%s %*s", line, format[name], "-")
	}
	fmt.Println(line)
}

func printHeader(format map[string]int) {
	line := fmt.Sprintf("%6s", " ")
	for _, name := range fieldOrder {
		line = fmt.Sprintf("%s %*s", line, format[name], name)
	}
	fmt.Println(line)
}

func main() {
	newRootCmd().Execute()
}

func contains[T comparable](sl []T, e T) bool {
	for _, i := range sl {
		if i == e {
			return true
		}
	}
	return false
}

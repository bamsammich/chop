package main

import (
	"bufio"
	"log"
	"os"

	"github.com/bamsammich/chop/config"
	"github.com/spf13/cobra"
)

// func init() {
// 	rootCmd.PersistentFlags().StringVarP(&messageField, "message-field", "m", defaultMessageField, "field containing log message")
// 	rootCmd.PersistentFlags().StringSliceVarP(&columns, "columns", "c", []string{defaultMessageField}, "field names to extract to columns")
// 	rootCmd.PersistentFlags().IntVarP(&maxColumnWidth, "max-column-width", "w", 60, "set maximum column width [0 unsets this]")
// 	rootCmd.PersistentFlags().IntVarP(&beforeLine, "before", "B", -1, "print lines before this count [-1 unsets this]")
// 	rootCmd.PersistentFlags().IntVarP(&afterLine, "after", "A", 0, "print lines after this count")
// 	rootCmd.PersistentFlags().BoolVarP(&excludeExtraFields, "exclude-extra-fields", "x", false, "exclude extra fields not defined by `columns`")
// 	rootCmd.PersistentFlags().SortFlags = true
// }

// var (
// 	beforeLine int
// 	afterLine  int

// 	maxColumnWidth int
// 	columns        []string

// 	countHeader         = "   #"
// 	fieldsHeader        = "fields"
// 	defaultMessageField = "message"
// 	messageField        string
// 	excludeExtraFields  bool
// )

// var rootCmd = &cobra.Command{
// 	Use:   "chop [path]",
// 	Short: "Write structured log in a human-readable way.",
// 	Args:  cobra.RangeArgs(0, 1),
// 	PreRun: func(cmd *cobra.Command, args []string) {
// 		if messageField != defaultMessageField {
// 			columns = replace(columns, defaultMessageField, messageField)
// 		}
// 	},
// 	RunE: func(cmd *cobra.Command, args []string) (err error) {
// 		var (
// 			tw      = tablewriter.NewWriter(os.Stdout)
// 			headers = append([]string{countHeader}, columns...)
// 		)

// 		tw.SetAutoWrapText(false)
// 		tw.SetAutoFormatHeaders(true)

// 		if !excludeExtraFields {
// 			headers = append(headers, fieldsHeader)
// 		}
// 		tw.SetHeader(headers)
// 		if len(args) == 0 {
// 			err = FromStdin(tw, headers...)
// 		} else {
// 			err = FromFile(tw, args[0], headers...)
// 		}
// 		tw.Render()
// 		return
// 	},
// }

// func FromStdin(tw *tablewriter.Table, headers ...string) error {
// 	stat, _ := os.Stdin.Stat()
// 	if (stat.Mode() & os.ModeCharDevice) != 0 {
// 		return fmt.Errorf("nothing passed to chop")
// 	}

// 	var (
// 		scanner = bufio.NewScanner(os.Stdin)
// 		count   = 0
// 	)

// 	for scanner.Scan() {
// 		if err := printLine(tw, count, scanner.Text(), columns...); err != nil {
// 			return err
// 		}
// 		count++
// 	}

// 	return nil
// }

// func FromFile(tw *tablewriter.Table, path string, headers ...string) error {
// 	fi, err := os.Stat(path)
// 	if err != nil {
// 		return err
// 	}
// 	if fi.IsDir() {
// 		return fmt.Errorf("must not be a directory")
// 	}
// 	file, err := os.Open(path)
// 	defer file.Close()
// 	if err != nil {
// 		return err
// 	}

// 	var (
// 		scanner = bufio.NewScanner(file)
// 		count   = 0
// 	)
// 	// printHeader(w, headers...)
// 	for scanner.Scan() {
// 		if err := printLine(tw, count, scanner.Text(), columns...); err != nil {
// 			return err
// 		}
// 		count++
// 	}
// 	return nil
// }

// func printLine(tw *tablewriter.Table, count int, line string, columns ...string) error {
// 	if beforeLine >= 0 && count > beforeLine {
// 		return nil
// 	}
// 	if count < afterLine {
// 		return nil
// 	}
// 	fieldMap := make(map[string]any)
// 	if !json.Valid([]byte(line)) {
// 		line = fmt.Sprintf("{%q:%q}", messageField, line)
// 	}
// 	if err := json.Unmarshal([]byte(line), &fieldMap); err != nil {
// 		return err
// 	}
// 	var row = []string{fmt.Sprintf("%4d", count)}
// 	for _, col := range columns {
// 		var (
// 			val string
// 			ok  bool
// 		)
// 		if val, ok = fieldMap[col].(string); !ok {
// 			row = append(row, " --- ")
// 			continue
// 		}
// 		if len(val) > maxColumnWidth && maxColumnWidth > 0 {
// 			val = fmt.Sprintf("%s...", string(val[:maxColumnWidth]))
// 		}
// 		row = append(row, val)
// 		delete(fieldMap, col)
// 	}
// 	if !excludeExtraFields {
// 		if len(fieldMap) > 0 {
// 			row = append(row, fmt.Sprintf("%+v", fieldMap))
// 		} else {
// 			row = append(row, " --- ")
// 		}
// 	}
// 	tw.Append(row)
// 	return nil
// }

// func printHeader(w io.Writer, columns ...string) {
// 	var line string
// 	for _, col := range columns {
// 		if line == "" {
// 			line = fmt.Sprintf("%s\t", col)
// 		} else {
// 			line = fmt.Sprintf("%s%s\t", line, col)
// 		}
// 	}
// 	fmt.Fprintf(w, "%s\n", line)
// }

// func replace[T comparable](l []T, old, new T) []T {
// 	for i, other := range l {
// 		if other == old {
// 			l[i] = new
// 		}
// 	}
// 	return l
// }

var path = "example/app.log"
var cfg *config.Config

func init() {
	var err error
	cfg, err = config.Load()
	if err != nil {
		log.Fatal(err)
	}
}

var rootCmd = &cobra.Command{}

func main() {
	// rootCmd.Execute()

	fi, err := os.Stat(path)
	if err != nil {
		log.Fatal(err)
	}
	if fi.IsDir() {
		log.Fatal("must not be a directory")
	}
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	var (
		scanner   = bufio.NewScanner(file)
		count     = 0
		formatter = cfg.Formats["default"]
	)

	formatter.PrintHeaders()
	for scanner.Scan() {
		if err := formatter.PrintLine(scanner.Text()); err != nil {
			log.Fatal(err)
		}
		count++
	}
}

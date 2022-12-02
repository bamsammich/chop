package config

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/creasty/defaults"
	"gopkg.in/go-playground/validator.v9"
	"gopkg.in/yaml.v3"
)

//go:embed default.yaml
var defaultRawYAML []byte

var defaultMessageField = "message"

type Config struct {
	Formats map[string]*Format `yaml:"formats" json:"formats" validate:"required,dive"`
}

type Column struct {
	// Name of the log field. If left empty, column name is assumed to be the field name.
	Field string `yaml:"field" json:"field"`
	// Header to print for column.
	Header string `yaml:"header" json:"header" validate:"required"`
	// Maximum column width
	MaxWidth  int `yaml:"maxWidth" json:"maxWidth" default:"60" validate:"gt=-1"`
	isMessage bool
}

func (col *Column) FieldName() string {
	name := col.Header
	if col.Field != "" {
		name = col.Field
	}
	return name
}

func (col *Column) FormatText(text string, wrap bool) []string {
	if col.MaxWidth == 0 {
		return []string{text}
	}
	if len(text) > col.MaxWidth && wrap {
		var multilines []string
		for len(text) > col.MaxWidth {
			multilines = append(multilines, text[:col.MaxWidth])
			text = text[col.MaxWidth:]
		}
		if len(text) > 0 {
			multilines = append(multilines, text)
		}

		return multilines
	}
	if len(text) > col.MaxWidth && col.MaxWidth > 0 && !wrap {
		return []string{fmt.Sprintf("%s...", string(text[:col.MaxWidth-3]))}
	}
	return []string{text}
}

type Format struct {
	LineNumbers      bool      `yaml:"lineNumbers" json:"line_numbers"`
	AutoWrapText     bool      `yaml:"autoWrapText" json:"auto_wrap_text"`
	PrintExtraFields bool      `yaml:"printExtraFields" json:"print_extra_fields"`
	ColumnSpacer     string    `yaml:"columnSpacer" json:"column_spacer" default:" "`
	Columns          []*Column `yaml:"columns" json:"columns" validate:"required,dive"`
}

func (f *Format) MessageColumn() *Column {
	for _, col := range f.Columns {
		if col.isMessage {
			return col
		}
	}
	return nil
}

func (f *Format) PrintHeaders() {
	var headers []string
	for _, col := range f.Columns {
		headers = append(headers,
			fmt.Sprintf("%-*s", col.MaxWidth, col.Header))
	}
	fmt.Println(strings.Join(headers, f.ColumnSpacer))
}

func (f *Format) PrintLine(line string) error {
	fields := make(map[string]any)
	if !json.Valid([]byte(line)) {
		messageField := defaultMessageField
		if msgCol := f.MessageColumn(); msgCol != nil {
			messageField = msgCol.FieldName()
		}
		line = fmt.Sprintf("{%q:%q}", messageField, line)
	}
	if err := json.Unmarshal([]byte(line), &fields); err != nil {
		return fmt.Errorf("failed to unmarshal JSON log: %w", err)
	}
	var (
		lines  = make([][]string, 1)
		cursor = 0
	)
	for _, col := range f.Columns {
		val, ok := fields[col.FieldName()]
		if !ok {
			val = strings.Repeat("-", col.MaxWidth/2)
		}
		fieldLines := col.FormatText(val.(string), f.AutoWrapText)
		for i := 0; i < len(fieldLines); i++ {
			line := fmt.Sprintf("%-*s", col.MaxWidth, fieldLines[i])
			if i == 0 {
				lines[i] = append(lines[i], line)
			} else {
				lines = append(lines, []string{strings.Repeat(" ", cursor), line})
			}
		}
		cursor += col.MaxWidth
		delete(fields, col.FieldName())
	}
	if f.PrintExtraFields {
		lines[0] = append(lines[0], fmt.Sprintf("%+v", fields))
	}
	for i := range lines {
		fmt.Println(strings.Join(lines[i], f.ColumnSpacer))
	}

	return nil
}

func Load(paths ...string) (*Config, error) {
	var cfg = Config{}
	if err := yaml.Unmarshal(defaultRawYAML, &cfg); err != nil {
		return nil, fmt.Errorf("failed to read default config: %w", err)
	}
	if err := defaults.Set(&cfg); err != nil {
		return nil, err
	}
	validate := validator.New()
	if err := validate.Struct(cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

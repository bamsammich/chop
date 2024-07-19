package log

import (
	"encoding/json"
	"fmt"
	"slices"
	"strings"

	"github.com/bamsammich/chop/internal/config"
)

type Formatter struct {
	fields    *Fields
	lastWidth int
}

func (f *Formatter) FromString(count int, line string) (string, error) {
	if !json.Valid([]byte(line)) {
		line = fmt.Sprintf(`{"%s":"%s"}`, config.DefaultField, line)
	}

	fields := make(map[string]any)
	if err := json.Unmarshal([]byte(line), &fields); err != nil {
		return "", err
	}

	var sb strings.Builder
	f.fields.AddMany(fields)

	lineWidth := f.fields.Width()
	if lineWidth > f.lastWidth {
		sb.WriteString("\n\n")
		sb.WriteString(f.Headers())
		f.lastWidth = lineWidth
	}

	for _, o := range f.fields.Order {
		sb.WriteString(f.format(o, fmt.Sprintf("%v", fields[o])))
	}

	return sb.String(), nil
}

func (f *Formatter) Headers() string {
	var (
		sb         strings.Builder
		dividerLen int
	)
	for _, o := range f.fields.Order {
		field, _ := f.fields.Get(o)
		h := f.format(o, o)
		if h == "" {
			continue
		}
		w := field.Width()
		dividerLen += w
		sb.WriteString(field.Format(o, w))
	}
	sb.WriteString("\n")
	sb.WriteString(strings.Repeat("-", dividerLen))
	sb.WriteString("\n")

	return sb.String()
}

func (f *Formatter) format(name string, value string) string {
	field := f.fields.Add(name, len(value)).Format(value, len(value))

	if slices.Contains(config.Exclude, name) || (len(config.Include) > 0 && !slices.Contains(config.Include, name)) {
		return ""
	}

	return field

}

func NewFormatter() *Formatter {
	fmtr := &Formatter{
		fields: NewFields(),
	}
	fmtr.fields.Add(config.DefaultField, 40)
	return fmtr
}

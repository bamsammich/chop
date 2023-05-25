package log

import (
	"encoding/json"
	"fmt"

	"github.com/bamsammich/chop/internal/config"
	"github.com/bamsammich/chop/internal/pointers"
)

type Formatter struct {
	order  []string
	widths map[string]int
}

func (f *Formatter) Add(name string, width int) {
	f.order = append(f.order, name)
	f.widths[name] = width
}

func (f *Formatter) format(field string, value any) string {
	strValue := fmt.Sprintf("%v", value)
	return fmt.Sprintf("%*s", f.widths[field], strValue)
}

func (f *Formatter) FromString(count int, line string) (*Log, error) {
	var (
		log    = &Log{Number: pointers.To(count), Fields: map[string]string{}, order: f.order}
		fields = make(map[string]any)
		extras = make(map[string]any)
	)
	if !json.Valid([]byte(line)) {
		log.Fields[config.DefaultField] = f.format(config.DefaultField, line)
	} else {
		if err := json.Unmarshal([]byte(line), &fields); err != nil {
			return nil, err
		}
		for k, v := range fields {
			if _, ok := f.widths[k]; !ok {
				extras[k] = v
				continue
			}
			log.Fields[k] = f.format(k, v)
		}
	}
	for k := range f.widths {
		if _, ok := log.Fields[k]; !ok {
			log.Fields[k] = f.format(k, "-")
		}
	}
	if config.PrintExtras {
		b, err := json.Marshal(extras)
		if err != nil {
			return nil, err
		}
		log.Fields[config.ExtraFieldsName] = f.format(config.ExtraFieldsName, string(b))
	}
	return log, nil
}

func (f *Formatter) Header() *Log {
	log := &Log{Fields: map[string]string{}, order: f.order}
	for _, field := range f.order {
		log.Fields[field] = f.format(field, field)
	}
	return log
}

func NewFormatter() *Formatter {
	return &Formatter{
		order:  []string{},
		widths: make(map[string]int),
	}
}

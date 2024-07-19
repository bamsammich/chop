package log

import (
	"fmt"
	"slices"
)

var widthBufferSize = 5

type Fields struct {
	fields map[string]*Field
	Order  []string
}

type Field struct {
	width int
}

func (f *Field) Width() int {
	return f.width
}

func (f *Field) Format(s string, lineWidth int) string {
	if len(s) > f.width {
		f.width = len(s) + widthBufferSize
	}
	return fmt.Sprintf(" %*s ", f.Width()-2, s)
}

func NewFields() *Fields {
	return &Fields{
		fields: make(map[string]*Field),
		Order:  make([]string, 0),
	}
}

func (f *Fields) Add(name string, width int) *Field {
	if _, ok := f.fields[name]; !ok {
		// add buffer to reduce header resizing
		f.fields[name] = &Field{width: width + widthBufferSize}
	}
	if !slices.Contains(f.Order, name) {
		f.Order = append(f.Order, name)
	}

	return f.fields[name]
}

func (f *Fields) Get(name string) (*Field, bool) {
	width, ok := f.fields[name]
	if !ok {
		return nil, false
	}
	return width, ok
}

func (f *Fields) AddMany(input map[string]any) {
	for k, v := range input {
		vs := fmt.Sprintf("%v", v)
		f.Add(k, len(vs))
	}
}

func (f *Fields) Width() int {
	var sum int
	for _, field := range f.fields {
		sum += field.Width()
	}
	return sum
}

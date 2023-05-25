package log

import "fmt"

type Log struct {
	order  []string
	Fields map[string]string
	Number *int
}

func (l *Log) Print() {
	var line string
	if l.Number != nil {
		line = fmt.Sprintf("%6d", *l.Number)
	}
	for _, name := range l.order {
		if _, ok := l.Fields[name]; ok {
			line = fmt.Sprintf("%s %s", line, l.Fields[name])
			continue
		}
		line = fmt.Sprintf("%s %s", line, "FOUND A FIELD WE DIDNT FORMAT")
	}
	fmt.Println(line)
}

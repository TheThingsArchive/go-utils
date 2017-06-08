package errors

import (
	"fmt"
	"regexp"
)

type ValueFormatter func(interface{}) string

func defaultValueFormatter(v interface{}) string {
	return fmt.Sprintf("%v", v)
}

// re finds names in the format string
var re = regexp.MustCompile("{(.+?)}")

// Format formats the values into the provided string
func Format(format string, values Attributes, formatValue ValueFormatter) string {
	return re.ReplaceAllStringFunc(format, func(name string) string {
		if values == nil {
			return "<nil>"
		}

		stripped := name[1 : len(name)-1]
		return formatValue(values[stripped])
	})
}

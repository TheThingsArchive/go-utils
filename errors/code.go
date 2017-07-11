package errors

import (
	"fmt"
	"strconv"
)

// String implmenents stringer
func (c Code) String() string {
	return fmt.Sprintf("%v", uint32(c))
}

// pareCode parses a string into a Code or returns 0 if the parse failed
func parseCode(str string) Code {
	code, err := strconv.Atoi(str)
	if err != nil {
		return Code(0)
	}
	return Code(code)
}

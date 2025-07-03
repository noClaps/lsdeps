package logger

import (
	"fmt"
	"os"
)

func Errorf(format string, a ...any) {
	fmt.Fprintf(os.Stderr, fmt.Sprintf("\033[31m%s\033[0m", format), a...)
}

package logger

import (
	"fmt"
	"os"
)

func Errorln(a ...any) {
	fmt.Fprint(os.Stderr, "\033[31mERROR:\033[0m ")
	fmt.Fprintln(os.Stderr, a...)
}

func Fatalln(a ...any) {
	fmt.Fprint(os.Stderr, "\033[31mERROR:\033[0m ")
	fmt.Fprintln(os.Stderr, a...)
	os.Exit(1)
}

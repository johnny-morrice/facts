package main

import (
	"fmt"
	"io"
)

func newline(w io.Writer) {
	fmt.Fprint(w, "\n")
}

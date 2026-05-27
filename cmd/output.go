package cmd

import (
	"fmt"
	"io"
)

type ErrorOutput struct {
	Command string
	Message string
	Cause   error
	Hint    string
}

func WriteError(w io.Writer, out ErrorOutput) {
	if out.Command != "" {
		fmt.Fprintf(w, "netdbg %s: %s\n", out.Command, out.Message)
	} else {
		fmt.Fprintf(w, "netdbg: %s\n", out.Message)
	}
	if out.Cause != nil {
		fmt.Fprintf(w, "cause: %v\n", out.Cause)
	}
	if out.Hint != "" {
		fmt.Fprintf(w, "hint: %s\n", out.Hint)
	}
}

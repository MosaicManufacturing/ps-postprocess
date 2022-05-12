package printerscript

import "strings"

func normalizeNewlines(input string) string {
	// replace CR LF \r\n (Windows) with LF \n (Unix)
	input = strings.ReplaceAll(input, "\r\n", "\n")
	// replace CF \r (Mac Classic) with LF \n (Unix)
	return strings.ReplaceAll(input, "\r", "\n")
}

func normalizeInput(input string) string {
	// remove directive if present
	if strings.HasPrefix(input, "@printerscript ") {
		input = input[strings.IndexRune(input, '\n'):]
	}
	// trim leading and trailing whitespace
	trimmed := strings.TrimSpace(input)
	// normalize all newlines to \n
	lfOnly := normalizeNewlines(trimmed)
	// add a trailing newline
	return lfOnly + "\n"
}

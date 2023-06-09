package printerscript

import "strings"

func normalizeNewlines(input string) string {
	// replace CR LF \r\n (Windows) with LF \n (Unix)
	input = strings.ReplaceAll(input, "\r\n", "\n")
	// replace CF \r (Mac Classic) with LF \n (Unix)
	return strings.ReplaceAll(input, "\r", "\n")
}

func Normalize(input string) string {
	// remove directive if present
	if strings.HasPrefix(input, "@printerscript ") {
		newlineIdx := strings.IndexRune(input, '\n')
		if newlineIdx < 0 {
			// empty script
			return "\n"
		} else {
			input = input[newlineIdx:]
		}
	}
	// trim leading and trailing whitespace
	trimmed := strings.TrimSpace(input)
	// normalize all newlines to \n
	lfOnly := normalizeNewlines(trimmed)
	// add a trailing newline
	return lfOnly + "\n"
}

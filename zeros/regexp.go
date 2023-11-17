package zeros

import "regexp"

var re = regexp.MustCompile("( [XYZEF]-?)\\.([0-9]+)")

// https://gist.github.com/elliotchance/d419395aa776d632d897
func replaceAllStringSubmatchFunc(str string, repl func([]string) string) string {
	result := ""
	lastIndex := 0

	for _, v := range re.FindAllSubmatchIndex([]byte(str), -1) {
		var groups []string
		for i := 0; i < len(v); i += 2 {
			groups = append(groups, str[v[i]:v[i+1]])
		}

		result += str[lastIndex:v[0]] + repl(groups)
		lastIndex = v[1]
	}

	return result + str[lastIndex:]
}

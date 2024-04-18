package embedded

import "strings"

type CodeSnippet []string

func (s CodeSnippet) String() string {
	return strings.Join([]string(s), "\n")
}

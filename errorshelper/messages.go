package errorshelper

import (
	"fmt"
	"strings"
)

func GenErrorMessage(global string, err error, i int, words []string, line int) string {
	s := ""
	for j := range i {
		s += words[j] + " "
	}
	l1 := len(s)
	s += words[i]
	for j := range len(words) - i - 1 {
		s += " " + words[j+i+1]
	}
	l2 := len(s)
	title := " " + global + " "
	displayLine := fmt.Sprintf(" (line %d)", line+1)
	after := ""
	maxErrorSize := 0
	for l := range strings.SplitSeq(err.Error(), "\n") {
		maxErrorSize = max(maxErrorSize, len(l))
	}
	size := max(maxErrorSize, l2+len(displayLine)) - len(title)
	if size > 0 {
		for n := range size {
			if n%2 == 0 {
				title = "=" + title
			} else {
				title += "="
			}
		}
	}
	for range len(title) {
		after += "="
	}
	s += "\n"
	if i == len(words)-1 && i != 0 {
		for range l2 - 1 {
			s += "-"
		}
		s += "^"
	} else {
		for range l1 {
			s += "-"
		}
		s += "^"
		for range l2 - l1 - 1 {
			s += "-"
		}
	}
	return fmt.Sprintf("%s\n%s%s\n\n%s\n%s\n\n", title, s, displayLine, err.Error(), after)
}

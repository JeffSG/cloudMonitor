package utils

import (
	"strings"
)

func ConcatStrings(items ...string) string {
	var sb strings.Builder
	for _, str := range items {
		sb.WriteString(str)
	}
	return sb.String()
}


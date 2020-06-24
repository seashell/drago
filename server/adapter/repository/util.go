package repository

import "strings"

func commaSeparatedStrToSlice(s string) []string {
	return strings.Fields(strings.Replace(s, ",", " ", -1))
}

package string

import "strings"

// SliceToString : Takes a slice containing strings and returns a comma-separated string
func SliceToString(s []string) string {
	return strings.Join(s, ",")
}

// StringToSlice : Takes a comma-separated string and returns a slice
func StringToSlice(s string) []string {
	return strings.Fields(strings.Replace(s, ",", " ", -1))
}

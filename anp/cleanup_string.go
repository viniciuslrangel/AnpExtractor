package anp

import "strings"

func cleanupName(city string) string {
	return strings.ReplaceAll(city, "@", " ")
}

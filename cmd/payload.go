package cmd

import "strings"

func looksLikePayloadArg(value string) bool {
	value = strings.TrimSpace(value)
	return strings.HasSuffix(strings.ToLower(value), ".json")
}

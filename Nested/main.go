package main

import (
	"strings"
)

func removeProfanity(message *string) {
	// ?
	dic := map[string]string{
		"fubb":  "****",
		"shiz":  "****",
		"witch": "*****",
	}
	for key, val := range dic {
		*message = strings.ReplaceAll(*message, key, val)

	}
}

package main

import (
	"strings"
)

type sms struct {
	id      string
	content string
	tags    []string
}

func tagMessages(messages []sms, tagger func(sms) []string) []sms {
	// ?
	for i, message := range messages {
		messages[i].tags = tagger(message)
	}
	return messages
}

func tagger(msg sms) []string {
	tags := []string{}
	tagList := map[string]bool{
		"Urgent": false,
		"Promo":  false,
	}
	for word := range strings.SplitSeq(msg.content, " ") {
		word = strings.ToLower(word)
		if strings.Contains(word, "urgent") {
			tagList["Urgent"] = true
		}
		if strings.Contains(word, "sale") {
			tagList["Promo"] = true
		}
	}
	for s, v := range tagList {
		if v {
			tags = append(tags, s)
		}
	}
	return tags
}

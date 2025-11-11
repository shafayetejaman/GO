package main

import (
	"errors"
)

func validateStatus(status string) error {
	// ?
	n := len(status)
	if n > 140 {
		return errors.New("status exceeds 140 characters")
	}

	if n == 0 {
		return errors.New("status cannot be empty")
	}
	return nil
}

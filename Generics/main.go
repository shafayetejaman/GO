package main

import (
	"encoding/json"
)

func marshalAll[T any](items []T) ([][]byte, error) {
	data := [][]byte{}
	for item := range items {
		byteData, err := json.Marshal(item)
		if err != nil {
			return nil, err
		}
		data = append(data, byteData)
	}

	return data, nil
}

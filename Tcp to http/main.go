package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

func getLinesChannel(f io.ReadCloser) <-chan string {
	channel := make(chan string)

	go func() {
		defer f.Close()

		var s string

		for {
			newBytes := make([]byte, 8)
			_, err := f.Read(newBytes)
			if err != nil {
				break
			}
			s += string(newBytes)
			if i := strings.Index(s, "\n"); i != -1 {

				channel <- s[:i]
				s = s[i+1:]
			}
			// time.Sleep(200 * time.Millisecond)
		}

		if len(s) != 0 {
			channel <- s
		}
		close(channel)
	}()

	return channel

}
func main() {
	file, err := os.Open("message.txt")

	if err != nil {
		log.Fatal(err)
	}
	channel := getLinesChannel(file)

	for line := range channel {
		fmt.Printf("read: %s\n", line)
	}

}

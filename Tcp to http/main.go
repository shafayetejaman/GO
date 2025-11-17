package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

func getLinesChannel(f io.ReadCloser) <-chan string {
	channel := make(chan string)

	go func() {
		defer f.Close()
		defer close(channel)

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
	}()

	return channel

}
func main() {
	listener, err := net.Listen("tcp", ":42069")
	defer listener.Close()

	if err != nil {
		log.Fatal(err)
	}

	for {
		cnn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("connection has been accepted!")
		channel := getLinesChannel(cnn)
		for line := range channel {
			fmt.Printf("read: %s\n", line)
		}
		cnn.Close()
	}

}

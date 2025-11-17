package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	// fmt.Println("I hope I get the job!")
	file, err := os.Open("message.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	var s string
	for {
		newBytes := make([]byte, 8)
		_, err := file.Read(newBytes)
		if err != nil {
			break
		}
		s += string(newBytes)
		if i := strings.Index(s, "\n"); i != -1 {

			fmt.Printf("read : %s\n", s[:i])
			s = s[i+1:]
		}
	}

	// if len(s) != 0 {
	// 	fmt.Printf("Read : %s\n", s)
	// }

}

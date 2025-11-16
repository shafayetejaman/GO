package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("I hope I get the job!")
	file, err := os.Open("message.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	newBytes := make([]byte, 8)
	data := make([]byte, 0)

	for {
		_, err := file.Read(newBytes)
		if err != nil {
			fmt.Println(err)
			break
		}
		data = append(data, newBytes...)
		fmt.Println(string(newBytes))
	}

	fmt.Print(string(data))

}

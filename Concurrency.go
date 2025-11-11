package main

import (
	"fmt"
	"sync"
	"time"
)

func sendEmail(message string) {
	time.Sleep(time.Millisecond * 250)
	fmt.Printf("Email received: '%s'\n", message)
	fmt.Printf("Email sent: '%s'\n", message)
}

func test(wg *sync.WaitGroup, message string, name string) {
	if wg != nil {
		defer wg.Done()
	}
	fmt.Println("sending email to: ", name)
	sendEmail(message)
	time.Sleep(time.Millisecond * 500)
	fmt.Println("========================")
}

func main() {
	start := time.Now()

	test(nil, "Hello there Kaladin!", "akash")
	test(nil, "Hi there Shallan!", "batash")
	test(nil, "Hey there Dalinar!", "abir")

	fmt.Printf("Total time: %v\n", time.Since(start))

	start = time.Now()

	var wg sync.WaitGroup
	wg.Add(3) // number of goroutines

	go test(&wg, "Hello there Kaladin!", "akashgo")
	go test(&wg, "Hi there Shallan!", "batashgo")
	go test(&wg, "Hey there Dalinar!", "abirgo")

	wg.Wait() // wait for all goroutines
	fmt.Printf("Total time: %v\n", time.Since(start))
}

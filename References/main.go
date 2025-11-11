package main

import "unicode/utf8"

func getNameCounts(names []string) map[rune]map[string]int {
	// Your code here
	dic := map[rune]map[string]int{}
	for _, name := range names {
		firstCar, _ := utf8.DecodeRuneInString(name)
		if _, ok := dic[firstCar]; ok {
			dic[firstCar][name]++
		} else {
			dic[firstCar] = map[string]int{name: 1}
		}
	}
	return dic
}

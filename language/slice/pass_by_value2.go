package main

import "fmt"

func modifySlice(s []int) {
	s = append(s, 2048)
	s[0] = 1024
}

func main() {
	var s []int
	for i := 0; i < 3; i++ {
		s = append(s, i)}
	modifySlice(s)
	fmt.Println(s)
}

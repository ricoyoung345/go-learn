package main

import "fmt"

func main() {
	m := make(map[int]int)
	modifyMap(m)
	fmt.Println(m)
}

func modifyMap(m map[int]int) {
	m[1] = 1
	m[2] = 2
}

package main

func main() {
	var s []int  // 0, 0
	s = append(s, 0) // 1, 1
	println("len:", len(s), ", cap:", cap(s))

	s = append(s, 1) // 2, 2
	println("len:", len(s), ", cap:", cap(s))
	s = append(s, 2) // 3, 4
	println("len:", len(s), ", cap:", cap(s))
	s = append(s, 3) // 4, 4
	println("len:", len(s), ", cap:", cap(s))
	s = append(s, 4) // 5, 8
	println("len:", len(s), ", cap:", cap(s))

	for i := 5; i < 1025; i++ {
		s = append(s, i)
	} // 1025, 1280
	println("len:", len(s), ", cap:", cap(s))

	s = []int{} // 0, 0
	println("len:", len(s), ", cap:", cap(s))
	s = append(s, 0, 1, 2) // 3, 3
	println("len:", len(s), ", cap:", cap(s))
	s = append(s, 3) // 4, 6
	println("len:", len(s), ", cap:", cap(s))
	s = append(s, 4) // 5, 6
	println("len:", len(s), ", cap:", cap(s))
	s = append(s, 5, 6, 7) // 8, 12
	println("len:", len(s), ", cap:", cap(s))
}

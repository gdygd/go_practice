package main

import "fmt"

func main() {
	a := []int{1, 2, 3}
	b := a[:0]

	b = append(b, 9) // b = [9], a는 [9, 2, 3]

	fmt.Printf("a: %v \n", a) //a: [9 2 3]
	fmt.Printf("b: %v \n", b) //b: [9]

}

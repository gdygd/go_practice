package main

import "fmt"

func main() {
	a := []int{1, 2, 3}
	b := a[:0]

	b = append(b, 9) // b = [9], a는 [9, 2, 3]

	fmt.Printf("a: %v \n", a) //a: [9 2 3]
	fmt.Printf("b: %v \n", b) //b: [9]

	fmt.Println("===============================================")

	a2 := []int{1, 2, 3, 4}
	b2 := a2[:0] // 길이는 0, cap은 4

	fmt.Printf("a2: %v, len: %d, cap: %d\n", a2, len(a2), cap(a2))
	fmt.Printf("b2: %v, len: %d, cap: %d\n", b2, len(b2), cap(b2))

	b2 = append(b2, 99)
	fmt.Println("b after append:", b2)
	fmt.Println("a after append:", a2)

}

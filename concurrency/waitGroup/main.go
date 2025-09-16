package main

import (
	"fmt"
	"sync"
)

func printSomething(s string, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println(s)
}

func test1() {
	words := []string{
		"aaa",
		"bbb",
		"ccc",
		"ddd",
		"eee",
		"fff",
		"ggg",
		"hhh",
		"iii",
	}

	var wg sync.WaitGroup

	wg.Add(len(words))

	for i, word := range words {
		go printSomething(fmt.Sprintf("%d : %s", i, word), &wg)
	}

	wg.Wait()
	wg.Add(1)
	printSomething("End..", &wg)

}
func main() {
	// test1()
	waitGrupTest()

}

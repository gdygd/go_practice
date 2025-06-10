package main

import (
	"fmt"
	"sync"
)

func sum(n int) int {
	s := 0
	for i := 0; i < n; i++ {
		s += i
	}
	return s
}

func main() {
	wg := sync.WaitGroup{}
	for n := 0; n < 5; n++ {

		wg.Add(1)
		go func(n int) {
			s := sum(n)
			fmt.Printf("sum(%d)= %d\n", n, s)
			wg.Done()
		}(n)
	}
	wg.Wait()

}

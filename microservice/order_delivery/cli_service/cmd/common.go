package cmd

import "fmt"

func printLine1(ch string, length int) {
	for i := 0; i < length; i++ {
		fmt.Printf("%s", ch)
	}
	fmt.Printf("\n")
}

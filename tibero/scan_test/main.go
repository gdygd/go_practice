// You can edit this code!
// Click here and start typing.
package main

import "fmt"

func main() {
	a := make([]string, 100)
	ttt := make([]interface{}, 100)

	for i, _ := range ttt {
		ttt[i] = &a[i]
	}

	for {
		fmt.Print(">>>> ")

		count, _ := fmt.Scanln(ttt...)
		fmt.Println("[count] : ", count)
		fmt.Println("[a] : ", a)

		if count == 0 {
			continue
		}
	}

}

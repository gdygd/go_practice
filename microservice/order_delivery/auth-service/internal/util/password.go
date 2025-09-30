package util

import "fmt"

func CheckPassword(reqPw string, userPw string) error {
	if reqPw == userPw {
		return nil
	}

	return fmt.Errorf("Invalid password...")
}

package main

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

func checkToken(strToken string) {

	hmacSampleSecret := []byte("AllYourBase")
	tokenString := strToken

	// Parse takes the token string and a function for looking up the key. The latter is especially
	// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
	// head of the token to identify which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return hmacSampleSecret, nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println(claims)
	} else {
		fmt.Println(err)
	}

}

func main() {
	mySigningKey := []byte("AllYourBase")

	type MyCustomClaims struct {
		Name string `json:"name"`
		jwt.StandardClaims
	}

	// Create the Claims
	claims := MyCustomClaims{
		"bar",
		jwt.StandardClaims{
			//ExpiresAt: time.Now().Add(time.Second * 60).Unix(),
			Issuer: "test",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(mySigningKey)

	fmt.Printf("%d\n", len(ss))
	fmt.Printf("%d\n", len("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"))

	fmt.Printf("%v %v\n", ss, err)

	checkToken(ss)
	//checkToken("userid,token-value")

	time.Sleep(time.Second * 5)
	checkToken(ss)
}

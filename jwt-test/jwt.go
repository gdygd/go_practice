package main

import (
	"fmt"

	jwt "github.com/dgrijalva/jwt-go"
)

func generateToken(name string) (string, error) {
	var mySigningKey = []byte("keymaker")
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["name"] = name
	claims["role"] = "redpill"
	//claims["exp"] = time.Now().Add(time.Minute * 1).Unix()
	//claims["exp"] = 9999

	//tm := time.Now().Add(time.Minute * 30).Unix()
	//log.Printf("exptm : %v, exp : %v,   %v", tm, claims["exp"], token)

	tokenString, err := token.SignedString(mySigningKey)

	if err != nil {
		fmt.Errorf("Something Went Wrong: %s", err.Error())
		return "", err
	}
	return tokenString, nil
}

func verifyToken(tokenString string) (jwt.Claims, error) {
	signingKey := []byte("keymaker")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return signingKey, nil
	})
	if err != nil {
		return nil, err
	}
	return token.Claims, err
}

func verifyToken2(tokenString string) (*jwt.Token, error) {
	var mySigningKey = []byte("keymaker")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("There was an error in parsing")
		}
		return mySigningKey, nil
	})

	return token, err
}

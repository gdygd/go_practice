package main

import (
	"encoding/json"
	"log"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
)

//------------------------------------------------------------------------------
// AuthMiddleware
//------------------------------------------------------------------------------
func AuthMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		reqUrl := r.Host
		reqPath := r.URL.Path

		log.Printf("Req url : %s,%s", reqUrl, reqPath)

		if reqPath == "/signin" {
			log.Print("url path is /signin..")
			next.ServeHTTP(w, r)
		} else {
			tokenString := r.Header.Get("Authorization")
			log.Print("header token : ", tokenString)
			if len(tokenString) == 0 {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Missing Authorization Header"))
				return
			}
			//tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
			token, err := verifyToken2(tokenString)

			log.Printf("token : %v", token)

			if err != nil {
				log.Printf("AuthMiddleware err  : %v", err)
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Error verifying JWT token: " + err.Error()))
				return
			}

			//--------------------------------------------
			claims, _ := token.Claims.(jwt.MapClaims)

			// name := claims.(jwt.MapClaims)["name"].(string)
			// role := claims.(jwt.MapClaims)["role"].(string)
			// exp := claims.(jwt.MapClaims)["exp"].(string)

			name := claims["name"]
			role := claims["role"]
			//exp := claims["exp"]

			//expTime, _ := claims["exp"].(int64)

			log.Printf("name : %v", name)
			log.Printf("role : %v", role)
			//log.Printf("exp : %v", exp)
			//log.Printf("expTime : %v", expTime)

			// r.Header.Set("name", name)
			// r.Header.Set("role", role)
			// r.Header.Set("exp", exp)

			next.ServeHTTP(w, r)
		}

	})
}

func SignUp(w http.ResponseWriter, r *http.Request) {

	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		//var err Error
		//err = SetError(err, "Error in reading body")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(err)
		return
	}

	log.Print("Sign up info : ", user.Email, user.Password)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func SignIn(w http.ResponseWriter, r *http.Request) {

	log.Print("SignIn..")

	var authdetails Authentication
	err := json.NewDecoder(r.Body).Decode(&authdetails)
	if err != nil {
		log.Print("SignIn err..", err)
		//var err Error
		//err = SetError(err, "Error in reading body")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(err)
		return
	}

	var authuser User

	log.Print("Sign in info : ", authdetails.Email, authdetails.Password)

	// if authuser.Email == "" {
	// 	var err Error
	// 	err = SetError(err, "Username or Password is incorrect")
	// 	w.Header().Set("Content-Type", "application/json")
	// 	json.NewEncoder(w).Encode(err)
	// 	return
	// }

	// check := CheckPasswordHash(authdetails.Password, authuser.Password)

	// if !check {
	// 	var err Error
	// 	err = SetError(err, "Username or Password is incorrect")
	// 	w.Header().Set("Content-Type", "application/json")
	// 	json.NewEncoder(w).Encode(err)
	// 	return
	// }

	// validToken, err := GenerateJWT(authuser.Email, authuser.Role)
	validToken, err := generateToken(authdetails.Email)
	// if err != nil {
	// 	var err Error
	// 	err = SetError(err, "Failed to generate token")
	// 	w.Header().Set("Content-Type", "application/json")
	// 	json.NewEncoder(w).Encode(err)
	// 	return
	// }

	var token Token
	token.Email = authuser.Email
	token.Role = authuser.Role
	token.Role = authuser.Role
	token.Exp = authuser.Exp
	token.TokenString = validToken
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(token)
}

func Apitest(w http.ResponseWriter, r *http.Request) {
	log.Print("Apitest..")

	w.Write([]byte("Apitest ok"))

}

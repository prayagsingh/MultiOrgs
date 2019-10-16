package rest

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
)

// JwtToken for storing tokenString
type JwtToken struct {
	Token string
}

// Create the JWT key used to create the signature
var jwtKey = []byte("my_secret_key")

var savedToken = make(map[string]string)

func (app *RestApp) processAuthentication(w http.ResponseWriter, key string) string {

	validToken, err := GetToken(key)
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	client := &http.Client{}
	r, _ := http.NewRequest("GET", "http://localhost:"+PORT+"/token_auth", nil)
	r.Header.Set("Token", validToken)

	_, err = client.Do(r)
	if err != nil {
		respondJSON(w, map[string]string{"error": "Auth Error - " + err.Error()})
	}

	savedToken["token"] = validToken

	return validToken
}

// GetToken : generating JWT
func GetToken(email string) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"authorized": true,
		"user":       email,
		"exp":        time.Now().Add(time.Minute * 30).Unix(),
	})

	// Get the complete signed token
	tokenString, err := token.SignedString(jwtKey)

	if err != nil {
		fmt.Printf("\nsomething went wrong: %s", err.Error())
		return "", err
	}

	return tokenString, nil

}

func (app *RestApp) hasSavedToken(endpoint func(http.ResponseWriter, *http.Request, string)) func(http.ResponseWriter, *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		if savedToken["token"] != "" {
			token := savedToken["token"]
			endpoint(w, r, token)
		} else {
			respondJSON(w, map[string]string{"error": "Not Authorized"})
		}
	}
}

func (app *RestApp) isAuthorized(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

		type Exception struct {
			Message string
		}

		fmt.Println("\n Value of req.header is: ", req.Header, "\n Header[Token] is: ", req.Header.Get("authorization"))
		autorizationHeader := req.Header.Get("authorization")
		fmt.Println("\n Value of autorization header is: ", autorizationHeader)

		if autorizationHeader != " " {
			bearerToken := strings.Split(autorizationHeader, " ")
			fmt.Println("\n bearer token is: ", bearerToken[1])
			token, err := jwt.Parse(bearerToken[1], func(token *jwt.Token) (interface{}, error) {

				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("There was an error")
				}
				return jwtKey, nil
			})
			if err != nil {
				json.NewEncoder(w).Encode(Exception{Message: "Error in parsing jwt.Token"})
				return
			}
			if token.Valid {
				log.Println("TOKEN WAS VALID")
				context.Set(req, "decoded", token.Claims)
				next(w, req)
			} else {
				json.NewEncoder(w).Encode(Exception{Message: "Invalid authorization token"})
			}
		} else {
			json.NewEncoder(w).Encode(Exception{Message: "An authorization header is required"})
		}
	})
}

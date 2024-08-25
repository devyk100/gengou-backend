package auth

import (
	"context"
	"fmt"
	"gengou-main-backend/internals/database"
	"gengou-main-backend/internals/redis"
	"github.com/golang-jwt/jwt/v4"
	"log"
	"net/http"
	"os"
	"time"
)

//type AuthenticationBody struct {
//	Token string `json:"token"`
//}

const (
	Authorized = "Authorized"
)

func AuthenticateUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//var body AuthenticationBody
		//err := json.NewDecoder(r.Body).Decode(&body)
		//if err != nil {
		//	fmt.Println(err.Error())
		//}
		//fmt.Println(body.Token)
		//fmt.Println(r.Header)
		//fmt.Println(r.Host)
		//fmt.Println(r.URL)
		//fmt.Println(r.Body)

		now := time.Now()

		tokenString := r.Header.Get("Authorization")
		secretKey := []byte(os.Getenv("GLOBAL_AUTH_SECRET"))
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return secretKey, nil
		})
		if err != nil {
			log.Fatalf("Failed to parse token: %v", err)
		}

		var userId string
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			fmt.Println("Token is valid")
			userId, ok = claims["userId"].(string)
			if !ok {
				log.Fatalf("Claim 'userId' is not of type string or does not exist")
			}
			fmt.Printf("User ID: %s\n", userId)
		} else {
			fmt.Println("Token is not valid")
		}

		ctx := context.WithValue(r.Context(), "userIdString", userId)

		val, err := redis.Instance.Get(userId)
		fmt.Println("The redis value is ", val)
		if err != nil {
			return
		}

		if val != Authorized {
			_, err = database.Queries.GetAUserWithUserId(context.Background(), userId)
			if err != nil {
				fmt.Printf("Failed to get user: %v", err.Error())
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			err = redis.Instance.Set(userId, Authorized, time.Minute*5)
			if err != nil {
				fmt.Println("Failed to set authorised", err.Error())
				return
			}
		}

		log.Println("Inside the authorisation middleware, took", time.Since(now))

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

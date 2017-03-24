package web

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/btshopng/btshopng/messages"
)

// Middlewares

//RecoverHandler catches all panics, so the serverdoesnt go down ocmpletely, jsust because of a panic, that could be in one handler request by one user, affecting every other user.
func RecoverHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic: %+v", err)
				messages.WriteError(w, messages.ErrInternalServer)
			}
		}()

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

//LoggingHandler Logs request time, method and duration of handler/request execution
func LoggingHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		t1 := time.Now()
		next.ServeHTTP(w, r)
		t2 := time.Now()
		log.Printf("[%s] %q %v\n", r.Method, r.URL.String(), t2.Sub(t1))
	}

	return http.HandlerFunc(fn)
}

func GetUserInfoFromToken(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		tokenValue := r.Header.Get("authorization")
		client := &http.Client{}

		param := make(map[string]string)
		param["id_token"] = strings.Split(tokenValue, " ")[1]
		log.Println(param)

		b := new(bytes.Buffer)
		err := json.NewEncoder(b).Encode(param)
		if err != nil {
			log.Println(err)
		}

		req, err := http.NewRequest("POST", "https://emikra.auth0.com/tokeninfo", b)
		if err != nil {
			log.Println(err)
		}

		req.Header.Set("Content-type", "application/json")

		res, err := client.Do(req)
		if err != nil {
			log.Println(err)
		}

		responseMap := make(map[string]interface{})
		err = json.NewDecoder(res.Body).Decode(&responseMap)
		if err != nil {
			log.Println(err)
		}
		log.Println(responseMap)

		ctx := context.WithValue(r.Context(), "username", responseMap["username"].(string))
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

//FrontAuthHandler checks if authenticated, and moves the user details from the JWT token into the request context
// func FrontAuthHandler(next http.Handler) http.Handler {
// 	ac := config.Get()
// 	fn := func(w http.ResponseWriter, r *http.Request) {
//
// 		// check if we have a cookie with out tokenName
// 		tokenValue := r.Header.Get("X-AUTH-TOKEN")
// 		//log.Println(tokenValue)
//
// 		// validate the token
// 		token, err := jwt.Parse(tokenValue, func(token *jwt.Token) (interface{}, error) {
// 			publicKey, err := jwt.ParseRSAPublicKeyFromPEM(ac.Encryption.Public)
//
// 			if err != nil {
// 				return publicKey, err
// 			}
// 			return publicKey, nil
// 		})
//
// 		// branch out into the possible error from signing
// 		switch err.(type) {
//
// 		case nil: // no error
// 			if !token.Valid { // but may still be invalid
// 				log.Println(err)
//
// 				messages.WriteError(w, messages.ErrBadToken)
// 				return
// 			}
//
// 			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
// 				ctx := context.WithValue(r.Context(), "User", claims["User"])
// 				r = r.WithContext(ctx)
// 				//context.Set(r, "User", claims["User"])
// 				//context.Set(r, "UserID", claims["UserID"])
// 			} else {
// 				log.Println(err)
// 				return
// 			}
//
// 			next.ServeHTTP(w, r)
//
// 		case *jwt.ValidationError: // something was wrong during the validation
// 			vErr := err.(*jwt.ValidationError)
//
// 			switch vErr.Errors {
// 			case jwt.ValidationErrorExpired:
// 				messages.WriteError(w, messages.ErrBadToken)
// 				return
// 			default:
// 				messages.WriteError(w, messages.ErrBadToken)
// 				log.Printf("ValidationError error: %+v\n", vErr.Errors)
// 				return
// 			}
//
// 		default: // something else went wrong
// 			messages.WriteError(w, messages.ErrBadToken)
// 			return
// 		}
// 	}
// 	return http.HandlerFunc(fn)
//
// }

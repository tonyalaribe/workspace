package web

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"
)

// Middlewares

//RecoverHandler catches all panics, so the serverdoesnt go down ocmpletely, jsust because of a panic, that could be in one handler request by one user, affecting every other user.
func RecoverHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic: %+v", err)
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

		/* map[string]interface {}{"picture":"https://s.gravatar.com/avatar/b87d99b15e2a5bc064e7a0ad44cf1e24?s=480&r=pg&d=https%3A%2F%2Fcdn.auth0.com%2Favatars%2Fan.png", "nickname":"tonyalaribe",
		    "updated_at":"2017-05-31T14:19:06.627Z",
		    "identities":[]interface {}{map[string]interface {}{
		    	"user_id":"58d3e3eec1002318c624c24d",
		   	"provider":"auth0",
		   	"connection":"Username-Password-Authentication",
		   	"isSocial":false
		   	}},
		   	"email":"anthonyalaribe@gmail.com",
		   	"name":"anthonyalaribe@gmail.com",
		   	"email_verified":true,
		   	"clientID":"yqZpzeiFgoapsnpczQHIz0t6XoZjvEjL", "user_id":"auth0|58d3e3eec1002318c624c24d", "created_at":"2017-03-23T15:04:14.544Z", "global_client_id":"KG2EO2ZSAH0qxSfaRB0Ru61sTKoFCqeF", "username":"tonyalaribe"}*/

		username := responseMap["username"].(string)
		user, err := User{}.Get(username)

		if err != nil {
			log.Println(err)
			//User doesnt exist, so create the user in the local store
			user.Username = responseMap["username"].(string)
			user.Email = responseMap["email"].(string)
			user.Name = responseMap["name"].(string)
			user.ProviderUserID = responseMap["user_id"].(string)
			err = user.Create()
			if err != nil {
				log.Println(err)
			}
		}

		ctx := context.WithValue(r.Context(), "username", username)
		ctx2 := context.WithValue(r.Context(), "user", user)

		r = r.WithContext(ctx)
		r = r.WithContext(ctx2)

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

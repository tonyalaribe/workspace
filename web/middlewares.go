package web

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"gitlab.com/middlefront/workspace/config"

	"github.com/Jeffail/gabs"
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

		jsonDecoder := json.NewDecoder(res.Body)
		responseObject, err := gabs.ParseJSONDecoder(jsonDecoder)
		if err != nil {
			log.Println(err)
		}

		username := responseObject.Path("username").Data().(string)

		db := config.Get().Database
		user, err := db.GetUser(username)
		if err != nil {
			log.Println(err)
			//User doesnt exist, so create the user in the local store
			user.Username = responseObject.Path("username").Data().(string)
			user.Email = responseObject.Path("email").Data().(string)
			user.Name = responseObject.Path("name").Data().(string)
			user.ProviderUserID = responseObject.Path("user_id").Data().(string)
			if responseObject.ExistsP("app_metadata.roles") {
				for _, role := range responseObject.Path("app_metadata.roles").Data().([]interface{}) {
					user.Roles = append(user.Roles, role.(string))
				}
			}
			log.Printf("%#v", user)
			err = db.CreateUser(user)
			if err != nil {
				log.Println(err)
			}
		}

		r = r.WithContext(context.WithValue(r.Context(), "user", user))
		r = r.WithContext(context.WithValue(r.Context(), "username", user.Username))

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

package web

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"github.com/rs/cors"
	"gitlab.com/middlefront/workspace/config"
)

// Router struct would carry the httprouter instance, so its methods could be verwritten and replaced with methds with wraphandler
type Router struct {
	*httprouter.Router
}

// Get is an endpoint to only accept requests of method GET
func (r *Router) Get(path string, handler http.Handler) {
	r.GET(path, wrapHandler(handler))
}

// Post is an endpoint to only accept requests of method POST
func (r *Router) Post(path string, handler http.Handler) {
	r.POST(path, wrapHandler(handler))
}

// Put is an endpoint to only accept requests of method PUT
func (r *Router) Put(path string, handler http.Handler) {
	r.PUT(path, wrapHandler(handler))
}

// Delete is an endpoint to only accept requests of method DELETE
func (r *Router) Delete(path string, handler http.Handler) {
	r.DELETE(path, wrapHandler(handler))
}

// NewRouter is a wrapper that makes the httprouter struct a child of the router struct
func NewRouter() *Router {
	return &Router{httprouter.New()}
}

func wrapHandler(h http.Handler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		ctx := context.WithValue(r.Context(), "params", ps)
		r = r.WithContext(ctx)
		h.ServeHTTP(w, r)
	}
}

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func App() {
	authMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			secret := []byte(config.Get().Auth0ClientSecret)

			if len(secret) == 0 {
				log.Fatal("AUTH0_CLIENT_SECRET is not set")
			}
			return secret, nil
		},
	})

	commonHandlers := alice.New(LoggingHandler)
	//RecoverHandler
	router := NewRouter()

	router.Post("/api/new_workspace", commonHandlers.Append(authMiddleware.Handler, GetUserInfoFromToken).ThenFunc(CreateWorkspaceHandler))
	router.Get("/api/workspaces", commonHandlers.Append(authMiddleware.Handler, GetUserInfoFromToken).ThenFunc(GetWorkspacesHandler))
	router.Get("/api/workspaces/:workspaceID", commonHandlers.Append(authMiddleware.Handler, GetUserInfoFromToken).ThenFunc(GetWorkspaceBySlugHandler))

	router.Post("/api/workspaces/:workspaceID/permissions", commonHandlers.Append(authMiddleware.Handler, GetUserInfoFromToken).ThenFunc(ChangeUserWorkspacePermission))
	router.Get("/api/users_and_workspaces", commonHandlers.ThenFunc(UsersAndWorkspaceRoles))
	router.Get("/api/users_in_workspace", commonHandlers.ThenFunc(GetWorkspaceUsersAndRolesHandler))

	router.Get("/api/workspaces/:workspaceID/forms", commonHandlers.Append(authMiddleware.Handler, GetUserInfoFromToken).ThenFunc(GetFormsHandler))
	router.Post("/api/workspaces/:workspaceID/new_form", commonHandlers.Append(authMiddleware.Handler, GetUserInfoFromToken).ThenFunc(CreateFormHandler))
	router.Get("/api/workspaces/:workspaceID/forms/:formID", commonHandlers.Append(authMiddleware.Handler, GetUserInfoFromToken).ThenFunc(GetFormBySlugHandler))

	router.Post("/api/workspaces/:workspaceID/forms/:formID/new_submission", commonHandlers.Append(authMiddleware.Handler, GetUserInfoFromToken).ThenFunc(NewFormSubmissionHandler))

	router.Get("/api/workspaces/:workspaceID/forms/:formID/submissions", commonHandlers.Append(authMiddleware.Handler, GetUserInfoFromToken).ThenFunc(GetSubmissionsHandler))

	router.Get("/api/workspaces/:workspaceID/forms/:formID/submissions/:submissionID/changelog", commonHandlers.Append(authMiddleware.Handler, GetUserInfoFromToken).ThenFunc(GetSubmissionChangelogHandler))

	router.Get("/api/workspaces/:workspaceID/forms/:formID/submissions/:submissionID", commonHandlers.Append(authMiddleware.Handler, GetUserInfoFromToken).ThenFunc(GetSubmissionInfoHandler))
	router.Put("/api/workspaces/:workspaceID/forms/:formID/submissions/:submissionID", commonHandlers.Append(authMiddleware.Handler, GetUserInfoFromToken).ThenFunc(UpdateSubmissionHandler))
	router.Delete("/api/workspaces/:workspaceID/forms/:formID/submissions/:submissionID", commonHandlers.Append(authMiddleware.Handler, GetUserInfoFromToken).ThenFunc(DeleteSubmissionHandler))

	//Triggers and Integrations
	router.Post("/api/workspaces/:workspaceID/forms/:formID/integrations", commonHandlers.Append(authMiddleware.Handler, GetUserInfoFromToken).ThenFunc(UpdateTriggerHandler))
	router.Get("/api/workspaces/:workspaceID/forms/:formID/integrations", commonHandlers.Append(authMiddleware.Handler, GetUserInfoFromToken).ThenFunc(GetFormTriggersHandler))

	router.Get("/", commonHandlers.ThenFunc(HomePageHandler))

	fileServer := http.FileServer(http.Dir("./ui/build/static"))
	router.GET("/static/*filepath", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		w.Header().Set("Vary", "Accept-Encoding")
		w.Header().Set("Cache-Control", "public, max-age=7776000")
		r.URL.Path = p.ByName("filepath")
		fileServer.ServeHTTP(w, r)
	})

	files, err := ioutil.ReadDir("./ui/build")
	if err != nil {
		fmt.Println(err)
	}
	for _, file := range files {
		filename := file.Name()
		log.Println(filename)
		router.Get("/"+filename, commonHandlers.ThenFunc(func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "./ui/build/"+filename)
		}))
	}

	uploadedFileServer := http.FileServer(http.Dir(config.Get().RootDirectory))

	router.GET("/"+config.Get().RootDirectory+"/*filepath", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		w.Header().Set("Vary", "Accept-Encoding")
		w.Header().Set("Cache-Control", "public, max-age=7776000")
		r.URL.Path = p.ByName("filepath")
		uploadedFileServer.ServeHTTP(w, r)
	})

	router.NotFound = commonHandlers.ThenFunc(HomePageHandler)

	PORT := os.Getenv("PORT")
	if PORT == "" {
		log.Println("No Global port has been defined, using default")
		PORT = "8080"
	}

	handler := cors.New(cors.Options{
		//		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "DELETE"},
		AllowCredentials: true,
		AllowedHeaders:   []string{"Accept", "Content-Type", "X-Auth-Token", "*"},
		Debug:            false,
	}).Handler(router)
	log.Println("serving ")
	log.Fatal(http.ListenAndServe(":"+PORT, handler))
}

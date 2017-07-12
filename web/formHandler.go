package web

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/metal3d/go-slugify"
	"gitlab.com/middlefront/workspace/config"
	"gitlab.com/middlefront/workspace/database"
)

func CreateFormHandler(w http.ResponseWriter, r *http.Request) {
	httprouterParams := r.Context().Value("params").(httprouter.Params)
	workspaceID := httprouterParams.ByName("workspaceID")

	formData := database.Form{}
	err := json.NewDecoder(r.Body).Decode(&formData)
	if err != nil {
		log.Println(err)
	}

	formData.Creator = r.Context().Value("username").(string)
	formData.ID = slugify.Marshal(formData.Name, true)

	conf := config.Get()
	err = conf.Database.CreateForm(workspaceID, formData)
	if err != nil {
		log.Println(err)
	}

	message := make(map[string]interface{})
	message["code"] = 200
	message["message"] = "success"
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(message)
	if err != nil {
		log.Println(err)
	}
}

func GetFormsHandler(w http.ResponseWriter, r *http.Request) {
	httprouterParams := r.Context().Value("params").(httprouter.Params)
	workspaceID := httprouterParams.ByName("workspaceID")

	conf := config.Get()
	forms, err := conf.Database.GetForms(workspaceID)
	if err != nil {
		log.Println(err)
	}

	w.Header().Set("Content-type", "application/json")
	err = json.NewEncoder(w).Encode(forms)
	if err != nil {
		log.Println(err)
	}
}

func GetFormBySlugHandler(w http.ResponseWriter, r *http.Request) {
	httprouterParams := r.Context().Value("params").(httprouter.Params)
	workspaceID := httprouterParams.ByName("workspaceID")
	formID := httprouterParams.ByName("formID")

	conf := config.Get()
	form, err := conf.Database.GetFormBySlug(workspaceID, formID)
	if err != nil {
		log.Println(err)
	}
	w.Header().Set("Content-type", "application/json")
	err = json.NewEncoder(w).Encode(form)
	if err != nil {
		log.Println(err)
	}
}

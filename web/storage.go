package web

import (
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	"gitlab.com/middlefront/workspace/storage"
)

func GetUploadedFile(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	w.Header().Set("Vary", "Accept-Encoding")
	w.Header().Set("Cache-Control", "public, max-age=7776000")
	fileURL := strings.TrimLeft(p.ByName("filepath"), "/")
	log.Printf("fileURL %#v", fileURL)

	item, err := storage.GetByURL(fileURL)
	if err != nil {
		log.Println(err)
	}

	readCloser, err := item.Open()
	if err != nil {
		log.Println(err)
	}

	_, err = io.Copy(w, readCloser)
	if err != nil {
		log.Println(err)
	}

}

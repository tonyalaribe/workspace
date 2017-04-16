package web

import "net/http"

func HomePageHandler(w http.ResponseWriter, r *http.Request) {
	// t := template.New("homepage") // Create a template.
	// var err error
	//
	// t, err = t.ParseFiles("./ui/build/index.html") // Parse template file.
	// if err != nil {
	// 	log.Println(err)
	// }
	//
	// data := struct {
	// 	Schema string
	// }{
	// 	Schema: JSONSchema,
	// }
	// t.ExecuteTemplate(os.Stdout, "index.html", data)
	//
	// t.ExecuteTemplate(w, "index.html", data) // merge.

	http.ServeFile(w, r, "./ui/build/index.html")

}

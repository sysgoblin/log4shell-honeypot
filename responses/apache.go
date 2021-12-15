package responses

import (
	"html/template"
	"log"
	"net/http"
)

var (
	apacheHeaders = map[string]string{
		"Server":     "Apache/2.4",
		"Connection": "close",
	}
)

func CreateApacheResponse(w http.ResponseWriter) {
	tmpl, err := template.ParseFiles("responses/form.html")
	if err != nil {
		log.Fatal("Failed to parse template: ", err)
	}

	for key, value := range apacheHeaders {
		w.Header().Set(key, value)
	}

	tmpl.Execute(w, nil)
	return
}

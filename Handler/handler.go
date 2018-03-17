package Handler

import (
	"html/template"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// Index is a data structure used to populate an indexTemplate.
type Index struct {
	Title string
	Body  string
}

//HomeHandler is Default Landing Handler
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	data := &Index{
		Title: "Image gallery",
		Body:  "Welcome to the image gallery.",
	}
	t := template.Must(template.ParseFiles("./template/welcome.tmpl"))
	if err := t.Execute(w, data); err != nil {
		log.Error(err)
	}
}

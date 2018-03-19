package Handler

import (
	"html/template"
	"math/rand"
	"net/http"
	"strconv"

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
		Title: "Image Gallery",
		Body:  "Welcome to the image gallery.",
	}
	t := template.Must(template.ParseFiles("./template/welcome.tmpl"))
	data.Title = data.Title + strconv.Itoa(rand.Intn(9999))
	if err := t.Execute(w, data); err != nil {
		log.Error(err)
	}
}
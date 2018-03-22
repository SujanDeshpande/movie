package Handler

import (
	"html/template"
	"math/rand"
	"movie/MovieDB"
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"
)

// Index is a data structure used to populate an indexTemplate.
type Index struct {
	Title string
	Body  string
}

//HomeHandler is Default landing Handler
func HomeHandler(db MovieDB.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		db.AllMovies()
		data := &Index{
			Title: "Image Gallery",
			Body:  "Welcome to the image gallery.",
		}
		t := template.Must(template.ParseFiles("./template/welcome.tmpl"))
		data.Title = data.Title + strconv.Itoa(rand.Intn(9999))
		if err := t.Execute(w, data); err != nil {
			log.Error(err)
		}
	})
}

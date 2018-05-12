package handle

import (
	"html/template"
	"math/rand"
	"movie/files"
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"
)

// Index is a data structure used to populate an indexTemplate.
type Index struct {
	Title string
	Body  string
}

//Home is Default landing Handler
func Home() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

//Test is Default landing Handler
func Test() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t := template.Must(template.ParseFiles("./template/listFilesTest.tmpl"))
		if err := t.Execute(w, nil); err != nil {
			log.Error(err)
		}
	})
}

//Sort is Default Sort API Handler
func Sort() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		buildPipe()
		log.Info("Sent Data")
		w.Write([]byte("\n[]byte\n\n"))
	})
}

//ListFileData something
type ListFileData struct {
	Headers   []string
	FileInfos []files.FileInfo
}

//ListFile is Default Sort API Handler
func ListFile() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		done := make(chan struct{})
		infoChan := walkfiles(done)
		var fileInfos []files.FileInfo
		for fileInfo := range infoChan {
			log.Debug("Appending " + fileInfo.Name)
			fileInfos = append(fileInfos, fileInfo)
		}
		data := ListFileData{
			Headers:   []string{"Name", "Size", "Mode", "ModTime", "IsDir", "From", "To"},
			FileInfos: fileInfos,
		}
		t := template.Must(template.ParseFiles("./template/listFiles.tmpl"))
		if err := t.Execute(w, data); err != nil {
			log.Error(err)
		}
	})
}

//GetAll is Default Sort API Handler
func GetAll() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	})
}

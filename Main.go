package main

import (
	"movie/Handler"
	"movie/MovieDB"
	"movie/Utils"
	"net/http"
	"os"

	"github.com/boltdb/bolt"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetOutput(os.Stdout)
}

func main() {
	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	dbase := MovieDB.DB{}
	dbase.Bolted = db
	movieInfo, _ := dbase.AllMovies()
	log.Info("movieInfo")
	log.Info(movieInfo)

	r := mux.NewRouter()
	r.HandleFunc("/", Handler.HomeHandler).Methods("GET")
	config := Utils.GetConfig()
	log.Fatal(http.ListenAndServe(":"+config.Port, r))
}

package main

import (
	"net/http"
	"os"

	"github.com/boltdb/bolt"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	"movie/Handler"
	"movie/MovieDB"
	"movie/Utils"
)

func init() {
	log.SetOutput(os.Stdout)
}

func main() {

	config := Utils.GetConfig()
	db, err := bolt.Open(config.Database.Name, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	MovieDB.CreateBucket(db)
	dbase := MovieDB.DB{Bolted: db}

	r := mux.NewRouter()
	r.Handle("/", Handler.HomeHandler(dbase)).Methods("GET")
	r.Handle("/sort", Handler.SortHandler()).Methods("GET")
	log.Fatal(http.ListenAndServe(":"+config.Port, r))
}

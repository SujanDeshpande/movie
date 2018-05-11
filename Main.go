package main

import (
	"movie/Utils"
	"movie/handler"
	"net/http"
	"os"

	"github.com/boltdb/bolt"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetOutput(os.Stdout)
}

type filer struct {
	Info        os.FileInfo `json:"info"`
	Destination string      `json:"destination"`
}

func main() {
	setupDB()
	r := mux.NewRouter()
	r.Handle("/", Handler.HomeHandler()).Methods("GET")
	r.Handle("/sort", Handler.SortHandler()).Methods("GET")
	log.Fatal(http.ListenAndServe(":"+Utils.GetConfig().Port, r))
}

func setupDB() {
	db, _ := bolt.Open(Utils.GetConfig().Database.Name, 0600, nil)
	db.Update(func(tx *bolt.Tx) error {
		tx.CreateBucketIfNotExists([]byte(Utils.GetConfig().Database.Bucket))
		return nil
	})
	log.Info("DB Setup Done")
}

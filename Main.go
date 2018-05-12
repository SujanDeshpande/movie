package main

import (
	"movie/database"
	"movie/handle"
	"movie/utils"
	"net/http"
	"os"

	"github.com/boltdb/bolt"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
	//log.SetFormatter(&log.JSONFormatter{})
}

type filer struct {
	Info        os.FileInfo `json:"info"`
	Destination string      `json:"destination"`
}

func main() {
	var err error
	database.DBCon, err = setupDB()
	if err != nil {
		log.Error("Error")
	}
	r := mux.NewRouter()

	fs := http.FileServer(http.Dir("static"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))
	r.Handle("/sort", handle.Sort()).Methods("GET")
	r.Handle("/list", handle.ListFile()).Methods("GET")
	r.Handle("/test", handle.Test()).Methods("GET")
	log.Fatal(http.ListenAndServe(":"+utils.GetConfig().Port, r))
}

func setupDB() (*bolt.DB, error) {
	db, err := bolt.Open(utils.GetConfig().Database.Name, 0600, nil)
	if err != nil {
		return nil, err
	}
	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(utils.GetConfig().Database.Bucket))
		if err != nil {
			return nil
		}
		log.WithFields(log.Fields{
			"DB":     utils.GetConfig().Database.Name,
			"bucket": utils.GetConfig().Database.Bucket,
		}).Debug("Created")
		return nil
	})
	log.Info("DB Setup Done")
	return db, nil
}

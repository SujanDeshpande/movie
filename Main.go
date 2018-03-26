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
	config := Utils.GetConfig()
	db, err := bolt.Open(config.Database.Name, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	createBucket(db)
	dbase := MovieDB.DB{Bolted: db}

	r := mux.NewRouter()
	r.Handle("/", Handler.HomeHandler(dbase)).Methods("GET")
	log.Fatal(http.ListenAndServe(":"+config.Port, r))
}

func createBucket(db *bolt.DB) {
	db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(Utils.GetDatabaseConfig().Bucket))
		if bucket != nil {
			log.WithField("bucket", Utils.GetDatabaseConfig().Bucket).Info("Bucket already exists")
		}
		if bucket == nil {
			_, berr := tx.CreateBucket([]byte(Utils.GetDatabaseConfig().Bucket))
			if berr != nil {
				log.WithError(berr).WithField("bucket", Utils.GetDatabaseConfig().Bucket).Fatal("Unable to create a bucket")
			}
			log.Info("Bucket created sucessfully")
		}
		return nil
	})
}

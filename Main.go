package main

import (
	"movie/Handler"
	"movie/Utils"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetOutput(os.Stdout)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", Handler.HomeHandler).Methods("GET")
	config := Utils.GetConfig()
	log.Fatal(http.ListenAndServe(":"+config.Port, r))
}

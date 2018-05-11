package Utils

import (
	"io/ioutil"

	log "github.com/sirupsen/logrus"
)

//ReadFile - reads the File from specified location and returns as a string
func ReadFile(filename string) string {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		log.WithError(err).Error()
	}
	str := string(file)
	return str
}

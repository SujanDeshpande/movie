package utils

import (
	"encoding/json"
	"io/ioutil"

	log "github.com/sirupsen/logrus"
)

//Config - Represent all the config information.
type Config struct {
	Port     string   `json:"port"`
	Database DBConfig `json:"db"`
	Location Location `json:"location"`
}

//DBConfig - Represent all the config information specific to database.
type DBConfig struct {
	Name   string `json:"name"`
	Bucket string `json:"bucket"`
}

//Location - Represent all the default Location information specific to files.
type Location struct {
	Src  string `json:"src"`
	Dest string `json:"dest"`
}

//Configuration - All configurations
var Configuration = &Config{}

//init - initialize the configurations.
func init() {
	configFile := "./utils/config.json"
	file, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.WithError(err).Error("Unable to read config file : " + configFile)
	}
	jerr := json.Unmarshal(file, &Configuration)
	if jerr != nil {
		log.WithError(err).Error("Unable to marshal config file : " + configFile)
	}
}

//GetConfig - all configurations
func GetConfig() *Config {
	return Configuration
}

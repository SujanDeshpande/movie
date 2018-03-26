package Utils

import (
	"encoding/json"
	"os"

	log "github.com/sirupsen/logrus"
)

//Config - Represent all the config information.
type Config struct {
	Port     string         `json:"port"`
	Database DBConfig       `json:"db"`
	Producer ProducerConfig `json:"producer"`
	Consumer ConsumerConfig `json:"consumer"`
}

//DBConfig - Represent all the config information specific to database.
type DBConfig struct {
	Name       string      `json:"name"`
	Permission os.FileMode `json:"permission"`
	Bucket     string      `json:"bucket"`
}

//ProducerConfig - Represent all the config information specific to Producer.
type ProducerConfig struct {
	Protocol string `json:"protocol"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	ClientID string `json:"clientId"`
	Username string `json:"username"`
	Password string `json:"password"`
	Topic    string `json:"topic"`
}

//ConsumerConfig - Represent all the config information specific to Consumer.
type ConsumerConfig struct {
	Protocol string `json:"protocol"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	ClientID string `json:"clientId"`
	Username string `json:"username"`
	Password string `json:"password"`
	Topic    string `json:"topic"`
}

//Configuration - All configurations
var Configuration = &Config{}

//init - initialize the configurations.
func init() {
	jsonString := ReadFile("./config.json")
	jsonBytes := []byte(jsonString)
	err := json.Unmarshal(jsonBytes, &Configuration)
	if err != nil {
		log.Panic(err)
		panic(err)
	}
	log.Debug(jsonString)
}

//GetProducerConfig -  the configuration specific for producer
func GetProducerConfig() ProducerConfig {
	return Configuration.Producer
}

//GetConsumerConfig - the configuration specific for consumer
func GetConsumerConfig() ConsumerConfig {
	return Configuration.Consumer
}

//GetDatabaseConfig - the configuration specific for database
func GetDatabaseConfig() DBConfig {
	return Configuration.Database
}

//GetConfig - all configurations
func GetConfig() *Config {
	return Configuration
}

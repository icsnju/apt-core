package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

var configPath = "config/master.json"

type MasConf struct {
	SharePath string
	Port      string
	DBUrl     string
	DBNAME    string
}

var configuration *MasConf = &MasConf{}

func InitConfig() {

	content, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatalln(err)
	}

	err = json.Unmarshal(content, configuration)
	if err != nil {
		log.Fatalln(err)
	}
}

func GetSharePath() string {
	return configuration.SharePath
}

func GetPort() string {
	return configuration.Port
}

func GetDBUrl() string {
	return configuration.DBUrl
}

func GetDBName() string {
	return configuration.DBNAME
}

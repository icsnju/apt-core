package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type SlaConf struct {
	SharePath string
	ServerIP  string
	AdbPath   string
}

var configPath = "slave.json"
var configuration *SlaConf = &SlaConf{}

func initConfig() {
	content, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatalln(err)
	}

	err = json.Unmarshal(content, configuration)
	if err != nil {
		log.Fatalln(err)
	}
}

func getSharePath() string {
	return configuration.SharePath
}

func getServerIP() string {
	return configuration.ServerIP
}

func getAdbPath() string {
	return configuration.AdbPath
}

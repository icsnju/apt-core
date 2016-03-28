package main

import (
	"apsaras/comm/framework"
	"log"
)

func main() {
	//init
	initConfig()
	//register in gob
	framework.RigisterGob()

	go deviceManager.loopUpdate()
	//start connet to master
	diaMaster(getServerIP())
	log.Println("Slave Over!")
}

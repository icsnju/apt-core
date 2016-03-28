package main

import (
	"apsaras/comm/filter"
	"apsaras/comm/framework"
	"apsaras/server/config"
	"apsaras/server/controllers"
	"apsaras/server/master"
	"apsaras/server/models"
	"log"

	"github.com/astaxie/beego"
)

func startServer() {
	//api
	beego.Router("/job", &controllers.JobController{}, "get:ListJobs;post:CreateJob")
	beego.Router("/job/:id", &controllers.JobController{}, "get:GetJob")
	beego.Router("/device", &controllers.DeviceController{}, "get:ListDevices")
	//html
	beego.Router("/*", &controllers.MainController{})

	beego.Run()
}

func main() {
	//register in gob
	framework.RigisterGob()
	filter.RigisterGob()

	//init config
	config.InitConfig()
	//init DB
	err := models.InitDB(config.GetDBUrl(), config.GetDBName())
	if err != nil {
		log.Fatalln(err)
	}
	defer models.CloseDB()

	go startServer()
	master.StartMaster(config.GetPort())
}

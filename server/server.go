package main

import (
	"apsaras/server/controllers"
	"apsaras/server/master"
	"apsaras/server/models"
	"fmt"
	"os"

	"github.com/astaxie/beego"
)

const DBURL = "localhost"

func startClient() {
	//api
	beego.Router("/job", &controllers.JobController{}, "get:ListJobs;post:CreateJob")
	beego.Router("/job/:id", &controllers.JobController{}, "get:GetJob")
	beego.Router("/device", &controllers.DeviceController{}, "get:ListDevices")
	//html
	beego.Router("/*", &controllers.MainController{})

	beego.Run()
}

func main() {

	//init DB
	err := models.InitDB(DBURL)
	if err != nil {
		fmt.Println("DB:", err)
		os.Exit(1)
	}
	defer models.CloseDB()

	go startClient()
	master.StartMaster()
}

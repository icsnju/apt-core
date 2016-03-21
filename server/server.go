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
	beego.Router("/", &controllers.MainController{})
	beego.Router("/submit", &controllers.MainController{})
	beego.Router("/job/", &controllers.JobController{}, "get:ListJobs;post:CreateJob")
	beego.Router("/device", &controllers.DeviceController{}, "get:ListDevices")

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

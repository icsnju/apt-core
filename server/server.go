package main

import (
	"apsaras/server/controllers"
	"apsaras/server/master"

	"github.com/astaxie/beego"
)

func startClient() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/submit", &controllers.MainController{})
	beego.Router("/table", &controllers.MainController{})
	beego.Router("/task/", &controllers.TaskController{}, "get:ListTasks;post:NewTask")
	beego.Router("/task/:id:int", &controllers.TaskController{}, "get:GetTask;put:UpdateTask")
	beego.Router("/", &controllers.MainController{})
	beego.Run()
}

func main() {
	go startClient()
	master.StartMaster()
}

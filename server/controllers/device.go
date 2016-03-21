package controllers

import (
	"apsaras/server/master"

	"github.com/astaxie/beego"
)

type DeviceController struct {
	beego.Controller
}

func (this *DeviceController) ListDevices() {
	devices := master.GetDevices()
	this.Data["json"] = devices
	this.ServeJSON()
}

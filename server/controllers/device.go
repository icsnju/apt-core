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

func (this *DeviceController) GetDevice() {
	ip := this.Ctx.Input.Param(":ip")
	id := this.Ctx.Input.Param(":id")
	device, err := master.GetDevice(ip, id)
	if err != nil {
		this.Ctx.Output.SetStatus(404)
		this.Ctx.Output.Body([]byte("Device not found"))
		return
	}
	this.Data["json"] = device
	this.ServeJSON()
}

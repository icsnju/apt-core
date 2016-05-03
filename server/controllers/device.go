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
	id := this.Ctx.Input.Param(":id")
	device, _, err := master.GetDevice(id)
	if err != nil {
		this.Ctx.Output.SetStatus(404)
		this.Ctx.Output.Body([]byte("Device not found"))
		return
	}
	this.Data["json"] = device
	this.ServeJSON()
}

func (this *DeviceController) GetDeviceIP() {
	id := this.Ctx.Input.Param(":id")
	_, ip, err := master.GetDevice(id)
	if err != nil {
		this.Ctx.Output.SetStatus(404)
		this.Ctx.Output.Body([]byte("Device not found"))
		return
	}
	this.Ctx.Output.Body([]byte(ip))
}

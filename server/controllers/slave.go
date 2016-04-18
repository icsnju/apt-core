package controllers

import (
	"apsaras/server/master"

	"github.com/astaxie/beego"
)

type SlaveController struct {
	beego.Controller
}

func (this *SlaveController) ListSlaves() {
	slaves := master.GetSlaveSketches()
	this.Data["json"] = slaves
	this.ServeJSON()
}

func (this *SlaveController) GetSlave() {
	id := this.Ctx.Input.Param(":id")
	slave, ex := master.GetSlave(id)
	if !ex {
		this.Ctx.Output.SetStatus(404)
		this.Ctx.Output.Body([]byte("Slave not found"))
		return
	}

	this.Data["json"] = slave
	this.ServeJSON()
}

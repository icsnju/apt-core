package controllers

import (
	"apsaras/server/models"
	"fmt"

	"github.com/astaxie/beego"
)

type JobController struct {
	beego.Controller
}

func (this *JobController) ListJobs() {
	jobs := models.GetJobsInDB()
	this.Data["json"] = jobs
	this.ServeJSON()
}

func (this *JobController) CreateJob() {
	fmt.Println(string(this.Ctx.Input.RequestBody))
	fmt.Println(this.Ctx.Input.IsUpload())
	fmt.Println(this.Ctx.Input.Data())
}

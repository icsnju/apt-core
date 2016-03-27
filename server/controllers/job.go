package controllers

import (
	"apsaras/comm"
	"apsaras/comm/framework"
	"apsaras/server/master"
	"apsaras/server/models"
	"log"
	"os"
	"path"

	"github.com/astaxie/beego"
)

type JobController struct {
	beego.Controller
}

//List all job sketch in db
func (this *JobController) ListJobs() {
	jobs := models.GetJobSketchesInDB()
	this.Data["json"] = jobs
	this.ServeJSON()
}

func (this *JobController) GetJob() {
	id := this.Ctx.Input.Param(":id")
	job, err := models.GetJobInDB(id)
	if err != nil {
		this.Ctx.Output.SetStatus(404)
		this.Ctx.Output.Body([]byte("Job not found"))
		log.Println(err)
		return
	}

	this.Data["json"] = job
	this.ServeJSON()
}

//Create a new job for client
func (this *JobController) CreateJob() {
	//create a subjob
	jobjson := this.GetString("job")
	subjob, err := models.ParserSubJobFromJson([]byte(jobjson))
	if err != nil {
		log.Println(err)
		return
	}

	//create a job
	job := master.CreateJob(subjob)
	//move test files
	err = moveTestFile(&job, this)
	if err != nil {
		log.Println(err)
		return
	}
	//save job in db
	err = models.SaveJobInDB(job)
	if err != nil {
		log.Println(err)
		return
	}
	err = models.SaveJobSketchInDB(job.ToSketch())
	if err != nil {
		log.Println(err)
		return
	}
	//add job in master
	master.AddJobInMaster(job)

}

func moveTestFile(job *models.Job, control *JobController) error {
	if job.JobInfo.FrameKind == framework.FRAME_MONKEY {
		//move file to share path
		dist := path.Join(master.GetSharePath(), job.JobId)
		err := os.Mkdir(dist, os.ModePerm)
		if err != nil {
			return err
		}
		dist = path.Join(dist, comm.APPNAME)
		err = control.SaveToFile("file", dist)
		if err != nil {
			return err
		}
		//update framework info
		monkey := job.JobInfo.Frame.(framework.MonkeyFrame)
		monkey.AppPath = dist
		job.JobInfo.Frame = monkey
	}
	return nil
}

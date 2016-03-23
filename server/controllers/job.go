package controllers

import (
	"apsaras/comm"
	"apsaras/comm/comp"
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

//List all job in db
func (this *JobController) ListJobs() {
	jobs := models.GetJobsInDB()
	this.Data["json"] = jobs
	this.ServeJSON()
}

//Create a new job for client
func (this *JobController) CreateJob() {
	//create a subjob
	jobjson := this.GetString("job")
	subjob, err := comp.ParserSubJobFromJson([]byte(jobjson))
	if err != nil {
		log.Println(err)
		return
	}

	//create a job
	job := master.CreateJob(subjob)
	err = moveTestFile(&job, this)
	if err != nil {
		log.Println(err)
		return
	}
	//save job in db
	err = models.SaveJobInDB(job.ToBrief())
	//add job in master
	master.AddJobInMaster(job)

}

func moveTestFile(job *comp.Job, control *JobController) error {
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

package controllers

import (
	"apsaras/comm"
	"apsaras/comm/framework"
	"apsaras/server/config"
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

//Get a job by id
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
	this.Ctx.Output.SetStatus(200)
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
		dist := path.Join(config.GetSharePath(), job.JobId)
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
	} else if job.JobInfo.FrameKind == framework.FRAME_ROBOT {
		//move file to share path
		dist := path.Join(config.GetSharePath(), job.JobId)
		err := os.Mkdir(dist, os.ModePerm)
		if err != nil {
			return err
		}
		distApp := path.Join(dist, comm.APPNAME)
		err = control.SaveToFile("app", distApp)
		if err != nil {
			return err
		}
		distTest := path.Join(dist, comm.TESTNAME)
		err = control.SaveToFile("test", distTest)
		if err != nil {
			return err
		}
		//update framework info
		robotium := job.JobInfo.Frame.(framework.RobotFrame)
		robotium.AppPath = distApp
		robotium.TestPath = distTest
		job.JobInfo.Frame = robotium
	}
	return nil
}

//Get job result by device id
func (this *JobController) GetTaskResult() {
	jid := this.GetString("jobid")
	did := this.GetString("deviceid")

	resultPath := path.Join(config.GetSharePath(), jid, did)
	zipPath := path.Join(config.GetSharePath(), jid, did+".zip")
	err := comm.Zipit(resultPath, zipPath)
	if err != nil {
		log.Println(err)
		return
	}
	name := jid + "_" + did + ".zip"
	this.Ctx.Output.Download(zipPath, name)
}

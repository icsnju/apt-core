package models

import (
	"apsaras/comm/comp"
	"fmt"

	"gopkg.in/mgo.v2/bson"
)

const (
	JOB_ID     = "jobid"
	JOB_STATUS = "status"
)

func SaveJobInDB(job comp.JobBrief) error {
	return jobCollection.Insert(job)

}

func GetJobsInDB() []comp.JobBrief {
	jobs := make([]comp.JobBrief, 0)
	var job comp.JobBrief
	iter := jobCollection.Find(nil).Iter()
	for iter.Next(&job) {
		jobs = append(jobs, job)
	}
	return jobs
}

func UpdateJobInDB(id string, update interface{}) {
	err := jobCollection.Update(bson.M{JOB_ID: id}, update)
	if err != nil {
		fmt.Println("job update err in db :", err)
	}
}

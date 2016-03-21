package models

import (
	"apsaras/comm/comp"
	"fmt"

	"gopkg.in/mgo.v2/bson"
)

func SaveJobInDB(job comp.JobBrief) {
	err := jobCollection.Insert(job)
	if err != nil {
		fmt.Println(err)
	}
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
	err := jobCollection.Update(bson.M{"JobId": id}, update)
	if err != nil {
		fmt.Println("job update err in db :", err)
	}
}

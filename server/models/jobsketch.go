package models

import (
	"log"

	"gopkg.in/mgo.v2/bson"
)

type JobSketch struct {
	JobId      string
	StartTime  string
	FrameKind  string
	FilterKind string
	Status     int
}

func SaveJobSketchInDB(job JobSketch) error {
	return jobSketchCollection.Insert(job)
}

func GetJobSketchesInDB() []interface{} {
	jobs := make([]interface{}, 0)
	var job interface{}
	iter := jobSketchCollection.Find(nil).Iter()
	for iter.Next(&job) {
		jobs = append(jobs, job)
	}
	return jobs
}

func UpdateJobSketchInDB(id string, update interface{}) {
	err := jobSketchCollection.Update(bson.M{JOB_ID: id}, update)
	if err != nil {
		log.Println("job sketch update err in db :", err)
	}
}

func DeleteJobSketchInDB(id string) {
	log.Println("Delete Job sketch", id)
	err := jobSketchCollection.Remove(bson.M{JOB_ID: id})
	if err != nil {
		log.Println("job sketch delete err in db :", err)
	}
}

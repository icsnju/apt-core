package models

import "gopkg.in/mgo.v2"

const (
	DBNAME      = "aptweb-dev"
	JOBSKETCH_C = "jobsketches"
	JOB_C       = "jobs"
)

const (
	JOB_ID     = "jobid"
	JOB_STATUS = "status"
)

var session *mgo.Session

var jobSketchCollection *mgo.Collection
var jobCollection *mgo.Collection

func InitDB(url string) error {
	session, err := mgo.Dial(url)
	if err != nil {
		return err
	}
	session.SetMode(mgo.Monotonic, true)
	db := session.DB(DBNAME)

	jobCollection = db.C(JOB_C)
	jobSketchCollection = db.C(JOBSKETCH_C)
	return nil
}

func CloseDB() {
	session.Close()
}

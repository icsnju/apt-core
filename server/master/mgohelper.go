package master

import (
	"apsaras/task"
	"fmt"

	"gopkg.in/mgo.v2"
)

const (
	DBNAME = "aptweb-dev"
	JOB_C  = "jobs"
)

var session *mgo.Session
var jobC *mgo.Collection

func initDB(url string) error {
	session, err := mgo.Dial(url)
	if err != nil {
		return err
	}
	session.SetMode(mgo.Monotonic, true)
	db := session.DB(DBNAME)

	jobC = db.C(JOB_C)
	return nil
}

func closeDB() {
	session.Close()
}

func saveJobInDB(job task.JobBrief) {
	err := jobC.Insert(job)
	if err != nil {
		fmt.Println(err)
	}
}

func updateJobInDB() {

}

package models

import "gopkg.in/mgo.v2"

const (
	DBNAME   = "aptweb-dev"
	JOB_C    = "jobs"
	DEVICE_C = "devices"
)

var session *mgo.Session

var jobCollection *mgo.Collection
var deviceCollection *mgo.Collection

func InitDB(url string) error {
	session, err := mgo.Dial(url)
	if err != nil {
		return err
	}
	session.SetMode(mgo.Monotonic, true)
	db := session.DB(DBNAME)

	jobCollection = db.C(JOB_C)
	deviceCollection = db.C(DEVICE_C)
	return nil
}

func CloseDB() {
	session.Close()
}

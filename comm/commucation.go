package comm

import "time"

const (
	HIMASTER     = "himaster"
	CHECKSLAVES  = "slaves"
	CHECKJOBS    = "jobs"
	CHECKDEVICES = "devices"
	SUBJOB       = "subjob"
	WEBSUBJOB    = "websubjob"
	JOBSTATE     = "job"

	HEATBEAT = "heartbeat"
	IMALIVE  = "imalive"
	GOODJOB  = "goodjob"

	WAITFORHEARTTIME = 5 * time.Minute
	WAITFORDIA       = 1 * time.Minute
	HEARTTIME        = 10 * time.Second
	//ALLOTTIME        = 10 * time.Second
	UPDATEDEVINFO = 10 * time.Minute
	SHIFTTIME     = 30 * time.Second
)

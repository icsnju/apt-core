package comp

import (
	"apsaras/comm/framework"
	"time"
)

const (
	TASK_WAIT = iota
	TASK_RUN
	TASK_COMPLETE
	TASK_FAIL
)

const (
	MAXWAIT      = 30 * time.Second
	MAXPRIOPRITY = 2
)

//Task
type Task struct {
	JobId      string
	DeviceId   string
	TargetId   string
	State      int
	FinishTime time.Time
	StartTime  time.Time
}

type TaskBrief struct {
	JobId    string
	TargetId string
	State    string
}

//A task runner
type RunTask struct {
	TaskInfo Task
	Frame    framework.FrameStruct
}

type RunTaskList struct {
	Tasks []RunTask
}

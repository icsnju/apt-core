package task

import (
	"apsaras/framework"
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

//A task runner
type RunTask struct {
	TaskInfo Task
	Frame    framework.FrameStruct
}

type RunTaskList struct {
	Tasks []RunTask
}

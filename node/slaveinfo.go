package node

import (
	"apsaras/andevice"
	"apsaras/task"
)

//IP
//DevicesStates,string:int
type SlaveInfo struct {
	IP           string
	DeviceStates map[string]andevice.Device
	TaskStates   map[string]task.Task
}

type SlaveMap struct {
	Map map[string]SlaveInfo
}

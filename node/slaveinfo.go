package node

import (
	"apsaras/device"
	"apsaras/task"
)

//IP
//DevicesStates,string:int
type SlaveInfo struct {
	IP           string
	DeviceStates map[string]device.Device
	TaskStates   map[string]task.Task
}

type SlaveMap struct {
	Map map[string]SlaveInfo
}

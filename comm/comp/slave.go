package comp

type SlaveInfo struct {
	IP           string
	DeviceStates map[string]Device
	TaskStates   map[string]Task
}

type SlaveMap struct {
	Map map[string]SlaveInfo
}

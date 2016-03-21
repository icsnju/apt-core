package master

import "apsaras/comm/comp"

//get all devices in all slaves
func GetDevices() []comp.DeviceInfo {
	devices := make([]comp.DeviceInfo, 0)
	for _, dev := range slaveManager.getDevices() {
		devices = append(devices, dev.Info)
	}
	return devices
}

func SubmitJob(message []byte) {
	jobManager.submitJob(message)
}

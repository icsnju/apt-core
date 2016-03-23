package master

import "apsaras/comm/comp"

/*Api of master*/

//get all devices in all slaves
func GetDevices() []comp.DeviceInfo {
	devices := make([]comp.DeviceInfo, 0)
	for _, dev := range slaveManager.getDevices() {
		devices = append(devices, dev.Info)
	}
	return devices
}

//Create a new job
func CreateJob(subjob comp.SubJob) comp.Job {

	job := jobManager.createJob(subjob)
	return job
}

//Add this job in master
func AddJobInMaster(job comp.Job) {
	jobManager.addJob(job)
}

func GetSharePath() string {
	return shareDirPath
}

package master

import (
	"apsaras/comm/comp"
	"apsaras/server/models"
)

/*Api of master*/

//get all devices in all slaves
func GetDevices() []comp.DeviceInfo {
	devices := make([]comp.DeviceInfo, 0)
	for _, dev := range slaveManager.getDevices() {
		devices = append(devices, dev.Info)
	}
	return devices
}

//get a device in a node
func GetDevice(ip, id string) (comp.Device, error) {
	return slaveManager.getDevice(ip, id)
}

//Create a new job
func CreateJob(subjob models.SubJob) models.Job {

	job := jobManager.createJob(subjob)
	return job
}

//Kill a job
func KillJob(id string) {
	jobManager.killJob(id)
}

//Add this job in master
func AddJobInMaster(job models.Job) {
	jobManager.addJob(job)
}

func IsFinished(id string) bool {
	return !jobManager.contain(id)
}

//Get sketches of slaves
func GetSlaveSketches() []models.SlaveSketch {
	return slaveManager.getSlaveSketches()
}

func GetSlave(IP string) (comp.SlaveInfo, bool) {
	return slaveManager.getSlave(IP)
}

package master

import (
	"apsaras/device"
	"apsaras/task"
	"apsaras/tools"
	"math/rand"
	"time"
)

//find the oldest job
func findOldJob(id string) string {

	var maxJobId string = "-1"
	oldTime := time.Now()

	//find a old job
	for ke, jo := range jobMap {
		if jo.Finished {
			continue
		}
		ts, ex := jo.TaskMap[id]
		if ex && ts.State == task.TASK_WAIT {
			if jo.StartTime.Before(oldTime) {
				maxJobId = ke
				oldTime = jo.StartTime
			}
		}
	}
	return maxJobId
}

//find the best job
func findBestJob(id string) string {
	var maxPriority float64 = -1
	var count1 int = 0
	var maxJob1 string = "-1"

	var maxWaitTime float64 = -1
	var count2 int = 0
	var maxJob2 string = "-1"

	//find a good job
	for ke, jo := range jobMap {
		if jo.Finished {
			continue
		}
		ts, ex := jo.TaskMap[id]
		if ex && ts.State == task.TASK_WAIT {
			p := jo.GetPriority()
			if p == -1 {
				if jo.GetWaitTime() > maxWaitTime {
					maxWaitTime = jo.GetWaitTime()
					maxJob2 = ke
				}
				count2++
			} else {
				if p > maxPriority {
					maxPriority = p
					maxJob1 = ke
				}
				count1++
			}
		}
	}

	sum := count1 + count2
	if sum <= 0 {
		return "-1"
	}

	rand.Seed(time.Now().UnixNano())
	lottery := rand.Intn(count1 + count2)
	if lottery < count1 {
		return maxJob1
	} else {
		return maxJob2
	}
}

//Shift it
func shift() {
	for {
		jobLock.Lock()
		for jid, job := range jobMap {
			if !job.Finished && job.GetWaitTime() > task.MAXWAIT.Seconds() {
				slavesLock.Lock()
				var devMap map[string]device.Device = make(map[string]device.Device)
				for _, si := range slavesMap {
					for k, v := range si.DeviceStates {
						devMap[k] = v
					}
				}

				var nTaskMap map[string]task.Task = make(map[string]task.Task)
				for tid, ts := range job.TaskMap {
					if ts.State == task.TASK_WAIT {
						//get similar device
						target := ts.TargetId
						idMap := device.FindSimilarDevs(target, devMap, 0)
						//TODO

						//remove repeated id
						var idList []string = make([]string, 0)
						for cid, _ := range idMap {
							_, ex := job.TaskMap[cid]
							if !ex || cid == target {
								idList = append(idList, cid)
							}
						}

						//get a redom id
						if len(idList) > 0 {
							//give me a replacement for this device
							rand.Seed(time.Now().UnixNano())
							r := rand.Intn(len(idList))
							ts.DeviceId = idList[r]
							tid = ts.DeviceId
						}
					}
					nTaskMap[tid] = ts
				}
				slavesLock.Unlock()
				job.TaskMap = nTaskMap
				jobMap[jid] = job
			}
		}
		jobLock.Unlock()
		time.Sleep(tools.SHIFTTIME)
	}
}

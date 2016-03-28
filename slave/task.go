package main

import (
	"apsaras/comm"
	"apsaras/comm/comp"
	"log"
	"path"
	"sync"
)

type TaskManager struct {
	taskMap  map[string]comp.RunTask
	taskLock *sync.Mutex
}

var taskManager *TaskManager = &TaskManager{make(map[string]comp.RunTask), new(sync.Mutex)}

func (m *TaskManager) getTaskInfo() map[string]comp.Task {
	tasks := make(map[string]comp.Task)

	m.taskLock.Lock()
	for ke, ts := range m.taskMap {
		tasks[ke] = ts.TaskInfo
	}

	//remove finished task
	for ke, ts := range tasks {
		if ts.State == comp.TASK_COMPLETE || ts.State == comp.TASK_FAIL {
			//TODO move file is time-consuming
			srcPath := path.Join(ts.JobId, ts.DeviceId)
			dstPath := path.Join(getSharePath(), ts.JobId)
			cmd := "cp -r " + srcPath + " " + dstPath
			comm.ExeCmd(cmd)
			delete(m.taskMap, ke)
		}
	}
	m.taskLock.Unlock()
	return tasks
}

func (m *TaskManager) addTask(ts comp.RunTask) {
	//If this task is in this slave
	key := ts.TaskInfo.JobId + ":" + ts.TaskInfo.DeviceId
	m.taskLock.Lock()
	_, ex := m.taskMap[key]
	m.taskLock.Unlock()
	if ex {
		log.Println("Error! Task have same ID!")
		return
	}
	m.taskLock.Lock()
	m.taskMap[key] = ts
	m.taskLock.Unlock()
}

func (m *TaskManager) updateTaskStates(task comp.RunTask) {
	key := task.TaskInfo.JobId + ":" + task.TaskInfo.DeviceId
	m.taskLock.Lock()
	m.taskMap[key] = task
	m.taskLock.Unlock()
}

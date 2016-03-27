package master

import (
	"apsaras/comm/comp"
	"sync"
)

type SlaveManager struct {
	slavesMap  map[string]comp.SlaveInfo
	slavesLock *sync.Mutex
}

var slaveManager *SlaveManager = &SlaveManager{make(map[string]comp.SlaveInfo), new(sync.Mutex)}

func (m *SlaveManager) getDevices() []comp.Device {
	devList := make([]comp.Device, 0)
	m.slavesLock.Lock()
	for _, si := range m.slavesMap {
		for _, v := range si.DeviceStates {
			devList = append(devList, v)
		}
	}
	m.slavesLock.Unlock()
	return devList
}

func (m *SlaveManager) updateSlave(s comp.SlaveInfo) {
	m.slavesLock.Lock()
	m.slavesMap[s.IP] = s
	m.slavesLock.Unlock()
}

func (m *SlaveManager) addSlave(slave comp.SlaveInfo) bool {
	ok := false
	m.slavesLock.Lock()
	_, ex := m.slavesMap[slave.IP]
	if !ex {
		m.slavesMap[slave.IP] = slave
		ok = true
	}
	m.slavesLock.Unlock()
	return ok
}

func (m *SlaveManager) getSlave(ip string) comp.SlaveInfo {
	return m.slavesMap[ip]
}

func (m *SlaveManager) removeSlave(ip string) {
	m.slavesLock.Lock()
	delete(m.slavesMap, ip)
	m.slavesLock.Unlock()
}

func (m *SlaveManager) updateSlaveDeviceState(ip, id string, state int) {
	m.slavesLock.Lock()
	slave := m.slavesMap[ip]
	dev := slave.DeviceStates[id]
	dev.State = state
	slave.DeviceStates[id] = dev
	m.slavesMap[ip] = slave
	m.slavesLock.Unlock()
}

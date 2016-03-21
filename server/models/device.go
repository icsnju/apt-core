package models

import (
	"apsaras/comm/comp"
	"fmt"

	"gopkg.in/mgo.v2/bson"
)

func SaveDeviceInDB(device comp.DeviceInfo) {
	err := deviceCollection.Insert(device)
	if err != nil {
		fmt.Println(err)
	}
}

func GetDevicesInDB() []comp.DeviceInfo {
	devices := make([]comp.DeviceInfo, 0)
	var device comp.DeviceInfo
	iter := deviceCollection.Find(nil).Iter()
	for iter.Next(&device) {
		devices = append(devices, device)
	}
	return devices
}

func DeleteDeviceInDB(id string) {
	deviceCollection.Remove(bson.M{"id": id})
}

package framework

import (
	"apsaras/andevice"
	"math/rand"
	"time"
)

//The kind of device filter
const (
	FILTER_SPECIFYDEVICES = "specify_devices"
	FILTER_COMPATIBILITY  = "compatibilty"
	FILTER_SPECIFYATTR    = "specify_attr"
)

type FilterInterface interface {
	GetDeviceSet(availDevs map[string]andevice.Device) []string
}

//Specify devices
type SpecifyDevFilter struct {
	IdList []string
}

func (sd SpecifyDevFilter) GetDeviceSet(availDevs map[string]andevice.Device) []string {
	if len(sd.IdList) == 1 && sd.IdList[0] == "*" {
		var nlist []string = make([]string, 0)
		for id, _ := range availDevs {
			nlist = append(nlist, id)
		}
		sd.IdList = nlist
	}

	return sd.IdList
}

//Compatibity rule
type CompatibilityFilter struct {
	Quantity int
	Dominate string
}

func (cf CompatibilityFilter) GetDeviceSet(availDevs map[string]andevice.Device) []string {
	var idList []string = make([]string, 0)
	clus := andevice.Kmedoids(availDevs, cf.Dominate)
	var idMap map[string]int = make(map[string]int)
	for len(idMap) < cf.Quantity {
		for _, clu := range clus {
			rand.Seed(time.Now().UnixNano())
			r := rand.Intn(len(clu.Member))
			idMap[clu.Member[r]] = 1
		}
	}
	return idList
}

//Specify attribute
type SpecifyAttrFilter struct {
	Quantity int
	Attr     string
	Value    string
}

func (cf SpecifyAttrFilter) GetDeviceSet(availDevs map[string]andevice.Device) []string {
	var idList []string = make([]string, 0)
	return idList
}

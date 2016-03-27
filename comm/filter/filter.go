package filter

import (
	"apsaras/comm/comp"
	"encoding/gob"
)

//The kind of device filter
const (
	FILTER_SPECIFYDEVICES = "specify_devices"
	FILTER_COMPATIBILITY  = "compatibilty"
	FILTER_SPECIFYATTR    = "specify_attr"
)

type FilterInterface interface {
	GetDeviceSet(availDevs []comp.Device) []string
}

//Specify devices
type SpecifyDevFilter struct {
	IdList []string
}

func (sd SpecifyDevFilter) GetDeviceSet(availDevs []comp.Device) []string {
	if len(sd.IdList) == 1 && sd.IdList[0] == "*" {
		var nlist []string = make([]string, 0)
		for _, dev := range availDevs {
			nlist = append(nlist, dev.Info.Id)
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

func (cf CompatibilityFilter) GetDeviceSet(availDevs []comp.Device) []string {
	var idList []string = make([]string, 0)
	//	clus := device.Kmedoids(availDevs, cf.Dominate)
	//	var idMap map[string]int = make(map[string]int)
	//	for len(idMap) < cf.Quantity {
	//		for _, clu := range clus {
	//			rand.Seed(time.Now().UnixNano())
	//			r := rand.Intn(len(clu.Member))
	//			idMap[clu.Member[r]] = 1
	//		}
	//	}
	return idList
}

//Specify attribute
type SpecifyAttrFilter struct {
	Quantity int
	Attr     string
	Value    string
}

func (cf SpecifyAttrFilter) GetDeviceSet(availDevs []comp.Device) []string {
	var idList []string = make([]string, 0)
	return idList
}

func RigisterGob() {
	gob.Register(SpecifyAttrFilter{})
	gob.Register(SpecifyDevFilter{})
	gob.Register(CompatibilityFilter{})
}

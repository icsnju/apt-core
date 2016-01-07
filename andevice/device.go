package andevice

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

//device state
const (
	DEVICE_RUN = iota
	DEVICE_FREE
)

const (
	DEV_ATTR_ID           = "Id"
	DEV_ATTR_MANU         = "Manu"
	DEV_ATTR_MODEL        = "Model"
	DEV_ATTR_API          = "Api"
	DEV_ATTR_BUILDVERSION = "BuildVersion"
	DEV_ATTR_CPUABI       = "CpuAbi"

	MAX_DIS = 0

	MAXVALUE = 1000
)

var ATTRCOUNT float64 = 5

//device infomation
//serialNumber, manufacturer, model, api, buildVersion, cpuAbi
type DeviceInfo struct {
	Id           string
	Manu         string
	Model        string
	Api          string
	BuildVersion string
	CpuAbi       string
}

//device detail
//device information, slave ip, state
type Device struct {
	Info DeviceInfo
	//IP    string
	State int
}

//device slice
type DeviceInfoSlice struct {
	DeviceInfos []DeviceInfo
}

type Cluster struct {
	Medoid string
	Member []string
}

//the distance between two devices
func (dv Device) Distance(other Device, attr string) float64 {
	return DistanceAttr(dv.Info, other.Info, attr)
}

//distance between two deviceinfo
func Distance(one, two DeviceInfo) float64 {
	var dist float64
	sim := 1/ATTRCOUNT*simManu(one.Manu, two.Manu) + 1/ATTRCOUNT*simModel(one.Model, two.Model) + 1/ATTRCOUNT*simApi(one.Api, two.Api) + 1/ATTRCOUNT*simBuildVersion(one.BuildVersion, two.BuildVersion) + 1/ATTRCOUNT*simCpuAbi(one.CpuAbi, two.CpuAbi)
	if sim == 0 {
		dist = MAXVALUE
	} else {
		dist = 1/sim - 1
	}
	return dist
}

//simlarity between two deviceinfo
func Simlarity(one, two DeviceInfo) float64 {

	sim := 1/ATTRCOUNT*simManu(one.Manu, two.Manu) + 1/ATTRCOUNT*simModel(one.Model, two.Model) + 1/ATTRCOUNT*simApi(one.Api, two.Api) + 1/ATTRCOUNT*simBuildVersion(one.BuildVersion, two.BuildVersion) + 1/ATTRCOUNT*simCpuAbi(one.CpuAbi, two.CpuAbi)

	return sim
}

func DistanceAttr(one, two DeviceInfo, attr string) float64 {
	var dist float64 = 0
	switch attr {
	case DEV_ATTR_MANU:
		sim := simManu(one.Manu, two.Manu)
		if sim == 0 {
			dist = MAXVALUE
		} else {
			dist = 1/sim - 1
		}
	case DEV_ATTR_MODEL:
		sim := simModel(one.Model, two.Model)
		if sim == 0 {
			dist = MAXVALUE
		} else {
			dist = 1/sim - 1
		}
	case DEV_ATTR_API:
		sim := simApi(one.Api, two.Api)
		if sim == 0 {
			dist = MAXVALUE
		} else {
			dist = 1/sim - 1
		}
	case DEV_ATTR_BUILDVERSION:
		sim := simBuildVersion(one.BuildVersion, two.BuildVersion)
		if sim == 0 {
			dist = MAXVALUE
		} else {
			dist = 1/sim - 1
		}
	case DEV_ATTR_CPUABI:
		sim := simCpuAbi(one.CpuAbi, two.CpuAbi)
		if sim == 0 {
			dist = MAXVALUE
		} else {
			dist = 1/sim - 1
		}
	default:
		dist = Distance(one, two)
	}
	return dist
}

func simSeri(m1, m2 string) float64 {
	if strings.ToLower(m1) == strings.ToLower(m2) {
		return 1
	}
	return 0
}

func simManu(m1, m2 string) float64 {
	if m1 == m2 {
		return 1
	}
	return 0
}

func simModel(m1, m2 string) float64 {
	if strings.ToLower(m1) == strings.ToLower(m2) {
		return 1
	}
	return 0
}

func simApi(m1, m2 string) float64 {
	i1, err1 := strconv.Atoi(m1)
	i2, err2 := strconv.Atoi(m2)
	if err1 != nil || err2 != nil {
		fmt.Println("Error in distBuildVersion")
		return 0
	}
	var dis float64
	if i1 > i2 {
		dis = float64(i1) - float64(i2)
	} else {
		dis = float64(i2) - float64(i1)
	}
	return 1 / (dis + 1)

}

func simBuildVersion(m1, m2 string) float64 {
	if m1 == m2 {
		return 1
	}
	return 0
}

func simCpuAbi(m1, m2 string) float64 {
	if m1 == m2 {
		return 1
	}
	return 0
}

//K-medoids algorithm
func Kmedoids(devMap map[string]Device, attr string) []Cluster {
	size := len(devMap)
	var bestJ float64 = -1
	var result []Cluster
	MAX_CLUSTER := 9
	if size <= MAX_CLUSTER {
		fmt.Println("Devices are not enough!")
		return result
	}

	var idList []string = make([]string, size)
	count := 0
	for key, _ := range devMap {
		idList[count] = key
		count++
	}

	TIMES := 50
	for k := 2; k < MAX_CLUSTER; k++ {
		for t := 0; t < TIMES; t++ {
			var clus []Cluster = make([]Cluster, k)

			var indexMap map[int]int = make(map[int]int)

			//random medoids
			for len(indexMap) < k {
				rand.Seed(time.Now().UnixNano())
				r := rand.Intn(size)
				indexMap[r] = 1
			}

			//find medoids
			count = 0
			for key, _ := range indexMap {
				clus[count].Medoid = idList[key]
				count++
			}

			change := true
			for change {
				//find cluster
				for _, id := range idList {

					//clear my member
					for _, cl := range clus {
						cl.Member = make([]string, 0)
					}
					mClu := 0
					var minDis float64 = -1
					dev1 := devMap[id]
					//find this id in which cluster
					for j, cl := range clus {
						if cl.Medoid == id {
							mClu = j
							break
						} else {
							dev2, ex := devMap[cl.Medoid]
							if !ex {
								fmt.Println("Error! Kmedoids")
								return result
							}
							dis := dev1.Distance(dev2, attr)
							if minDis == -1 || dis < minDis {
								minDis = dis
								mClu = j
							}
						}
					}
					clus[mClu].Member = append(clus[mClu].Member, id)
				}

				change = false
				//udpate medoid
				for i, clu := range clus {
					var minDisSum float64 = -1
					newMed := clu.Medoid
					for _, id1 := range clu.Member {
						var mSum float64 = 0
						dv1 := devMap[id1]
						for _, id2 := range clu.Member {
							dv2 := devMap[id2]
							mSum += dv1.Distance(dv2, attr)
						}
						if minDisSum == -1 || mSum < minDisSum {
							minDisSum = mSum
							newMed = id1
						}
					}

					if newMed != clu.Medoid {
						clus[i].Medoid = newMed
						change = true
					}
				}
			}
			//calculate the mJ
			var disSum float64 = 0
			for _, clu := range clus {
				dv1 := devMap[clu.Medoid]
				for _, m := range clu.Member {
					if m != clu.Medoid {
						dv2 := devMap[m]
						disSum += dv1.Distance(dv2, attr)
					}
				}
			}
			disSum = disSum / float64(size-k)
			if bestJ == -1 || bestJ > disSum {
				bestJ = disSum
				result = clus
			}
		}
	}

	return result
}

//find device close to me
func FindSimilarDevs(id string, devMap map[string]Device, maxDis float64) map[string]int {
	me, _ := devMap[id]
	idMap := make(map[string]int, 0)

	for key, dev := range devMap {
		if Simlarity(me.Info, dev.Info) >= maxDis {
			idMap[key] = 1
		}
	}
	return idMap
}

//K-medoids algorithm
func Kmedoids2(devMap map[string]Device, clunum int) []Cluster {
	size := len(devMap)
	var bestJ float64 = -1
	var result []Cluster
	if size <= clunum {
		fmt.Println("Devices are not enough!")
		return result
	}

	var idList []string = make([]string, size)
	count := 0
	for key, _ := range devMap {
		idList[count] = key
		count++
	}

	TIMES := 50

	for t := 0; t < TIMES; t++ {
		var clus []Cluster = make([]Cluster, clunum)

		var indexMap map[int]int = make(map[int]int)

		//random medoids
		for len(indexMap) < clunum {
			rand.Seed(time.Now().UnixNano())
			r := rand.Intn(size)
			indexMap[r] = 1
		}

		//find medoids
		count = 0
		for key, _ := range indexMap {
			clus[count].Medoid = idList[key]
			count++
		}

		change := true
		for change {
			//find cluster
			for _, id := range idList {

				//clear my member
				for _, cl := range clus {
					cl.Member = make([]string, 0)
				}
				mClu := 0
				var minDis float64 = -1
				dev1 := devMap[id]
				//find this id in which cluster
				for j, cl := range clus {
					if cl.Medoid == id {
						mClu = j
						break
					} else {
						dev2, ex := devMap[cl.Medoid]
						if !ex {
							fmt.Println("Error! Kmedoids")
							return result
						}
						dis := dev1.Distance(dev2, "")
						if minDis == -1 || dis < minDis {
							minDis = dis
							mClu = j
						}
					}
				}
				clus[mClu].Member = append(clus[mClu].Member, id)
			}

			change = false
			//udpate medoid
			for i, clu := range clus {
				var minDisSum float64 = -1
				newMed := clu.Medoid
				for _, id1 := range clu.Member {
					var mSum float64 = 0
					dv1 := devMap[id1]
					for _, id2 := range clu.Member {
						dv2 := devMap[id2]
						mSum += dv1.Distance(dv2, "")
					}
					if minDisSum == -1 || mSum < minDisSum {
						minDisSum = mSum
						newMed = id1
					}
				}

				if newMed != clu.Medoid {
					clus[i].Medoid = newMed
					change = true
				}
			}
		}
		//calculate the mJ
		var disSum float64 = 0
		for _, clu := range clus {
			dv1 := devMap[clu.Medoid]
			for _, m := range clu.Member {
				if m != clu.Medoid {
					dv2 := devMap[m]
					disSum += dv1.Distance(dv2, "")
				}
			}
		}
		disSum = disSum / float64(size-clunum)
		if bestJ == -1 || bestJ > disSum {
			bestJ = disSum
			result = clus
		}
	}

	return result
}

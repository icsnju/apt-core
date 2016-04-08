package framework

import (
	"apsaras/comm"
	"fmt"
	"os"
	"path"
)

const OUTPATH = "out.txt"

//Robotium
//a test framework based on Robotium
type MonkeyFrame struct {
	AppPath string
	PkgName string
	Argu    string
}

//Roborium executor
func (mf MonkeyFrame) TaskExecutor(jobId, deviceId, sdkPath string) {
	outPath := path.Join(jobId, deviceId, OUTPATH)

	file, err := os.Create(outPath)
	if err != nil {
		fmt.Println("InstallFrame create out file err!")
		fmt.Println(err)
	}
	var adb = path.Join(sdkPath, ADB_PATH)
	var out string = ""

	cmd := adb + " -s" + deviceId + " uninstall " + mf.PkgName
	comm.ExeCmd(cmd)
	cmd = adb + " -s " + deviceId + " install " + mf.AppPath
	out += comm.ExeCmd(cmd)
	//	cmd = "sleep 120"
	//	comm.ExeCmd(cmd)
	cmd = adb + " -s " + deviceId + " shell monkey -p " + mf.PkgName + " " + mf.Argu
	out += comm.ExeCmd(cmd)
	cmd = adb + " -s " + deviceId + " uninstall " + mf.PkgName
	comm.ExeCmd(cmd)

	file.WriteString(out)
	file.Sync()
	file.Close()
}

//move test file to target file
func (mf MonkeyFrame) MoveTestFile(disPath string) FrameStruct {
	jobPath := path.Join(disPath, comm.APPNAME)
	cmd := "cp " + mf.AppPath + " " + jobPath
	comm.ExeCmd(cmd)
	return mf
}

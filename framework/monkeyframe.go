package framework

import (
	"apsaras/tools"
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
func (mf MonkeyFrame) TaskExecutor(jobId, deviceId string) {
	outPath := path.Join(jobId, deviceId, OUTPATH)

	file, err := os.Create(outPath)
	if err != nil {
		fmt.Println("InstallFrame create out file err!")
		fmt.Println(err)
	}

	var out string = ""

	cmd := "adb -s" + deviceId + " uninstall " + mf.PkgName
	tools.ExeCmd(cmd)
	cmd = "adb -s " + deviceId + " install " + mf.AppPath
	tools.ExeCmd(cmd)
	cmd = "sleep 120"
	tools.ExeCmd(cmd)
	cmd = "adb -s " + deviceId + " shell monkey -p " + mf.PkgName + " " + mf.Argu
	out += tools.ExeCmd(cmd)
	cmd = "adb -s " + deviceId + " uninstall " + mf.PkgName
	tools.ExeCmd(cmd)

	file.Sync()
	file.Close()
}

//move test file to target file
func (mf MonkeyFrame) MoveTestFile(disPath string) FrameStruct {
	jobPath := path.Join(disPath, tools.APPNAME)
	cmd := "cp " + mf.AppPath + " " + jobPath
	tools.ExeCmd(cmd)
	mf.AppPath = jobPath

	return mf
}

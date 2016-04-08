package framework

import (
	"apsaras/comm"
	"path"
)

//Robotium
//a test framework based on Robotium
type RobotFrame struct {
	AppPath  string
	TestPath string
}

//Roborium executor
func (bf RobotFrame) TaskExecutor(jobId, deviceId, sdkPath string) {
	outPath := path.Join(jobId, deviceId)
	cmd := "java  -Djava.awt.headless=true -jar spoon-runner.jar --apk " + bf.AppPath
	cmd += " --test-apk " + bf.TestPath + " --output " + outPath + " --sdk " + sdkPath
	cmd += " -serial " + deviceId
	comm.ExeCmd(cmd)
}

//move test file to target file
func (bf RobotFrame) MoveTestFile(disPath string) FrameStruct {
	jobPath := path.Join(disPath, comm.APPNAME)
	cmd := "cp " + bf.AppPath + " " + jobPath
	comm.ExeCmd(cmd)
	bf.AppPath = jobPath

	jobPath = path.Join(disPath, comm.TESTNAME)
	cmd = "cp " + bf.TestPath + " " + jobPath
	comm.ExeCmd(cmd)
	bf.TestPath = jobPath
	return bf
}

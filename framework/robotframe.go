package framework

import (
	"nata/tools"
	"path"
)

//Robotium
//a test framework based on Robotium
type RobotFrame struct {
	AppPath  string
	TestPath string
}

//Roborium executor
func (bf RobotFrame) TaskExecutor(jobId, deviceId string) {
	outPath := path.Join(jobId, deviceId)
	cmd := "java  -Djava.awt.headless=true -jar spoon-runner.jar --apk " + bf.AppPath
	cmd += " --test-apk " + bf.TestPath + " --output " + outPath + " --sdk /opt/android-sdk/ "
	cmd += "-serial " + deviceId
	tools.ExeCmd(cmd)
}

//move test file to target file
func (bf RobotFrame) MoveTestFile(disPath string) FrameStruct {
	jobPath := path.Join(disPath, tools.APPNAME)
	cmd := "cp " + bf.AppPath + " " + jobPath
	tools.ExeCmd(cmd)
	bf.AppPath = jobPath

	jobPath = path.Join(disPath, tools.TESTNAME)
	cmd = "cp " + bf.TestPath + " " + jobPath
	tools.ExeCmd(cmd)
	bf.TestPath = jobPath
	return bf
}

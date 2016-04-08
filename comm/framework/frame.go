package framework

import "encoding/gob"

const (
	FRAME_ROBOT   = "robotium"
	FRAME_MONKEY  = "monkey"
	FRAME_INSTALL = "install"
	ADB_PATH      = "platform-tools/adb"
)

//Framework
//framework interface
type FrameStruct interface {
	TaskExecutor(jobId, deviceId, sdkPath string)
	MoveTestFile(disPath string) FrameStruct
}

func RigisterGob() {
	gob.Register(RobotFrame{})
	gob.Register(MonkeyFrame{})
	gob.Register(InstallFrame{})
}

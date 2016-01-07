package framework

import (
	"encoding/gob"
)

const (
	FRAME_ROBOT   = "robotium"
	FRAME_MONKEY  = "monkey"
	FRAME_INSTALL = "install"
)

//Framework
//framework interface
type FrameStruct interface {
	TaskExecutor(jobId, deviceId string)
	MoveTestFile(disPath string) FrameStruct
}

func RigisterGob() {
	gob.Register(RobotFrame{})
	gob.Register(MonkeyFrame{})
	gob.Register(InstallFrame{})
	gob.Register(SpecifyAttrFilter{})
	gob.Register(SpecifyDevFilter{})
	gob.Register(CompatibilityFilter{})
}

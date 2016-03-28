package comm

import (
	"log"
	"os"
	"os/exec"
	"strings"
)

const (
	APPNAME  = "App.apk"
	TESTNAME = "Test.apk"

	MASTER = "master"
	SLAVE  = "slave"
)

//Execute cmdline
func ExeCmd(cmd string) string {
	log.Println("command is ", cmd)
	// splitting head => g++ parts => rest of the command
	parts := strings.Fields(cmd)
	head := parts[0]
	parts = parts[1:len(parts)]

	out, err := exec.Command(head, parts...).Output()
	var sout string = ""
	if err != nil {
		log.Println(err)
		sout += err.Error() + "\n"
	}
	return string(out)
}

// exists returns whether the given file or directory exists or not
func FileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

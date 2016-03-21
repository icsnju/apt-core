package comm

import (
	"fmt"
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

//Check if err is nil
func CheckError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}

//Execute cmdline
func ExeCmd(cmd string) string {
	fmt.Println("command is ", cmd)
	// splitting head => g++ parts => rest of the command
	parts := strings.Fields(cmd)
	head := parts[0]
	parts = parts[1:len(parts)]

	out, err := exec.Command(head, parts...).Output()
	var sout string = ""
	if err != nil {
		fmt.Printf("%s", err)
		sout += err.Error() + "\n"
	}
	fmt.Printf("%s", out)

	sout += string(out)
	return sout
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

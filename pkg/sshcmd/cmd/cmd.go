package cmd

import (
	"bytes"
	"os"
	"os/exec"

	"github.com/currycan/supkube/pkg/logger"
)

//Cmd is exec on os ,no return
func Cmd(name string, arg ...string) {
	logger.Info("[os]exec cmd is : ", name, arg)
	cmd := exec.Command(name, arg[:]...)
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		logger.Error("[os]os call error.", err)
	}
}

//String is exec on os , return result
func String(name string, arg ...string) string {
	logger.Info("[os]exec cmd is : ", name, arg)
	cmd := exec.Command(name, arg[:]...)
	cmd.Stdin = os.Stdin
	var b bytes.Buffer
	cmd.Stdout = &b
	cmd.Stderr = &b
	if err := cmd.Run(); err != nil {
		logger.Error("[os]os call error.", err)
		return ""
	}
	return b.String()
}

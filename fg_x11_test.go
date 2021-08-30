package goforeground

import (
	"os/exec"
	"testing"
	"time"
)

func TestActivateByPID(t *testing.T) {
	//Functional test done on debian with XDE
	t.Skip()
	cmd := exec.Command("kate")
	cmd.Start()
	time.Sleep(5 * time.Second)
	err := activate(cmd.Process.Pid)
	if err != nil {
		t.Error("expected no error")
	}
}

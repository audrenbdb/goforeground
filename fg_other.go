//+build !linux,!windows

package goforeground

import (
	"os"
)

func activate(windowName string, pid int, callback func() error) error {
	process, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	process.Kill()
	return callback()
}


//+build !linux,!windows

package goforeground

func activateByWindowTitle(windowName string) error {
	return nil
}

func activateByWindowPID(pid int) error {
	return nil
}


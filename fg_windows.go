//+build windows

package goforeground

import (
	"fmt"
	"syscall"
	"unsafe"
)

const (
	//Activates and displays the window.
	//If the window is minimized or maximized,
	//the system restores it to its original size and position
	swRestore = 9
)

var (
	user32                  = syscall.MustLoadDLL("user32.dll")
	procEnumWindows         = user32.MustFindProc("EnumWindows")
	procGetWindowPID = user32.MustFindProc("GetWindowThreadProcessId")
	procIsWindowVisible = user32.MustFindProc("IsWindowVisible")
	procGetWindow = user32.MustFindProc("GetWindow")


	procSetForegroundWindow = user32.MustFindProc("SetForegroundWindow")
	procShowWindow          = user32.MustFindProc("ShowWindow")

)

func isMainWindow(hwnd syscall.Handle) bool {
	isVisible, _, _ := procIsWindowVisible.Call(uintptr(hwnd))
	gwOwner := uintptr(4)
	isOwned, _, _ := procGetWindow.Call(uintptr(hwnd), gwOwner)
	return isVisible == 1 && isOwned == 0
}

func enumWindows(enumFunc uintptr, lparam uintptr) (err error) {
	r1, _, e1 := syscall.Syscall(procEnumWindows.Addr(), 2, uintptr(enumFunc), uintptr(lparam), 0)
	if r1 == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}


func getWindowThreadProcessId(hwnd syscall.Handle) (int, error) {
	var pid uintptr = 0
	_, _, err := procGetWindowPID.Call(uintptr(hwnd), uintptr(unsafe.Pointer(&pid)))
	return int(pid), err
}

func findWindow(pid int) (syscall.Handle, error) {
	var hwnd syscall.Handle
	cb := syscall.NewCallback(func(h syscall.Handle, p uintptr) uintptr {
		if pidMatch(h, pid) && isMainWindow(h) {
			hwnd = h
			return 0
		}
		return 1
	})
	enumWindows(cb, 0)
	if hwnd == 0 {
		return 0, fmt.Errorf("No window with pid %d found", pid)
	}
	return hwnd, nil
}

func pidMatch(h syscall.Handle, pid int) bool {
	wpid, _ := getWindowThreadProcessId(h)
	return pid == wpid && wpid != 0
}

func setForeground(h syscall.Handle) error {
	procShowWindow.Call(uintptr(h), swRestore)
	procSetForegroundWindow.Call(uintptr(h))
	return nil
}

func activate(pid int) error {
	h, err := findWindow(pid)
	if err != nil {
		return err
	}
	return setForeground(h)
}

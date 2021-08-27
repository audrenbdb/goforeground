//+build windows

package goforeground

import (
	"fmt"
	"strings"
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
	procGetWindowTextW      = user32.MustFindProc("GetWindowTextW")
	procGetWindowPID = user32.MustFindProc("GetWindowThreadProcessId")
	procSetForegroundWindow = user32.MustFindProc("SetForegroundWindow")
	procShowWindow          = user32.MustFindProc("ShowWindow")
)

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

func getWindowText(hwnd syscall.Handle, str *uint16, maxCount int32) (len int32, err error) {
	r0, _, e1 := syscall.Syscall(procGetWindowTextW.Addr(), 3, uintptr(hwnd), uintptr(unsafe.Pointer(str)), uintptr(maxCount))
	len = int32(r0)
	if len == 0 {
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

func findWindow(title string, pid int) (syscall.Handle, error) {
	var hwnd syscall.Handle
	cb := syscall.NewCallback(func(h syscall.Handle, p uintptr) uintptr {
		if pidMatch(h, pid) {
			if titleMatch(h, title) {
				hwnd = h
				return 0
			}
		}
		return 1
	})
	enumWindows(cb, 0)
	if hwnd == 0 {
		return 0, fmt.Errorf("No window with pid %d found", pid)
	}
	return hwnd, nil
}

func titleMatch(h syscall.Handle, title string) bool {
	b := make([]uint16, 200)
	_, err := getWindowText(h, &b[0], int32(len(b)))
	if err != nil {
		return false
	}
	return strings.Contains(syscall.UTF16ToString(b), title)
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

func activate(title string, pid int) error {
	h, err := findWindow(title, pid)
	if err != nil {
		return err
	}
	return setForeground(h)
}

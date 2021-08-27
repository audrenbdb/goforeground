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

func findWindowByTitle(title string) (syscall.Handle, error) {
	var hwnd syscall.Handle
	cb := syscall.NewCallback(func(h syscall.Handle, p uintptr) uintptr {
		b := make([]uint16, 200)
		_, err := getWindowText(h, &b[0], int32(len(b)))
		if err != nil {
			// ignore the error
			return 1 // continue enumeration
		}
		if syscall.UTF16ToString(b) == title {
			// note the window
			hwnd = h
			return 0 // stop enumeration
		}
		return 1 // continue enumeration
	})
	enumWindows(cb, 0)
	if hwnd == 0 {
		return 0, fmt.Errorf("No window with title '%s' found", title)
	}
	return hwnd, nil
}

func findWindowByPID(pid int) (syscall.Handle, error) {
	var hwnd syscall.Handle
	cb := syscall.NewCallback(func(h syscall.Handle, p uintptr) uintptr {
		wpid, _ := getWindowThreadProcessId(h)
		if pid == wpid && wpid != 0 {
			hwnd = h
			return 0 // stop enumeration
		}
		return 1 // continue enumeration
	})
	enumWindows(cb, 0)
	if hwnd == 0 {
		return 0, fmt.Errorf("No window with pid %d found", pid)
	}
	return hwnd, nil
}

func setForeground(h syscall.Handle) error {
	procShowWindow.Call(uintptr(h), swRestore)
	procSetForegroundWindow.Call(uintptr(h))
	return nil
}

func activateByWindowTitle(name string) error {
	h, err := findWindowByTitle(name)
	if err != nil {
		return err
	}
	return setForeground(h)
}

func activateByPID(pid int) error {
	h, err := findWindowByPID(pid)
	if err != nil {
		return err
	}
	return setForeground(h)
}

//+build !linux,!windows

package goforeground

func activate(title string, pid int) error {
	return nil
}

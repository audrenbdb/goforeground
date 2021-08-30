//+build !linux,!windows

package goforeground

func activate(pid int) error {
	return nil
}

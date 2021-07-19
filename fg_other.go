//+build !linux,!windows

package goforeground

func activate(windowName string) error {
	return nil
}


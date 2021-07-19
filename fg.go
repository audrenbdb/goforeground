package goforeground

//Activate places a window foreground
func Activate(windowName string, pid int, callback func() error) error {
	return activate(windowName, pid, callback)
}

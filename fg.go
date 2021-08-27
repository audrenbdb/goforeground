package goforeground

//ActivateByWindowTitle attempts to set a window foreground by its title
func ActivateByWindowTitle(windowName string) error {
	return activateByWindowTitle(windowName)
}

//ActivateByPID attempts to set a window foreground by its process ID
func ActivateByPID(pid int) error {
	return activateByPID(pid)
}

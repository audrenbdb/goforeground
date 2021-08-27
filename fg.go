package goforeground

//Activate attempts to set a window foreground by its title and its process id
func Activate(title string, pid int) error {
	return activate(title, pid)
}

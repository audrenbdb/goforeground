package goforeground

//Activate attemps to place a window foreground.
func Activate(windowName string) error {
	return activate(windowName)
}

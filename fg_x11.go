//+build linux,!android,!nox11 freebsd openbsd !windows

package goforeground

/*
#cgo openbsd CFLAGS: -I/usr/X11R6/include -I/usr/local/include
#cgo openbsd LDFLAGS: -L/usr/X11R6/lib -L/usr/local/lib
#cgo freebsd openbsd LDFLAGS: -lX11 -lxkbcommon -lxkbcommon-x11 -lX11-xcb
#cgo linux pkg-config: x11

#include <stdlib.h>
#include <string.h>
#include <X11/Xlib.h>			// `apt-get install libx11-dev`
#include <X11/Xatom.h>

Window *getDisplayWindows (Display *disp, unsigned long *len);
char *winame (Display *disp, Window win);
void activateWindow(Display *display, Window window);

void activateWindowName(Display *disp, char *title) {
    int i;
    unsigned long len;
    Window *windows;
    char *name;

    windows = (Window*)getDisplayWindows(disp,&len);
    for (i=0;i<(int)len;i++) {
        name = winame(disp,windows[i]);
		if (strcmp(name, title) == 0) {
			free(name);
			activateWindow(disp, windows[i]);
			break;
		}
        free(name);
    }
	XFree(windows);
	return;
}

void activateWindow(Display *display, Window window) {
	XWindowAttributes attr = { 0 };
	XGetWindowAttributes(display, window, &attr);
	int s = XScreenNumberOfScreen(attr.screen);

	Atom prop = XInternAtom(display,"_NET_ACTIVE_WINDOW",False), type;

	XClientMessageEvent e = { 0 };
	e.window = window;
	e.format = 32;
	e.message_type = prop;
	e.display = display;
	e.type = ClientMessage;
	e.data.l[0] = 2;
	e.data.l[1] = CurrentTime;
	XSendEvent(display, XRootWindow(display, s), False, SubstructureNotifyMask | SubstructureRedirectMask,
			(XEvent*) &e);
	return;
}

Window *getDisplayWindows (Display *disp, unsigned long *len) {
    Atom prop = XInternAtom(disp,"_NET_CLIENT_LIST",False), type;
    int form;
    unsigned long remain;
    unsigned char *list;

    if (XGetWindowProperty(disp,XDefaultRootWindow(disp),prop,0,1024,False,XA_WINDOW,
                &type,&form,len,&remain,&list) != Success) {
        return 0;
    }

    return (Window*)list;
}


char *winame (Display *disp, Window win) {
    Atom prop = XInternAtom(disp,"WM_NAME",False), type;
    int form;
    unsigned long remain, len;
    unsigned char *list;
 	if (XGetWindowProperty(disp,win,prop,0,1024,False,AnyPropertyType,
                &type,&form,&len,&remain,&list) != Success) {
        return NULL;
    }

    return (char *)list;
}
*/
import "C"

func activate(windowName string, pid int, callback func() error) error {
	display := C.XOpenDisplay(nil)
	defer C.XCloseDisplay(display)
	C.activateWindowName(display, C.CString(windowName))
	return nil
}
//go:build (linux && !android && !nox11) || freebsd || openbsd || !windows
// +build linux,!android,!nox11 freebsd openbsd !windows

package goforeground

/*
#cgo openbsd CFLAGS: -I/usr/X11R6/include -I/usr/local/include
#cgo openbsd LDFLAGS: -L/usr/X11R6/lib -L/usr/local/lib
#cgo freebsd openbsd LDFLAGS: -lX11 -lxkbcommon -lxkbcommon-x11 -lX11-xcb
#cgo linux pkg-config: x11

#include <stdio.h>

#include <stdlib.h>
#include <string.h>
#include <X11/Xlib.h>			// `apt-get install libx11-dev`
#include <X11/Xatom.h>

Window *getDisplayWindows (Display *disp, unsigned long *len);
void activateWindow(Display *display, Window window);
int getWindowPID (Display *disp, Window win);

void activateWindowByPID(Display *disp, int pid) {
    int i;
    unsigned long len;
    Window *windows;
    int wpid;
    windows = (Window*)getDisplayWindows(disp,&len);
    for (i=0;i<(int)len;i++) {
        wpid = getWindowPID(disp, windows[i]);
		if (pid == wpid) {
			activateWindow(disp, windows[i]);
			break;
		}
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

int getWindowPID (Display *disp, Window win) {
    Atom prop = XInternAtom(disp,"_NET_WM_PID", True);

    Atom actual_type_return;
    int actual_format_return;
    unsigned long remaining_bytes_in_prop_return, nitems_return;
    unsigned char *result = NULL;

    if (XGetWindowProperty(disp, win, prop,
            0,                      // in: long_offset (in counts of 32 bits)
            1,                      // in: long_length (in counts of 32 bits)
            False,                  // in: delete
            AnyPropertyType,        // in: req_type
            &actual_type_return,                // out: actual_type_return
            &actual_format_return,              // out: actual_format_return
            &nitems_return,                     // out: nitems_return
            &remaining_bytes_in_prop_return,    // out: bytes_after_return (remaining bytes in property)
            &result                             // out: prop_return
        ) != Success) {
        // Just in case XGetWindowProperty allocated data despite failure, the doc doesn't specify clearly
        if (result != NULL) {
            XFree(result);
        }
        return 0;
    }

    // XGetWindowProperty sometimes returns null with Success
    // No idea why, We can't do anything in this case, so just assume failure
    if (result == NULL) {
        return 0;
    }

    int pid;
    pid = result[1] * 256;
    pid += result[0];

    // XGetWindowProperty: The function returns Success if it executes successfully. To free the resulting data, use XFree.
    XFree(result);

    return pid;
}
*/
import "C"

func activate(pid int) error {
	display := C.XOpenDisplay(nil)
	defer C.XCloseDisplay(display)
	C.activateWindowByPID(display, C.int(pid))
	return nil
}

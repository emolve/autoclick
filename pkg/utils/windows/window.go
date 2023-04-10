package windows

import (
	"fmt"
	"syscall"
	"unsafe"
)

var (
	user32           = syscall.NewLazyDLL("user32.dll")
	enumWindows      = user32.NewProc("EnumWindows")
	getWindowText    = user32.NewProc("GetWindowTextW")
	getWindowTextLen = user32.NewProc("GetWindowTextLengthW")
	getWindowClass   = user32.NewProc("GetClassNameW")
)

func EnumWindows() {
	enumWindows.Call(syscall.NewCallback(enumWindowsProc), 0)
}

func enumWindowsProc(hwnd syscall.Handle, lParam uintptr) uintptr {
	const bufLen = 512
	var buf [bufLen]uint16
	_, _, _ = getWindowTextLen.Call(uintptr(hwnd))
	getWindowText.Call(
		uintptr(hwnd),
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(bufLen))
	windowText := syscall.UTF16ToString(buf[:])
	if len(windowText) > 0 {
		var classNameBuf [bufLen]uint16
		_, _, _ = getWindowClass.Call(uintptr(hwnd), uintptr(unsafe.Pointer(&classNameBuf[0])), uintptr(bufLen))
		className := syscall.UTF16ToString(classNameBuf[:])
		fmt.Printf("Window Name: %s, Class Name: %s\n", windowText, className)
	}
	return 1
}

var (
	findWindow = user32.NewProc("FindWindowW")
	showWindow = user32.NewProc("ShowWindow")
)

// 在windows包中定义的常量
const (
	SW_HIDE            = 0
	SW_SHOWNORMAL      = 1
	SW_SHOWMINIMIZED   = 2
	SW_SHOWMAXIMIZED   = 3
	SW_SHOWNOACTIVATE  = 4
	SW_SHOW            = 5
	SW_MINIMIZE        = 6
	SW_SHOWMINNOACTIVE = 7
	SW_SHOWNA          = 8
	SW_RESTORE         = 9
	SW_SHOWDEFAULT     = 10
	SW_FORCEMINIMIZE   = 11
)

func HideWindow() {
	className := syscall.StringToUTF16Ptr("StandardFrame_DingTalk")
	windowName := syscall.StringToUTF16Ptr("钉钉")
	hwnd, _, _ := findWindow.Call(uintptr(unsafe.Pointer(className)), uintptr(unsafe.Pointer(windowName)))
	if hwnd != 0 {
		showWindow.Call(hwnd, SW_HIDE)
	}
}

func ShowWindow() {
	className := syscall.StringToUTF16Ptr("StandardFrame_DingTalk")
	windowName := syscall.StringToUTF16Ptr("钉钉")
	hwnd, _, _ := findWindow.Call(uintptr(unsafe.Pointer(className)), uintptr(unsafe.Pointer(windowName)))
	if hwnd != 0 {
		showWindow.Call(hwnd, SW_SHOW)
	}
}

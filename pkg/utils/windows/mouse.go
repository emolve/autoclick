package windows

import (
	"fmt"
	"syscall"
)

var (
	m_user32          = syscall.MustLoadDLL("user32.dll")
	kernel32          = syscall.MustLoadDLL("kernel32.dll")
	procSetCursorPos  = m_user32.MustFindProc("SetCursorPos")
	procMouseLeftDown = m_user32.MustFindProc("mouse_event")
	procMouseLeftUp   = m_user32.MustFindProc("mouse_event")
	procGetLastError  = kernel32.MustFindProc("GetLastError")
)

func Click(x, y int) {
	// 设置光标位置
	_, _, err := procSetCursorPos.Call(
		uintptr(x),
		uintptr(y),
	)

	if err != nil && err.Error() != "The operation completed successfully." {
		fmt.Println("SetCursorPos error:", err)
		return
	}

	// 鼠标左键按下
	_, _, err = procMouseLeftDown.Call(
		uintptr(0x0002), // MOUSEEVENTF_LEFTDOWN
		0,
		0,
		0,
		0,
	)

	if err != nil && err.Error() != "The operation completed successfully." {
		fmt.Println("mouse_event error:", err)
		return
	}

	// 鼠标左键抬起
	_, _, err = procMouseLeftUp.Call(
		uintptr(0x0004), // MOUSEEVENTF_LEFTUP
		0,
		0,
		0,
		0,
	)

	if err != nil && err.Error() != "The operation completed successfully." {
		fmt.Println("mouse_event error:", err)
		return
	}

	// 检查错误
	_, _, err = procGetLastError.Call()
	if err != nil && err.Error() != "The operation completed successfully." {
		fmt.Println("GetLastError error:", err)
		return
	}
}

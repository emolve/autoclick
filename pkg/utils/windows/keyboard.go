package windows

import (
	"fmt"
	"syscall"
	"time"
)

var (
	k_user32           = syscall.MustLoadDLL("user32.dll")
	K_kernel32         = syscall.MustLoadDLL("kernel32.dll")
	procKeybdEvent     = k_user32.MustFindProc("keybd_event")
	k_procGetLastError = K_kernel32.MustFindProc("GetLastError")
)

func Space(seconds int) {
	// 模拟键盘空格按下事件
	_, _, err := procKeybdEvent.Call(
		uintptr(0x20), // VK_SPACE
		uintptr(0),
		uintptr(0),
		0,
	)

	if err != nil && err.Error() != "The operation completed successfully." {
		fmt.Println("keybd_event error:", err)
		return
	}

	// 等待 seconds 秒钟
	time.Sleep(time.Duration(seconds) * time.Second)

	// 模拟键盘空格抬起事件
	_, _, err = procKeybdEvent.Call(
		uintptr(0x20), // VK_SPACE
		uintptr(0),
		uintptr(0x0002), // KEYEVENTF_KEYUP
		0,
	)

	if err != nil && err.Error() != "The operation completed successfully." {
		fmt.Println("keybd_event error:", err)
		return
	}

	// 检查错误
	_, _, err = k_procGetLastError.Call()
	if err != nil && err.Error() != "The operation completed successfully." {
		fmt.Println("GetLastError error:", err)
		return
	}
}

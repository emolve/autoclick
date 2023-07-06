package adb

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

const (
	ImagesPath = "D:\\Project\\emolve\\autoclick\\storage\\images\\auto_click.png"
)

func GetDeviceID() (*string, error) {

	// 执行 cmd /c adb devices 命令
	cmd := exec.Command("cmd", "/c", "adb devices")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	// 解析输出结果，提取设备ID
	lines := strings.Split(string(output), "\n")
	if len(lines) < 2 {
		return nil, errors.New("device not found")
	}

	deviceLine := strings.Fields(lines[1])
	if len(deviceLine) < 1 {
		return nil, err
	}

	deviceID := deviceLine[0]
	return &deviceID, nil
}

func Screen(deviceID string) error {
	fmt.Println("deviceID:", deviceID)
	cmd := exec.Command("cmd", "/c", fmt.Sprintf("adb -s %s shell screencap -p /sdcard/Pictures/auto_click.png", deviceID))
	output, err := cmd.Output()
	fmt.Println(string(output))
	cmd = exec.Command("cmd", "/c", fmt.Sprintf("adb.exe -s %s pull /sdcard/Pictures/auto_click.png %s", deviceID, ImagesPath))
	output, err = cmd.Output()
	fmt.Println(string(output))
	if err != nil {
		return err
	}
	return nil
}

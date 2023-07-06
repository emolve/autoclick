package adb

import (
	"fmt"
	"testing"
)

func TestGetDeviceID(t *testing.T) {
	id, err := GetDeviceID()
	if err != nil {
		t.Failed()
	}
	fmt.Println("device id:", *id)
}

func TestScreen(t *testing.T) {
	id, err := GetDeviceID()
	if err != nil {
		t.Failed()
	}
	err = Screen(*id)
	if err != nil {
		t.Failed()
	}
}

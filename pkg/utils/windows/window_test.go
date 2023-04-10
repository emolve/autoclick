package windows

import "testing"

func TestName(t *testing.T) {
	EnumWindows()
}

func TestHideWindow(t *testing.T) {
	HideWindow()
}

func TestShowWindow(t *testing.T) {
	ShowWindow()
}

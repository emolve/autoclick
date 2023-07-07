package main

import (
	"autoclick/global"
	"fmt"
	"testing"
)

func TestSetupSetting(t *testing.T) {
	err := setupSetting()
	if err != nil {
		t.Failed()
	}
	fmt.Println(global.AppSetting.RunMode)
	fmt.Println(global.NotificationSetting.PlusPlus)
	fmt.Println(global.NotificationSetting.Mail)
}

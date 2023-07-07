package main

import (
	"autoclick/global"
	"fmt"
	"testing"
)

func TestPlusPlus(t *testing.T) {

	//plusplus()

}
func TestMail(t *testing.T) {

	//err := Send163Mail()
	//if err != nil {
	//	fmt.Println(err)
	//}

}

func TestGetFormattedName(t *testing.T) {
	name := "autoTest.exe"
	formattedName := getFormattedName(name)
	fmt.Println(formattedName)
}

func TestSetupSetting(t *testing.T) {
	setupSetting()
	fmt.Println(global.AppSetting.RunMode)
	fmt.Println(global.NotificationSetting.PlusPlus)
	fmt.Println(global.NotificationSetting.Mail)
}

package main

import (
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

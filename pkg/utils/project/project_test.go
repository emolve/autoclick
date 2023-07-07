package project

import (
	"fmt"
	"testing"
)

func TestGetFormattedName(t *testing.T) {
	name := "autoTest.exe"
	formattedName := GetFormattedName(name)
	fmt.Println(formattedName)
}

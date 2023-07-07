package project

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func GetProgramName() string {
	// 获取命令行参数
	args := os.Args
	// 第一个参数是程序的名称
	programPath := args[0]
	// 提取文件的基本名称
	name := filepath.Base(programPath)
	return name
}

func GetFormattedName(name string) string {
	fileNameWithoutExt := strings.TrimSuffix(name, ".exe")
	// 使用正则表达式将驼峰式命名转换为下划线格式
	reg := regexp.MustCompile("([a-z0-9])([A-Z])")
	formattedName := reg.ReplaceAllString(fileNameWithoutExt, "${1}_${2}")
	formattedName = strings.ToLower(formattedName)
	return formattedName
}

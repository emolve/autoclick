package main

import (
	"autoclick/global"
	"autoclick/internal/service"
	"autoclick/pkg/setting"
	"autoclick/pkg/utils/project"
	win "autoclick/pkg/windows"
	"fmt"
	"github.com/sirupsen/logrus"
	"golang.org/x/sys/windows/svc"
	"log"
	"os"
)

const (
	LogDir = "C:\\Users\\Administrator\\Documents\\AutoClick\\"
)

var (
	localLogger *logrus.Logger
)

func init() {
	setupLogger()

	err := setupSetting()
	if err != nil {
		log.Fatalf("init.setupSetting err: %v", err)
	}
}

func main() {

	localLogger.Info("-----------------------Starting service-----------------------")

	// 注册服务
	err := svc.Run("", win.InitService(global.AppSetting.RunMode, localLogger, service.DingTask))
	if err != nil {
		localLogger.Errorf("Failed to register service: %v", err)
	}

}

func setupSetting() error {
	setter, err := setting.NewSetting()
	if err != nil {
		return err
	}
	err = setter.ReadSection("Notification", &global.NotificationSetting)
	if err != nil {
		return err
	}
	err = setter.ReadSection("App", &global.AppSetting)
	if err != nil {
		return err
	}
	return nil
}

func setupLogger() *logrus.Logger {
	// get logFile Name
	logFile := fmt.Sprintf("log_%s.txt", project.GetFormattedName(project.GetProgramName()))

	// Create a new logger
	localLogger = logrus.New()

	err := os.MkdirAll(LogDir, 0755)
	if err != nil {
		// 处理创建目录时出现的错误
	}
	// Open a file for writing the log output
	file, err := os.OpenFile(fmt.Sprintf("%s%s", LogDir, logFile), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		localLogger.Fatal(err)
	}
	//defer file.Close()

	// Set the logger output to the file
	localLogger.SetOutput(file)
	return localLogger
}

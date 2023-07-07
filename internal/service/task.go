package service

import (
	"autoclick/global"
	"autoclick/pkg/adb"
	"autoclick/pkg/notification"
	"autoclick/pkg/utils"
	"autoclick/pkg/utils/project"
	"github.com/sirupsen/logrus"
	"os/exec"
	"time"
)

const (
	//  对应构建名称

	ModelMorning = "auto_on"
	ModelNight   = "auto_off"

	DEBUG  = "debug"
	NORMAL = "normal"
	ACTUAL = "actual"

	MailTo = "13735599246@163.com"
)

func DingTask(mode string, logger *logrus.Logger) {
	var ticker *time.Ticker
	var targetHour, targetMinute int
	// 执行任务的逻辑

	runModel := project.GetFormattedName(project.GetProgramName())
	switch runModel {
	case ModelMorning:
		// 指定每天 8:40 执行任务
		targetHour = 8
		targetMinute = 40
		logger.Info("执行上班任务")
	case ModelNight:
		// 指定每天 18:10 执行任务
		targetHour = 18
		targetMinute = 10
		logger.Info("执行下班任务")
	default:
		inAMinute := time.Now().Add(time.Second * 60)
		// 获取小时和分钟
		targetHour = inAMinute.Hour()
		targetMinute = inAMinute.Minute()
		logger.Info("执行其他任务,将在一分钟后执行")
	}

	// 获取当前时间
	now := time.Now()

	// 计算距离下一次执行任务还有多长时间
	targetTime := time.Date(now.Year(), now.Month(), now.Day(), targetHour, targetMinute, 0, 0, time.Local)
	if now.After(targetTime) {
		targetTime = targetTime.Add(24 * time.Hour)
	}
	for targetTime.Weekday() == time.Saturday || targetTime.Weekday() == time.Sunday {
		targetTime = targetTime.AddDate(0, 0, 1)
	}

	duration := targetTime.Sub(now)

	// 设置偏移量
	off, actual := utils.AddRandomTime(duration)

	// 等待指定时间后执行任务

	switch mode {
	case DEBUG:
		ticker = time.NewTicker(time.Second * 5) // 延迟5s开始
		logger.Info("距离下一次执行：", time.Second*5)
		logger.Info("最终执行时间：", now.Add(time.Second*5))
	case NORMAL:
		ticker = time.NewTicker(duration) // 整点
		logger.Info("距离下一次执行：", duration)
		logger.Info("最终执行时间：", targetTime)
	case ACTUAL:
		ticker = time.NewTicker(actual) // 偏移
		logger.Info("本次偏移时间为：", off)
		logger.Info("距离下一次执行：", actual)
		logger.Info("最终执行时间：", targetTime.Add(off))
	}

	select {
	case <-ticker.C:

		var err error
		var deviceID string
		var mailSubject string

		id, err := adb.GetDeviceID()
		if err != nil {
			logger.Error("get devices id failed:", err)
		}
		deviceID = *id

		// 网络检测
		for i := 0; i < 10; i++ {

			err = exec.Command("cmd", "/c", "adb shell curl www.baidu.com").Run() // net check
			if err == nil {
				break
			}
			if err != nil && err.Error() == "exit status 6" {
				logger.Errorf("network checking failed :%v, try again later, times:%d", err, i)
				time.Sleep(time.Second * 3)
			}
		}
		if err != nil {
			mailSubject = "[failed] network failed"
			notification.Plusplus(global.NotificationSetting.PlusPlus, mailSubject)
			notification.Send163Mail(global.NotificationSetting.Mail, MailTo, mailSubject)
		}

		// 发送点亮屏幕信号
		logger.Info("1. power on")
		err = exec.Command("cmd", "/c", "adb  shell input keyevent 26").Run()
		// todo: 查看屏幕状态
		// adb shell dumpsys window policy
		logger.Info("2. back to home")
		err = exec.Command("cmd", "/c", "adb shell input tap 544 2270").Run() // home
		time.Sleep(time.Second * 5)
		logger.Info("3. click app")
		err = exec.Command("cmd", "/c", "adb shell input tap 172 847").Run() // app
		time.Sleep(time.Second * 10)
		logger.Info("4. screen")
		err = adb.Screen(deviceID)
		logger.Info("5. back to home")
		err = exec.Command("cmd", "/c", "adb shell input tap 544 2270").Run() // home
		time.Sleep(time.Second * 5)
		logger.Info("6. view app history")
		err = exec.Command("cmd", "/c", "adb  shell input tap  276 2286").Run() // used app
		time.Sleep(time.Second * 5)
		logger.Info("7. delete app history")
		err = exec.Command("cmd", "/c", "adb shell input tap 533 2044").Run() //close

		// 输出执行结果
		if err != nil {
			logger.Error("adb operation failed", err)
			return
		} else {
			mailSubject = "[successful] adb operation done"
			logger.Info("adb operation done")
		}

		notification.Plusplus(global.NotificationSetting.PlusPlus, mailSubject)
		notification.Send163Mail(global.NotificationSetting.Mail, MailTo, mailSubject)
	}

}

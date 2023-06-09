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

	// 设置偏移量
	randomTime, difference, executionTime, ticker := getTrigger(logger, mode)

	logger.Info("随机偏移时间：", randomTime)
	logger.Infof("最终执行时间：%s，相距：%s", executionTime, difference)

	select {
	case <-ticker.C:
		adbOperation(logger)
	}

}

func getTrigger(logger *logrus.Logger, mode string) (
	randomTime time.Duration,
	difference time.Duration,
	executionTime time.Time,
	ticker *time.Ticker) {
	var targetHour, targetMinute int

	runTime := project.GetFormattedName(project.GetProgramName())
	switch runTime {
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
		logger.Info("未指定执行时间,采用默认策略")
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

	// 目标时间与现在的时间差
	duration := targetTime.Sub(now)

	// 随机偏差， 计划相隔时间
	off, actual := utils.AddRandomTime(duration)

	switch mode {
	case DEBUG:
		return time.Second * 5, time.Second * 5, now.Add(time.Second * 5), time.NewTicker(time.Second * 5)
	case NORMAL:
		return 0, duration, targetTime, time.NewTicker(duration)
	case ACTUAL:
		return off, actual, targetTime.Add(off), time.NewTicker(actual)
	default:
		return time.Second * 5, time.Second * 5, now.Add(time.Second * 5), time.NewTicker(time.Second * 5)
	}
}

func adbOperation(logger *logrus.Logger) {
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
		err = notification.Plusplus(global.NotificationSetting.PlusPlus, mailSubject)
		if err != nil {
			logger.Error("[notification] plus failed", err)
			return
		}
		err = notification.Send163Mail(global.NotificationSetting.Mail, MailTo, mailSubject)
		if err != nil {
			logger.Error("[notification] mail failed", err)
			return
		}
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

	err = notification.Plusplus(global.NotificationSetting.PlusPlus, mailSubject)
	if err != nil {
		logger.Error("[notification] plus failed", err)
		return
	}
	err = notification.Send163Mail(global.NotificationSetting.Mail, MailTo, mailSubject)
	if err != nil {
		logger.Error("[notification] mail failed", err)
		return
	}
}

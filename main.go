package main

import (
	"autoclick/pkg/adb"
	"autoclick/pkg/notification"
	win "autoclick/pkg/windows"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"golang.org/x/sys/windows/svc"
	"io"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

var (
	runModel = getFormattedName(getProgramName())
)

const (
	LogDir = "C:\\Users\\Administrator\\Documents\\AutoClick\\"

	//  对应构建名称

	ModelMorning = "auto_on"
	ModelNight   = "auto_off"

	DEBUG  = "DEBUG"
	NORMAL = "NORMAL"
	ACTUAL = "ACTUAL"
)

func myTask(mode string, logger *logrus.Logger) {
	var ticker *time.Ticker
	var targetHour, targetMinute int
	// 执行任务的逻辑
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
	off, actual := addRandomTime(duration)

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
			plusplus(mailSubject)
			send163Mail(mailSubject)
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

		// screen()
		plusplus(mailSubject)
		send163Mail(mailSubject)
	}

}

func main() {

	// get logFile Name
	programName := getProgramName()
	formattedName := getFormattedName(programName)
	logFile := fmt.Sprintf("log_%s.txt", formattedName)

	// Create a new logger
	localLogger := logrus.New()

	err := os.MkdirAll(LogDir, 0755)
	if err != nil {
		// 处理创建目录时出现的错误
	}
	// Open a file for writing the log output
	file, err := os.OpenFile(fmt.Sprintf("%s%s", LogDir, logFile), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		localLogger.Fatal(err)
	}
	defer file.Close()

	// Set the logger output to the file
	localLogger.SetOutput(file)
	localLogger.Info("-----------------------Starting service-----------------------")
	localLogger.Info("logging at: ", fmt.Sprintf("%s%s", LogDir, logFile))

	// 注册服务
	err = svc.Run("", win.InitService(ACTUAL, localLogger, myTask))

	//err = svc.Run("", &myService{logger: localLogger, mode: ACTUAL})
	if err != nil {
		fmt.Printf("Failed to register service: %v", err)
	}

}
func getProgramName() string {
	// 获取命令行参数
	args := os.Args
	// 第一个参数是程序的名称
	programPath := args[0]
	// 提取文件的基本名称
	name := filepath.Base(programPath)
	return name
}

func getFormattedName(name string) string {
	fileNameWithoutExt := strings.TrimSuffix(name, ".exe")
	// 使用正则表达式将驼峰式命名转换为下划线格式
	reg := regexp.MustCompile("([a-z0-9])([A-Z])")
	formattedName := reg.ReplaceAllString(fileNameWithoutExt, "${1}_${2}")
	formattedName = strings.ToLower(formattedName)
	return formattedName
}

func addRandomTime(duration time.Duration) (randomTime time.Duration, actual time.Duration) {
	if duration < 5*time.Minute {
		// Add a random time between 0 and 3 minutes
		rand.Seed(time.Now().UnixNano())
		randomTime = time.Duration(rand.Int63n(int64(2*time.Minute))) + 3*time.Minute
		duration += randomTime
	} else {
		// Add or subtract a random time between 0 and 3 minutes
		rand.Seed(time.Now().UnixNano())
		randomTime = time.Duration(rand.Int63n(int64(2*time.Minute))) + 3*time.Minute
		if rand.Intn(2) == 0 {
			duration += randomTime
		} else {
			duration -= randomTime
			randomTime = 0 - randomTime
		}
	}
	return randomTime, duration
}

func plusplus(subject string) {
	s := struct {
		Token    string `json:"token"`
		Title    string `json:"title"`
		Content  string `json:"content"`
		Template string `json:"template"`
	}{
		Token:    "xxxxx",
		Title:    subject,
		Content:  time.Now().String() + "打卡成功",
		Template: "html",
	}
	jsonStr, err := json.Marshal(s)
	url := "http://www.pushplus.plus/send"

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := io.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
}

func send163Mail(subject string) error {
	//err := SendMail("xxxx@163.com", "xxxx", "smtp.163.com", "25", "xxxx@163.com", "xxxx@163.com", subject, "11111")
	err := notification.SendMail("13735599246@163.com", "xxx", "smtp.163.com", "25", "13735599246@163.com", "13735599246@163.com", subject, "11111")
	return err
}

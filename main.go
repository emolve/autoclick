package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/kbinani/screenshot"
	"github.com/sirupsen/logrus"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/eventlog"
	"gopkg.in/gomail.v2"
	"image"
	"image/png"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"sync"
	"time"
)

type myService struct {
	logger *logrus.Logger
	mutex  sync.Mutex
	mode   string
}

func (m *myService) Execute(args []string, r <-chan svc.ChangeRequest, changes chan<- svc.Status) (ssec bool, errno uint32) {

	m.logger.Info("start mode: ", m.mode)
	// Open the Windows event log
	el, err := eventlog.Open("autoClickService")
	if err != nil {
		m.logger.WithError(err).Error("failed to open event log")
		return
	}
	defer el.Close()

	// 初始化服务状态
	changes <- svc.Status{State: svc.StartPending}

	// 启动服务
	changes <- svc.Status{State: svc.Running, Accepts: svc.AcceptStop | svc.AcceptShutdown}

	// 处理系统事件
	for {
		select {
		case c := <-r:
			switch c.Cmd {
			case svc.Stop, svc.Shutdown:
				// 停止服务
				changes <- svc.Status{State: svc.StopPending}
				m.logger.Info("-----------------------Stopping service-----------------------")
				return
			default:
				//debug.Log.(fmt.Sprintf("unexpected control request #%d", c))
			}
		case <-time.After(10 * time.Second):
			// 每隔10秒钟执行一次任务
			go func() {
				if !m.mutex.TryLock() {
					// 锁获取失败，直接返回
					return
				}
				defer m.mutex.Unlock()
				m.doTask()
			}()
		}
	}
}

func (m *myService) doTask() {
	var ticker *time.Ticker

	// 执行任务的逻辑
	// 指定每天8:40执行任务
	targetHour := 8
	targetMinute := 40

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

	switch m.mode {
	case DEBUG:
		ticker = time.NewTicker(time.Second * 5) // 延迟5s开始
		m.logger.Info("距离下一次执行：", time.Second*5)
		m.logger.Info("最终执行时间：", now.Add(time.Second*5))
	case NORMAL:
		ticker = time.NewTicker(duration) // 整点
		m.logger.Info("距离下一次执行：", duration)
		m.logger.Info("最终执行时间：", targetTime)
	case ACTUAL:
		ticker = time.NewTicker(actual) // 偏移
		m.logger.Info("本次偏移时间为：", off)
		m.logger.Info("距离下一次执行：", actual)
		m.logger.Info("最终执行时间：", targetTime.Add(off))
	}

	select {
	case <-ticker.C:

		var mailSubject string
		// 网络检测
		var err error
		for i := 0; i < 10; i++ {

			err = exec.Command("cmd", "/c", "adb shell curl www.baidu.com").Run() // net check
			if err == nil {
				break
			}
			if err != nil && err.Error() == "exit status 6" {
				m.logger.Errorf("network checking failed :%v, try again later, times:%d", err, i)
				time.Sleep(time.Second * 3)
			}
		}
		if err != nil {
			mailSubject = "[failed] network failed"
			plusplus(mailSubject)
			send163Mail(mailSubject)
		}

		// 发送点亮屏幕信号
		m.logger.Info("1. power on")
		err = exec.Command("cmd", "/c", "adb  shell input keyevent 26").Run()
		// todo: 查看屏幕状态
		// adb shell dumpsys window policy
		m.logger.Info("2. back to home")
		err = exec.Command("cmd", "/c", "adb shell input tap 544 2270").Run() // home
		time.Sleep(time.Second * 5)
		m.logger.Info("3. click app")
		err = exec.Command("cmd", "/c", "adb shell input tap 172 847").Run() // app
		time.Sleep(time.Second * 10)
		m.logger.Info("4. back to home")
		err = exec.Command("cmd", "/c", "adb shell input tap 544 2270").Run() // home
		time.Sleep(time.Second * 5)
		m.logger.Info("5. view app history")
		err = exec.Command("cmd", "/c", "adb  shell input tap  276 2286").Run() // used app
		time.Sleep(time.Second * 5)
		m.logger.Info("6. delete app history")
		err = exec.Command("cmd", "/c", "adb shell input tap 533 2044").Run() //close

		// 输出执行结果
		if err != nil {
			m.logger.Error("adb operation failed", err)
			return
		} else {
			mailSubject = "[successful] adb operation done"
			m.logger.Info("adb operation done")
		}

		// screen()
		plusplus(mailSubject)
		send163Mail(mailSubject)
	}

}

const (
	logDir  = "C:\\Users\\Administrator\\Documents\\AutoClick\\"
	logFile = "log.txt"

	DEBUG  = "DEBUG"
	NORMAL = "NORMAL"
	ACTUAL = "ACTUAL"
)

func main() {

	// Create a new logger
	localLogger := logrus.New()

	err := os.MkdirAll(logDir, 0755)
	if err != nil {
		// 处理创建目录时出现的错误
	}
	// Open a file for writing the log output
	file, err := os.OpenFile(fmt.Sprintf("%s%s", logDir, logFile), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		localLogger.Fatal(err)
	}
	defer file.Close()

	// Set the logger output to the file
	localLogger.SetOutput(file)
	localLogger.Info("-----------------------Starting service-----------------------")
	localLogger.Info("logging at: ", fmt.Sprintf("%s%s", logDir, logFile))

	// 注册服务
	err = svc.Run("", &myService{logger: localLogger, mode: ACTUAL})
	if err != nil {
		fmt.Printf("Failed to register service: %v", err)
	}

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
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
}

func send163Mail(subject string) error {
	err := SendMail("xxxx@163.com", "xxxx", "smtp.163.com", "25", "xxxx@163.com", "xxxx@163.com", subject, "11111")
	return err
}

func SendMail(userName, authCode, host, portStr, mailTo, sendName string, subject, body string) error {
	port, _ := strconv.Atoi(portStr)
	m := gomail.NewMessage()
	m.SetHeader("From", m.FormatAddress(userName, sendName))
	m.SetHeader("To", mailTo)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)
	//m.Embed("C:\\Users\\Administrator\\Pictures\\uToolsWallpapers\\wallhaven-yje5gg.jpg") // 图片路径
	//m.SetBody("text/html", `<img src="cid:wallhaven-yje5gg.jpg" alt="My image" />`)       //设置邮件正文

	d := gomail.NewDialer(host, port, userName, authCode)
	err := d.DialAndSend(m)
	return err
}

func screen() {
	//使用 GetDisplayBounds获取指定屏幕显示范围，全屏截图
	bounds := screenshot.GetDisplayBounds(0)
	img, err := screenshot.CaptureRect(bounds)
	if err != nil {
		panic(err)
	}
	//拼接图片名
	t := time.Now().Unix()
	tt := strconv.Itoa(int(t)) + ".png"

	save(img, tt)
}

// save *image.RGBA to filePath with PNG format.
func save(img *image.RGBA, filePath string) {
	file, err := os.Create(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	png.Encode(file, img)
}

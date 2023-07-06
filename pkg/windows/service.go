package windows

import (
	"github.com/sirupsen/logrus"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/eventlog"
	"sync"
	"time"
)

type Task func(mode string, logger *logrus.Logger)

type Service struct {
	logger  *logrus.Logger
	mutex   sync.Mutex
	mode    string
	crontab Task
}

func (m *Service) Execute(args []string, r <-chan svc.ChangeRequest, changes chan<- svc.Status) (ssec bool, errno uint32) {

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
				m.crontab(m.mode, m.logger)
			}()
		}
	}
}

func InitService(mode string, logger *logrus.Logger, task Task) *Service {

	return &Service{logger: logger, mode: mode, crontab: task}
}

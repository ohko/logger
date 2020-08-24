package logger

import (
	"testing"
	"time"
)

// go test -run TestNewMonitor -v -count=1
func TestNewMonitor(t *testing.T) {
	o := NewMonitor(&MonitorOption{
		ID:         123,               // 标识符
		LogPath:    "./log",           // 要监控的日志目录
		MaxSize:    100 * 1024 * 1024, // 100MB
		NotifyRate: time.Minute,       // 监控频率

		// 钉钉Webhook通知
		DingDing: "https://oapi.dingtalk.com/robot/send?access_token=xxxxx",

		// MailAddr: "smtp.exmail.qq.com:465",
		// MailUser: "xxx@qq.com",
		// MailPass: "123",
		// MailName: "XXX",
		// ToAddr:   "to@qq.com",
	})
	t.Log(o.GetSize("./log"))
	time.Sleep(time.Minute)
	// o.NotifyCallback(1, 123*1024)
}

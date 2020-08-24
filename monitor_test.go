package logger

import (
	"testing"
	"time"
)

// go test -run TestNewMonitor -v -count=1
func TestNewMonitor(t *testing.T) {
	o := NewMonitor(&MonitorOption{
		ID:         123,
		LogPath:    "./log",
		MaxSize:    1024,
		NotifyRate: time.Minute,

		MailAddr: "smtp.exmail.qq.com:465",
		MailUser: "xxx@qq.com",
		MailPass: "123",
		MailName: "XXX",
		ToAddr:   "to@qq.com",
	})
	t.Log(o.GetSize("./log"))
	time.Sleep(time.Minute)
	// o.NotifyCallback(1, 123*1024)
}

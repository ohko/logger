package logger

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/mail"
	"net/smtp"
	"strings"
	"time"

	"github.com/ohko/logger/email"
)

// MonitorOption ...
type MonitorOption struct {
	ID             int           // 标识符
	LogPath        string        // 日志目录
	MaxSize        int64         // 日志目录最大磁盘占用字节数
	NotifyRate     time.Duration // 通知频率
	NotifyCallback func() error  // 达到最大占用数量时，回调通知函数

	// Email setting
	MailAddr string // 邮件服务器SSL地址
	MailUser string // 发件人账号
	MailPass string // 发件人密码
	MailName string // 发件人名字
	ToAddr   string // 收件人地址
}

// Monitor ...
type Monitor struct {
	option *MonitorOption
}

// NewMonitor ...
func NewMonitor(option *MonitorOption) *Monitor {
	o := &Monitor{option: option}
	go o.monitor()
	return o
}

func (o *Monitor) monitor() {
	rate := o.option.NotifyRate
	if rate < time.Minute {
		rate = time.Minute
	}
	for {
		time.Sleep(rate)
		size := o.GetSize(o.option.LogPath)
		if size > o.option.MaxSize {
			if o.option.NotifyCallback != nil {
				o.option.NotifyCallback()
			} else {
				o.NotifyCallback(o.option.ID, size)
			}
		}
	}
}

// GetSize ...
func (o *Monitor) GetSize(dirPath string) int64 {
	dirSize := int64(0)
	flist, e := ioutil.ReadDir(dirPath)
	if e != nil {
		return 0
	}
	for _, f := range flist {
		if f.IsDir() {
			dirSize = o.GetSize(dirPath+"/"+f.Name()) + dirSize
		} else {
			dirSize = f.Size() + dirSize
		}
	}
	return dirSize
}

// NotifyCallback ...
func (o *Monitor) NotifyCallback(id int, size int64) error {
	html := true
	mailAddr := o.option.MailAddr
	mailUser := o.option.MailUser
	mailPass := o.option.MailPass
	mailName := o.option.MailName
	toAddr := o.option.ToAddr
	subject := "logger monitor notify"
	body := fmt.Sprintf("ID: %d, Size: %.3fMB", id, float64(size)/1024/1024)
	attachFile := ""

	// 发送通知邮件
	hp := strings.Split(mailAddr, ":")
	to := mail.Address{Name: "", Address: toAddr}
	from := mail.Address{Name: mailName, Address: mailUser}
	auth := smtp.PlainAuth("", mailUser, mailPass, hp[0])

	var m *email.Message
	if html {
		m = email.NewHTMLMessage(subject, body)
	} else {
		m = email.NewMessage(subject, body)
	}
	m.From = from
	m.To = []string{toAddr}
	if attachFile != "" {
		if err := m.Attach(attachFile); err != nil {
			return err
		}
	}

	// get SSL connection
	conn, err := tls.Dial("tcp", mailAddr, &tls.Config{InsecureSkipVerify: true})
	if err != nil {
		return err
	}
	// create new SMTP client
	smtpClient, err := smtp.NewClient(conn, hp[0])
	if err != nil {
		return err
	}
	defer smtpClient.Quit()
	// auth the smtp client
	err = smtpClient.Auth(auth)
	if err != nil {
		return err
	}
	// set To && From address, note that from address must be same as authorization user.
	err = smtpClient.Mail(from.Address)
	if err != nil {
		return err
	}
	err = smtpClient.Rcpt(to.Address)
	if err != nil {
		return err
	}
	// Get the writer from SMTP client
	writer, err := smtpClient.Data()
	if err != nil {
		return err
	}
	// compose message body
	// write message to recp
	_, err = writer.Write(m.Bytes())
	if err != nil {
		return err
	}
	// close the writer
	err = writer.Close()
	if err != nil {
		return err
	}

	return nil
}

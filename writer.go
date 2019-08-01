package logger

import (
	"errors"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

// DefaultWriter ...
type DefaultWriter struct {
	cur   io.Writer
	clone io.Writer
}

// NewDefaultWriter ...
func NewDefaultWriter(clone io.Writer) *DefaultWriter {
	o := new(DefaultWriter)
	o.clone = clone
	o.next()

	go func() {
		for {
			// 等待明天
			t1 := time.Now()
			t2, _ := time.Parse("2006-01-02 -0700", t1.Add(time.Hour*24).Format("2006-01-02 -0700"))
			<-time.After(t2.Sub(t1))

			o.next()
		}
	}()

	return o
}

func (o *DefaultWriter) next() {
	f := "./log/" + time.Now().Format("2006/2006-01-02") + ".log"
	os.MkdirAll(filepath.Dir(f), 0755)
	nc, err := os.OpenFile(f, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Println(err)
		return
	}

	if o.clone != nil {
		o.cur = io.MultiWriter(nc, o.clone)
	} else {
		o.cur = nc
	}
}

func (o *DefaultWriter) Write(p []byte) (n int, err error) {
	if o.cur == nil {
		return 0, errors.New("io nil error")
	}
	return o.cur.Write(p)
}

package logger

import (
	"log"
	"strings"
	"testing"
	"time"
)

func TestNewLogger(t *testing.T) {
	l1 := NewLogger("/tmp/1[%v].log")
	l1.Printf("%v", "Printf")

	defer func() {
		if err := recover(); err != nil {
			l1.LogCalldepth(5, err)
		}
	}()

	l1.SetPrefix("[prefix] ")
	l1.Printf("Prefix")

	l1.SetPrefix("")
	l1.Println("No Prefix")

	l2 := NewLogger("/tmp/2[%v].log")
	l2.setFile("/tmp/2016-01-02.log")
	l2.Printf("New File 1")
	l2.setFile("/tmp/2016-01-03.log")
	l2.Printf("New File 2")

	l2.SetFlag(log.Ltime)
	l2.Printf("%v", "Only Time")

	l2.SetFlag(log.Ltime | log.Lshortfile)
	l2.Println("Time | ShortFile")

	panic("panic")
}

func TestPipe(t *testing.T) {
	l := NewLogger("/tmp/l-%v.log")
	l.SetPipe("/tmp/abc")
	for {
		s := strings.Repeat(time.Now().Format("2006-01-02 15:04:05"), 100)
		l.Printf(s)
		time.Sleep(time.Second)
	}
}

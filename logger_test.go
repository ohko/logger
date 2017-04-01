package logger

import (
	"log"
	"testing"
)

func TestNewLogger(t *testing.T) {
	l1 := NewLogger(0, "", "1-%v.log")
	l2 := NewLogger(0, "", "2-%v.log")
	l1.Log0Debug("0:%v", "Info")

	l1.SetPrefix("prefix")
	l1.Log1Warn("1:Warning")

	l1.SetPrefix("")
	l1.Log2Error("2:Error")
	l1.Log3Fatal("3:Fatal")

	l2.SetFile("2016-01-02.log")
	l2.Log4Trace("4:Trace")
	l2.SetFile("2016-01-03.log")
	l2.Log5NoFormat("5:NoFormat")

	l2.SetFlag(log.Ltime)
	l2.Printf("%v", "Printf")

	l2.SetFlag(log.Ltime | log.Lshortfile)
	l2.Println("Println")
}

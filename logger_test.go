package logger

import (
	"log"
	"testing"
)

func TestNewLogger(t *testing.T) {
	Log0Info("0:%v", "Info")

	SetPrefix("prefix")
	Log1Warn("1:Warning")

	SetPrefix("")
	Log2Error("2:Error")
	Log3Fatal("3:Fatal")

	SetFile("2016-01-02.log")
	Log4Trace("4:Trace")
	SetFile("2016-01-03.log")
	Log5NoFormat("5:NoFormat")

	SetFlag(log.Ltime)
	Printf("%v", "Printf")

	SetFlag(log.Ltime | log.Lshortfile)
	Println("Println")
}

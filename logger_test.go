package logger

import (
	"fmt"
	"log"
	"testing"
	"time"
)

func TestNewLogger(t *testing.T) {
	// l1 := NewLogger(NewDefaultWriter(os.Stdout))
	// l1 := NewLogger(NewDefaultWriter(nil))
	// l1 := NewLogger(os.Stdout)
	l1 := NewLogger(nil)
	// l1.SetLevel(LoggerLevel3Fatal)
	l1.SetPrefix("L1")

	for {
		l1.Log0Debug(fmt.Sprintf("0:%v", "Debug"))
		l1.Log1Warn("1:Warning")
		l1.Log2Error("2:Error")
		l1.SetPrefix("l1")
		// l1.Log3Fatal("3:Fatal")
		l1.Log4Trace("4:Trace")
		time.Sleep(time.Second)
	}

	l1.SetColor(false)
	l1.Log1Warn("no color")

	l1.SetPrefix("")
	l1.Log1Warn("no prefix")

	l1.SetFlags(log.Lshortfile)
	l1.Log4Trace("log.Lshortfile")

	l1.SetLevel(LoggerLevel5Off)
	l1.Log4Trace("LoggerLevelOff")
}

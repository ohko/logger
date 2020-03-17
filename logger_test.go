package logger

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"
	"time"
)

// go test github.com/ohko/logger -run TestNewLogger -v -count=1
func TestNewLogger(t *testing.T) {
	log.SetFlags(log.Llongfile)

	l1 := NewLogger(NewDefaultWriter(&DefaultWriterOption{Clone: os.Stdout, Path: "./log", Label: "lable", Name: "name_"}))
	// l1 := NewLogger(NewDefaultWriter(nil))
	// l1 := NewLogger(os.Stdout)
	// l1 := NewLogger(nil)
	// l1.SetLevel(LoggerLevel3Fatal)
	l1.SetPrefix("L1")

	l1.Log0Debug(fmt.Sprintf("0:%v", "Debug"))
	l1.Log1Warn("1:Warning")
	l1.SetPrefix("l1")
	l1.Log2Error("2:Error")
	// l1.Log3Fatal("3:Fatal")
	l1.Log4Trace("4:Trace")

	l1.SetColor(true)
	l1.Log2Error("has color")
	l1.SetColor(false)
	l1.Log1Warn("no color")

	l1.SetPrefix("")
	l1.Log1Warn("no prefix")

	l2 := l1.Fork("fork")
	l2.Log4Trace("new fork")

	l1.SetFlags(log.Lshortfile)
	l1.Log4Trace("log.Lshortfile")

	l1.SetLevel(LoggerLevel5Off)
	l1.Log4Trace("!!! DOT NOT SEE THIS !!! LoggerLevelOff")

	l2.Log4Trace("!!! DOT NOT SEE THIS !!!")
	// l1.Listen(":8080")
}

func TestCompressMonth(t *testing.T) {
	os.MkdirAll("./log/2019/01/", 0755)
	os.MkdirAll("./log/2019/02/", 0755)
	os.MkdirAll("./log/2019/03/", 0755)
	ioutil.WriteFile("./log/2019/01/2019-01-01.log", []byte("2019-01-01"), 0644)
	ioutil.WriteFile("./log/2019/02/2019-02-01.log", []byte("2019-02-01"), 0644)
	ioutil.WriteFile("./log/2019/02/2019-02-02.log", []byte("2019-02-02"), 0644)

	if err := compressAndRemoveDir("./log/2019/01/", "./log/2019/2019-01.zip"); err != nil {
		log.Println(err)
	}

	if err := compressAndRemoveDir("./log/2019/02/", "./log/2019/2019-02.zip"); err != nil {
		log.Println(err)
	}

	if err := os.RemoveAll("./log/2019/2019-01.zip"); err != nil {
		log.Println(err)
	}
}

func TestCompressDay(t *testing.T) {
	os.MkdirAll("./log/2019/01/", 0755)
	os.MkdirAll("./log/2019/02/", 0755)
	os.MkdirAll("./log/2019/03/", 0755)
	ioutil.WriteFile("./log/2019/02/2019-02-01.log", []byte("2019-02-01"), 0644)
	ioutil.WriteFile("./log/2019/02/2019-02-02.log", []byte("2019-02-02"), 0644)

	logFile := "./log/2019/02/2019-02-01.log"
	zipFile := "./log/2019/02/2019-02-01.zip"
	if err := compressAndRemoveFile(logFile, zipFile); err != nil {
		log.Println(err)
	}
}

func TestSubMonth(t *testing.T) {
	t0, _ := time.Parse("2006-01-02 -0700", "2019-06-15 +0800")
	t1, _ := time.Parse("2006-01-02 -0700", "2019-05-01 +0800")
	t2, _ := time.Parse("2006-01-02 -0700", "2019-04-01 +0800")
	t3, _ := time.Parse("2006-01-02 -0700", "2019-03-01 +0800")
	log.Println(t0, t1, t2, t3)
	if subMoth(t0, 1).Sub(t1) != 0 {
		t.Fail()
	}
	if subMoth(t0, 2).Sub(t2) != 0 {
		t.Fail()
	}
	if subMoth(t0, 3).Sub(t3) != 0 {
		t.Fail()
	}
}

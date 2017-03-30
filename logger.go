package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"
)

const (
	// LoggerLevelInfo 一般信息
	LoggerLevelInfo = iota
	// LoggerLevelWarning 警告信息
	LoggerLevelWarning
	// LoggerLevelError 错误信息
	LoggerLevelError
	// LoggerLevelFatal 严重信息
	LoggerLevelFatal
	// LoggerLevelTrace 打印信息
	LoggerLevelTrace
	// LoggerLevelNoFormat 无格式信息
	LoggerLevelNoFormat
	// LoggerLevelOff 关闭信息
	LoggerLevelOff
)

// Logger ...
type Logger struct {
	l          *log.Logger
	fileName   string
	fileReg    string
	fileHandle *os.File
	level      int
	prefix     string
	lock       sync.Mutex
}

// var ll = NewLogger(LoggerLevelInfo, "", "/tmp/%v.log")
var ll = NewLogger(LoggerLevelInfo, "", "")

// NewLogger ...
// eg: prefix="demo", file="./log/%v.log"
func NewLogger(level int, prefix, file string) *Logger {

	if prefix != "" {
		prefix = "[" + prefix + "] "
	}

	var l *log.Logger
	var fileName string
	var logFile *os.File
	var err error
	if file != "" {
		os.MkdirAll(filepath.Dir(file), 0755)
		fileName = fmt.Sprintf(file, time.Now().Format("2006-01-02"))
		logFile, err = os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
		if err != nil {
			panic(err)
		}
		l = log.New(io.MultiWriter(logFile, os.Stdout), prefix, log.Ltime|log.Lshortfile)
	} else {
		l = log.New(os.Stdout, prefix, log.Ltime|log.Lshortfile)
	}

	return &Logger{l: l, fileName: fileName, fileReg: file, fileHandle: logFile, level: level, prefix: prefix}
}

func (o *Logger) nextLogFile() {
	_f := fmt.Sprintf(o.fileReg, time.Now().Format("2006-01-02"))
	// 跨日
	if o.fileName != "" && _f != o.fileName {
		o.lock.Lock()
		defer o.lock.Unlock()

		oldFileName := o.fileName

		o.fileName = _f
		logFile, err := os.OpenFile(_f, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
		if err == nil {
			// 先关闭前一个日志
			o.fileHandle.Close()

			// TODO: 压缩前一月日志
			// 压缩前一天日志
			exec.Command("tar", "czf", oldFileName+".tar.gz", oldFileName).Run()
			os.Remove(oldFileName)

			// 赋值新日志
			o.fileHandle = logFile
			o.l.SetOutput(io.MultiWriter(logFile, os.Stdout))
		}
	}
}

// LogCalldepth ...
func (o *Logger) LogCalldepth(calldepth int, level int, msg ...interface{}) {
	if level < o.level {
		return
	}
	o.nextLogFile()
	o.lock.Lock()
	defer o.lock.Unlock()
	switch level {
	case LoggerLevelInfo:
		o.l.SetPrefix("\033[32m" + o.prefix)
	case LoggerLevelWarning:
		o.l.SetPrefix("\033[33m" + o.prefix)
	case LoggerLevelError:
		o.l.SetPrefix("\033[31m" + o.prefix)
	case LoggerLevelFatal:
		o.l.SetPrefix("\033[31;1;5;7m" + o.prefix)
	case LoggerLevelTrace:
		o.l.SetPrefix("\033[37m" + o.prefix)
	case LoggerLevelNoFormat:
		o.l.SetPrefix(o.prefix)
		o.l.Output(calldepth, fmt.Sprint(msg...))
		return
	default:
		o.l.SetPrefix(o.prefix)
	}

	o.l.Output(calldepth, fmt.Sprint(msg...)+"\033[m")
}

// SetFlag ...
func SetFlag(flag int) {
	ll.l.SetFlags(flag)
}

// SetLevel ...
func SetLevel(level int) {
	ll.level = level
}

// SetPrefix ...
func SetPrefix(prefix string) {
	if prefix != "" {
		ll.prefix = "[" + prefix + "] "
	} else {
		ll.prefix = prefix
	}
}

// Println ...
func Println(v ...interface{}) {
	ll.LogCalldepth(3, LoggerLevelNoFormat, fmt.Sprintln(v...))
}

// Printf ...
func Printf(format string, v ...interface{}) {
	ll.LogCalldepth(3, LoggerLevelNoFormat, fmt.Sprintf(format, v...))
}

// Log0Info ...
func Log0Info(format string, v ...interface{}) {
	ll.LogCalldepth(3, LoggerLevelInfo, fmt.Sprintf(format, v...))
}

// Log1Warn ...
func Log1Warn(format string, v ...interface{}) {
	ll.LogCalldepth(3, LoggerLevelWarning, fmt.Sprintf(format, v...))
}

// Log2Error ...
func Log2Error(format string, v ...interface{}) {
	ll.LogCalldepth(3, LoggerLevelError, fmt.Sprintf(format, v...))
}

// Log3Fatal ...
func Log3Fatal(format string, v ...interface{}) {
	ll.LogCalldepth(3, LoggerLevelFatal, fmt.Sprintf(format, v...))
}

// Log4Trace ...
func Log4Trace(format string, v ...interface{}) {
	ll.LogCalldepth(3, LoggerLevelTrace, fmt.Sprintf(format, v...))
}

// Log5NoFormat ...
func Log5NoFormat(format string, v ...interface{}) {
	ll.LogCalldepth(3, LoggerLevelNoFormat, fmt.Sprintf(format, v...))
}

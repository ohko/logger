package logger

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"
)

// Logger ...
type Logger struct {
	l          *log.Logger // log对象
	fileName   string      // 日志文件名
	fileReg    string      // 日志文件名格式
	fileHandle *os.File    // 日志文件handle
	lock       sync.Mutex  // 日志锁

	// fifo
	pipe     *os.File
	logCache chan string

	// 是否在记录日志
	logged bool
}

// NewLogger ...
// eg: ll := NewLogger("")
// eg: ll := NewLogger("./log/%v.log")
func NewLogger(file string) *Logger {
	var err error
	var l *log.Logger
	var fileName string
	var logFile *os.File

	if file != "" {
		os.MkdirAll(filepath.Dir(file), 0755)
		fileName = fmt.Sprintf(file, time.Now().Format("2006-01-02"))
		logFile, err = os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
		if err != nil {
			panic(err)
		}
		l = log.New(io.MultiWriter(logFile, os.Stdout), "", log.Ltime|log.Lshortfile)
	} else {
		l = log.New(os.Stdout, "", log.Ltime|log.Lshortfile)
	}

	return &Logger{l: l, fileName: fileName, fileReg: file, fileHandle: logFile, logged: true}
}

func (o *Logger) nextLogFile() {
	o.lock.Lock()
	defer o.lock.Unlock()

	var _nextFileName string
	if strings.Contains(o.fileReg, "%v") {
		_nextFileName = fmt.Sprintf(o.fileReg, time.Now().Format("2006-01-02"))
	} else {
		_nextFileName = o.fileReg
	}

	if o.fileReg == "" || _nextFileName == o.fileName {
		return
	}

	logFile, err := os.OpenFile(_nextFileName, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		return
	}

	if o.fileHandle != nil {
		// 先关闭前一个日志
		o.fileHandle.Close()

		// TODO: 压缩前一月日志
		// 压缩前一天日志
		exec.Command("tar", "czf", o.fileName+".tar.gz", o.fileName).Run()
		os.Remove(o.fileName)
	}

	// 赋值新日志
	o.fileName = _nextFileName
	o.fileHandle = logFile
	o.l.SetOutput(io.MultiWriter(logFile, os.Stdout))
}

// LogCalldepth ...
func (o *Logger) LogCalldepth(calldepth int, msg ...interface{}) {
	if !o.logged {
		return
	}
	if o.pipe != nil {
		select {
		case o.logCache <- fmt.Sprint(msg...):
		default:
		}
	}
	o.nextLogFile()
	o.l.Output(calldepth, fmt.Sprint(msg...))
}

func (o *Logger) setFile(fileReg string) {
	os.MkdirAll(filepath.Dir(fileReg), 0755)
	o.fileReg = fileReg
}

// SetFlag ...
func (o *Logger) SetFlag(flag int) {
	o.l.SetFlags(flag)
}

// SetPrefix ...
func (o *Logger) SetPrefix(prefix string) {
	o.l.SetPrefix(prefix)
}

// Println ...
func (o *Logger) Println(v ...interface{}) {
	if !o.logged {
		return
	}
	o.LogCalldepth(3, fmt.Sprintln(v...))
}

// Printf ...
func (o *Logger) Printf(format string, v ...interface{}) {
	if !o.logged {
		return
	}
	o.LogCalldepth(3, fmt.Sprintf(format, v...))
}

// SetPipe ...
func (o *Logger) SetPipe(filePath string) error {
	if o.pipe != nil {
		return errors.New("pipe exists")
	}
	if err := syscall.Unlink(filePath); err != nil {
		return err
	}
	if err := syscall.Mkfifo(filePath, 0666); err != nil {
		return err
	}
	// syscall.Mknod(filePath, syscall.S_IFIFO|0666, 0)

	logFile, err := os.OpenFile(filePath, os.O_RDWR|os.O_SYNC, os.ModeNamedPipe)
	if err != nil {
		return err
	}

	o.pipe = logFile
	o.logCache = make(chan string, 10000)
	go func() {
		for {
			if _, err := o.pipe.WriteString(<-o.logCache); err != nil {
				break
			}
		}
	}()
	return nil
}

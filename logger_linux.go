package logger

import (
	"errors"
	"os"
	"syscall"
)

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

package log

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"runtime"
	"time"
)

type LogLevel = int

//日志等级
const (
	DEBUG_N LogLevel = iota
	INFO_N
	WARN_N
	ERROR_N
	FATAL_N
	TRASH
)

const (
	sDEBUG = "DEBUG"
	sINFO  = "INFO"
	sWARN  = "WARN"
	sERROR = "ERROR"
	sFATAL = "FATAL"

	defaultBufSize = 1024 * 1024
)

var (
	logLevel = DEBUG_N

	levelText = map[int]string{DEBUG_N: sDEBUG, INFO_N: sINFO, WARN_N: sWARN, ERROR_N: sERROR, FATAL_N: sFATAL}

	fileOutput *os.File
	fileBuffer *bufio.Writer
	chOutput   = make(chan string, 1000)
)

func SetFile(fileName string) bool {
	os.MkdirAll(path.Dir(fileName), os.ModeDir)
	fileOutput, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println(err)
		return false
	}
	fileBuffer = bufio.NewWriterSize(fileOutput, defaultBufSize)
	return true
}

func SetLevel(level LogLevel) bool {
	if level < DEBUG_N || level >= TRASH {
		return false
	}

	logLevel = level
	return true
}

func Close() {
	for {
		if len(chOutput) == 0 {
			time.Sleep(time.Millisecond * 5)
			break
		}
		time.Sleep(time.Millisecond * 10)
	}
	if fileBuffer != nil {
		fileBuffer.Flush()
		fileOutput.Sync()
	}
}

func init() {
	go func() {
		var output string
		write := os.Stdout.WriteString
		for {
			select {
			case output = <-chOutput:
				write(output)
				if fileBuffer != nil {
					fileBuffer.WriteString(output)
				}
			case <-time.After(time.Millisecond * 100):
				if fileBuffer != nil {
					fileBuffer.Flush()
					fileOutput.Sync()
				}
			}
		}
	}()
}

func timeFmt() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func fileLine(skip int) string {
	_, file, line, _ := runtime.Caller(skip)
	i, count := len(file)-4, 0
	for ; i > 0; i-- {
		if file[i] == '/' {
			count++
			if count == 2 {
				break
			}
		}
	}
	return fmt.Sprintf("%s:%d", file[i+1:], line)
}

func writeBufferf(level LogLevel, skip int, format string, data ...interface{}) {
	if level < logLevel {
		return
	}
	chOutput <- fmt.Sprintf("%s %s [%s] - %s\n", timeFmt(), levelText[level], fileLine(skip), fmt.Sprintf(format, data...))
}

func writeBuffer(level LogLevel, skip int, data ...interface{}) {
	if level < logLevel {
		return
	}
	chOutput <- fmt.Sprintf("%s %s [%s] - %s\n", timeFmt(), levelText[level], fileLine(skip), fmt.Sprint(data...))
}

func Debug(data ...interface{}) {
	writeBuffer(DEBUG_N, 3, data...)
}

func Debugf(format string, data ...interface{}) {
	writeBufferf(DEBUG_N, 3, format, data...)
}

func Info(data ...interface{}) {
	writeBuffer(INFO_N, 3, data...)
}

func Infof(format string, data ...interface{}) {
	writeBufferf(INFO_N, 3, format, data...)
}

func Warn(data ...interface{}) {
	writeBuffer(WARN_N, 3, data...)
}

func Warnf(format string, data ...interface{}) {
	writeBufferf(WARN_N, 3, format, data...)
}

func Error(data ...interface{}) {
	writeBuffer(ERROR_N, 3, data...)
}

func Errorf(format string, data ...interface{}) {
	writeBufferf(ERROR_N, 3, format, data...)
}

func Fatal(data ...interface{}) {
	writeBuffer(FATAL_N, 3, data...)
}

func Fatalf(format string, data ...interface{}) {
	writeBufferf(FATAL_N, 3, format, data...)
}

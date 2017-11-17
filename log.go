package log

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"runtime"
	"time"
)

//日志等级
const (
	TRASH = iota
	FATAL_N
	ERROR_N
	WARN_N
	INFO_N
	DEBUG_N
)

const (
	sDEBUG = "DEBUG"
	sINFO  = "INFO"
	sWARN  = "WARN"
	sERROR = "ERROR"
	sFATAL = "FATAL"
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
	fileBuffer = bufio.NewWriter(fileOutput)
	return true
}

func SetLevel(level int) bool {
	if level <= TRASH || level > DEBUG_N {
		return false
	}

	logLevel = level
	return true
}

func init() {
	go func(chOutput chan string) {
		var output string
		write := os.Stdout.WriteString
		for {
			select {
			case output = <-chOutput:
				write(output)
				if fileBuffer != nil {
					fileBuffer.WriteString(output)
				}
			case <-time.After(time.Millisecond * 200):
				if fileBuffer != nil {
					fileBuffer.Flush()
					fileOutput.Sync()
				}
			}
		}
	}(chOutput)
}

func timeFmt() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func fileLine() string {
	_, file, line, _ := runtime.Caller(3)
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

func writeBufferf(level int, format string, data ...interface{}) {
	if logLevel < level {
		return
	}
	chOutput <- fmt.Sprintf("%s %s [%s] - %s\n", timeFmt(), levelText[level], fileLine(), fmt.Sprintf(format, data...))
}

func writeBuffer(level int, data ...interface{}) {
	if logLevel < level {
		return
	}
	chOutput <- fmt.Sprintf("%s %s [%s] - %s\n", timeFmt(), levelText[level], fileLine(), fmt.Sprint(data...))
}

func Debug(data ...interface{}) {
	writeBuffer(DEBUG_N, data...)
}

func Debugf(format string, data ...interface{}) {
	writeBufferf(DEBUG_N, format, data...)
}

func Info(data ...interface{}) {
	writeBuffer(INFO_N, data...)
}

func Infof(format string, data ...interface{}) {
	writeBufferf(INFO_N, format, data...)
}

func Warn(data ...interface{}) {
	writeBuffer(WARN_N, data...)
}

func Warnf(format string, data ...interface{}) {
	writeBufferf(WARN_N, format, data...)
}

func Error(data ...interface{}) {
	writeBuffer(ERROR_N, data...)
}

func Errorf(format string, data ...interface{}) {
	writeBufferf(ERROR_N, format, data...)
}

func Fatal(data ...interface{}) {
	writeBuffer(FATAL_N, data...)
}

func Fatalf(format string, data ...interface{}) {
	writeBufferf(FATAL_N, format, data...)
}

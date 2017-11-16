package log

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"
)

const (
	TRASH = iota
	FATAL_N
	ERROR_N
	WARN_N
	INFO_N
	DEBUG_N
)

const (
	DEBUG = "DEBUG"
	INFO  = "INFO"
	WARN  = "WARN"
	ERROR = "ERROR"
	FATAL = "FATAL"
)

var (
	LOG_LEVEL = DEBUG_N

	LevelText = map[int]string{DEBUG_N: DEBUG, INFO_N: INFO, WARN_N: WARN, ERROR_N: ERROR, FATAL_N: FATAL}

	fileOutput *os.File
	fileBuffer *bufio.Writer
	chOutput   = make(chan string, 10000)
)

func SetFile(file *os.File) {
	fileOutput = file
	fileBuffer = bufio.NewWriter(file)

}

func Start() {
	go func(chOutput chan string) {
		go func() {
			for {
				time.Sleep(time.Second * 2)
				if fileOutput != nil {
					fileBuffer.Flush()
					fileOutput.Sync()
				}
			}
		}()

		var output string
		for {
			output = <-chOutput
			os.Stdout.WriteString(output)
			if fileOutput != nil {
				fileBuffer.WriteString(output)
			}
		}
	}(chOutput)
}

func Close() {
	if fileOutput != nil {
		fileBuffer.Flush()
		fileOutput.Sync()
	}
}

func Time() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func FileLine() string {
	var f string

	_, file, line, ok := runtime.Caller(3)
	if ok == true {
		files := strings.Split(file, "/src/")
		if len(files) >= 2 {
			f = files[len(files)-1]
		} else {
			f = file
		}

		f = fmt.Sprintf("[%s(%v)]", f, line)
	}

	return f
}

func writeBuffer(level int, format string, data ...interface{}) {
	if LOG_LEVEL < level {
		return
	}
	levelText := LevelText[level]
	chOutput <- fmt.Sprintf("%s %s %s %s\n", Time(), levelText, FileLine(), fmt.Sprintf(format, data...))
}

func Debug(format string, data ...interface{}) {
	writeBuffer(DEBUG_N, format, data...)
}

func Info(format string, data ...interface{}) {
	writeBuffer(INFO_N, format, data...)
}

func Warn(format string, data ...interface{}) {
	writeBuffer(WARN_N, format, data...)
}

func Error(format string, data ...interface{}) {
	writeBuffer(ERROR_N, format, data...)
}

func Fatal(format string, data ...interface{}) {
	writeBuffer(FATAL_N, format, data...)
}

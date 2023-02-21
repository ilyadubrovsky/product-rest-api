package logging

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path"
	"runtime"
)

var e *logrus.Entry

type Logger struct {
	*logrus.Entry
}

// GetLogger returns the global logger entry
func GetLogger() *Logger {
	return &Logger{e}
}

// writeHook custom hook which implements logrus hook interface
type writerHook struct {
	Writer    []io.Writer
	LogLevels []logrus.Level
}

// Levels returns log levels
func (hook *writerHook) Levels() []logrus.Level {
	return hook.LogLevels
}

// Fire write received log info to custom writers: file and ostream
func (hook *writerHook) Fire(entry *logrus.Entry) error {
	line, err := entry.String()
	if err != nil {
		return err
	}

	for _, w := range hook.Writer {
		w.Write([]byte(line))
	}

	return err
}

// init initializes logger from logrus
func init() {
	l := logrus.New()
	l.SetReportCaller(true)
	l.Formatter = &logrus.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			filename := path.Base(frame.File)
			return fmt.Sprintf("%s()", frame.Func), filename
		},
	}

	if err := os.Mkdir("logs", 0644); err != nil {
		panic(err)
	}

	file, err := os.OpenFile("logs/log.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0640)
	if err != nil {
		panic(err)
	}

	l.SetOutput(io.Discard)

	l.AddHook(&writerHook{
		Writer:    []io.Writer{os.Stdout, file},
		LogLevels: logrus.AllLevels,
	})
}

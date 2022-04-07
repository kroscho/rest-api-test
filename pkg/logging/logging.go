package logging

import (
	"fmt"
	"io"
	"os"
	"path"
	"runtime"

	"github.com/sirupsen/logrus"
)

type writerHook struct {
	Writer    []io.Writer
	LogLevels []logrus.Level
}

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

func (hook *writerHook) Levels() []logrus.Level {
	return hook.LogLevels
}

var e *logrus.Entry

type Logger struct {
	*logrus.Entry
}

func GetLogger() *Logger {
	return &Logger{e}
}

func (l *Logger) GetLogerWithField(k string, v interface{}) *Logger {
	return &Logger{l.WithField(k, v)}
}

func init() {
	l := logrus.New()
	l.SetReportCaller(true)
	l.Formatter = &logrus.TextFormatter{
		// передается frame
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			// в frame есть информация о файле, в котором происходит логирование строчки
			filename := path.Base(frame.File)
			// также в нем есть инфа о line на которой находимся и о функции внутри которой находимся
			// возвращается функция, имя файла и номер строчки
			return fmt.Sprintf("%s()", frame.Function), fmt.Sprintf("%s %d", filename, frame.Line)
		},
		// отключаем цвета
		DisableColors: false,
		FullTimestamp: true,
	}

	// создаем папку logs
	err := os.MkdirAll("logs", 0644)
	if err != nil {
		panic(err)
	}

	allFile, err := os.OpenFile("logs/all.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0640)
	if err != nil {
		panic(err)
	}

	// чтоб по умолчанию логи никуда не уходили
	l.SetOutput(io.Discard)

	l.AddHook(&writerHook{
		Writer:    []io.Writer{allFile, os.Stdout},
		LogLevels: logrus.AllLevels,
	})

	l.SetLevel(logrus.TraceLevel)

	e = logrus.NewEntry(l)
}

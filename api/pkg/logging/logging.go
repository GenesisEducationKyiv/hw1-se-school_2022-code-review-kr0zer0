package logging

import (
	"github.com/sirupsen/logrus"
	"io"
)

type writerHook struct {
	Writers   []io.Writer
	LogLevels []logrus.Level
}

func (hook *writerHook) Fire(entry *logrus.Entry) error {
	line, err := entry.String()
	if err != nil {
		return err
	}
	for _, w := range hook.Writers {
		w.Write([]byte(line))
	}

	return err
}

func (hook *writerHook) Levels() []logrus.Level {
	return hook.LogLevels
}

func Init(writers []io.Writer) *logrus.Logger {
	l := logrus.New()
	l.SetReportCaller(true)
	l.Formatter = &logrus.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	}

	l.SetOutput(io.Discard)

	l.AddHook(&writerHook{
		Writers:   writers,
		LogLevels: logrus.AllLevels,
	})

	return l
}

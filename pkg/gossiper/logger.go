package gossiper

import (
	"github.com/sirupsen/logrus"
	"os"
)

type logger struct {
	log *logrus.Logger
}

func Newlogger() *logger {
	return &logger{log: logrus.New()}
}

//hint and content can be ignored
func (ml *logger) Info(msg, hint string, content ...string) {
	if content != nil {
		ml.log.WithFields(logrus.Fields{hint: content}).Info(msg)
	} else {
		ml.log.Info(msg)
	}
}

func (ml *logger) Debug(msg, hint string, content ...string) {
	if content != nil {
		ml.log.WithFields(logrus.Fields{hint: content}).Debug(msg)
	} else {
		ml.log.Debug(msg)
	}
}

func (ml *logger) Warn(msg, hint string, content ...string) {
	if content != nil {
		ml.log.WithFields(logrus.Fields{hint: content}).Warn(msg)
	} else {
		ml.log.Warn(msg)
	}
}

func (ml *logger) Error(msg, hint string, content ...string) {
	if content != nil {
		ml.log.WithFields(logrus.Fields{hint: content}).Error(msg)
	} else {
		ml.log.Error(msg)
	}
}

func (ml *logger) Fatal(msg, hint string, content ...string) {
	if content != nil {
		ml.log.WithFields(logrus.Fields{hint: content}).Fatal(msg)
	} else {
		ml.log.Fatal(msg)
	}
}

func (ml *logger) Panic(msg, hint string, content ...string) {
	if content != nil {
		ml.log.WithFields(logrus.Fields{hint: content}).Panic(msg)
	} else {
		ml.log.Panic(msg)
	}
}

func (ml *logger) SaveToFile(path string) {
	ml.log.Out = os.Stdout
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		ml.log.Out = file
	} else {
		logrus.Info("Failed to log to file, using default stderr")
	}
}

//transfer to default from save to file
func (ml *logger) DefaultSet() {
	ml.log.Out = os.Stderr
}

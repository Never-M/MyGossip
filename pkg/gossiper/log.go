package gossiper

import (
	"github.com/sirupsen/logrus"
	"os"
)

type mylog struct {
	log *logrus.Logger
}

func Newlog() *mylog {
	return &mylog{log: logrus.New()}
}

//hint and content can be ignored
func (ml *mylog) Info(msg, hint string, content ...string) {
	if content != nil {
		ml.log.WithFields(logrus.Fields{hint: content}).Info(msg)
	} else {
		ml.log.Info(msg)
	}
}

func (ml *mylog) Debug(msg, hint string, content ...string) {
	if content != nil {
		ml.log.WithFields(logrus.Fields{hint: content}).Debug(msg)
	} else {
		ml.log.Debug(msg)
	}
}

func (ml *mylog) Warn(msg, hint string, content ...string) {
	if content != nil {
		ml.log.WithFields(logrus.Fields{hint: content}).Warn(msg)
	} else {
		ml.log.Warn(msg)
	}
}

func (ml *mylog) Error(msg, hint string, content ...string) {
	if content != nil {
		ml.log.WithFields(logrus.Fields{hint: content}).Error(msg)
	} else {
		ml.log.Error(msg)
	}
}

func (ml *mylog) Fatal(msg, hint string, content ...string) {
	if content != nil {
		ml.log.WithFields(logrus.Fields{hint: content}).Fatal(msg)
	} else {
		ml.log.Fatal(msg)
	}
}

func (ml *mylog) Panic(msg, hint string, content ...string) {
	if content != nil {
		ml.log.WithFields(logrus.Fields{hint: content}).Panic(msg)
	} else {
		ml.log.Panic(msg)
	}
}

func (ml *mylog) SaveToFile(path string) {
	ml.log.Out = os.Stdout
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		ml.log.Out = file
	} else {
		logrus.Info("Failed to log to file, using default stderr")
	}
}

//transfer to default from save to file
func (ml *mylog) Defaultset() {
	ml.log.Out = os.Stderr
}

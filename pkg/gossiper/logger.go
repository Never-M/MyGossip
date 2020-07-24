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
func (ml *logger) Info(msg ...string) {
	if len(msg) != 1 {
		ml.log.WithFields(logrus.Fields{msg[1]: msg[2]}).Info(msg[0])
	} else {
		ml.log.Fatal(msg[0])
	}
}

func (ml *logger) Debug(msg ...string) {
	if len(msg) != 1 {
		ml.log.WithFields(logrus.Fields{msg[1]: msg[2]}).Debug(msg[0])
	} else {
		ml.log.Fatal(msg[0])
	}
}

func (ml *logger) Warn(msg ...string) {
	if len(msg) != 1 {
		ml.log.WithFields(logrus.Fields{msg[1]: msg[2]}).Warn(msg[0])
	} else {
		ml.log.Fatal(msg[0])
	}
}

func (ml *logger) Error(msg ...string) {
	if len(msg) != 1 {
		ml.log.WithFields(logrus.Fields{msg[1]: msg[2]}).Error(msg[0])
	} else {
		ml.log.Fatal(msg[0])
	}
}

func (ml *logger) Fatal(msg ...string) {
	if len(msg) != 1 {
		ml.log.WithFields(logrus.Fields{msg[1]: msg[2]}).Fatal(msg[0])
	} else {
		ml.log.Fatal(msg[0])
	}
}

func (ml *logger) Panic(msg ...string) {
	if len(msg) != 1 {
		ml.log.WithFields(logrus.Fields{msg[1]: msg[2]}).Panic(msg[0])
	} else {
		ml.log.Fatal(msg[0])
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

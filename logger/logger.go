package logger

import (
	"MinecraftServerLogger/config"
	"github.com/hpcloud/tail"
	"log"
	"sync"
)

const (
	historySize = 1000
)

type LoggerReceiver interface {
	ReceiveHistory([]string)
	ReceiveMessage(string)
}

type Logger struct {
	conf    *config.RootConfig
	history []string

	loggerReceivers []LoggerReceiver
	lrmu            *sync.Mutex
}

func NewLogger(conf *config.RootConfig) *Logger {
	return &Logger{
		conf:            conf,
		history:         make([]string, 0, historySize),
		loggerReceivers: []LoggerReceiver{},
		lrmu:            &sync.Mutex{},
	}
}

func (l *Logger) Start() error {
	t, err := tail.TailFile(l.conf.LogLocation, tail.Config{
		Follow: true,
		ReOpen: true,
	})

	if err != nil {
		return err
	}

	go func(t *tail.Tail) {
		for line := range t.Lines {
			if line.Err != nil {
				log.Fatalln(line.Err)
			}

			l.history = append(l.history, line.Text)

			if len(l.history) > historySize {
				l.history = l.history[1:]
			}

			l.lrmu.Lock()
			for _, receiver := range l.loggerReceivers {
				receiver.ReceiveMessage(line.Text)
			}
			l.lrmu.Unlock()
		}
	}(t)

	return nil
}

func (l *Logger) RegisterNewReceiver(receiver LoggerReceiver) {
	// At first, send then log history
	historyLocal := make([]string, 0, len(l.history))
	copy(historyLocal, l.history)

	receiver.ReceiveHistory(historyLocal)

	l.lrmu.Lock()
	defer l.lrmu.Unlock()
	l.loggerReceivers = append(l.loggerReceivers, receiver)
}

func (l *Logger) DeleteReceiver(receiver LoggerReceiver) {
	l.lrmu.Lock()
	defer l.lrmu.Unlock()

	for i, listReceiver := range l.loggerReceivers {
		if listReceiver == receiver {
			l.loggerReceivers = append(l.loggerReceivers[:i], l.loggerReceivers[i+1:]...)
			return
		}
	}
}

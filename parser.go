package ansiterm

import (
	"os"

	"github.com/Sirupsen/logrus"
)

var parser *AnsiParser
var logFile *os.File
var log *logrus.Logger

func init() {
	filename := "parse.txt"

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		logFile, err = os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0x0666)
		if err != nil {
			panic(err)
		}
	}

	log = &logrus.Logger{
		Out:       logFile,
		Formatter: new(logrus.TextFormatter),
		Level:     logrus.InfoLevel,
	}
}

type AnsiParser struct {
	state        State
	eventHandler AnsiEventHandler
	context      *AnsiContext
}

func CreateParser(initialState State, evtHandler AnsiEventHandler) *AnsiParser {
	log.Infof("CreateParser")

	parser = &AnsiParser{
		state:        initialState,
		eventHandler: evtHandler,
		context:      &AnsiContext{},
	}

	return parser
}

func (ap *AnsiParser) Parse(bytes []byte) {
	for _, b := range bytes {
		ap.handle(b)
	}
}

func (ap *AnsiParser) handle(b byte) {
	newState, _ := ap.state.Handle(b)

	if newState == nil {
		log.Warning("newState is nil")
		return
	}

	if newState != ap.state {
		ap.changeState(newState)
	}
}

func (ap *AnsiParser) changeState(newState State) error {
	log.Infof("ChangeState %s --> %s", ap.state.Name(), newState.Name())

	// Exit old state
	err := ap.state.Exit()
	if err != nil {
		return err
	}

	// Perform transition action
	err = ap.state.Transition(newState)

	// Enter new state
	err = newState.Enter()
	if err != nil {
		return err
	}

	ap.state = newState
	return nil
}

package ansiterm

import (
	"errors"
	"fmt"
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

func (ap *AnsiParser) Parse(bytes []byte) (int, error) {
	for i, b := range bytes {
		if err := ap.handle(b); err != nil {
			return i, err
		}
	}

	return len(bytes), nil
}

func (ap *AnsiParser) handle(b byte) error {
	newState, err := ap.state.Handle(b)
	if err != nil {
		return err
	}

	if newState == nil {
		log.Warning("newState is nil")
		return errors.New(fmt.Sprintf("New state of 'nil' is invalid."))
	}

	if newState != ap.state {
		if err := ap.changeState(newState); err != nil {
			return err
		}
	}

	return nil
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

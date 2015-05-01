package ansiterm

import (
	log "github.com/Sirupsen/logrus"
)

type CsiParamState struct {
	BaseState
}

func (csiState CsiParamState) Handle(b byte) (s State, e error) {
	log.Infof("CsiParam::Handle %#x", b)
	parser.context.currentChar = b

	nextState, err := csiState.BaseState.Handle(b)
	if nextState != nil || err != nil {
		return nextState, err
	}

	switch {
	case sliceContains(Alphabetics, b):
		return Ground, nil
	case sliceContains(CsiCollectables, b):
		parser.collectParam(parser.context.currentChar)
		return CsiParam, nil
	}

	return csiState, nil
}

func (csiState CsiParamState) Transition(s State) error {
	log.Infof("CsiParam::Transition %s --> %s", csiState.Name(), s.Name())

	switch s {
	case Ground:
		parser.context.finalChar = parser.context.currentChar
		return parser.csiDispatch()
	}

	return nil
}

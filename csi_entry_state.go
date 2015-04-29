package ansiterm

import (
	log "github.com/Sirupsen/logrus"
)

type CsiEntryState struct {
	BaseState
}

func (csiState CsiEntryState) Handle(b byte) (s State, e error) {
	log.Infof("CsiEntry::Handle %#x", b)
	parser.context.currentChar = b

	nextState, err := csiState.BaseState.Handle(b)
	if nextState != nil || err != nil {
		return nextState, err
	}

	switch {
	case sliceContains(AllCase, b):
		return Ground, nil
	case sliceContains(CsiCollectables, b):
		return CsiParam, nil
	}

	return csiState, nil
}

func (csiState CsiEntryState) Transition(s State) error {
	log.Infof("CsiEntry::Transition %s --> %s", csiState.Name(), s.Name())

	switch s {
	case Ground:
		parser.context.finalChar = parser.context.currentChar
		return parser.csiDispatch()
	case CsiParam:
		switch {
		case sliceContains(CsiParams, parser.context.currentChar):
			parser.collectParam(parser.context.currentChar)
		}
	}

	return nil
}

package ansiterm

type EscapeState struct {
	BaseState
}

func (escState EscapeState) Handle(b byte) (s State, e error) {
	nextState, err := escState.BaseState.Handle(b)
	if nextState != nil || err != nil {
		return nextState, err
	}

	switch b {
	case ANSI_ESCAPE_SECONDARY:
		return CsiEntry, nil
	case ANSI_OSC_STRING_ENTRY:
		return OscString, nil
	}

	return escState, nil
}

func (escState EscapeState) Enter() error {
	parser.clear()
	return nil
}

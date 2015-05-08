package ansiterm

type OscStringState struct {
	BaseState
}

func (csiState OscStringState) Handle(b byte) (s State, e error) {
	log.Infof("OscString::Handle %#x", b)
	parser.context.currentChar = b

	nextState, err := csiState.BaseState.Handle(b)
	if nextState != nil || err != nil {
		return nextState, err
	}

	if isOscStringTerminator(b) {
		return Ground, nil
	}

	return csiState, nil
}

// See below for OSC string terminators for linux
// http://man7.org/linux/man-pages/man4/console_codes.4.html
func isOscStringTerminator(b byte) bool {

	if b == ANSI_BEL || b == 0x5C {
		return true
	}

	return false
}

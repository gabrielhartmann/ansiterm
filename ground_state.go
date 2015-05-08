package ansiterm

type GroundState struct {
	BaseState
}

func (gs GroundState) Handle(b byte) (s State, e error) {
	log.Infof("Ground::Handle %#x", b)
	parser.context.currentChar = b

	nextState, err := gs.BaseState.Handle(b)
	if nextState != nil || err != nil {
		return nextState, err
	}

	if sliceContains(Printables, b) {
		return gs, parser.print()
	}

	if sliceContains(C0Control, b) {
		return gs, parser.execute()
	}

	return gs, nil
}

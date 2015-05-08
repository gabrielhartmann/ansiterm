package ansiterm

type StateId int

const (
	csiEntry StateId = iota
	csiParam
	dcsEntry
	escape
	errorId
	ground
	oscString
	stateCount
	invalid = -1
)

var stateMap = map[StateId]State{
	csiEntry:  CsiEntry,
	csiParam:  CsiParam,
	dcsEntry:  DcsEntry,
	escape:    Escape,
	errorId:   Error,
	ground:    Ground,
	oscString: OscString,
}

type State interface {
	Enter() error
	Exit() error
	Handle(byte) (State, error)
	Name() string
	Transition(State) error
}

var CsiEntry = CsiEntryState{BaseState{name: "CsiEntry", id: csiEntry}}
var CsiParam = CsiParamState{BaseState{name: "CsiParam", id: csiParam}}
var DcsEntry = DcsEntryState{BaseState{name: "DcsEntry", id: dcsEntry}}
var Escape = EscapeState{BaseState{name: "Escape", id: escape}}
var Error = ErrorState{BaseState{name: "Error", id: errorId}}
var Ground = GroundState{BaseState{name: "Ground", id: ground}}
var OscString = OscStringState{BaseState{name: "OscString", id: oscString}}

type BaseState struct {
	id   StateId
	name string
}

func (base BaseState) Enter() error {
	return nil
}

func (base BaseState) Exit() error {
	return nil
}

func (base BaseState) Handle(b byte) (s State, e error) {

	switch b {
	case CSI_ENTRY:
		return CsiEntry, nil
	case DCS_ENTRY:
		return DcsEntry, nil
	case ANSI_ESCAPE_PRIMARY:
		return Escape, nil
	case OSC_STRING:
		return OscString, nil
	}

	return nil, nil
}

func (base BaseState) Name() string {
	return base.name
}

func (base BaseState) Transition(State) error {
	return nil
}

type DcsEntryState struct {
	BaseState
}

type ErrorState struct {
	BaseState
}

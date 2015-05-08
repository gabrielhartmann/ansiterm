package ansiterm

import (
	"fmt"
	"testing"
)

func TestStateTransitions(t *testing.T) {
	stateTransitionHelper(t, CsiEntry, Ground, Alphabetics)
	stateTransitionHelper(t, CsiEntry, CsiParam, CsiCollectables)
	stateTransitionHelper(t, Escape, CsiEntry, []byte{ANSI_ESCAPE_SECONDARY})
}

func TestAnyToX(t *testing.T) {
	anyToXHelper(t, []byte{ANSI_ESCAPE_PRIMARY}, Escape)
	anyToXHelper(t, []byte{DCS_ENTRY}, DcsEntry)
	anyToXHelper(t, []byte{OSC_STRING}, OscString)
	anyToXHelper(t, []byte{CSI_ENTRY}, CsiEntry)
}

func TestCollectCsiParams(t *testing.T) {
	parser, _ := createTestParser(CsiEntry)
	parser.Parse(CsiCollectables)

	buffer := parser.context.paramBuffer
	bufferCount := len(buffer)

	if bufferCount != len(CsiCollectables) {
		t.Errorf("Buffer:    %v", buffer)
		t.Errorf("CsiParams: %v", CsiCollectables)
		t.Errorf("Buffer count failure: %d != %d", bufferCount, len(CsiParams))
		return
	}

	for i, v := range CsiCollectables {
		if v != buffer[i] {
			t.Errorf("Buffer:    %v", buffer)
			t.Errorf("CsiParams: %v", CsiParams)
			t.Errorf("Mismatch at buffer[%d] = %d", i, buffer[i])
		}
	}
}

func TestParseParams(t *testing.T) {
	parseParamsHelper(t, []byte{}, []string{})
	parseParamsHelper(t, []byte{';'}, []string{})
	parseParamsHelper(t, []byte{';', ';'}, []string{})
	parseParamsHelper(t, []byte{'7'}, []string{"7"})
	parseParamsHelper(t, []byte{'7', ';'}, []string{"7"})
	parseParamsHelper(t, []byte{'7', ';', ';'}, []string{"7"})
	parseParamsHelper(t, []byte{'7', ';', ';', '8'}, []string{"7", "8"})
	parseParamsHelper(t, []byte{'7', ';', '8', ';'}, []string{"7", "8"})
	parseParamsHelper(t, []byte{'7', ';', ';', '8', ';', ';'}, []string{"7", "8"})
	parseParamsHelper(t, []byte{'7', '8'}, []string{"78"})
	parseParamsHelper(t, []byte{'7', '8', ';'}, []string{"78"})
	parseParamsHelper(t, []byte{'7', '8', ';', '9', '0'}, []string{"78", "90"})
	parseParamsHelper(t, []byte{'7', '8', ';', ';', '9', '0'}, []string{"78", "90"})
	parseParamsHelper(t, []byte{'7', '8', ';', '9', '0', ';'}, []string{"78", "90"})
	parseParamsHelper(t, []byte{'7', '8', ';', '9', '0', ';', ';'}, []string{"78", "90"})
}

func TestCursor(t *testing.T) {
	cursorSingleParamHelper(t, 'A', "CUU")
	cursorSingleParamHelper(t, 'B', "CUD")
	cursorSingleParamHelper(t, 'C', "CUF")
	cursorSingleParamHelper(t, 'D', "CUB")
	cursorSingleParamHelper(t, 'E', "CNL")
	cursorSingleParamHelper(t, 'F', "CPL")
	cursorSingleParamHelper(t, 'G', "CHA")
	cursorTwoParamHelper(t, 'H', "CUP")
	cursorTwoParamHelper(t, 'f', "HVP")
	funcCallParamHelper(t, []byte{'?', '2', '5', 'h'}, Ground, []string{"DECTCEM([true])"})
	funcCallParamHelper(t, []byte{'?', '2', '5', 'l'}, Ground, []string{"DECTCEM([false])"})
}

func TestErase(t *testing.T) {
	// Erase in Display
	eraseHelper(t, 'J', "ED")

	// Erase in Line
	eraseHelper(t, 'K', "EL")
}

func TestSelectGraphicRendition(t *testing.T) {
	funcCallParamHelper(t, []byte{'m'}, Ground, []string{"SGR([0])"})
	funcCallParamHelper(t, []byte{'0', 'm'}, Ground, []string{"SGR([0])"})
	funcCallParamHelper(t, []byte{'0', ';', '1', 'm'}, Ground, []string{"SGR([0 1])"})
	funcCallParamHelper(t, []byte{'0', ';', '1', ';', '2', 'm'}, Ground, []string{"SGR([0 1 2])"})
}

func TestPan(t *testing.T) {
	panHelper(t, 'S', "SU")
	panHelper(t, 'T', "SD")
}

func TestPrint(t *testing.T) {
	parser, evtHandler := createTestParser(Ground)
	parser.Parse(Printables)
	validateState(t, parser.state, Ground)

	for i, v := range Printables {
		expectedCall := fmt.Sprintf("Print([%s])", string(v))
		actualCall := evtHandler.FunctionCalls[i]
		if actualCall != expectedCall {
			t.Errorf("Actual != Expected: %v != %v at %d", actualCall, expectedCall, i)
		}
	}
}

func TestClear(t *testing.T) {
	p, _ := createTestParser(Ground)
	fillContext(p.context)
	p.clear()
	validateEmptyContext(t, p.context)
}

func TestClearOnStateChange(t *testing.T) {
	clearOnStateChangeHelper(t, Ground, Escape, []byte{ANSI_ESCAPE_PRIMARY})
	clearOnStateChangeHelper(t, Ground, CsiEntry, []byte{CSI_ENTRY})
}

func TestOscStringToGround(t *testing.T) {
	stateTransitionHelper(t, OscString, Ground, []byte{ANSI_BEL})
	stateTransitionHelper(t, OscString, Ground, []byte{0x5C})
}

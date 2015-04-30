package ansiterm

import (
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

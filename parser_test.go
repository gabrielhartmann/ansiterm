package ansiterm

import (
	"fmt"
	"testing"
)

func TestStateTransitions(t *testing.T) {
	stateTransitionHelper(t, CsiEntry, Ground, AllCase)
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
	parser.Parse(CsiParams)

	buffer := parser.context.parameterBuffer
	bufferCount := len(buffer)

	if bufferCount != len(CsiParams) {
		t.Errorf("%v", parser.context.parameterBuffer)
		t.Errorf("Buffer count failure: %d != %d", bufferCount, len(CsiParams))
		return
	}

	for i, v := range CsiParams {
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

func cursorHelper(t *testing.T, command byte, funcName string) {
	funcCallParamHelper(t, []byte{command}, Ground, []string{fmt.Sprintf("%s([])", funcName)})
	funcCallParamHelper(t, []byte{'2', command}, Ground, []string{fmt.Sprintf("%s([2])", funcName)})
	funcCallParamHelper(t, []byte{'2', '3', command}, Ground, []string{fmt.Sprintf("%s([23])", funcName)})
	funcCallParamHelper(t, []byte{'2', ';', '3', command}, Ground, []string{fmt.Sprintf("%s([2 3])", funcName)})
	funcCallParamHelper(t, []byte{'2', ';', '3', ';', '4', command}, Ground, []string{fmt.Sprintf("%s([2 3 4])", funcName)})
}

func TestCursor(t *testing.T) {
	cursorHelper(t, 'A', "CUU")
	cursorHelper(t, 'B', "CUD")
	cursorHelper(t, 'C', "CUF")
	cursorHelper(t, 'D', "CUB")
	cursorHelper(t, 'E', "CNL")
	cursorHelper(t, 'F', "CPL")
	cursorHelper(t, 'G', "CHA")
	cursorHelper(t, 'H', "CUP")
	cursorHelper(t, 'f', "HVP")
}

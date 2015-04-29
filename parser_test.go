package ansiterm

import (
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
	parseParamsHelper(t, []byte{}, []int{})
	parseParamsHelper(t, []byte{';'}, []int{})
	parseParamsHelper(t, []byte{';', ';'}, []int{})
	parseParamsHelper(t, []byte{'7'}, []int{7})
	parseParamsHelper(t, []byte{'7', ';'}, []int{7})
	parseParamsHelper(t, []byte{'7', ';', ';'}, []int{7})
	parseParamsHelper(t, []byte{'7', ';', ';', '8'}, []int{7, 8})
	parseParamsHelper(t, []byte{'7', ';', '8', ';'}, []int{7, 8})
	parseParamsHelper(t, []byte{'7', ';', ';', '8', ';', ';'}, []int{7, 8})
	parseParamsHelper(t, []byte{'7', '8'}, []int{78})
	parseParamsHelper(t, []byte{'7', '8', ';'}, []int{78})
	parseParamsHelper(t, []byte{'7', '8', ';', '9', '0'}, []int{78, 90})
	parseParamsHelper(t, []byte{'7', '8', ';', ';', '9', '0'}, []int{78, 90})
	parseParamsHelper(t, []byte{'7', '8', ';', '9', '0', ';'}, []int{78, 90})
	parseParamsHelper(t, []byte{'7', '8', ';', '9', '0', ';', ';'}, []int{78, 90})
}

func TestCursor(t *testing.T) {
	funcCallParamHelper(t, []byte{'A'}, Ground, []string{"CUU([])"})
	funcCallParamHelper(t, []byte{'2', 'A'}, Ground, []string{"CUU([2])"})
	funcCallParamHelper(t, []byte{'2', '3', 'A'}, Ground, []string{"CUU([23])"})
	funcCallParamHelper(t, []byte{'2', ';', '3', 'A'}, Ground, []string{"CUU([2 3])"})
	funcCallParamHelper(t, []byte{'2', ';', '3', ';', '4', 'A'}, Ground, []string{"CUU([2 3 4])"})

	funcCallParamHelper(t, []byte{'B'}, Ground, []string{"CUD([])"})
	funcCallParamHelper(t, []byte{'2', 'B'}, Ground, []string{"CUD([2])"})
	funcCallParamHelper(t, []byte{'2', '3', 'B'}, Ground, []string{"CUD([23])"})
	funcCallParamHelper(t, []byte{'2', ';', '3', 'B'}, Ground, []string{"CUD([2 3])"})
	funcCallParamHelper(t, []byte{'2', ';', '3', ';', '4', 'B'}, Ground, []string{"CUD([2 3 4])"})

	funcCallParamHelper(t, []byte{'C'}, Ground, []string{"CUF([])"})
	funcCallParamHelper(t, []byte{'2', 'C'}, Ground, []string{"CUF([2])"})
	funcCallParamHelper(t, []byte{'2', '3', 'C'}, Ground, []string{"CUF([23])"})
	funcCallParamHelper(t, []byte{'2', ';', '3', 'C'}, Ground, []string{"CUF([2 3])"})
	funcCallParamHelper(t, []byte{'2', ';', '3', ';', '4', 'C'}, Ground, []string{"CUF([2 3 4])"})

	funcCallParamHelper(t, []byte{'D'}, Ground, []string{"CUB([])"})
	funcCallParamHelper(t, []byte{'2', 'D'}, Ground, []string{"CUB([2])"})
	funcCallParamHelper(t, []byte{'2', '3', 'D'}, Ground, []string{"CUB([23])"})
	funcCallParamHelper(t, []byte{'2', ';', '3', 'D'}, Ground, []string{"CUB([2 3])"})
	funcCallParamHelper(t, []byte{'2', ';', '3', ';', '4', 'D'}, Ground, []string{"CUB([2 3 4])"})
}

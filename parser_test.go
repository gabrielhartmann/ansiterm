package ansiterm

import (
	"testing"
)

func anyToXHelper(t *testing.T, bytes []byte, expectedState State) {
	for i := 0; i < int(stateCount); i++ {
		s := StateId(i)
		parser, _ := createTestParser(stateMap[s])
		parser.Parse(bytes)
		validateState(t, parser.state, expectedState)
	}
}

func TestAnyToX(t *testing.T) {
	anyToXHelper(t, []byte{ANSI_ESCAPE_PRIMARY}, Escape)
	anyToXHelper(t, []byte{DCS_ENTRY}, DcsEntry)
	anyToXHelper(t, []byte{OSC_STRING}, OscString)
	anyToXHelper(t, []byte{CSI_ENTRY}, CsiEntry)
}

func TestEscapeToCsiEntry(t *testing.T) {
	parser, _ := createTestParser(Escape)
	parser.Parse([]byte{ANSI_ESCAPE_SECONDARY})
	validateState(t, parser.state, CsiEntry)
}

func stateTransitionHelper(t *testing.T, start State, end State, bytes []byte) {
	for _, b := range bytes {
		bytes := []byte{byte(b)}
		parser, _ := createTestParser(start)
		parser.Parse(bytes)
		validateState(t, parser.state, end)
	}
}

func TestCsiEntryToX(t *testing.T) {
	stateTransitionHelper(t, CsiEntry, Ground, AllCase)
	stateTransitionHelper(t, CsiEntry, CsiParam, CsiCollectables)
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

func parseParamsHelper(t *testing.T, bytes []byte, expectedParams []int) {
	params, err := parseParams(bytes)

	if err != nil {
		t.Errorf("Parameter parse error: %v", err)
		return
	}

	if len(params) != len(expectedParams) {
		t.Errorf("Parsed   parameters: %v", params)
		t.Errorf("Expected parameters: %v", expectedParams)
		t.Errorf("Parameter length failure: %d != %d", len(params), len(expectedParams))
		return
	}

	for i, v := range expectedParams {
		if v != params[i] {
			t.Errorf("Parsed   parameters: %v", params)
			t.Errorf("Expected parameters: %v", expectedParams)
			t.Errorf("Parameter parse failure: %d != %d at position %d", v, params[i], i)
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

func funcCallParamHelper(t *testing.T, bytes []byte, expectedState State, expectedCalls []string) {
	parser, evtHandler := createTestParser(CsiEntry)
	parser.Parse(bytes)
	validateState(t, parser.state, expectedState)
	validateFuncCalls(t, evtHandler.FunctionCalls, expectedCalls)
}

func TestCursorParam(t *testing.T) {
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

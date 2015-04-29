package ansiterm

import (
	"testing"
)

func validateState(t *testing.T, actualState State, expectedState State) {
	actualName := "Nil"
	expectedName := "Nil"

	if actualState != nil {
		actualName = actualState.Name()
	}

	if expectedState != nil {
		expectedName = expectedState.Name()
	}

	if actualState != expectedState {
		t.Errorf("Invalid State: '%s' != '%s'", actualName, expectedName)
	}
}

func validateFuncCalls(t *testing.T, actualCalls []string, expectedCalls []string) {
	actualCount := len(actualCalls)
	expectedCount := len(expectedCalls)

	if actualCount != expectedCount {
		t.Errorf("Actual   calls: %v", actualCalls)
		t.Errorf("Expected calls: %v", expectedCalls)
		t.Errorf("Call count error: %d != %d", actualCount, expectedCount)
	}

	for i, v := range actualCalls {
		if v != expectedCalls[i] {
			t.Errorf("Actual   calls: %v", actualCalls)
			t.Errorf("Expected calls: %v", expectedCalls)
			t.Errorf("Mismatched calls: %s != %s", v, expectedCalls[i])
		}
	}
}

func createTestParser(s State) (AnsiParser, *TestAnsiEventHandler) {
	evtHandler := CreateTestAnsiEventHandler()
	parser := CreateParser(s, &evtHandler)

	return parser, &evtHandler
}

func anyToStateHelper(t *testing.T, bytes []byte, expectedState State) {
	for i := 0; i < int(stateCount); i++ {
		s := StateId(i)
		parser, _ := createTestParser(stateMap[s])
		parser.Parse(bytes)
		validateState(t, parser.state, expectedState)
	}
}

func TestAnyToEscapeTransition(t *testing.T) {
	bytes := []byte{ANSI_ESCAPE_PRIMARY}
	expectedState := Escape
	anyToStateHelper(t, bytes, expectedState)
}

func TestAnyToDcsEntryTransition(t *testing.T) {
	bytes := []byte{DCS_ENTRY}
	expectedState := DcsEntry
	anyToStateHelper(t, bytes, expectedState)
}

func TestAnyToOcsStringTransition(t *testing.T) {
	bytes := []byte{OSC_STRING}
	expectedState := OscString
	anyToStateHelper(t, bytes, expectedState)
}

func TestAnyToCsiEntryTransition(t *testing.T) {
	bytes := []byte{CSI_ENTRY}
	expectedState := CsiEntry
	anyToStateHelper(t, bytes, expectedState)
}

func TestEscapeToCsiEntry(t *testing.T) {
	bytes := []byte{ANSI_ESCAPE_SECONDARY}
	parser, _ := createTestParser(Escape)
	parser.Parse(bytes)
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

func TestCsiEntryToGround(t *testing.T) {
	stateTransitionHelper(t, CsiEntry, Ground, AllCase)
}

func TestCsiEntryToCsiParam(t *testing.T) {
	stateTransitionHelper(t, CsiEntry, CsiParam, CsiCollectables)
}

func csiToGroundNoParamHelper(t *testing.T, b byte, funcCall string) {
	bytes := []byte{b}
	parser, evtHandler := createTestParser(CsiEntry)
	parser.Parse(bytes)
	validateState(t, parser.state, Ground)
	validateFuncCalls(t, evtHandler.FunctionCalls, []string{funcCall})
}

func TestCUU(t *testing.T) {
	csiToGroundNoParamHelper(t, 'A', "CUU")
}

func TestCUD(t *testing.T) {
	csiToGroundNoParamHelper(t, 'B', "CUD")
}

func TestCUF(t *testing.T) {
	csiToGroundNoParamHelper(t, 'C', "CUF")
}

func TestCUB(t *testing.T) {
	csiToGroundNoParamHelper(t, 'D', "CUB")
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

func TestCUUParam(t *testing.T) {
	funcCallParamHelper(t, []byte{'2', 'A'}, Ground, []string{"CUU([2])"})
	funcCallParamHelper(t, []byte{'2', '3', 'A'}, Ground, []string{"CUU([23])"})
	funcCallParamHelper(t, []byte{'2', ';', '3', 'A'}, Ground, []string{"CUU([2 3])"})
	funcCallParamHelper(t, []byte{'2', ';', '3', ';', '4', 'A'}, Ground, []string{"CUU([2 3 4])"})

	funcCallParamHelper(t, []byte{'2', 'B'}, Ground, []string{"CUD([2])"})
	funcCallParamHelper(t, []byte{'2', '3', 'B'}, Ground, []string{"CUD([23])"})
	funcCallParamHelper(t, []byte{'2', ';', '3', 'B'}, Ground, []string{"CUD([2 3])"})
	funcCallParamHelper(t, []byte{'2', ';', '3', ';', '4', 'B'}, Ground, []string{"CUD([2 3 4])"})

	funcCallParamHelper(t, []byte{'2', 'C'}, Ground, []string{"CUF([2])"})
	funcCallParamHelper(t, []byte{'2', '3', 'C'}, Ground, []string{"CUF([23])"})
	funcCallParamHelper(t, []byte{'2', ';', '3', 'C'}, Ground, []string{"CUF([2 3])"})
	funcCallParamHelper(t, []byte{'2', ';', '3', ';', '4', 'C'}, Ground, []string{"CUF([2 3 4])"})

	funcCallParamHelper(t, []byte{'2', 'D'}, Ground, []string{"CUB([2])"})
	funcCallParamHelper(t, []byte{'2', '3', 'D'}, Ground, []string{"CUB([23])"})
	funcCallParamHelper(t, []byte{'2', ';', '3', 'D'}, Ground, []string{"CUB([2 3])"})
	funcCallParamHelper(t, []byte{'2', ';', '3', ';', '4', 'D'}, Ground, []string{"CUB([2 3 4])"})
}

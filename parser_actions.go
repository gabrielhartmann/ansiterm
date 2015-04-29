package ansiterm

import (
	log "github.com/Sirupsen/logrus"
)

func (ap *AnsiParser) collectParam(b byte) error {
	log.Infof("AnsiParser::collectParam %#x", b)
	ap.context.parameterBuffer = append(ap.context.parameterBuffer, b)
	return nil
}

func parseParams(bytes []byte) ([]int, error) {
	paramBuff := make([]byte, 0, 0)
	params := []int{}

	for _, v := range bytes {
		if v == ';' {
			if len(paramBuff) > 0 {
				// Completed parameter, append it to the list
				i := convertBytesToInteger(paramBuff)
				params = append(params, i)
				paramBuff = make([]byte, 0, 0)
			}
		} else {
			paramBuff = append(paramBuff, v)
		}
	}

	// Last parameter may not be terminated with ';'
	if len(paramBuff) > 0 {
		i := convertBytesToInteger(paramBuff)
		params = append(params, i)
	}

	log.Infof("Parsed params: %v with length: %d", params, len(params))
	return params, nil
}

func (ap *AnsiParser) csiDispatch() error {
	params, _ := parseParams(ap.context.parameterBuffer)

	switch ap.context.finalChar {
	case 'A':
		return ap.eventHandler.CUU(params)
	case 'B':
		return ap.eventHandler.CUD(params)
	case 'C':
		return ap.eventHandler.CUF(params)
	case 'D':
		return ap.eventHandler.CUB(params)
	}

	return nil
}

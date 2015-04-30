package ansiterm

import (
	"errors"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"strconv"
)

func (ap *AnsiParser) collectParam(b byte) error {
	log.Infof("AnsiParser::collectParam %#x", b)
	ap.context.paramBuffer = append(ap.context.paramBuffer, b)
	return nil
}

func (ap *AnsiParser) collectInter(b byte) error {
	log.Infof("AnsiParser::collectInter %#x", b)
	ap.context.interBuffer = append(ap.context.interBuffer, b)
	return nil
}

func parseParams(bytes []byte) ([]string, error) {
	paramBuff := make([]byte, 0, 0)
	params := []string{}

	for _, v := range bytes {
		if v == ';' {
			if len(paramBuff) > 0 {
				// Completed parameter, append it to the list
				s := string(paramBuff)
				params = append(params, s)
				paramBuff = make([]byte, 0, 0)
			}
		} else {
			paramBuff = append(paramBuff, v)
		}
	}

	// Last parameter may not be terminated with ';'
	if len(paramBuff) > 0 {
		s := string(paramBuff)
		params = append(params, s)
	}

	log.Infof("Parsed params: %v with length: %d", params, len(params))
	return params, nil
}

func parseCmd(context AnsiContext) (string, error) {
	return string(context.finalChar), nil
}

func getInt(params []string, dflt int) int {
	i := getInts(params, 1, dflt)[0]
	log.Infof("getInt: %v", i)
	return i
}

func getInts(params []string, minCount int, dflt int) []int {
	ints := []int{}

	for _, v := range params {
		i, _ := strconv.Atoi(v)
		ints = append(ints, i)
	}

	if len(ints) < minCount {
		remaining := minCount - len(ints)
		for i := 0; i < remaining; i++ {
			ints = append(ints, dflt)
		}
	}

	log.Infof("getInts: %v", ints)

	return ints
}

func (ap *AnsiParser) hDispatch(params []string) error {
	if len(params) == 1 && params[0] == "?25" {
		return ap.eventHandler.DECTCEM(true)
	}

	return nil
}

func (ap *AnsiParser) lDispatch(params []string) error {
	if len(params) == 1 && params[0] == "?25" {
		return ap.eventHandler.DECTCEM(false)
	}

	return nil
}

func getEraseParam(params []string) int {
	param := getInt(params, 0)
	if param < 0 || 3 < param {
		param = 0
	}

	return param
}

func (ap *AnsiParser) csiDispatch() error {
	cmd, _ := parseCmd(*ap.context)
	params, _ := parseParams(ap.context.paramBuffer)

	switch cmd {
	case "A":
		return ap.eventHandler.CUU(getInt(params, 1))
	case "B":
		return ap.eventHandler.CUD(getInt(params, 1))
	case "C":
		return ap.eventHandler.CUF(getInt(params, 1))
	case "D":
		return ap.eventHandler.CUB(getInt(params, 1))
	case "E":
		return ap.eventHandler.CNL(getInt(params, 1))
	case "F":
		return ap.eventHandler.CPL(getInt(params, 1))
	case "G":
		return ap.eventHandler.CHA(getInt(params, 1))
	case "H":
		ints := getInts(params, 2, 1)
		x, y := ints[0], ints[1]
		return ap.eventHandler.CUP(x, y)
	case "J":
		param := getEraseParam(params)
		return ap.eventHandler.ED(param)
	case "K":
		param := getEraseParam(params)
		return ap.eventHandler.EL(param)
	case "f":
		ints := getInts(params, 2, 1)
		x, y := ints[0], ints[1]
		return ap.eventHandler.HVP(x, y)
	case "h":
		return ap.hDispatch(params)
	case "l":
		return ap.lDispatch(params)
	case "m":
		return ap.eventHandler.SGR(getInts(params, 1, 0))
	}

	return errors.New(fmt.Sprintf("%v", ap.context))
}

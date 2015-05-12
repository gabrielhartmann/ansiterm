package ansiterm

import (
	"errors"
	"fmt"
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

func (ap *AnsiParser) csiDispatch() error {
	cmd, _ := parseCmd(*ap.context)
	params, _ := parseParams(ap.context.paramBuffer)

	log.Infof("csiDispatch: %v(%v)", cmd, params)

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
	case "S":
		return ap.eventHandler.SU(getInt(params, 0))
	case "T":
		return ap.eventHandler.SD(getInt(params, 0))
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
	case "r":
		ints := getInts(params, 2, 1)
		top, bottom := ints[0], ints[1]
		return ap.eventHandler.DECSTBM(top, bottom)
	default:
		return errors.New(fmt.Sprintf("Unsupported CSI command: '%s', with full context:  %v", cmd, ap.context))
	}

}

func (ap *AnsiParser) print() error {
	log.Infof("AnsiParser::print %#x", ap.context.currentChar)
	return ap.eventHandler.Print(ap.context.currentChar)
}

func (ap *AnsiParser) clear() error {
	ap.context = &AnsiContext{}
	return nil
}

func (ap *AnsiParser) execute() error {
	log.Infof("AnsiParser::execute %#x", ap.context.currentChar)

	return ap.eventHandler.Execute(ap.context.currentChar)

}

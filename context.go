package ansiterm

type AnsiContext struct {
	currentChar     byte
	finalChar       byte
	parameterBuffer []byte
}

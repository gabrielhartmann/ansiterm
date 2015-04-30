package ansiterm

type AnsiContext struct {
	currentChar byte
	finalChar   byte
	paramBuffer []byte
	interBuffer []byte
}

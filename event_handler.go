package ansiterm

type AnsiEventHandler interface {
	// CUrsor Up
	CUU([]int) error

	// CUrsor Down
	CUD([]int) error

	//CUrsor Forward
	CUF([]int) error

	//CUrsor Backward
	CUB([]int) error

	//Cursor to Next Line
	CNL([]int) error

	//Cursor to Previous Line
	CPL([]int) error
}

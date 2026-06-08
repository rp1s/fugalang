package reporter

import (
	"fmt"
	"fugu/pkg/token"
	"strings"
	"sync/atomic"
)

type Report interface {
	Input() string
}

type Reporter struct {
	Source Report
	input  string
	lines  []string
	err    chan Err

	isInit  atomic.Bool
	isClose atomic.Bool
}

type Err struct {
	FileName string
	Msg      string
	Start    int
	End      int
	Pos      token.Position
}

// Init делать не надо при создание через конструктор
func New(source Report, fileName string) *Reporter {
	rp := &Reporter{
		Source: source,
	}
	rp.Init()
	return rp
}

func (rp *Reporter) Init() {
	if !rp.isInit.Load() {
		rp.input = rp.Source.Input()
		rp.lines = strings.Split(rp.input, "\n")
		rp.err = make(chan Err, 64)
		rp.isInit.Store(true)

		go rp.outputer()
	} else {
		panic("Cannot initialize twice")
	}
}

func (rp *Reporter) Close() {
	if !rp.isClose.Load() {
		rp.isClose.Store(true)
		close(rp.err)
	} else {
		panic("Cant close it twice")
	}
}

func (rp *Reporter) Send(err Err) {
	if !rp.isClose.Load() {
		rp.err <- err
	} else {
		panic("You cant write to a closed reporter")
	}
}

func (rp *Reporter) SSend(msg string, tk token.Token) {
	rp.Send(Err{
		FileName: tk.Pos.FileName,
		Msg:      msg,
		Start:    tk.Start,
		End:      tk.End,
		Pos:      tk.Pos,
	})
}

func (rp *Reporter) outputer() {
	for err := range rp.err {
		rp.print(err)
	}
}

func (rp *Reporter) print(err Err) {
	fmt.Println(BoldCyan(fmt.Sprintf("%s:%d:%d:", err.FileName, err.Pos.Line, err.Pos.Column)))

	arrowsLen := err.End - err.Start
	if arrowsLen <= 0 {
		arrowsLen = 1
	}

	rawLines := rp.getLine(err)
	if rawLines == "" {
		fmt.Printf("%s %s\n", ErrorLabel("error:"), err.Msg)
		fmt.Printf("%s%s \n", Gray(fmt.Sprintf("%2d ", err.Pos.Line)), Gray("|"))

		padding := strings.Repeat(" ", 6+(err.Pos.Column-1))
		fmt.Printf("%s%s\n\n", padding, BoldRed(strings.Repeat("^", arrowsLen)))
		return
	}

	errorLines := strings.Split(rawLines, "\n")
	for i, line := range errorLines {
		fmt.Printf("%s%s %s\n", Gray(fmt.Sprintf("%2d ", err.Pos.Line+i)), Gray("|"), line)
	}

	if len(errorLines) > 1 {
		arrowsLen = len(errorLines) - (err.Pos.Column - 1)
		if arrowsLen <= 0 {
			arrowsLen = 1
		}
	}

	padding := strings.Repeat(" ", 6+(err.Pos.Column-1))
	fmt.Printf("%s%s\n%s%s\n\n", padding, BoldRed(strings.Repeat("^", arrowsLen)), padding, BoldRed(err.Msg))
}

func (rp *Reporter) getLine(err Err) string {
	if err.Pos.Line-1 < 0 || err.Pos.Line-1 >= len(rp.lines) {
		return ""
	}

	tokenText := rp.input[err.Start:err.End]
	lineCount := strings.Count(tokenText, "\n")

	if lineCount == 0 {
		return rp.lines[err.Pos.Line-1]
	}

	endLine := err.Pos.Line - 1 + lineCount
	if endLine >= len(rp.lines) {
		endLine = len(rp.lines) - 1
	}

	return strings.Join(rp.lines[err.Pos.Line-1:endLine+1], "\n")
}

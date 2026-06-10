package reporter

import (
	"bytes"
	"fmt"
	"fugu/pkg/token"
	"strings"
	"sync/atomic"
)

type Report interface {
	Input() *[]byte
}

type Reporter struct {
	Source  Report
	lines   [][]byte
	err     chan Err
	isInit  atomic.Bool
	isClose atomic.Bool

	IsUse bool // чтобы знать были ли ошибки на прошлом этапе
}

type Msg interface {
	Code() string
	String() string
	Notes() []string
	Arrow() string
}

type Err struct {
	Code     Msg
	FileName string
	Msg      Msg
	ArrowMsg Msg
	Start    int
	End      int
	Pos      token.Position
	Notes    Msg
}

func New(source Report, fileName string) *Reporter {
	rp := &Reporter{
		Source: source,
	}
	rp.Init()
	return rp
}

func (rp *Reporter) Init() {
	if !rp.isInit.Load() {
		rp.lines = SplitLines(*rp.Source.Input())
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
	if !rp.IsUse {
		rp.IsUse = true
	}
	if !rp.isClose.Load() {
		rp.err <- err
	} else {
		panic("You cant write to a closed reporter")
	}
}

func (rp *Reporter) SendTk(msg Msg, tk token.Token) {
	rp.Send(Err{
		Code:     msg,
		FileName: tk.Pos.FileName,
		Msg:      msg,
		ArrowMsg: msg,
		Notes:    msg,
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
	label := "error"
	if err.Code.Code() != "" {
		label = fmt.Sprintf("error[%s]", err.Code.Code())
	}
	fmt.Printf("%s: %s\n", BoldRed(label), err.Msg.String())
	fmt.Printf("%s %s:%d:%d\n", BoldCyan("  -->"), err.FileName, err.Pos.Line, err.Pos.Column)

	arrowsLen := err.End - err.Start
	if arrowsLen <= 0 {
		arrowsLen = 1
	}

	maxLine := err.Pos.Line
	rawLines := rp.getLine(err)
	var errorLines [][]byte
	if rawLines != nil {
		errorLines = SplitLines(rawLines)
		maxLine = err.Pos.Line + len(errorLines) - 1
	}

	width := len(fmt.Sprintf("%d", maxLine))
	if width < 2 {
		width = 2
	}

	emptyPrefix := fmt.Sprintf("%s%s ", strings.Repeat(" ", width), Gray("|"))

	if rawLines == nil {
		fmt.Printf("%s%s \n", Gray(fmt.Sprintf("%*d", width, err.Pos.Line)), Gray("|"))
		prefixLen := width + 2
		padding := strings.Repeat(" ", prefixLen+(err.Pos.Column-1))
		if err.ArrowMsg.Arrow() != "" {
			fmt.Printf("%s%s %s\n", padding, BoldRed(strings.Repeat("^", arrowsLen)), BoldRed(err.ArrowMsg.Arrow()))
		} else {
			fmt.Printf("%s%s\n", padding, BoldRed(strings.Repeat("^", arrowsLen)))
		}
	} else {
		fmt.Println(emptyPrefix)
		for i, line := range errorLines {
			fmt.Printf("%s%s %s\n", Gray(fmt.Sprintf("%*d", width, err.Pos.Line+i)), Gray("|"), line)
		}
		if len(errorLines) > 1 {
			arrowsLen = len(errorLines) - (err.Pos.Column - 1)
			if arrowsLen <= 0 {
				arrowsLen = 1
			}
		}
		prefixLen := width + 2
		padding := strings.Repeat(" ", prefixLen+(err.Pos.Column-1))
		if err.ArrowMsg.Arrow() != "" {
			fmt.Printf("%s%s %s\n", padding, BoldRed(strings.Repeat("^", arrowsLen)), BoldRed(err.ArrowMsg.Arrow()))
		} else {
			fmt.Printf("%s%s\n", padding, BoldRed(strings.Repeat("^", arrowsLen)))
		}
	}

	if len(err.Notes.Notes()) > 0 {
		fmt.Println(emptyPrefix)
		for _, note := range err.Notes.Notes() {
			fmt.Printf("%s%s %s\n", strings.Repeat(" ", width), BoldCyan("="), note)
		}
	}
	fmt.Println()
}

func (rp *Reporter) getLine(err Err) []byte {
	lineIdx := err.Pos.Line - 1
	if lineIdx < 0 || lineIdx >= len(rp.lines) {
		return []byte{}
	}

	tokenText := (*rp.Source.Input())[err.Start:err.End]
	if !bytes.Contains(tokenText, []byte{'\n'}) {
		return rp.lines[lineIdx]
	}

	lineCount := bytes.Count(tokenText, []byte{'\n'})

	endLine := lineIdx + lineCount
	if endLine >= len(rp.lines) {
		endLine = len(rp.lines) - 1
	}

	return bytes.Join(rp.lines[lineIdx:endLine+1], []byte{'\n'})
}

func SplitLines(data []byte) [][]byte {
	if len(data) == 0 {
		return nil
	}
	return bytes.Split(data, []byte{'\n'})
}

package reporter

import (
	"bytes"
	"fmt"
	"fugu/pkg/token"
	"strings"
	"sync"
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

	wg sync.WaitGroup
}

type Msg interface {
	Code() string
	Msg() string
	Notes() []string
	Arrow() string
	IsUseBlock() bool
}

type Err struct {
	Code       string
	FileName   string
	Msg        string
	ArrowMsg   string
	Start      int
	End        int
	Pos        token.Position
	Notes      []string
	IsUseBlock bool
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
		rp.wg.Add(1)
		go rp.outputer()
	} else {
		panic("Cannot initialize twice")
	}
}

func (rp *Reporter) Close() {
	if !rp.isClose.Load() {
		rp.isClose.Store(true)
		close(rp.err)
		rp.wg.Wait()
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
		Code:       msg.Code(),
		FileName:   tk.Pos.FileName,
		Msg:        msg.Msg(),
		ArrowMsg:   msg.Arrow(),
		Notes:      msg.Notes(),
		IsUseBlock: msg.IsUseBlock(),
		Start:      tk.Start,
		End:        tk.End,
		Pos:        tk.Pos,
	})
}

func (rp *Reporter) outputer() {
	defer rp.wg.Done()
	for err := range rp.err {
		fmt.Println(rp.buildMsg(err).String())
	}
}

func (rp *Reporter) buildMsg(err Err) *strings.Builder {
	var out *strings.Builder

	label := "error"
	if err.Code != "" {
		label = fmt.Sprintf("error[%s]", err.Code)
	}
	out.WriteString(fmt.Sprintf("%s: %s\n", BoldRed(label), err.Msg))
	out.WriteString(fmt.Sprintf("%s %s:%d:%d\n", BoldYellow(" -->"), err.FileName, err.Pos.Line, err.Pos.Column))

	if err.IsUseBlock {
		arrowLen := err.End - err.Start
		if arrowLen <= 0 {
			arrowLen = 1
		}
		rawLines := rp.getLine(err)

		maxLine := err.Pos.Line
		width := len(fmt.Sprintf("%d", maxLine))
		if width < 2 {
			width = 2
		}

		if len(rawLines) == 0 {
			out.WriteString(fmt.Sprintf("%s%s \n", Gray(fmt.Sprintf("%*d", width, err.Pos.Line)), Gray("|")))
			padding := strings.Repeat(" ", width+3+(err.Pos.Column-1))
			if err.ArrowMsg != "" {
				out.WriteString(fmt.Sprintf("%s%s %s\n", padding, BoldRed(strings.Repeat("^", arrowLen)), BoldRed(err.ArrowMsg)))
			} else {
				out.WriteString(fmt.Sprintf("%s%s\n", padding, BoldRed(strings.Repeat("^", arrowLen))))
			}
		} else {
			prefix := fmt.Sprintf("%s%s ", strings.Repeat(" ", width), Gray("|"))
			out.WriteString(prefix + "\n")

			startLineNum := err.Pos.Line - len(rawLines) + 1
			if err.Pos.Line <= len(rawLines) {
				startLineNum = 1
			}

			var targetLine []byte

			for i, line := range rawLines {
				lineNum := startLineNum + i
				if lineNum == err.Pos.Line {
					targetLine = line
				}
				out.WriteString(fmt.Sprintf("%s%s %s\n", Gray(fmt.Sprintf("%*d", width, lineNum)), Gray("|"), string(line)))
			}

			lineStrPrefix := fmt.Sprintf("%*d| ", width, err.Pos.Line)
			basePadding := len(lineStrPrefix)

			// символы перед целевой линией
			cbtl := 0
			for i := 0; i < err.Pos.Line-1 && i < len(rp.lines); i++ {
				cbtl += len(rp.lines[i]) + 1
			}

			//  bbtl - байты перед токеном в строке
			bbtl := err.Start - cbtl
			if bbtl < 0 {
				bbtl = 0
			}
			if bbtl > len(targetLine) {
				bbtl = len(targetLine)
			}

			prefixBytes := targetLine[:bbtl]
			prefixRunes := []rune(string(prefixBytes))

			codePadding := 0
			for _, r := range prefixRunes {
				if r == '\t' {
					codePadding += 4
				} else {
					codePadding += 1
				}
			}

			padding := strings.Repeat(" ", basePadding+codePadding)

			if err.ArrowMsg != "" {
				out.WriteString(fmt.Sprintf("%s%s %s\n", padding, BoldRed(strings.Repeat("^", arrowLen)), BoldRed(err.ArrowMsg)))
			} else {
				out.WriteString(fmt.Sprintf("%s%s\n", padding, BoldRed(strings.Repeat("^", arrowLen))))
			}
		}
	}
	if len(err.Notes) > 0 {
		width := 4
		for _, note := range err.Notes {
			out.WriteString(fmt.Sprintf("%s%s %s\n", strings.Repeat(" ", width), BoldGreen("="), note))
		}
	}
	out.WriteString("\n")

	return out
}

func (rp *Reporter) getLine(err Err) [][]byte {
	if len(rp.lines) == 0 {
		return nil
	}
	lidx := err.Pos.Line - 1
	if lidx < 0 || lidx >= len(rp.lines) {
		return nil
	}
	start := lidx - 3
	if start < 0 {
		start = 0
	}
	end := lidx + 1
	if end > len(rp.lines) {
		end = len(rp.lines)
	}
	if start >= end {
		return nil
	}
	return rp.lines[start:end]
}

func SplitLines(data []byte) [][]byte {
	if len(data) == 0 {
		return nil
	}
	return bytes.Split(data, []byte{'\n'})
}

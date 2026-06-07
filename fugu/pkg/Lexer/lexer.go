package lexer

import (
	"fugu/pkg/token"
	"unicode"
	"unicode/utf8"
)

type Lexer struct {
	input string
	rn    rune // текущая rune

	curPos         int // абсолютное смещение c начала файла
	tokStart       int // абсолютное смещение до начала токена который разбираеться прямо сейчас
	tokStartLine   int // номер строки начала токена
	tokStartColumn int // номер колонки начала токена
	pos            token.Position
}

func (lex *Lexer) Reset() {
	lex = New(lex.input, lex.pos.FileName)
}

func New(input, fileName string) *Lexer {
	lex := &Lexer{
		input:  input,
		curPos: 0,
		pos: token.Position{
			FileName: fileName,
			Line:     1,
			Column:   1,
			Offset:   0,
		},
	}
	lex.advance()
	return lex
}

func (lex *Lexer) NextToken() token.Token {
	lex.tokStart = lex.pos.Offset
	lex.tokStartLine = lex.pos.Line
	lex.tokStartColumn = lex.pos.Column

	if unicode.IsSpace(lex.rn) {
		for unicode.IsSpace(lex.rn) {
			lex.advance()
		}
		return lex.NewToken(token.SPACING)
	}

	switch lex.rn {
	case '/':
		if lex.peekRn() == '/' {
			return lex.readLineComment()
		} else if lex.peekRn() == '*' {
			return lex.readMultiLineComment()
		} else {
			lex.NewToken(token.DIVIDE)
		}

	}
	// TODO: разбор операторов и тд
	lex.advance()
	return lex.NewToken(token.ILLEGAL)
}

func (lex *Lexer) readLineComment() token.Token {
	// Пропускаем '//'
	lex.advance()
	lex.advance()

	// останавливаемся перед '\n'
	for lex.rn != '\n' && lex.rn != 0 {
		lex.advance()
	}

	return lex.NewToken(token.COMMENT)
}

func (lex *Lexer) readMultiLineComment() token.Token {
	// пропуск '/*'
	lex.advance()
	lex.advance()

	for {
		// TODO: надо выкинуть в верх ошибку: файл закончился, а комментарий так и не был закрыт
		if lex.rn == 0 {
			return lex.NewToken(token.ILLEGAL)
		}

		if lex.rn == '*' && lex.peekRn() == '/' {
			lex.advance() // '*'
			lex.advance() // '/'
			break
		}

		lex.advance()
	}

	return lex.NewToken(token.COMMENT)
}

//
// вспомогательный функционал
//

func (lex *Lexer) advance() {
	if lex.curPos >= len(lex.input) {
		lex.rn = 0 // \x00
		lex.pos.Offset = lex.curPos
		return
	}

	r, size := utf8.DecodeRuneInString(lex.input[lex.curPos:])

	lex.rn = r
	lex.pos.Offset = lex.curPos
	lex.curPos += size

	if lex.rn == '\n' {
		lex.pos.Line++
		lex.pos.Column = 1
	} else {
		lex.pos.Column++
	}
}

// 0 - следущий симвл после Lexer.curPos
func (lex *Lexer) peekRn() rune {
	if lex.curPos >= len(lex.input) {
		return 0
	}

	r, _ := utf8.DecodeRuneInString(lex.input[lex.curPos:])

	return r
}

func (lex *Lexer) NewToken(kind token.TokenKind) token.Token {
	return token.Token{
		Kind: kind,
		Pos: token.Position{
			FileName: lex.pos.FileName,
			Line:     lex.tokStartLine,
			Column:   lex.tokStartColumn,
			Offset:   lex.tokStart,
		},
		Start: lex.tokStart,
		End:   lex.pos.Offset,
	}
}

func (lex *Lexer) LiteralToken(tk token.Token) string {
	return lex.input[tk.Start:tk.End]
}

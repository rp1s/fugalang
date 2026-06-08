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

// для интерфейса для возможности получить Literal коректно
func (lex *Lexer) Input() string {
	return lex.input
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

	if lex.rn == 0 {
		return lex.NewToken(token.EOF)
	}

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

	case '.':
		if lex.peekRn() == '.' {
			lex.advance() // едим первую .
			if lex.peekRn() == '=' {
				lex.advance().advance()
				return lex.NewToken(token.RANGE_INCL)
			} else if lex.peekRn() == '<' {
				lex.advance().advance()
				return lex.NewToken(token.RANGE_HALF_OPEN)
			}
			lex.advance()
			return lex.NewToken(token.OP_RANGE)
		}
		lex.advance()
		return lex.NewToken(token.DOT)
	}

	// TODO: разбор операторов и тд
	lex.advance()
	return lex.NewToken(token.ILLEGAL)
}

func (lex *Lexer) readLineComment() token.Token {
	// Пропускаем '//'
	lex.advance().advance()

	// останавливаемся перед '\n'
	for lex.rn != '\n' && lex.rn != 0 {
		lex.advance()
	}

	return lex.NewToken(token.COMMENT)
}

func (lex *Lexer) readMultiLineComment() token.Token {
	// пропуск '/*'
	lex.advance().advance()

	for {
		// TODO: надо выкинуть в верх ошибку: файл закончился, а комментарий так и не был закрыт
		if lex.rn == 0 {
			return lex.NewToken(token.ILLEGAL)
		}

		if lex.rn == '*' && lex.peekRn() == '/' {
			lex.advance().advance() // '*', '/'
			break
		}

		lex.advance()
	}

	return lex.NewToken(token.M_COMMENT)
}

//
// вспомогательный функционал
//

func (lex *Lexer) advance() *Lexer {
	if lex.curPos >= len(lex.input) {
		lex.rn = 0 // \x00
		lex.pos.Offset = lex.curPos
		return lex
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

	return lex
}

// возвращает следущий симвл после Lexer.curPos
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

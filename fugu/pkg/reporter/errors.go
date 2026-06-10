package reporter

type Code uint16

const (
	NoError Code = iota
	TestError

	// NoClosing
	LexerNoClosing // не закрыт блок и тд
)

// TODO: надо будет очев сделать чтобы возвращал или на англ или на ру
func (c Code) String() string {
	switch c {
	case NoError:
		return "нету"
	case TestError:
		return "тестовая"
	case LexerNoClosing:
		return "пропущен закрывающий символ"
	default:
		return "неизвестная"
	}
}

func (c Code) Notes() []string {
	switch c {
	case NoError:
		return []string{
			"Ой, кажется, я сломался изнутри! Этого не должно было быть напечатано.",
			"Пожалуйста, создайте баг-репорт:",
			"  https://github.com/fugalang/fugu/issues",
			"Если вам не трудно, опишите в репорте сценарий, при котором вы обнаружили это недоразумение.",
		}
	case TestError:
		return []string{
			"Ой, кажется, я сломался изнутри! Это внутренняя отладочная ошибка.",
			"Я не должен был спотыкаться на этом месте при обычной работе.",
			"Пожалуйста, создайте об этом баг-репорт:",
			"  https://github.com/fugalang/fugu/issues",
			"Если вам не трудно, опишите в репорте сценарий, при котором вы обнаружили эту проблему.",
		}

	case LexerNoClosing:
		return []string{""}
	default:
		return []string{}
	}
}

func (c Code) Code() string {
	switch c {
	case NoError:
		return "NOERR"
	case TestError:
		return "TEST"
	case LexerNoClosing:
		return "LNC1"
	default:
		return "I!K"
	}
}

func (c Code) Arrow() string {
	switch c {
	case NoError:
		return ""
	case TestError:
		return ""
	case LexerNoClosing:
		return "закрой за собой!"
	default:
		return ""
	}
}

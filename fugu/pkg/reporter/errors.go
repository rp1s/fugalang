package reporter

type Code uint16

const (
	NoError Code = iota
	TestError

	// NoClosing
	LexerNoClosing // не закрыт блок и тд

	ParserCantStartWork
	StateDoesNotToken // состояние не расчитывает на этот токен
)

// TODO: надо будет очев сделать чтобы возвращал или на англ или на ру
func (c Code) Msg() string {
	switch c {
	case NoError:
		return "не найденна"
	case TestError:
		return "тестовая ошибка"
	case LexerNoClosing:
		return "пропущен закрывающий символ"
	case StateDoesNotToken:
		return "ошибка при работе с таблицой состояний"
	case ParserCantStartWork:
		return "нету возможности запучтить разбор"
	default:
		return "неизвестная ошибка"
	}
}

func (c Code) Notes() []string {
	switch c {
	case NoError:
		return []string{
			"Внутренняя ошибка отладки: сбой компонента.",
			"Некорректное состояние, не ожидаемое при штатной работе.",
			"Пожалуйста, сообщите об ошибке по адресу:",
			"  https://github.com/fugalang/fugu/issues",
			"По возможности приложите описание сценария воспроизведения.",
		}
	case TestError:
		return []string{
			"Внутренняя ошибка отладки: сбой компонента.",
			"Некорректное состояние, не ожидаемое при штатной работе.",
			"Пожалуйста, сообщите об ошибке по адресу:",
			"  https://github.com/fugalang/fugu/issues",
			"По возможности приложите описание сценария воспроизведения.",
		}

	case StateDoesNotToken:
		return []string{
			"Внутренняя ошибка отладки: сбой компонента.",
			"Некорректное состояние, не ожидаемое при штатной работе.",
			"Пожалуйста, сообщите об ошибке по адресу:",
			"  https://github.com/fugalang/fugu/issues",
			"По возможности приложите описание сценария воспроизведения.",
		}

	case ParserCantStartWork:
		return []string{
			"Исправьте прошлые ошибки, чтобы парсер отработал корректно.",
		}

	default:
		return []string{}
	}
}

func (c Code) Code() string {
	switch c {
	case NoError:
		return "NoError"
	case TestError:
		return "TestError"
	case LexerNoClosing:
		return "LexerNoClosing"
	case ParserCantStartWork:
		return "ParserCantStartWork"
	case StateDoesNotToken:
		return "StateDoesNotToken"
	default:
		return "NoError"
	}
}

func (c Code) Arrow() string {
	switch c {
	case LexerNoClosing:
		return "закрой за собой!"
	default:
		return ""
	}
}

func (c Code) IsUseBlock() bool {
	switch c {
	case NoError:
		return false
	case StateDoesNotToken:
		return false
	case ParserCantStartWork:
		return false
	default:
		return true
	}
}

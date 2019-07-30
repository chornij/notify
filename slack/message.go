package slack

type Message struct {
	Text  string
	Color string
}

type NotifyType string

func Info(text string) Message {
	return Message{
		Text:  text,
		Color: "good",
	}
}

func Warning(text string) Message {
	return Message{
		Text:  text,
		Color: "warning",
	}
}

func Error(text string) Message {
	return Message{
		Text:  text,
		Color: "danger",
	}
}

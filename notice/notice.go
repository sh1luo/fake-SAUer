package notice

type Notifier interface {
	Notice(to, subject, body string) error
}

func NewNotifier(notifierName string, notifierStruct interface{}) Notifier {
	if notifierName == "Email" {
		return NewEmailNotifier(notifierStruct)
	}

	return nil
}

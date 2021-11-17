package notice

type Notifier interface {
	Notice(to, subject, body string) error
}

func NewNotifier(notifierName string, args ...interface{}) Notifier {
	if notifierName == "email" {
		return NewEmailNotifier(args)
	}

	return nil
}

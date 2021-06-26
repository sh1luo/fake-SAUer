package notice

type Notifier interface {
	Notice(to, subject, body string) error
}

func NewNotifier(notifierName, Account, Token, Host string, Post int) Notifier {
	if notifierName == "email" {
		return NewEmailNotifier(Account, Token, Host, Post)
	}
	
	return nil
}

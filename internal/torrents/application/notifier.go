package application

type Link struct {
	Text string
	Url  string
}

type Notification struct {
	Text string
	Link *Link
}

type Notifier interface {
	Notify(notification Notification) error
	NotifyText(text string, args ...interface{}) error
	NotifyLink(text string, link *Link) error
}

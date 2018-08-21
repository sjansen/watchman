package connection

type WatchmanError struct {
	msg string
}

func (e *WatchmanError) Error() string {
	return e.msg
}

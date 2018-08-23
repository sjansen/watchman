package protocol

// WatchmanError is returned when the Watchman server responds to a
// request with an error instead of a normal response.
type WatchmanError struct {
	msg string
}

func (e *WatchmanError) Error() string {
	return e.msg
}

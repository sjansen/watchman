package watchman

import "errors"

var MarshalTypeError = errors.New("watchman: Marshal must be called with a struct or pointer to struct")

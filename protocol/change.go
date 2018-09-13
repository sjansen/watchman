package protocol

// A Change describes how a File's state has changed.
type Change int

const (
	Created Change = 1 << iota
	Removed
	Updated
)

func (c Change) String() string {
	switch c {
	case Created:
		return "created"
	case Removed:
		return "removed"
	case Updated:
		return "updated"
	}
	return "unknown"
}

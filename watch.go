package watchman

type Watch interface {
	Clock(timeout int) (string, error)
}

type watch struct {
	conn          *Connection
	root          string
	relative_path string
}

func (w *watch) Clock(timeout int) (value string, err error) {
	var result object
	if timeout > 0 {
		result, err = w.conn.command("clock", w.root, map[string]int{"sync_timeout": timeout})
	} else {
		result, err = w.conn.command("clock", w.root)
	}
	if err != nil {
		return
	}
	value = result["clock"].(string)
	return
}

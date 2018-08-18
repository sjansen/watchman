package watchman

import (
	"bytes"
	"fmt"
	"reflect"
)


type GetConfig struct {
	Path string
}

type WatchList struct{}

type WatchProject struct {
	Path string
}

func Marshal(cmd interface{}) ([]byte, error) {
	v := reflect.ValueOf(cmd)
	k := v.Kind()
	t := v.Type()
	if k == reflect.Ptr {
		v = v.Elem()
		k = v.Kind()
		t = t.Elem()
	}
	if k != reflect.Struct {
		return nil, MarshalTypeError
	}

	b := bytes.Buffer{}
	b.WriteRune('[')

	// Command
	name := t.Name()
	b.WriteRune('"')
	for i, r := range name {
		if r >= 'A' && r <= 'Z' {
			if i > 0 {
				b.WriteRune('-')
			}
			b.WriteRune(r + 32)
		} else {
			b.WriteRune(r)
		}
	}
	b.WriteRune('"')

	// Args
	for i := 0; i < v.NumField(); i++ {
		b.WriteString(", ")
		f := v.Field(i)
		fmt.Fprintf(&b, "%q", f)
	}

	b.WriteString("]\n")
	return b.Bytes(), nil
}

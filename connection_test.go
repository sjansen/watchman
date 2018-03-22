package watchman

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type source int

const (
	client source = iota
	server
)

type step struct {
	src  source
	data string
}

type testcase struct {
	script      []step
	results     []event
	unilaterals []event
}

var testcases = map[string]testcase{
	"simple": testcase{
		script: []step{
			{client, `["version"]`},
			{server, `{"version":"4.9.0"}`},
			{client, `["list-capabilities"]`},
			{server, `{"capabilities":["relative_root","cmd-subscribe"],"version":"4.9.0"}`},
		},
		results: []event{
			{
				"version": "4.9.0",
			},
			{
				"version": "4.9.0",
				"capabilities": []interface{}{
					"relative_root", "cmd-subscribe",
				},
			},
		},
	},
}

func connect(t *testing.T, label string, script []step) (c *connection) {
	commands := make(chan string)
	events := make(chan string)
	c = &connection{
		commands: commands,
		events:   events,
	}

	go func() {
		defer close(commands)
		defer close(events)
		for i, step := range script {
			switch step.src {
			case client:
				actual := <-commands
				if step.data != actual {
					t.Errorf("[%s] step %d expected: %#v actual: %#v", label, i, step.data, actual)
				}
			case server:
				events <- step.data
			}
		}
	}()

	return
}

func TestEventLoop(t *testing.T) {
	assert := assert.New(t)

	for label, tc := range testcases {
		c := connect(t, label, tc.script)
		l := loop(c)

		n := 0
		for _, step := range tc.script {
			switch step.src {
			case client:
				l.commands <- step.data
				actual := <-l.results
				expected := tc.results[n]
				n += 1
				if !assert.Equal(expected, actual, label) {
					break
				}
			case server:
				continue
			}
		}
	}
}

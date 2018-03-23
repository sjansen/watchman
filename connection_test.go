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
	u8l  bool
	data string
}

type testcase struct {
	script      []step
	results     []object
	unilaterals []object
}

var testcases = map[string]testcase{
	"simple": testcase{
		script: []step{
			{client, false, `["version"]`},
			{server, false, `{"version":"4.9.0"}`},
			{client, false, `["list-capabilities"]`},
			{server, false, `{"capabilities":["relative_root","cmd-subscribe"],"version":"4.9.0"}`},
		},
		results: []object{
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
	"log-level": testcase{
		script: []step{
			{client, false, `["log-level", "error"]`},
			{server, false, `{"log_level":"error","version":"4.9.0"}`},
			{server, true, `{"level":"error","unilateral":true,"log":"2018-03-22T01:18:52,901: [client=0x7ffe1dc035d8:stm=0x7ffe1dc03460:pid=0] test message\n"}`},
		},
		results: []object{
			{
				"version":   "4.9.0",
				"log_level": "error",
			},
		},
		unilaterals: []object{
			{
				"unilateral": true,
				"level":      "error",
				"log":        "2018-03-22T01:18:52,901: [client=0x7ffe1dc035d8:stm=0x7ffe1dc03460:pid=0] test message\n",
			},
		},
	},
	"subscribe": testcase{
		script: []step{
			{client, false, `["subscribe", "/tmp", "sub1", {"fields": ["name"]}]`},
			{server, false, `{"version":"4.9.0","clock":"c:1521588867:575:1:2","subscribe":"sub1"}`},
			{server, true, `{"version":"4.9.0","unilateral":true,"subscription":"sub1","clock":"c:1521588867:575:1:2","root":"/tmp","files":["foo"],"is_fresh_instance":true}`},
			{server, true, `{"version":"4.9.0","unilateral":true,"subscription":"sub1","clock":"c:1521588867:575:1:3","since":"c:1521588867:575:1:2","root":"/tmp","files":["bar"],"is_fresh_instance":false}`},
		},
		results: []object{
			{
				"version":   "4.9.0",
				"clock":     "c:1521588867:575:1:2",
				"subscribe": "sub1",
			},
		},
		unilaterals: []object{
			{
				"version":           "4.9.0",
				"unilateral":        true,
				"subscription":      "sub1",
				"clock":             "c:1521588867:575:1:2",
				"root":              "/tmp",
				"files":             []interface{}{"foo"},
				"is_fresh_instance": true,
			},
			{
				"version":           "4.9.0",
				"unilateral":        true,
				"subscription":      "sub1",
				"clock":             "c:1521588867:575:1:3",
				"since":             "c:1521588867:575:1:2",
				"root":              "/tmp",
				"files":             []interface{}{"bar"},
				"is_fresh_instance": false,
			},
		},
	},
	"watch-project": testcase{
		script: []step{
			{client, false, `["watch-project", "/tmp"]`},
			{server, false, `{"version":"4.9.0","watcher":"fsevents","watch":"/tmp"}`},
			{client, false, `["clock", "/tmp"]`},
			{server, false, `{"version":"4.9.0","clock":"c:1521588867:575:1:5"}`},
			{client, false, `["subscribe", "/tmp", "sub2", {"since":"c:1521588867:575:1:5", "fields":["name"]}]`},
			{server, false, `{"version":"4.9.0","clock":"c:1521588867:575:1:5","subscribe":"sub2"}`},
			{server, true, `{"version":"4.9.0","unilateral":true,"subscription":"sub2","clock":"c:1521588867:575:1:6","since":"c:1521588867:575:1:5","root":"/tmp","files":["baz"],"is_fresh_instance":false}`},
			{server, true, `{"version":"4.9.0","unilateral":true,"subscription":"sub2","clock":"c:1521588867:575:1:7","since":"c:1521588867:575:1:6","root":"/tmp","files":["qux"],"is_fresh_instance":false}`},
		},
		results: []object{
			{
				"version": "4.9.0",
				"watch":   "/tmp",
				"watcher": "fsevents",
			},
			{
				"version": "4.9.0",
				"clock":   "c:1521588867:575:1:5",
			},
			{
				"version":   "4.9.0",
				"clock":     "c:1521588867:575:1:5",
				"subscribe": "sub2",
			},
		},
		unilaterals: []object{
			{
				"version":           "4.9.0",
				"unilateral":        true,
				"subscription":      "sub2",
				"clock":             "c:1521588867:575:1:6",
				"since":             "c:1521588867:575:1:5",
				"root":              "/tmp",
				"files":             []interface{}{"baz"},
				"is_fresh_instance": false,
			},
			{
				"version":           "4.9.0",
				"unilateral":        true,
				"subscription":      "sub2",
				"clock":             "c:1521588867:575:1:7",
				"since":             "c:1521588867:575:1:6",
				"root":              "/tmp",
				"files":             []interface{}{"qux"},
				"is_fresh_instance": false,
			},
		},
	},
}

func connect(t *testing.T, script []step) (c *connection) {
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
					t.Errorf("step %d expected: %#v actual: %#v", i, step.data, actual)
				}
			case server:
				events <- step.data
			}
		}
	}()

	return
}

func TestEventLoop(t *testing.T) {
	for label, tc := range testcases {
		t.Run(label, func(t *testing.T) {
			assert := assert.New(t)

			c := connect(t, tc.script)
			l := loop(c)

			results := 0
			unilaterals := 0
			for _, step := range tc.script {
				switch step.src {
				case client:

					expected := tc.results[results]
					results += 1
					l.commands <- step.data
					actual := <-l.results
					if !assert.Equal(expected, actual) {
						break
					}

				case server:

					if step.u8l {
						expected := tc.unilaterals[unilaterals]
						unilaterals += 1
						actual := <-l.unilaterals
						if !assert.Equal(expected, actual) {
							break
						}
					}

				}
			}
		})
	}
}

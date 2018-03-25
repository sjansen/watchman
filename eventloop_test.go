package watchman

import (
	"context"
	"testing"

	"github.com/fortytw2/leaktest"
	"github.com/stretchr/testify/assert"
)

type source int

const (
	CLIENT source = iota
	SERVER
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
			{CLIENT, false, `["version"]`},
			{SERVER, false, `{"version":"4.9.0"}`},
			{CLIENT, false, `["list-capabilities"]`},
			{SERVER, false, `{"capabilities":["relative_root","cmd-subscribe"],"version":"4.9.0"}`},
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
			{CLIENT, false, `["log-level", "error"]`},
			{SERVER, false, `{"log_level":"error","version":"4.9.0"}`},
			{SERVER, true, `{"level":"error","unilateral":true,"log":"2018-03-22T01:18:52,901: [client=0x7ffe1dc035d8:stm=0x7ffe1dc03460:pid=0] test message\n"}`},
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
			{CLIENT, false, `["subscribe", "/tmp", "sub1", {"fields": ["name"]}]`},
			{SERVER, false, `{"version":"4.9.0","clock":"c:1521588867:575:1:2","subscribe":"sub1"}`},
			{SERVER, true, `{"version":"4.9.0","unilateral":true,"subscription":"sub1","clock":"c:1521588867:575:1:2","root":"/tmp","files":["foo"],"is_fresh_instance":true}`},
			{SERVER, true, `{"version":"4.9.0","unilateral":true,"subscription":"sub1","clock":"c:1521588867:575:1:3","since":"c:1521588867:575:1:2","root":"/tmp","files":["bar"],"is_fresh_instance":false}`},
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
			{CLIENT, false, `["watch-project", "/tmp"]`},
			{SERVER, false, `{"version":"4.9.0","watcher":"fsevents","watch":"/tmp"}`},
			{CLIENT, false, `["clock", "/tmp"]`},
			{SERVER, false, `{"version":"4.9.0","clock":"c:1521588867:575:1:5"}`},
			{CLIENT, false, `["subscribe", "/tmp", "sub2", {"since":"c:1521588867:575:1:5", "fields":["name"]}]`},
			{SERVER, false, `{"version":"4.9.0","clock":"c:1521588867:575:1:5","subscribe":"sub2"}`},
			{SERVER, true, `{"version":"4.9.0","unilateral":true,"subscription":"sub2","clock":"c:1521588867:575:1:6","since":"c:1521588867:575:1:5","root":"/tmp","files":["baz"],"is_fresh_instance":false}`},
			{SERVER, true, `{"version":"4.9.0","unilateral":true,"subscription":"sub2","clock":"c:1521588867:575:1:7","since":"c:1521588867:575:1:6","root":"/tmp","files":["qux"],"is_fresh_instance":false}`},
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

func serverFromScript(ctx context.Context, t *testing.T, script []step) (s *server) {
	/* SHUTDOWN
	s.commands: closed by caller
	s.events:   closed locally
	*/

	commands := make(chan string)
	events := make(chan []byte)
	s = &server{
		commands: commands,
		events:   events,
	}

	go func() {
		defer func() {
			close(events)
			for _ = range <-commands {
				continue
			}
		}()
		for i, step := range script {
			switch step.src {
			case CLIENT:
				select {
				case actual, ok := <-commands:
					if !ok {
						t.Errorf("step %d commands channel closed early", i)
					} else if step.data != actual {
						t.Errorf("step %d expected: %#v actual: %#v", i, step.data, actual)
					}
				case <-ctx.Done():
					return
				}
			case SERVER:
				select {
				case events <- []byte(step.data):
				case <-ctx.Done():
					return
				}
			}
		}
	}()

	return
}

func TestEventLoop(t *testing.T) {
	/* SHUTDOWN
	l.commands:      closed by stop()
	l.results:       closed by *eventloop
	l.unilaterals:   closed by *eventloop
	server.commands: closed by *eventloop
	server.events:   closed by *server
	*/

	for label, tc := range testcases {
		t.Run(label, func(t *testing.T) {
			assert := assert.New(t)

			ctx, cancelFunc := context.WithCancel(context.Background())
			server := serverFromScript(ctx, t, tc.script)
			defer cancelFunc()

			l, stop := loop(server)
			defer func() {
				// delayClose=true so that script steps aren't skipped
				stop(true)
				leaktest.Check(t)()
			}()

			results := 0
			unilaterals := 0
			for _, step := range tc.script {
				switch step.src {

				case CLIENT:
					expected := tc.results[results]
					results += 1
					l.commands <- step.data
					actual := <-l.results
					if !assert.Equal(expected, actual) {
						break
					}

				case SERVER:
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

func TestEventLoopCancellation(t *testing.T) {
	/* SHUTDOWN
	l.commands:      closed by stop()
	l.results:       closed by *eventloop
	l.unilaterals:   closed by *eventloop
	server.commands: closed by *eventloop
	server.events:   closed by *server
	*/

	script := testcases["simple"].script

	for _, label := range []string{"before", "after"} {
		t.Run(label, func(t *testing.T) {
			ctx, cancelFunc := context.WithCancel(context.Background())
			server := serverFromScript(ctx, t, script)

			l, stop := loop(server)
			defer func() {
				// delayClose=true so that shutdown is triggered
				// by cancellation instead of some other event
				stop(true)
				leaktest.Check(t)()
			}()

			if label == "before" {
				cancelFunc()
			} else {
				l.commands <- script[0].data
				cancelFunc()
			}

		})
	}
}

.PHONY:  default  examples  test  test-coverage

default: test

examples:
	@for I in examples/*/main.go; do \
		echo ; \
		echo $$I; \
		pushd `dirname "$$I"` >/dev/null; \
		echo ----------; \
		go run *.go; \
		echo ==========; \
		popd >/dev/null; \
		echo ; \
	done

test:
	mkdir -p dist
	go test -coverpkg ./... -coverprofile=dist/coverage.txt -tags integration ./...
	@echo ========================================
	go vet ./...
	golint -set_exit_status ./ ./protocol/...
	gocyclo -over 17 *.go protocol/
	@echo ========================================
	@git grep TODO  -- '**.go' || true
	@git grep FIXME -- '**.go' || true

test-coverage: test
	go tool cover -html=dist/coverage.txt

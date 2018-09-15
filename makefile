default: test

examples:
	@for I in examples/*/main.go; do \
	  echo ; \
	  echo $$I; \
	  pushd `dirname "$$I"` >/dev/null; \
	  echo ----------; \
	  echo go run *.go; \
	  echo ==========; \
	  popd >/dev/null; \
	  echo ; \
	done

test:
	go test -tags integration ./...
	@echo ========================================
	go vet ./...
	golint -set_exit_status ./ ./protocol/...
	gocyclo -over 17 *.go protocol/
	@echo ========================================
	@git grep TODO  -- '**.go' || true
	@git grep FIXME -- '**.go' || true

.PHONY: default examples test

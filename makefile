default: test

demos:
	@for I in demos/*/main.go; do \
	  echo ; \
	  echo $$I; \
	  pushd `dirname "$$I"` >/dev/null; \
	  echo ----------; \
	  echo '1+2' | go run *.go; \
	  echo ==========; \
	  popd >/dev/null; \
	  echo ; \
	done

test:
	go test -tags integration ./...
	@echo ========================================
	go vet ./...
	golint -set_exit_status ./ ./protocol/...
	gocyclo -over 16 *.go protocol/
	@echo ========================================
	@git grep TODO  -- '**.go' || true
	@git grep FIXME -- '**.go' || true

.PHONY: default demos test

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
	gocyclo -over 15 *.go protocol/
	@echo ========================================
	@git grep TODO  || true
	@git grep FIXME || true

.PHONY: default demos test

.PHONY:  default  coverage  examples  test  test-coverage

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

coverage:
	go tool cover -html=dist/coverage.txt

test:
	scripts/run-all-tests
	@echo ========================================
	@git grep TODO  -- '**.go' || true
	@git grep FIXME -- '**.go' || true

test-coverage: test coverage

include Python.mk

GIT_ROOT := $(shell git rev-parse --show-toplevel)

GO_TAGS = $(if $(PYGOLO_TAGS),-tags "$(PYGOLO_TAGS)")
GO_VERBOSE := $(if $(filter-out 0,$(V)),-v)
GO_FLAGS = $(strip $(GO_TAGS) $(GO_VERBOSE) $(PYGOLO_FLAGS))

PYTEST_VERBOSE := $(if $(filter-out 0,$(V)),-v -s,-q)

EXTENSIONS := test/ext.so
EXAMPLES := $(basename $(wildcard examples/*.go))

all: test

license-check:
	@bash scripts/check-license.sh

spell-check:
ifeq ($(CHECK_ONLY),1)
	test -z "`hunspell -a -d en_GB,en_US -p $(GIT_ROOT)/.hunspell_dict *.md docs/*.md | grep '^&' | tee /dev/stderr`"
else
	hunspell -d en_GB,en_US -p $(GIT_ROOT)/.hunspell_dict *.md docs/*.md
	sort $(GIT_ROOT)/.hunspell_dict > $(GIT_ROOT)/.hunspell_dict.new
	mv $(GIT_ROOT)/.hunspell_dict.new $(GIT_ROOT)/.hunspell_dict
endif

prereq:
	go install golang.org/x/lint/golint@latest
	go install honnef.co/go/tools/cmd/staticcheck@latest

lint: PATH := $(PATH):$(shell go env GOPATH)/bin
lint:
	test -z "$$(gofmt -d . | tee /dev/stderr)"
	test -z "$$(golint . | grep -vf .golint-ignore | tee /dev/stderr)"
	go vet $(GO_FLAGS) -cgocall=false .
	staticcheck $(GO_TAGS) .

test-matrix-check: GO_VER := $(shell go version | cut -d" " -f3)
test-matrix-check: PY_VER := $(shell $(PYTHON) -V | cut -d" " -f2)
test-matrix-check:
	$(PYTHON) scripts/check-test-matrix.py "$(TEST_MATRIX_NAME)" "$(GO_VER)" "$(PY_VER)"

test-embed:
	$(embed-python)
	go test $(GO_FLAGS) ./test

test-extend: $(EXTENSIONS)
	$(PYTHON) -m pytest $(PYTEST_VERBOSE) --benchmark-skip test

test: test-embed test-extend

benchmark-embed:
	$(embed-python)
	go test $(GO_FLAGS) -run='^$$' -bench=. ./test

benchmark-extend: $(EXTENSIONS)
	$(PYTHON) -m pytest $(PYTEST_VERBOSE) --benchmark-only test

benchmark: benchmark-embed benchmark-extend

examples: $(EXAMPLES)
ifneq ($(filter-out 0,$(RUN)),)
	@for E in $(notdir $(sort $^)); do \
		echo "### $$E ###"; \
		(cd examples && ./$$E); \
	done
endif

%: %.go FORCE
	$(embed-python)
	go build $(GO_FLAGS) -o $@ ./$<

%.so: % FORCE
	$(extend-python)
	go build $(GO_FLAGS) -o $@ ./$<

clean:
	rm -rf $(EXAMPLES)
	rm -rf $(EXTENSIONS) $(EXTENSIONS:.so=.h)

mrproper: clean
	go clean -cache -testcache

FORCE:

.PHONY: FORCE test

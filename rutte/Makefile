
SELF := github.com/GeenPeil/stem/rutte
CMDS := $(SELF)/cmd/rutte
LDFLAGS := -X $(SELF)/version.version=$(shell git describe --tags)

## first target is default target
.PHONY: build
build:
	@mkdir -p "dist/bin"
	GOBIN="$(CURDIR)/dist/bin" go install -ldflags '$(LDFLAGS)' $(CMDS)
	dist/bin/rutte -v >/dev/null

.PHONY: all
all: build

.PHONY: clean
clean:
	rm -f dist/bin/rutte

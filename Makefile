TARGET = tapper
SRC = $(shell find . -type f -name '*.go' -not -path "./vendor/*")

VERSION = $(shell git describe --tags --exact-match 2> /dev/null || echo development)
COMMIT = $(shell git rev-parse --short=10 HEAD)
DATE = $(shell date -u +"%Y-%m-%dT%H:%M:%S")

DESTDIR=
prefix=$(HOME)
bindir=$(prefix)/bin
INSTALL=install

BUILDDATA = \
	-X "$(TARGET)/version.Version=$(VERSION)" \
	-X "$(TARGET)/version.BuildDate=$(DATE)" \
	-X "$(TARGET)/version.BuildRef=$(COMMIT)"

LDFLAGS = -ldflags '$(BUILDDATA)'

all: $(TARGET)

$(TARGET): .checked $(SRC) go.sum
	go build $(LDFLAGS) -o $(TARGET)

check: .check-lint .check-test

.check-lint: $(SRC) go.sum
	golangci-lint run | tee $@

.check-test: $(SRC) go.sum
	go test $(shell go list ./... | grep -v cmd/) | tee $@

go.sum: go.mod
	go mod tidy && touch $@

install: $(TARGET)
	$(INSTALL) -m 755 $(TARGET) $(DESTDIR)$(bindir)/$(TARGET)

clean:
	rm -f $(TARGET) .checked

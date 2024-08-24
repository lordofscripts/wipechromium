
# Set Compiler to either of: go, gocolor OR gopretty
MODE=gopretty

# Compilers
GO=go
GOCOLOR=~/go/bin/colorgo
GOPRETTY="$(HOME)/go/bin/gofilter"
ifeq ($(MODE),gocolor)
        GO=$(GOCOLOR)
endif

# Sources
CMD_CLI=./cmd/wiper

# Outputs
BIN=./bin
BIN_OUT_CLI=$(BIN)/wipechromium


all: wiper

# ---------------------------------------------------
wiper:
	clear
ifeq ($(MODE),gopretty)
	$(GO) build -v -o $(BIN_OUT_CLI) -tags=none $(CMD_CLI)/*.go 2>&1 | $(GOPRETTY) -color -width 75 -version
else
	$(GO) build -v -o $(BIN_OUT_CLI) -tags=none $(CMD_CLI)/*.go
endif


# ---------------------------------------------------
update:
	go get -u all

testall:
	go test ./...

testfull:
	go test -v test/*_test.go

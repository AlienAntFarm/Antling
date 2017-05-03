
ALL_GO_FILES = $(wildcard *.go)
ALL_BIN_FILES = $(foreach file,$(patsubst %.go,%,$(ALL_GO_FILES)),$(GOBIN)/$(file))

all: $(ALL_BIN_FILES)

$(GOBIN)/%: %.go
	go install $<

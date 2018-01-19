.PHONY: all
all: deps build

.PHONY: build
build:
	go build getsystems.go

.PHONY: deps
deps:
	true

.PHONY: clean
clean:
	- rm *.exe
	- rm getsystems

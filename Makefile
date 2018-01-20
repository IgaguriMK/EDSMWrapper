.PHONY: build
build:
	go build terraformable.go
	go build getstartype.go

.PHONY: deps
deps:
	true

.PHONY: clean
clean:
	- rm *.exe

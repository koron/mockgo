.PHONY: build
build:
	go build

.PHONY: test
test:
	go test ./...

.PHONY: tags
tags:
	gotags -f tags -R .

.PHONY: lint
lint:
	golint ./...

.PHONY: clean
clean:
	go clean
	rm -f tags

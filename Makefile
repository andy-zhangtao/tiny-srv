
.PHONY: build
name = tiny-srv
version = v1.2.1

build:
	go build -ldflags "-X main._VERSION_=$(shell date +%Y%m%d_%H%M%S)" -o $(name)

run: build
	./$(name)

release: *.go *.md
	GOOS=linux GOARCH=amd64 go build -ldflags "-X main._VERSION_=$(shell date +%Y%m%d)" -a -o $(name)
	docker build -t vikings/$(name):$(version) -f Dockerfile.alpine .
	docker push vikings/$(name):$(version)

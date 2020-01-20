build:
	GOOS=linux go build -v ./cmd/techno-db-forum

clean:
	rm -rf ./apiserver

.DEFAULT_GOAL := build

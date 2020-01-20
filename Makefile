build:
	GOOS=linux go build -v -o ./db-forum-kosenkov ./cmd/techno-db-forum

clean:
	rm -rf ./db-forum-kosenkov

.DEFAULT_GOAL := build

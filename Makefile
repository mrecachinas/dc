all: dc

run:
	go run -v ./cmd/dc

dc:
	CGO_ENABLED=0 go build -tags production -o ./bin/dc -ldflags="-s -w" ./cmd/dc

clean:
	rm -f ./bin/*

.PHONY: dc

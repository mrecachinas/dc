all: ui dc

run:
	go run -v ./cmd/dc

dc:
	CGO_ENABLED=0 go build -tags production -o ./bin/dc -ldflags="-s -w" ./cmd/dc

ui:
	npm --prefix ./ui/webapp install ./ui/webapp
	npm --prefix ./ui/webapp run build

clean:
	rm -f ./bin/*
	rm -rf ./ui/webapp/node_modules

.PHONY: dc ui

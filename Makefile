.PHONY: build dist clean

BIN := tt
CMD := ./cmd/tt

build:
	go build -o $(BIN) $(CMD)

dist:
	mkdir -p dist
	GOOS=darwin  GOARCH=arm64 go build -o dist/tt-darwin-arm64  $(CMD)
	GOOS=darwin  GOARCH=amd64 go build -o dist/tt-darwin-amd64  $(CMD)
	GOOS=linux   GOARCH=amd64 go build -o dist/tt-linux-amd64   $(CMD)
	GOOS=linux   GOARCH=arm64 go build -o dist/tt-linux-arm64   $(CMD)

clean:
	rm -f $(BIN)
	rm -rf dist/

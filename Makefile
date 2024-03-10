build: tidy fmt test
	CGO_ENABLED=0 go build -o out/crypto-discord-bot . && chmod +x out/crypto-discord-bot

test:
	go test ./... -cover -v

fmt:
	gofmt -s -w .

tidy:
	go mod tidy -v

all:
	build run

build:
	GOOS=linux GOARCH=amd64 go build -o ./dist/shield_linux_amd64 cmd/shield/main.go
	GOOS=darwin GOARCH=amd64 go build -o ./dist/shield_darwin_amd64 cmd/shield/main.go

run:
	go run cmd/shield/main.go

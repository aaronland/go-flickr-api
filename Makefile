cli:
	go build -mod vendor -o bin/api cmd/api/main.go
	go build -mod vendor -o bin/authorize cmd/authorize/main.go

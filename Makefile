cli:
	go build -mod vendor -o bin/api cmd/api/main.go
	go build -mod vendor -o bin/auth-cli cmd/auth-cli/main.go
	go build -mod vendor -o bin/auth-www cmd/auth-www/main.go

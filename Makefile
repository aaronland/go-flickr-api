GOMOD=$(shell test -f "go.work" && echo "readonly" || echo "vendor")
LDFLAGS=-s -w

cli:
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/api cmd/api/main.go
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/upload cmd/upload/main.go
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/replace cmd/replace/main.go
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/auth-cli cmd/auth-cli/main.go
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/auth-www cmd/auth-www/main.go

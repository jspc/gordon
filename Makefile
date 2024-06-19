BINARY := sample-app/gordon
TYPES := types/page.go \
	request.mint

default: $(TYPES) $(BINARY)

$(BINARY): *.go go.mod go.sum sample-app/*.go $(TYPES)
	CGO_ENABLED=0 go build -ldflags="-s -w" -trimpath -o $@ ./sample-app
	upx $@

$(TYPES): $(wildcard mint/*.mint)
	mint generate mint/

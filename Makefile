NAME = pwg
VERSION = $(shell git describe --tags --abbrev=0)
PREFIX = /usr/local/bin
LDFLAGS =-w -s -X 'main.Version="$(VERSION)"' -X 'main.Name="$(NAME)"'

.PHONY: build clean install uninstall test

build: clean
	@go build -ldflags="$(LDFLAGS)" -o $(NAME) ./cmd/$(NAME)/*.go

clean:
	@$(RM) $(NAME)

install:
	@cp -i $(NAME) $(PREFIX)

uninstall:
	@rm -i $(PREFIX)/$(NAME)

test:
	@go test -v ./...
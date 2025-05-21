build:
	@templ generate components
	@go build -o datastarui.exe main.go

tailwind:
	@tailwindcss -i static/css/index.css -o static/css/build.css --watch

watch:
	air

run: build
	@./datastarui.exe

install:
	@go install github.com/cosmtrek/air@latest
	@go install github.com/a-h/templ/cmd/templ@latest
	@go get ./...
	@go mod tidy
	@go mod download
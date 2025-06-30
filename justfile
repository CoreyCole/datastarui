build:
	@templ generate
	@go build -o datastarui main.go

tailwind:
    @tailwindcss -i static/css/index.css -o static/css/build.css --watch --content "./components/**/*.templ" --content "./pages/**/*.templ" --content "./layouts/**/*.templ"

watch:
	air

install:
	@go install github.com/air-verse/air@latest
	@go install github.com/a-h/templ/cmd/templ@latest
	@go get ./...
	@go mod tidy
	@go mod download

# Docker development commands
up:
	@echo "🚀 Starting DatastarUI development environment..."
	@docker compose up --build

down:
	@echo "🛑 Stopping DatastarUI development environment..."
	@docker compose down

docker-shell service="app":
	@echo "🐚 Opening shell in {{service}} container..."
	@docker compose exec {{service}} sh

docker-tail service lines="20":
	@echo "📋 Last {{lines}} lines for {{service}} service..."
	@docker logs --tail {{lines}} datastarui-{{service}}-1

docker-follow service:
  @docker logs --follow datastarui-{{service}}
  @echo "📋 Following logs for {{service}} service (Ctrl+C to stop)..."

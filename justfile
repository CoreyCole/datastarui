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
# Builds, (re)creates, and starts containers
@up *args:
	just _compose local up '-d --build --remove-orphans' {{args}}
alias u := up

@down *args:
	just _compose local down {{args}}
alias d := down

# Generic Command
_compose environment cmd opts='' *args='':
	GOCACHE=$(go env GOCACHE) BUILDKIT_PROGRESS=plain \
	docker compose \
		--profile {{env_var_or_default("UP_PROFILES", "all")}} \
		-p datastarui-{{environment}} \
		-f $(if [ '{{environment}}' = 'local' ]; then echo 'docker-compose.yml'; else echo "docker-compose.{{environment}}.yml"; fi) \
		{{cmd}} {{opts}} {{args}}

# Docker shell access
docker-shell service="app":
	@echo "🐚 Opening shell in {{service}} container..."
	just _compose local exec '' {{service}} sh

# Docker logs viewing
docker-tail service lines="50":
	@echo "📋 Last {{lines}} lines for {{service}} service..."
	@docker logs --tail {{lines}} datastarui-local-{{service}}-1


build:
  @templ generate
  @go build -o datastarui main.go
  @just build-tailwind

build-tailwind:
    @echo "🎨 Building Tailwind CSS..."
    @pnpm tailwindcss -i static/css/index.css -o static/css/out.css --content "./components/**/*" --content "./pages/**/*" --content "./layouts/**/*"
    @if [ -f static/css/out.css ]; then \
        echo "📝 Generating CSS hash..."; \
        HASH=$(sha256sum static/css/out.css | cut -d' ' -f1 | head -c8); \
        echo "🔖 Hash: $HASH"; \
        rm -f static/css/out.*.css; \
        cp static/css/out.css static/css/out.$HASH.css; \
        echo "✅ Created static/css/out.$HASH.css"; \
    fi

tailwind:
  @tailwindcss -i static/css/index.css -o static/css/out.css --watch --content "./components/**/*" --content "./pages/**/*" --content "./layouts/**/*"

watch:
  air

install:
  @pnpm install
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


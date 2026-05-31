# Development stage - only sets up tools, no file copying
FROM golang:1.25-alpine AS development

# Install dependencies
RUN apk update && apk add --no-cache git curl nodejs npm bash coreutils && \
    rm -rf /var/cache/apk/*

# Install templ version pinned by go.mod-generated code compatibility.
RUN go install github.com/a-h/templ/cmd/templ@v0.3.977

# Install buf for proto generation
RUN go install github.com/bufbuild/buf/cmd/buf@latest

# Install air for live reload
RUN go install github.com/air-verse/air@latest

# Install just command runner
RUN curl --proto '=https' --tlsv1.2 -sSf https://just.systems/install.sh -o /tmp/install-just.sh && \
    bash /tmp/install-just.sh --to /usr/local/bin && \
    rm /tmp/install-just.sh

# Install pnpm for non-interactive Docker build steps.
RUN npm install -g pnpm@11.5.0

# Set working directory
WORKDIR /app

# Copy package files and install dependencies
COPY package.json pnpm-lock.yaml ./
RUN pnpm install --dangerously-allow-all-builds && pnpm store prune

# Expose port
EXPOSE 4242


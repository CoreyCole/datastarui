# Development stage - only sets up tools, no file copying
FROM golang:1.24-alpine AS development

# Install dependencies
RUN apk update && apk add --no-cache git curl nodejs npm bash coreutils && \
    rm -rf /var/cache/apk/*

# Install templ
RUN go install github.com/a-h/templ/cmd/templ@latest

# Install air for live reload
RUN go install github.com/air-verse/air@latest

# Install just command runner
RUN curl --proto '=https' --tlsv1.2 -sSf https://just.systems/install.sh -o /tmp/install-just.sh && \
    bash /tmp/install-just.sh --to /usr/local/bin && \
    rm /tmp/install-just.sh

# Install pnpm
ENV SHELL=/bin/bash
RUN curl -fsSL https://get.pnpm.io/install.sh | bash -
ENV PATH="/root/.local/share/pnpm:$PATH"

# Set working directory
WORKDIR /app

# Copy package files and install dependencies
COPY package.json pnpm-lock.yaml ./
RUN pnpm install && pnpm store prune

# Expose port
EXPOSE 4242


# Development stage - only sets up tools, no file copying
FROM golang:1.24-alpine AS development

# Install dependencies
RUN apk add --no-cache git curl nodejs npm bash

# Install templ
RUN go install github.com/a-h/templ/cmd/templ@latest

# Install air for live reload
RUN go install github.com/air-verse/air@latest

# Install just command runner
RUN curl --proto '=https' --tlsv1.2 -sSf https://just.systems/install.sh -o /tmp/install-just.sh && \
    bash /tmp/install-just.sh --to /usr/local/bin && \
    rm /tmp/install-just.sh

# Set working directory
WORKDIR /app

# Expose port
EXPOSE 4242


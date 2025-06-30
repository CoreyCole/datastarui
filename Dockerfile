# Development stage - only sets up tools, no file copying
FROM golang:1.24-alpine AS development

# Install dependencies
RUN apk add --no-cache git curl nodejs npm

# Install templ
RUN go install github.com/a-h/templ/cmd/templ@latest

# Install air for live reload
RUN go install github.com/air-verse/air@latest

# Install just command runner
RUN curl --proto '=https' --tlsv1.2 -sSf https://just.systems/install.sh | sh -s -- --to /usr/local/bin

# Install glibc compatibility for Alpine and Tailwind CSS standalone CLI
RUN apk add --no-cache gcompat && \
    curl -L https://github.com/tailwindlabs/tailwindcss/releases/download/v4.1.11/tailwindcss-linux-x64 -o /usr/local/bin/tailwindcss && \
    chmod +x /usr/local/bin/tailwindcss

# Set working directory
WORKDIR /app

# Expose port
EXPOSE 4242

# Default command for development
CMD ["air"]

# Production stage
FROM alpine:latest AS production

RUN apk --no-cache add ca-certificates
WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/datastarui .
COPY --from=builder /app/static ./static

EXPOSE 4242

CMD ["./datastarui"]
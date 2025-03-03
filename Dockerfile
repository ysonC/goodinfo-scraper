################################################
# 1) Build the Go scraper
################################################
FROM golang:1.24 AS builder
WORKDIR /app

# Copy mod files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build a static Go binary named "scraper"
RUN CGO_ENABLED=0 GOOS=linux go build -o /scraper .

################################################
# 2) Install the Playwright-Go driver
################################################
FROM golang:1.24 AS driver-installer
# playwright-go provides a helper CLI under "cmd/playwright"
RUN go install github.com/playwright-community/playwright-go/cmd/playwright@latest

# Pre-install the matching version of the driver
# (Installs it in /root/.cache/ms-playwright-go by default)
RUN /go/bin/playwright install --with-deps

################################################
# 3) Final image: Use official MS Playwright runtime
################################################
FROM mcr.microsoft.com/playwright:v1.50.1-jammy
WORKDIR /app

# Copy the Go binary from the builder stage
COPY --from=builder /scraper /app/scraper

# Copy the "playwright" driver CLI
COPY --from=driver-installer /go/bin/playwright /usr/local/bin/playwright

# Also copy the driver caches that were installed
COPY --from=driver-installer /root/.cache/ms-playwright-go /root/.cache/ms-playwright-go

# Default entrypoint: run the Go scraper
ENTRYPOINT ["/app/scraper"]


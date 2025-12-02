################################################
# 1) Builder: build your Go binary + fetch CLI
################################################
FROM golang:1.24 AS builder
WORKDIR /src

# Cache deps
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest
COPY . .

# Build the scraper (adjust path if different)
# → if your main is in cmd/scraper, build that path
RUN CGO_ENABLED=0 GOOS=linux \
    go build -trimpath -ldflags="-s -w" -o /scraper ./cmd/scraper

# Also build the playwright-go CLI here so final image stays Go-free
RUN go install github.com/playwright-community/playwright-go/cmd/playwright@v0.5001.0

################################################
# 2) Final: Playwright runtime with preinstalled browsers
################################################
FROM mcr.microsoft.com/playwright:v1.50.1-jammy
WORKDIR /app

# Use browsers pre-shipped in this image; skip downloading more
ENV PLAYWRIGHT_BROWSERS_PATH=/ms-playwright \
    PLAYWRIGHT_SKIP_BROWSER_DOWNLOAD=1

# Copy binary and CLI
COPY --from=builder /scraper /app/scraper
COPY --from=builder /go/bin/playwright /usr/local/bin/playwright


# Prepare mount points and permissions for pwuser (provided by base image)
RUN mkdir -p /app/data/downloaded_stock /app/data/final_output /app/data/input_stock /app/data/failed_stock \
 && chown -R pwuser:pwuser /app

# Preload inputstock
COPY resources/stocks.txt /app/data/input_stock/stocks.txt

# Install the Playwright-Go driver (tiny) against preinstalled browsers
USER pwuser
ENV HOME=/home/pwuser
RUN playwright install

ENTRYPOINT ["/app/scraper"]

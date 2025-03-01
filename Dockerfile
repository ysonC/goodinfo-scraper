# ---- Stage 1: Build the binary ----
FROM golang:1.24.0 AS builder
WORKDIR /app

# Copy go.mod and go.sum, then download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the full source code into the container
COPY . .

# Build the binary. We disable CGO for a static binary.
RUN CGO_ENABLED=0 GOOS=linux go build -o scraper .

# ---- Final Stage: Runtime image ----
FROM node:20-bookworm
ENV DEBIAN_FRONTEND=noninteractive

# Install dependencies required by Chromium/Playwright
RUN apt-get update && apt-get install -y \
    curl \
    wget \
    ca-certificates \
 && rm -rf /var/lib/apt/lists/*

WORKDIR /app
RUN mkdir -p input_stock output_stock

# Copy the built binary from the builder stage
COPY --from=builder /app/scraper .

# Set the browser installation path so the driver is installed in a known location
ENV PLAYWRIGHT_BROWSERS_PATH=/app/.playwright-browsers

# Install Playwright along with its browser binaries/driver
RUN npx -y playwright@1.50.1 install --with-deps

# (Optional) Add the directory to PATH if your app expects to find the executable there
ENV PATH=/app/.playwright-browsers:$PATH

# Set the entrypoint to run your scraper
ENTRYPOINT ["./scraper"]


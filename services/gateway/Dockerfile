FROM rust:1.90-alpine AS builder

# Install required build tools and static OpenSSL
RUN apk add --no-cache musl-dev pkgconfig openssl-dev openssl-libs-static

# Create app directory
WORKDIR /app

# Copy Cargo manifests first
COPY Cargo.toml Cargo.lock ./

# Make stub main.rs for caching deps
RUN mkdir src && echo "fn main() {}" > src/main.rs

# Fetch the dependencies
RUN cargo fetch

# Copy the full source
COPY src ./src

# Build for release
RUN cargo build --release

FROM alpine:3.20

# Install CA certificates
RUN apk add --no-cache ca-certificates

# Copy the binary from builder
COPY --from=builder /app/target/release/gateway /usr/local/bin/gateway

# Set environment variables
ENV PORT=1000
EXPOSE ${PORT}

# Run the API
CMD ["/usr/local/bin/gateway"]
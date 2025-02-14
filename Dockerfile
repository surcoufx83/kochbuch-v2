# Build Rust backend
FROM rust:latest AS backend-builder
WORKDIR /app
RUN apt update && apt install -y musl-tools
RUN rustup target add x86_64-unknown-linux-musl
COPY backend/ .
RUN cargo build --release --target=x86_64-unknown-linux-musl
RUN strip /app/target/x86_64-unknown-linux-musl/release/kochbuch

# Build Angular frontend
FROM node:latest AS frontend-builder
WORKDIR /app
COPY frontend/ .
RUN npm install && npm run build -- --output-path=/dist

# Serve with nginx
FROM nginx:alpine
COPY --from=frontend-builder /dist/ /usr/share/nginx/html
COPY --from=backend-builder /app/target/x86_64-unknown-linux-musl/release/kochbuch /usr/local/bin/kochbuch
CMD ["sh", "-c", "/usr/local/bin/kochbuch & nginx -g 'daemon off;'"]

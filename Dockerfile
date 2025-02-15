# Build Go backend
FROM golang:latest AS backend-builder
WORKDIR /app
COPY backend/ .
RUN go mod tidy && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o kochbuch

# Build Angular frontend
FROM node:latest AS frontend-builder
WORKDIR /app
COPY frontend/ .
RUN npm install && npm run build -- --output-path=/dist

# Serve frontend and proxy backend
FROM nginx:alpine
COPY --from=frontend-builder /dist/browser/ /usr/share/nginx/kochbuch/html
COPY --from=backend-builder /app/kochbuch /usr/local/bin/kochbuch
COPY nginx.conf /etc/nginx/conf.d/default.conf
CMD ["sh", "-c", "/usr/local/bin/kochbuch & nginx -g 'daemon off;'"]

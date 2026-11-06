# syntax=docker/dockerfile:1

FROM golang:1.22-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# Fail the image build early if any content file is invalid.
RUN go test ./internal/content
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/web ./cmd/web

FROM gcr.io/distroless/static-debian12:nonroot
WORKDIR /app
# Templates and static assets are embedded in the binary; content stays editable.
COPY --from=build /out/web /app/web
COPY content /app/content
ENV PORT=8080 ENVIRONMENT=production CONTENT_DIR=/app/content
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/app/web"]

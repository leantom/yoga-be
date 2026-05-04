FROM golang:1.22-bookworm AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -ldflags="-s -w" -o /out/yoga-api ./cmd/api

FROM gcr.io/distroless/static-debian12:nonroot

ENV FIRESTORE_DATABASE=(default)
ENV PORT=8080

WORKDIR /app
COPY --from=builder /out/yoga-api ./yoga-api

USER nonroot:nonroot
EXPOSE 8080

ENTRYPOINT ["/app/yoga-api"]

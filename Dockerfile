FROM golang:1.22-bookworm AS build

WORKDIR /src
COPY go.mod go.sum* ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o /bin/yoga-api ./cmd/api

FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=build /bin/yoga-api /yoga-api
USER nonroot:nonroot
EXPOSE 8080
ENTRYPOINT ["/yoga-api"]

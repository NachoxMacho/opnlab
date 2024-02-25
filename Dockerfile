FROM golang:1.21.6 as build-stage

WORKDIR /app

# Dependencies
COPY go.mod go.sum ./
RUN go mod download

# TODO: just copy .go files while keeping structure
COPY . ./

# Disable CGO so we can run without glibc
RUN CGO_ENABLED=0 GOOS=linux go build -o /docker-go

FROM gcr.io/distroless/base-debian11  AS release-stage

WORKDIR /app

COPY --from=build-stage /docker-go /app/docker-go

COPY ./views /app/views
COPY ./css /app/css
COPY ./migrations /app/migrations

EXPOSE 42069

USER nonroot:nonroot

ENTRYPOINT ["/app/docker-go"]

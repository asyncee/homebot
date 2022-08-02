##
## Build
##
FROM golang:1.18-buster AS build

WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./
RUN go build -o /homebot ./cmd/homebot/main.go

##
## Deploy
##
FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /homebot /homebot
USER nonroot:nonroot
ENTRYPOINT ["/homebot"]


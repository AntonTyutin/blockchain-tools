FROM golang:1.23 AS build

WORKDIR /build

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

ENV GOCACHE=/root/.cache/go-build
RUN --mount=type=cache,target="/root/.cache/go-build" go build -o bin/blockchain-tools cmd/blockchain-tools/main.go

FROM debian

COPY --from=build /build/bin/blockchain-tools /blockchain-tools

ENTRYPOINT ["/blockchain-tools"]
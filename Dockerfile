FROM golang:1.16-alpine AS build

WORKDIR /go/src/app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o mimic ./cmd

FROM golang:1.16-alpine

COPY --from=build /go/src/app/mimic /mimic

ENTRYPOINT ["/mimic"]

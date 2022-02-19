FROM golang:1.17 AS build

WORKDIR /go/src/github.com/dtan4/grpc-private-bff-example

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /server github.com/dtan4/grpc-private-bff-example/cmd/server

FROM gcr.io/distroless/static:nonroot

COPY --from=build /server /server

EXPOSE 8080

CMD ["/server"]

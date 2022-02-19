FROM golang:1.17 AS build

WORKDIR /go/src/github.com/dtan4/grpc-private-bff-example

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /bff github.com/dtan4/grpc-private-bff-example/cmd/bff

# nonroot image doesn't work well with volume mount
FROM gcr.io/distroless/static

COPY --from=build /bff /bff

EXPOSE 8080

CMD ["/bff"]

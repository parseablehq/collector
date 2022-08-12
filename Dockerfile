FROM golang:1.17 as builder

WORKDIR /workspace

COPY . .

RUN go mod tidy
RUN go fmt ./...
RUN go vet ./...
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o collector main.go

FROM golang:1.17-alpine
WORKDIR /
COPY --from=builder /workspace/collector .
ENTRYPOINT ["/collector"]

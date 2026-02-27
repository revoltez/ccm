FROM golang:1.25-alpine AS builder

WORKDIR /src

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o ccm ./cmd/ccm/

FROM scratch

COPY --from=builder /src/ccm /ccm

USER 65534:65534

ENTRYPOINT ["/ccm"]

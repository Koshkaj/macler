FROM golang:1.19.5-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /go/bin/app .

FROM gcr.io/distroless/static-debian11 as prod

COPY --from=builder /go/bin/app /

USER 1000:1000

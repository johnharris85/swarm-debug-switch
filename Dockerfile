FROM golang:alpine as builder

WORKDIR /app 

COPY . /app/

RUN go build -o node-debug-switch

FROM alpine:3.6

WORKDIR /app

COPY --from=builder /app/node-debug-switch .

CMD ["./node-debug-switch"]

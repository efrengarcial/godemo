FROM golang:1.11.1-stretch as builder

WORKDIR /godemo

COPY . ./

# Building using -mod=vendor, which will utilize the v
RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor  -o app ./users/cmd

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /godemo/app .

CMD ["./app"]
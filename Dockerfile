FROM golang:alpine as builder
ADD . /go-aws
WORKDIR /go-aws
RUN go get -t -v ./...
RUN go build -o app;

FROM alpine:latest
WORKDIR /root/
# can be handled from secrets as volume mount
# COPY ./credentials ./.aws/credentials
COPY --from=builder /go-aws/app .
EXPOSE 9090
#RUN chmod -r 777
CMD ["./app"]
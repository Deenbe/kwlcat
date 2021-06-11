FROM golang:1.16.3-alpine3.13

RUN mkdir -p /home/ots
WORKDIR /home/ots
COPY go.mod /home/ots/
COPY go.sum /home/ots/
COPY main.go /home/ots/
COPY awsap/ /home/ots/awsap/
COPY cmd/ /home/ots/cmd/
COPY global/ /home/ots/global/

RUN ls /home/ots/awsap/dynamodb/
ENV GOPROXY=https://goproxy.io,direct

RUN go get
RUN go build -o build/ots .


FROM alpine:3.13
COPY --from=0 /home/ots/build/ots /usr/local/bin/

ENTRYPOINT ["ots"]

FROM golang:1.16.3-alpine3.13

RUN mkdir -p /home/gopi

WORKDIR /home/gopi
COPY go.mod /home/gopi/
COPY go.sum /home/gopi/
COPY main.go /home/gopi/
COPY awsap/ /home/gopi/awsap/
COPY cmd/ /home/gopi/cmd/
COPY global/ /home/gopi/global/

RUN ls /home/gopi/awsap/dynamodb/
ENV GOPROXY=https://goproxy.io,direct

RUN go get
RUN go build -o build/gopi .


FROM alpine:3.13
COPY --from=0 /home/gopi/build/gopi /usr/local/bin/
RUN mkdir -p /etc/otel
COPY otel-collector-config.yaml /etc/otel/config.yaml
VOLUME /etc/otel

ENTRYPOINT ["gopi"]

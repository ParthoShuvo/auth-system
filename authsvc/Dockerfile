FROM golang:1.17-alpine3.15

WORKDIR /home/app
COPY . /home/app

RUN go get .
RUN go build

CMD ./authsvc $SVC_CONFIG_PATH
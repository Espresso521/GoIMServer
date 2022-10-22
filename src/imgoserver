

FROM golang:1.19.2-alpine

WORKDIR /Users/huze/Desktop/gotest/GoIMServer/src

ADD . .

RUN go build -o imserver

CMD [ "./imserver" ]
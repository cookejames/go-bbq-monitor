FROM golang:1.16.2

COPY ./src/ /go/src/
COPY ./certs /go/certs/

WORKDIR /go/src/

ENV LOGXI=main=INF
# RUN go get -d -v ./...
RUN go install -v ./...

CMD ["go-bbq-monitor"]
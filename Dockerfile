FROM golang:alpine

RUN mkdir -p /go/src/instatasks

WORKDIR /go/src/instatasks

COPY . /go/src/instatasks

RUN go install instatasks

CMD /go/bin/instatasks

EXPOSE 8080
FROM golang:alpine

RUN mkdir -p /go/src/instatasks

WORKDIR /go/src/instatasks

COPY . /go/src/instatasks

RUN CGO_ENABLED=0 go install instatasks

CMD /go/bin/instatasks


FROM golang:alpine

RUN apk update;apk add git
RUN apk add --no-cache gcc musl-dev
RUN go get -u github.com/golang/dep/cmd/dep

CMD ["/bin/sh"]
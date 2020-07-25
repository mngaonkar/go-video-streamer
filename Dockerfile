FROM golang:alpine

WORKDIR /dist

ADD public ./public
ADD server ./server
COPY main.go .
COPY go.mod .
COPY go.sum .

RUN go mod download
RUN go build main.go

EXPOSE 8080

CMD ["/dist/main"]





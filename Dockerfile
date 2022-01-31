FROM golang:alpine
WORKDIR /go-docker
ADD . .
COPY go.mod go.sum ./
RUN go mod download
RUN go build .
RUN go build main.go 
EXPOSE 8088
ENTRYPOINT ["./go-docker", "--port=8088"]

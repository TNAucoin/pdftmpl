FROM golang:alpine

WORKDIR /myapp
COPY . .
RUN go mod tidy
RUN go build -o ./bin/api ./cmd/api/

CMD ["/myapp/bin/api"]
EXPOSE 4000
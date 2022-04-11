FROM golang:1.18.0-alpine3.15
WORKDIR /md2htmlapiv
COPY . .
RUN go build -ldflags="-w -s" .
RUN rm -rf *.go && rm -rf go.*
CMD ["./md2htmlapi"]

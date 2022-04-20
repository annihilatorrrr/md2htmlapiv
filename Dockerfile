FROM golang:1.18.1-alpine3.15 as builder
WORKDIR /md2htmlapiv
COPY . .
RUN go build -ldflags="-w -s" .
RUN rm -rf *.go && rm -rf go.*
RUN ls
FROM alpine:latest
COPY --from=builder /md2htmlapiv/md2htmlapiv /md2htmlapiv
ENTRYPOINT ["/md2htmlapiv"]

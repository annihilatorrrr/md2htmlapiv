FROM golang:1.18.1-alpine3.15 as builder
WORKDIR /md2htmlapiv
COPY . .
RUN go build -ldflags="-w -s" .
RUN rm -rf *.go && rm -rf go.*
FROM alpine:3.15.4
COPY --from=builder /md2htmlapiv/md2htmlapi /md2htmlapiv
ENTRYPOINT ["/md2htmlapiv"]

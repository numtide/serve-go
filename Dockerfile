# Build
FROM golang:1.17-alpine as build
WORKDIR /build
COPY --chown=0:0 go.* *.go /build/
COPY --chown=0:0 ./spa/ /build/spa/
RUN go build -o serve-go .

# Deploy
FROM alpine:latest
COPY --from=build /build/serve-go /bin/serve-go
RUN mkdir /site
EXPOSE 3000
CMD ["serve-go", "/site"]

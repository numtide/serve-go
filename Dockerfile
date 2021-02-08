# Build
FROM golang:1.15-alpine as build
WORKDIR /build
COPY go.mod *.go /build/
RUN go build -o serve-go .

# Deploy
FROM alpine:latest
COPY --from=build /build/serve-go /bin/serve-go
RUN mkdir /site
EXPOSE 3000
CMD ["serve-go", "/site"]

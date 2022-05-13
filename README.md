# serve-go - HTTP web server for SPA

A mostly drop-in replacement to https://www.npmjs.com/package/serve but
written in Go. For when it's time to go to production.

It's dead simple. Not configurable (happy to get PRs).

It serves all the static files from the current directory. And falls backs to
serving the `index.html` if the file is not found. Perfect for Single-Page
Applications (SPA).

## Usage

```
dead-simple application that serves static files from the current directory
Usage: serve-go [options] [<work-dir>]

Options:
  -listen: Port to listen to (default 3000)
  -oembed-url: Sets the oEmbed Link header if set (env: $SERVEGO_OEMBED_URL) (default )
  <work-dir>: Folder to serve (default to current directory)
```

## Content Encoding

The server now also handles gzip and brotli content-encoding.

For each file that exists in the work directory, it will also look for a .br
or .gz file, and if it exists, and the client accepts the encoding, serve that
file instead of the original one.

## Example

Here is how to integrate it in a Docker image:

Dockerfile:
```dockerfile
FROM node:15.7.0-alpine3.12 as builder
WORKDIR /app/
# Install dependencies
COPY package.json package-lock.json /app/
RUN npm install
# Build the frontend
COPY . /app/
RUN npm run build

# Create a serve container
FROM ghcr.io/numtide/serve-go:v1.3.0
COPY --from=builder /app/build /site
```


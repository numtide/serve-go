# serve-go - HTTP web server for SPA

A mostly drop-in replacement to https://www.npmjs.com/package/serve but
written in Go. For when it's time to go to production.

It's dead simple. Not configurable (happy to get PRs).

It serves all the static files from the current directory. And falls backs to
serving the `index.html` if the file is not found. Perfect for Single-Page
Applications (SPA).

## Usage

```
Usage: serve [options] [<work-dir>]

Options:
* -listen: Port to listen to (default 3000)
* <work-dir>: Folder to serve (default to current directory)
```

## Example

Here is how to integrate it in a Docker image:

`Dockerfile`
```
FROM node:15.7.0-alpine3.12 as builder
WORKDIR /app/
# Install dependencies
COPY package.json package-lock.json /app/
RUN npm install
# Build the frontend
COPY . /app/
RUN npm run build

# Create a serve container
FROM numtide/serve-go:1.0.0
COPY --from=builder /app/build /site
```


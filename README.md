<div  align="center">
  <img src="./assets/echo-logo.jpg" alt="Echo Logo" width="200">
  <h2>HTTP Echo</h2>
  <p>A simple HTTP server that echoes back request information.</p>

  ![license](https://img.shields.io/badge/license-MIT-green.svg)
  ![go](https://img.shields.io/badge/golang-1.23-blue.svg?logo=go)
  [![gh-stars](https://img.shields.io/github/stars/rellyson/http-echo?logo=github)](https://github.com/rellyson/http-echo/stargazers)
  [![docker-pulls](https://img.shields.io/docker/pulls/rellyson/http-echo.svg?logo=docker)](https://hub.docker.com/repository/docker/rellyson/http-echo)
</div>

This project is a simple HTTP server that echoes back request information. It is useful
for testing and debugging deployments and networking configurations by providing a simple
interface to inspect requests. 

## Features

The HTTP Echo server provides the following features:

- Accepts HTTP methods `GET`, `POST`, `PUT`, `PATCH` and `DELETE`.
- Echoes back request information in `json` format:
  - **host**: `hostname` and `ip`
  - **headers**: request headers
  - **queries**: query parameters
  - **params**: path parameters
  - **body**: request body
- Returns a custom status code for a request (provided via query) - **TBD**.
- Exports metrics in `prometheus` format for monitoring.

## Usage

The server accepts the following parameters:

- `-listen`: the address to listen on (default: `:3000`)
- `-metrics`: the path to expose metrics on (default: `/metrics`)

## Running the server

You can run the server using via the following methods:

### Local using `go`

Yuu can run the server locally using the following command:

```bash
# execute via go run
$ go run cmd/server/main.go

# or build and execute
$ go build -o htt-echo cmd/server/main.go
$ ./http-echo
```

### Docker

There is also a [docker image](https://hub.docker.com/repository/docker/rellyson/http-echo)
available on _Docker Hub_.

```bash
# run using docker
$ docker run -p 3000:3000 rellyson/http-echo:latest
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

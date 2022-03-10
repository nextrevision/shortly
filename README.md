shortly
=======

`shortly` is a PoC url shortner written in go.

## Getting Started

Requirements:

- docker
- docker-compose

Start the server and local dependencies:

```shell
make local
```

Issue a request against the endpoint:

```shell
make shorten-url
```

Take the output of that command and browse to the location:

```shell
curl -vvv -L localhost:8000/[id]
```

## Endpoints

| Path         | Method | Description                                                                                                                                                                       |
|--------------|--------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `/`          | `POST` | Shortens a supplied URL and returns a plain text "shortened" URL. Accepts `application/x-www-form-urlencoded` content with the `url` form key set to the URL you wish to shorten. |
| `/{id}`      | `GET`  | Resolves a shortened URL to it's intended destination and responds with a redirect to the location.                                                                               |
| `/echo/{id}` | `GET`  | Used to test the shorten and redirect workflows. It merely echos the `id` supplied.                                                                                               |
| `/metrics`    | `GET`   | Prometheus metrics scraping endpoint                                                                                                                                              |


## Test and Build

Tests can be run with:

```shell
make test
```

and `shortly` can be built with:

```shell
make build
```

## Deployment

This service can be deployed on AWS via Terraform. All terraform files are in the `terraform` directory and can be run with:

```shell
make tf-plan
make tf-apply
```
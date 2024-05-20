# Micro-Services Template

This template is intended to be used by teams or companies, to create enterprise-scale microservice architectures quickly with a reasonable set of defaults. The architecture is built on top of [Kratos](https://go-kratos.dev/en/), a framework for rolling out microservices. A growing team will quickly find need to convert sections of the app into submodules for organization and repo size management - this is expected and the template code is structured to not lose organization as systems grow.

This template uses gRPC as well as REST. They are configured together using protofiles. Both gRPC and HTTP servers are run concurrently for each service. Services are modular are isolated - they can be configured to run independently on many machines, horizontally scaled, auto-scaled, or otherwise deployed as needed.

The main (core) service can be run using `make execute`. This service, by default, handles user operations and some financial services. End users should choose what services they will need.

Google's [Wire](https://github.com/google/wire/tree/main) tool is used for compile-time dependency injection. See Wire's [User Guide](https://github.com/google/wire/blob/main/docs/guide.md), [FAQ](https://github.com/google/wire/blob/main/docs/faq.md), and [Best Practices](https://github.com/google/wire/blob/main/docs/best-practices.md) document for further reference on how Wire works.

Caching and Kafka streaming are provided via [Upstash](https://upstash.com/). Instrumentation is provided using [Jaegar](https://www.jaegertracing.io/), and the [Prometheus](https://prometheus.io/) stack.

## Setup

You will need a `.env` file. Ask a team member for a copy.

## Working with the code

The workflow to create new services is roughly as follows:

- Define a proto service: 
  - If defining an internal service: `kratos proto add api/v1/my_service.proto`
  - For a public-facing service: `kratos proto add api/v1/my_namespace/my_service.proto`
- Convert the protofile to Go code: `kratos proto client api/v1/my_namespace/my_service.proto`
- Define a service
  - If defining an internal service: `kratos proto server api/v1/my_service.proto -t internal/service`
  - For a public-facing service: `kratos proto server api/v1/my_namespace/my_service.proto -t pkg/my_service/service`
  - Follow the structure of packages such as the `lodging` and `finance` packages for public packages
  - For internal services, follow the example of the Users service

Further documentation is ongoing. Check various directories for documentation of specific services and internal tooling.

## Install Kratos

```
go install github.com/go-kratos/kratos/cmd/kratos/v2@latest
```

## Create a service

```shell
# Add a proto template
kratos proto add api/v1/foobar.proto
# Generate the proto code
kratos proto client api/v1/foobar.proto
# Generate the source code of service by proto file
kratos proto server api/v1/foobar.proto -t internal/service

make api # Build protofiles
make build # Build Go code
```

Make sure to update the Makefile, adding a build step for the new service.

Add a build step for the `foobar` service:
```dart
.PHONY: foobar
fin:
	./bin/foobar &
```

## Generate auxiliary files

```shell
# Generate API files (include: pb.go, http, grpc, validate, swagger) by proto file
make api
# Build all files
make compile
```

## Automated Initialization (wire)
```shell
# install wire
go get github.com/google/wire/cmd/wire

# generate wire
cd cmd/server
wire
```

Make sure to add a wire build step to your make file when a service is ready to be deployed.

```go
build:
	mkdir -p bin/ && go build -ldflags "-X main.Version=$(VERSION)" -o ./bin/ ./...
	cd app/finance && wire
	cd app/b2b && wire
	cd app/lodging && wire
  cd app/your_pkg && wire
```

This will break builds if the service isn't defined.

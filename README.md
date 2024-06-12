# Micro-Services Template

This template allows for the to quick generation of enterprise-scale microservice architectures with a reasonable set of defaults. 

## Features

### Service Discovery and Configuration:

- etcd: Distributed key-value store for service discovery and configuration management.

### Framework and Libraries:

- Kratos: Go microservices framework with built-in support for gRPC and other features.
- gRPC: High-performance RPC framework for inter-service communication.
- gorm: ORM library for database interactions.
- JWT: JSON Web Tokens for authentication and authorization.


### Databases:

- Redis: In-memory data store for caching and quick data access.
- CockroachDB: Distributed SQL database for persistent data storage.


### Messaging and Streaming:

- Kafka: Distributed streaming platform for event-driven communication and data pipelines.


### Observability and Monitoring:

- Prometheus: Monitoring and alerting system for collecting metrics.
- Jaeger: Distributed tracing system for monitoring and troubleshooting microservices.
- OpenTelemetry: Observability framework for distributed tracing, metrics, and logging.


### API Analytics and Monitoring:

- Moesif: API analytics and monitoring platform for insights into API usage and customer behavior.


### Potential Additions (TBA):

- Sentry or Datadog: Error tracking and monitoring platforms for identifying and resolving issues.
- Docker: Containerization platform for packaging and deploying microservices.
- Kubernetes: Container orchestration system for automating deployment, scaling, and management of containerized applications.

## Internal Dependencies

- Go v. 1.22.2
- Docker
- Etcd
- Kubernetes

## External Dependencies

- CockroachDB
- InfluxDB (via InfluxDB Cloud)
- Redis (via Upstash)
- Kafka (via Upstash)

## Architecture

The architecture is built on top of [Kratos](https://go-kratos.dev/en/), a framework for rolling out microservices. A growing codebase will eventually bring about the need to convert sections of the app into submodules for organization and repo size management - the folder structure is designed to easily support this scenario.

When a service is started, it runs gRPC and REST servers concurrently. They are configured jointly using protofiles (see: [protobuf docs](api)). Both gRPC and HTTP servers are run concurrently for each service. Services are modular are isolated - they can be configured to run independently on many machines, horizontally scaled, auto-scaled, or otherwise deployed as needed.

The main (core) service can be run using `make execute`. This service handles user operations and (optionally) some others. End users should choose what services they will need.

Google's [Wire](https://github.com/google/wire/tree/main) tool is used for compile-time dependency injection. See Wire's [User Guide](https://github.com/google/wire/blob/main/docs/guide.md), [FAQ](https://github.com/google/wire/blob/main/docs/faq.md), and [Best Practices](https://github.com/google/wire/blob/main/docs/best-practices.md) document for further reference on how Wire works.

Caching and Kafka streaming are provided via [Upstash](https://upstash.com/). Instrumentation is provided using [Jaegar](https://www.jaegertracing.io/), and the [Prometheus](https://prometheus.io/) stack.

## Setup

You will need a `.env` file. Ask a team member for a copy or make your own based on [.env.example](.env.example). 

### Installation

Build a Docker image.

## Working with the code

The workflow to create new services is roughly as follows:

- Define a proto service: 
  - If defining an internal service: `kratos proto add api/v1/my_service.proto`
  - For a namespaced or public-facing service: `kratos proto add api/v1/my_namespace/my_service.proto`
- Convert the protofile to Go code: `kratos proto client api/v1/my_namespace/my_service.proto`
- Define a service
  - If defining an internal service: `kratos proto server api/v1/my_service.proto -t internal/service`
  - For a publicly accessible service: `kratos proto server api/v1/my_namespace/my_service.proto -t pkg/my_service/service`

## Install Kratos

```
go install github.com/go-kratos/kratos/cmd/kratos/v2@latest
```

## Create a service

```shell
kratos proto add api/v1/<filename>.proto # Add a proto template
make proto # Generate the proto code
```

To generate the code for a new service, use the following.

If creating an internal service (such as a user auth service, system service, or transaction processing service) then use the following:

```shell
# Generate the code layout for the service
kratos proto server api/v1/<filename>.proto -t internal/service
```

For a publicly available service, create a [package](https://www.golang-book.com/books/intro/11) for it. Make a new folder for your package in the `pkg` directory. Add a subfolder called `service`. Then run the following:

```shell
# Generate the code layout for the service
kratos proto server api/v1/<filename>.proto -t pkg/<pkg_folder_name>/service
```


Make sure to update the Makefile, adding an available build step for the new service. See following example.

Add a build step for the `foobar` service:
```dart
.PHONY: foobar
fin:
	./bin/foobar &
```

## Compile services

Services cannot currently be compiled independently.

```shell
make api # Build protofiles
make build # Build Go code
```

## Generate auxiliary files

```shell
# Generate API files (include: pb.go, http, grpc, validate, swagger) by proto file
make api
```

## Build all files

```shell
make compile # Build all files
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

# Tutorial: Adding A Service

Let's add a service for consultants. It will expose a gRPC and HTTP server through which CRUD operations can be run on consultant data.

("Consultant" here is a generic term and does not hold specialized in-context meaning)

There is more than one way to add a service. For critical internal services, we can add them to our core servers. There are several examples of this in action in the server package at `internal/server/` in the *NewCoreGRPCServer* and *NewCoreHTTPServer* functions.

For other services, especially publicly exposed or re-usable ones, we want to define a [package](https://www.golang-book.com/books/intro/11). 

First, we create our protofile:

```shell
kratos proto add api/v1/consultants/consultants.proto
```

This will generate a blank Protofile:

```proto
syntax = "proto3";

// this package now available like so:
//   import consultantsV1 "core/api/v1/consultants"
// use package api.v1 for internal services
package api.v1.consultants;

option go_package = "core/api/v1/consultants;consultant";
option java_multiple_files = true;
option java_package = "api.v1.consultants";

service Consultants {
	rpc CreateConsultants (CreateConsultantsRequest) returns (CreateConsultantsReply);
	rpc UpdateConsultants (UpdateConsultantsRequest) returns (UpdateConsultantsReply);
	rpc DeleteConsultants (DeleteConsultantsRequest) returns (DeleteConsultantsReply);
	rpc GetConsultants (GetConsultantsRequest) returns (GetConsultantsReply);
	rpc ListConsultants (ListConsultantsRequest) returns (ListConsultantsReply);
}

message CreateConsultantsRequest {}
message CreateConsultantsReply {}

message UpdateConsultantsRequest {}
message UpdateConsultantsReply {}

message DeleteConsultantsRequest {}
message DeleteConsultantsReply {}

message GetConsultantsRequest {}
message GetConsultantsReply {}

message ListConsultantsRequest {}
message ListConsultantsReply {}
```

Fill in your protofile. Note that you can import from other protofiles in the `api` directory as well as the third party vendor protofiles in the `third_party` directory. 

For example, to use timestamps we can import Google's *timestamp.proto* specification. Other useful data types in Google's specifications include Duration, Empty, Any, and HTTP.

Our Consultants protofile might now look like this:

```proto
syntax = "proto3";

package api.v1.consultant;

import "google/protobuf/timestamp.proto";
import "google/api/annotations.proto";
import "v1/users.proto";

// ...

service Consultants {
  // ...
}

message Consultant {
	string id = 1;
	string user_id = 2;
  api.v1.User user = 3; // imported from users protofile
}
```

(Note that all fields are considered nullable, so a `nil` user won't trigger errors in a consultant's data record.)

Run the makefile command to generate code from protofiles:

```shell
make proto
```

## Creating our package

A package is made of several sub-packages:

- biz/: Define business logic and establish the data schema. Think of this as the "model" if coming from an MVC context.
- data/: For DB operations and low-level requests
- internal/: (optional) Additional includes or definitions
- server/: (optional) Define HTTP and gRPC servers if exposing the service. 
- service/: Define request-level logic. Think of this as the "controller" if coming from an MVC paradigm.

For our package, we'll be defining packages such as `consultant_biz`, `consultant_data`, and so on.

### Business Layer

In our `consultant_biz` sub-package, we define a [gORM](https://gorm.io/) model and declare what actions we are able to run on it. We declare this by defining a "repo" of available actions such as saving or deleting records. Lastly, we define public functions for each action we will make publicly available.

For our example, we will define `Get` and `Save` functions.

```go
package consultants_biz

type Consultant struct {
	gorm.Model
	ID                string                 `gorm:"primaryKey" protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	User              *biz.User              `protobuf:"bytes,2,opt,name=user,proto3" json:"user,omitempty"`
	...
	CreatedAt         *timestamppb.Timestamp `protobuf:"bytes,10,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	UpdatedAt         *timestamppb.Timestamp `protobuf:"bytes,11,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
}

// See: GORM Hooks - https://gorm.io/docs/hooks.html
func (c *Consultant) BeforeCreate(tx *gorm.DB) error {
	if c.ID == "" {
		c.ID = uuid.New().String()
	}
	return nil
}

type ConsultantRepo interface {
	Get(context.Context, string) (*Consultant, error)
	Save(context.Context, *Consultant) (*Consultant, error)
}

type ConsultantAction struct {
	repo ConsultantRepo
	log  *log.Helper
}

func NewConsultantAction(repo ConsultantRepo, logger log.Logger) *ConsultantAction {
	return &ConsultantAction{repo: repo, log: log.NewHelper(logger)}
}

// Public function, calls our repo internally
func (uc *ConsultantAction) GetConsultant(ctx context.Context, id string) (*Consultant, error) {
	uc.log.WithContext(ctx).Infof("GetConsultant: %s", id)
	consultant, err := uc.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return consultant, nil
}

// Public function, calls our private repo.Save() function
func (uc *ConsultantAction) CreateConsultant(ctx context.Context, c *Consultant) (*Consultant, error) {
	uc.log.WithContext(ctx).Infof("CreateConsultant: %s", c.ID)
	res, err := uc.repo.Save(ctx, c)
	if err != nil {
		fmt.Println("error creating consultant: ", err)
	}
	fmt.Println("create consultant result: ", res)
	return res, err
}
```

With our business logic code finished, we need to configure what parts of the sub-package we want to make available in other packages. We do this by defining a `ProviderSet` for each package (and sub-package).

Expose the Consultant's available business actions:

```go
package consultants_biz

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	NewConsultantAction,
)
```

### Data Layer

We will define the lower-level actions behind our business logic here. Database transactions, caching requests, and external network requests live here.

Our `pkg/consultants/data/consultants.go` file may look like this:

```go
package consultants_data

// ...

type consultantRepo struct {
	data *Data
	log  *log.Helper
}

func NewConsultantRepo(data *Data, logger log.Logger) consultants_biz.ConsultantRepo {
	return &consultantRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *consultantRepo) Get(ctx context.Context, id string) (*consultants_biz.Consultant, error) {
	var consultant *consultants_biz.Consultant
	err := server.DB.First(&consultant, id).Error
	if err != nil {
		return nil, err
	}

	return consultant, nil
}

func (r *consultantRepo) Save(ctx context.Context, c *consultants_biz.Consultant) (*consultants_biz.Consultant, error) {
	if c.ID != "" {
		if err := server.DB.Save(&c).Error; err != nil {
			return nil, err
		} else {
			return c, nil
		}
	}

	if err := server.DB.FirstOrCreate(&c).Error; err != nil {
		return nil, err
	}

	return c, nil
}
```

Again we expose this sub-package using wire.

```go
package consultants_data

var ProviderSet = wire.NewSet(
	NewConsultantRepo, NewData,
)

// The following are used by convention

// Data .
type Data struct {
	// wrapped database client
}

// NewData .
func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}

	return &Data{}, cleanup, nil
}
```

### Service Layer

This may be familiar to those with backgrounds using MVC frameworks such as Rails or Django. The service layer is comparable to the "controller" concept in MVC. Here we process a request, trigger any business logic, and reply with the appropriate response data.

```go
package consultants_service

import (
	"context"

	consultantsV1 "core/api/v1/consultant"
	consultants_biz "core/pkg/consultants/biz"
)

type ConsultantService struct {
	consultantsV1.UnimplementedConsultantsServer

	action *consultants_biz.ConsultantAction
}

func NewConsultantService(action *consultants_biz.ConsultantAction) *ConsultantService {
	return &ConsultantService{action: action}
}

func (s *ConsultantService) GetConsultant(ctx context.Context, req *consultantsV1.GetConsultantRequest) (*consultantsV1.GetConsultantReply, error) {
	consultant, err := s.action.GetConsultant(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &consultantsV1.GetConsultantReply{
		Ok:         err == nil,
		Consultant: consultants_biz.ConsultantToProtoData(consultant),
	}, nil
}

func (s *ConsultantService) CreateConsultant(ctx context.Context, req *consultantsV1.CreateConsultantRequest) (*consultantsV1.CreateConsultantReply, error) {
	consultant := consultants_biz.ProtoToConsultantData(req.Consultant)
	createdConsultant, err := s.action.CreateConsultant(ctx, consultant)
	if err != nil {
		return nil, err
	}
	return &consultantsV1.CreateConsultantReply{
		Ok:         err == nil,
		Consultant: consultants_biz.ConsultantToProtoData(createdConsultant),
	}, nil
}
```

### Server Layer

This step is relatively simple. We'll set up our GRPC and HTTP servers, and register the relevant endpoints on them. We can do this in less than 40 lines of code.

Create gRPC and HTTP servers for the consulting service:

```go
package consultants_server

import (
	consultantsV1 "core/api/v1/consultant"
	"core/internal/conf"
	"core/internal/server"
	consultantsService "core/pkg/consultants/service"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
)

func NewConsultantsGrpcServer(
	c *conf.Server,
	logger log.Logger,
	consultant *consultantsService.ConsultantService,
) *grpc.Server {
	srv := server.GRPCServerFactory("consultants", c, logger)
	consultantsV1.RegisterConsultantsServer(srv, consultant)

	return srv
}

func NewConsultantsHTTPServer(
	c *conf.Server,
	logger log.Logger,
	consultant *consultantsService.ConsultantService,
) *http.Server {
	srv := server.HTTPServerFactory("consultants", c, logger)
	consultantsV1.RegisterConsultantsHTTPServer(srv, consultant)

	server.StartPrometheus(srv)
	return srv
}
```

As before, in the sub-package we register these functions as publicly available to use at compile time.

`pkg/consultants/server/server.go`:
```go
package consultants_server

import (
	"github.com/google/wire"
)

// Declare server types to run concurrently at runtime
var ProviderSet = wire.NewSet(
	NewConsultantsGrpcServer, NewConsultantsHTTPServer,
)

```

## Compilation

To build and compile our service, add a folder to the `app/` directory. We will need two files in order to configure our service.

First, we'll need a main entrypoint for our service to run from. You'll add something like the following in `app/consultants/main.go`:

```go
package main

import (
	"os"

	"core/internal/server"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	// "google.golang.org/grpc"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name string = "consultants"
	// Version is the version of the compiled software.
	Version string
	// flagconf is the config flag.
	flagconf string

	id, _ = os.Hostname()

	KafkaTopics = []string{"consultants", "consultants/cdc"}
)

func init() {
	server.InitEnv(id, Name, &flagconf, KafkaTopics)
}

func newConsultantsApp(logger log.Logger, gs *grpc.Server, hs *http.Server) *kratos.App {
	return server.NewApp(Name, id, Version, logger, gs, hs)
}

func main() {
	server.RunApp(Name, Version, flagconf, wireApp)
}

```

Meanwhile, we'll need to pull in the code from our sub-packages using Wire to compile all code exposed in our sub-packages' Provider Sets at build time.

In `app/consultants/wire.go`, add the following:

```go
//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	...
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// wireApp init kratos application.
func wireApp(*conf.Server, *conf.Data, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(
		consultants_server.ProviderSet, consultants_data.ProviderSet,
		consultants_biz.ProviderSet, consultants_service.ProviderSet, newConsultantsApp,
	))
}

```

(Note that `newConsultantsApp` was just defined in the previous step in our main package.)

Don't forget to add any new DB tables to the autoMigration list like so:

```go
func automigrateDBTables(*gorm.DB) {
	DB.AutoMigrate(&consultants_biz.Consultant{})
}
```

An important final step is to configure the endpoint the service will will run on. This can be configured in `configs/config.yaml` by adding something like the following:

```yaml
consultants:
  http:
    addr: 0.0.0.0:8103
    timeout: 1s
  grpc:
    addr: 0.0.0.0:9103
    timeout: 1s
  database:
    driver: mysql
    source: root:root@tcp(127.0.0.1:3306)/test?parseTime=True&loc=Local
  redis:
    addr: 127.0.0.1:6379
    read_timeout: 0.2s
    write_timeout: 0.2s
```

(Note: the `redis` and `database` fields aren't used as this architecture doesn't run redis locally on-machine. However they could be activated if using split DBs per namespace or other such techniques.)

To build our new service, simply run the following commands:

```shell
make build # if needed
make consultants
```

# ProtoBuffer Definitions

The first step in adding a service is defining protocol buffers.

For each service, we define services (which establish available requests) and messages (which establish the structure of data that will be passed). Services define both gRPC and HTTP endpoints. 

## Example 

The following example define a `Users` service and `User` message, along with messages that define what parameters we expect in the Request and Reply. The service holds a request called `GetUser` - which can be accessed via REST at the endpoint `"<server_url>/users/{id}"`.

```grpc
service Users {
  rpc GetUser (GetUserRequest) returns (GetUserReply) {
    option (google.api.http) = {
      get: "/users/{id}"
    };
  };
}

message GetUserRequest {
  optional string id = 1; // we can retrieve the user by ID
  optional string email = 3; // or retrieve by email
}
message GetUserReply {
  User user = 1;
}

message User {
  string id = 1;
  string username = 2;
  string email = 3;
  // ...
}
```

## Creating ProtoBuffers

DO NOT manually create proto files.

Use the following command:

```shell
kratos proto add api/v1/my_service.proto
```

Replace "my_service" above with the name of your service. Prefer plural spellings, such as "users" over singular, such as "user" where applicable.

## Compiling ProtoBuffers

```shell
make proto
```


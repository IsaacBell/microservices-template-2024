# Leads Service

An API for B2B lead gen and customer data. This can be used for lead gen data gen/consumption, or for tools such as email workflow management (to store and manage business contact data).

This package also includes access to financial data related to USA spending and Senate lobbying.

## Package Structure

The leads package consists of the following files:

- leads_service.go: Contains the implementation of the lead service.
- finance_util.go: Provides utility functions for initializing and accessing the Finnhub API client.

## Usage


```go
action := leads_biz.NewLeadAction(repo, logger)
leadService := leads_service.NewLeadService(action)
```

The LeadService depends on the LeadAction from the leads_biz package, which should be initialized with the necessary dependencies (repository and logger).

### Methods

```go
GetLead(ctx context.Context, req *leadsV1.GetLeadRequest) (*leadsV1.GetLeadReply, error): Retrieves a lead by its ID.
GetUSASpending(ctx context.Context, req *v1.GetUSASpendingRequest) (*v1.GetUSASpendingReply, error): Retrieves USA spending data (not yet implemented).
GetSenateLobbying(ctx context.Context, req *v1.GetSenateLobbyingRequest) (*v1.GetSenateLobbyingReply, error): Retrieves Senate lobbying data (not yet implemented).
```

### Finance Utility

The finance_util.go file provides utility functions for initializing and accessing the Finnhub API client.
Usage

## Configuration

The Finnhub API client requires an API token to be set in the configuration. The token should be provided as an environment variable named FINNHUB_API_TOKEN.

Make sure to set the FINNHUB_API_TOKEN environment variable with your Finnhub API token before using the finance utility functions.

## Dependencies

The leads package has the following dependencies:

```
microservices-template-2024/api/v1: Provides the generated gRPC service definitions.
microservices-template-2024/api/v1/b2b: Provides the generated gRPC service definitions specific to B2B.
microservices-template-2024/pkg/leads/biz: Provides the business logic for leads.
github.com/Finnhub-Stock-API/finnhub-go/v2: Provides the Finnhub API client for retrieving financial data.
```

Make sure to properly install and configure these dependencies before using the leads package.

### TODO

- Implement the logic for retrieving USA spending data in the GetUSASpending method.
- Implement the logic for retrieving Senate lobbying data in the GetSenateLobbying method.
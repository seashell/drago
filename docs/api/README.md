# HTTP API

## Prefix

All API routes are prefixed with `/api/`.

## Addressing and Ports
Drago binds to a specific set of addresses and ports. The HTTP API is served via the HTTP address and port. The default port for the Drago HTTP API is 8081. This can be overridden via the Drago configuration block. Here is an example curl request to query the status of a Drago server:

```bash
$ curl http://127.0.0.1:8081/api/status
```

## Data Model and Layout
There are four primary nouns in Drago:

- nodes
- networks
- interfaces
- connections

## ACLs
Several endpoints in Drago use or require ACL tokens to operate. The tokens are used to authenticate the request and determine if the request is allowed based on the associated authorizations. Tokens are specified per-request by using the X-Drago-Token request header set to the Secret of an ACL Token.

## Authentication
When ACLs are enabled, a Drago token should be provided to API requests using the X-Drago-Token header.

Here is an example using curl:

```bash
$ curl \
    --header "X-Drago-Token: aa534e09-6a07-0a45-2295-a7f77063d429" \
    https://localhost:4646/api/nodes
```

## Formatted JSON Output
By default, the output of all HTTP API requests is JSON.

## HTTP Methods
Drago's API aims to be RESTful, although there might be some exceptions. The API responds to the standard HTTP verbs GET, POST, PUT, and DELETE. Each API method will clearly document the verb(s) it responds to and the generated response. The same path with different verbs may trigger different behavior. For example:

```
POST /v1/networks/
GET /v1/networks/
```

Even though these share a path, the POST operation creates a new network whereas the GET operation reads all networks.

## HTTP Response Codes
Individual API's will contain further documentation in the case that more specific response codes are returned but all clients should handle the following:

- 200 and 204 as success codes.
- 400 indicates a validation failure and if a parameter is modified in the request, it could potentially succeed.
- 403 marks that the client isn't authenticated for the request.
- 404 indicates that the resource targeted does not exist.
- 5xx means that the client should not expect the request to succeed if retried. Whenever a status 5xx is returned, more details about the error will be contained within the response body.
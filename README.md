# apisonator-cache
An experimental PoC of providing a caching layer in front of 3scale Apisonator

## Usage

This server provides a scaled down version of the 3scale Service Management APIs.
It provides both the auth and authrep endpoints and accepts requests in a format that matches
the expectations of current implementation of [Apisonator](https://github.com/3scale/apisonator#apisonator-listener).

In order to keep the server as performant as possible, it does not respond in the XML format.
Instead, this server enforces the use of the following two extensions:

1. [no_body](https://github.com/3scale/apisonator/blob/master/docs/extensions.md#rejection_reason_header-boolean) - 
The server will respond with an empty body.
2. [rejection_reason_header](https://github.com/3scale/apisonator/blob/v2.100.0/docs/extensions.md#rejection_reason_header-boolean) - 
Non 2xx status code responses will store the rejection reason in the `3scale-rejection-reason` header.

## Using the application

Currently, the caching layer has been designed to work with [APIcast](https://github.com/3scale/APIcast) only.
The application presents itself to APIcast as Apisonator through configuration and forwards requests to 3scale during cache
warm-up and later, periodically to flush the cache and fetch updated state.

The application takes two flags:
1. `port` - the port to listen and serve the HTTP server on which is not required and defaults to `3000` if not provided.
2. `upstream` - the URL of 3scale backend which is required.

To start the server locally and link it against 3scale SaaS backend you can do as follows:
`go run main.go --upstream=https://su1.3scale.net`
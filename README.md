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

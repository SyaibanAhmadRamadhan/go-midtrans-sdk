# Go Midtrans SDK

This SDK is a Golang library designed to integrate payments via Midtrans, with additional features such as OpenTelemetry support, automatic retry for API requests, and type-safe request structures.

## Features

- **OpenTelemetry Support**: This SDK includes OpenTelemetry for monitoring and tracing transactions sent to Midtrans.
- **Retry Mechanism**: Automatically retries requests if the initial request to the Midtrans API fails.
- **Type-Safe Requests**: All requests are type-safe, minimizing the risk of type-related errors.
- **Validation Compliance**: Built-in validation ensures that request parameters adhere to Midtrans requirements, preventing errors from invalid fields or values.
- **Customizable Logging**: Allows custom logging configurations, with OpenTelemetry as the default option, to fit diverse logging and monitoring needs.

## OpenTelemetry Support
With OpenTelemetry, each request can be traced end-to-end within a trace. Ensure OpenTelemetry is initialized in your application to enable this feature.
if you use otel, it will automatically inject metadata when sending it to the midtrans api, so when there is a notifications callback from midtrans you can get the traceparent from the midtrans request body, and you can set propagators for distributed tracing
and this traceparent is taken from the context, if you use http you can set the traceparent context value.
For examples or references, you can check the http library that I have created https://github.com/SyaibanAhmadRamadhan/http-wrapper/blob/main/otel_Trace.go#L138

## Installation

Use `go get` to install the SDK into your Go project:

```shell
go get github.com/SyaibanAhmadRamadhan/go-midtrans-sdk@v1.241028.2152
```
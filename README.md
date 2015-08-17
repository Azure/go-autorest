# go-autorest

[![GoDoc](https://godoc.org/github.com/azure/go-autorest/autorest?status.png)](https://godoc.org/github.com/azure/go-autorest/autorest)

## Usage
This package implements an HTTP request pipeline suitable for use across multiple go-routines
and provides the shared routines relied on by AutoRest (see https://github.com/azure/autorest/)
generated Go code.

The package breaks sending and responding to HTTP requests into three phases: Preparing, Sending,
and Responding. A typical pattern is:

```
  req, err := Prepare(&http.Request{},
    WithAuthorization())

  resp, err := Send(req,
    WithLogging(logger),
    DoErrorIfStatusCode(500),
    DoCloseIfError(),
    DoRetryForAttempts(5, time.Second))

  err = Respond(resp,
    ByClosing())
```

Each phase relies on decorators to modify and / or manage processing. Decorators may first modify
and then pass the data along, pass the data first and then modify the result, or wrap themselves
around passing the data (such as a logger might do). Decorators run in the order provided. For
example, the following:

```
  req, err := Prepare(&http.Request{},
    WithBaseURL("https://microsoft.com/"),
    WithPath("a"),
    WithPath("b"),
    WithPath("c"))
```

will set the URL to:

```
  https://microsoft.com/a/b/c
```

Preparers and Responders may be shared and re-used (assuming the underlying decorators support
sharing and re-use). Performant use is obtained by creating one or more Preparers and Responders
shared among multiple go-routines, and a single Sender shared among multiple sending go-routines,
all bound together by means of input / output channels.

Decorators hold their passed state within a closure (such as the path components in the example
above). Be careful to share Preparers and Responders only in a context where such held state
applies. For example, it may not make sense to share a Preparer that applies a query string from a
fixed set of values. Similarly, sharing a Responder that reads the response body into a passed
struct (e.g., ByUnmarshallingJson) is likely incorrect.

Lastly, the Swagger specification (https://swagger.io) that drives AutoRest
(https://github.com/azure/autorest/) precisely defines two date forms: date and date-time. The
github.com/azure/go-autorest/autorest/date package provides time.Time derivations to ensure
correct parsing and formatting.

See the included examples for more detail. For details on the suggested use of this package by
generated clients, see the Client described below.


## Install

```bash
go get github.com/azure/go-autorest/autorest
```

## License

See LICENSE file.

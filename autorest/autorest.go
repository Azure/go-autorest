/*
  This package provides the shared routines used by Go code generated through the AutoRest code generator
  (see https://github.com/azure/autorest/).

  The package breaks sending and responding to HTTP requests into three phases: Preparing, Sending, and Responding.
  A typical pattern is:

    req, err := Prepare(&http.Request{}, WithBearerAuthorization("SECRET TOKEN"))
    resp, err := Send(req, WithRetry(5, 2), WithLogging(logger))
    err = Respond(resp, ByClosing())

  Each phase relies on decorators to modify and / or manage its processing. Decorators may first modify and then pass
  the data along, pass the data first and then modify the result, or wrap themselves around passing the data (such as
  a logger might do). Decorators run in the order provided (from the perspective of a decorator that modifies and then
  passes the data). For example, the following:

    req, err := Prepare(&http.Request{},
      WithBaseURL("https://microsoft.com/"),
      WithPath("a"),
      WithPath("b"),
      WithPath("c"))

  will set the URL to:

    https://microsoft.com/a/b/c

  Preparers (normally created by calling CreatePreparer, passing a sequence of PrepareDecorators) and Responders
  (normally created by calling CreateResponder, passing a sequence of RespondDecorators) may be shared and re-used
  (assuming the underlying decorators support sharing and re-use). Performant use is obtained by creating one or more
  Preparers and Responders shared among multiple go-routines, and a single Sender shared among multiple sending
  go-routines, all bound together by means of input / output channels.

  The decorators used to create Preparers and Responders hold on to the passed stated (such as the path components in
  the example above). Be careful to share Preparers and Responders only in a context where such held state applies.
  For example, it may not make sense to share a Preparer that applies a query string from a fixed set of values.
  Similarly, sharing a Responder that reads the response body into a passed struct (e.g., ByUnmarshallingJsonAndClosing)
  is likely incorrect.

  Lastly, the Swagger specification (https://swagger.io) that drives AutoRest (https://github.com/azure/autorest/)
  precisely defines two date forms (i.e., date and date-time). The two sub-packages -- github.com/azure/go-autorest/autorest/date
  and github.com/azure/go-autorest/autorest/datetime -- provide time.Time derivations to ensure correct parsing and
  formatting.

  See the included examples for more detail.

*/
package autorest

# Developer Notes

## Background

After experimenting with multiple approaches, the current design was
chosen because it should make supporting both BSER and JSON encoding
easier. It also provides mechanisms to add support for additional
commands outside the module.

Watchman's protocol is documented here:
https://facebook.github.io/watchman/docs/socket-interface.html

## Code Organization

The focus of this package is providing low-level building blocks
to interact with the Watchman server. In other words: connecting,
marshaling, and unmarshaling messages. It should leave decisions
about channels, goroutines, and such to other packages. Likewise,
high-level interfaces and abstractions should not be included in
this package.

Guidelines for this package include:

* One Watchman command per file.
* Each command file should include three items:
  * Request struct
  * Response struct
  * A ResponseTranslator to convert from raw ResponsePDU
    to Response struct.
* Request structs should expose their member values, and
  implement the Request interface,
* Response structs should hide their member values, and
  provide accessor methods.
* Keep the mapping from PDU to Request or Response struct
  as simple as possible.
* Each command file should have a matching unit test file.
* Each command should also be exercised by one or more
  integration tests.

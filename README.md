# The Things Network Utilities for Go

[![Build Status](https://travis-ci.org/TheThingsNetwork/go-utils.svg?branch=master)](https://travis-ci.org/TheThingsNetwork/go-utils)

## Utilities

- `backoff`: backoff algorithm extracted from [`github.com/grpc/grpc-go`](https://github.com/grpc/grpc-go)
- `encoding`: encoding and decoding between a `struct` and `map[string]string`
- `grpc/interceptor`: gRPC interceptor that logs RPCs
- `grpc/restartstream`: gRPC interceptor that restart streams when the underlying connection breaks and restores
- `handlers/cli`: CLI logger for [`github.com/apex/log`](https://github.com/apex/log)
- `handlers/elasticsearch`: [Elasticsearch](https://www.elastic.co/products/elasticsearch) logger for [`github.com/apex/log`](https://github.com/apex/log)
- `log`: log wrapper
- `random` and `pseudorandom`: wrappers for (pseudo)random functions
- `queue`: implementations of queues and schedules
- `rate`: rate counting and rate limiting
- `roots`: CA Root certificates that are used when the OS doesn't supply them

## License

Source code is released under the MIT License, which can be found in the [LICENSE](LICENSE) file. A list of authors can be found in the [AUTHORS](AUTHORS) file.

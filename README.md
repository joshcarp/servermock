# Mirror

A grpc/http dynamic stub

## Usage

1. run a mirror server ./mirror

2a. use the go API to set a response for a traceID

```go
SetResponse("http://localhost:8000", "/path/whatever", "{'Foo':'Bar'}")
SetGRPCResponse("http://localhost:8000", "/foo.service.bar.SomethingAPI/GetWhatever", proto.Message)
SetJsonResponse("http://localhost:8000", "/path/whatever", interface{})

```

2b. OR use the `Set-Data:` Header to post data for a traceid

```bash
curl -d '{"Request": "1"}' -H "Set-Data: true" localhost:8000/path/whatever -X POST
curl -d '{"Request": "2"}' -H "Set-Data: true" localhost:8000/path/whatever -X POST
curl -d '{"Request": "3"}' -H "Set-Data: true" localhost:8000/path/whatever -X POST
```

3. Set your downstream to point at `http://localhost:8000`

```bash
curl -H localhost:8000/path/whatever
> {"Request": "1"}
curl localhost:8000/path/whatever
> {"Request": "2"}
curl localhost:8000/path/whatever
> {"Request": "3"}
```

or whatever the equivalent grpc request would be Note: Requests are set FIFO for every TraceID, and a request is cleared
from memory once a request is completed


# Mirror
A grpc/http dynamic stub

## Usage

1. run a mirror server
./mirror
   
2a. use the go API to set a response for a traceID
```go
SetResponse("http://localhost:8000", "traceid1234", "{'Foo':'Bar'}")
SetGRPCResponse("http://localhost:8000", "traceid1234", proto.Message)
SetJsonResponse("http://localhost:8000", "traceid1234", interface{})

```
2b. Use the `Set-Data:` Header to post data for a traceid
```bash
curl -d '{"Request": "1"}' -H "Set-Data: 5678" -H "X-B3-TraceId: 99999" localhost:8000 -X POST
curl -d '{"Request": "2"}' -H "Set-Data: 5678" -H "X-B3-TraceId: 99999" localhost:8000 -X POST
curl -d '{"Request": "3"}' -H "Set-Data: 5678" -H "X-B3-TraceId: 99999" localhost:8000 -X POST
```

3. Set your downstream to point at `http://localhost:8000`

```bash
curl -H "X-B3-TraceId: 99999" localhost:8000
> {"Request": "1"}
curl -H "X-B3-TraceId: 99999" localhost:8000
> {"Request": "2"}
curl -H "X-B3-TraceId: 99999" localhost:8000
> {"Request": "3"}
```

or whatever the equivilent grpc request would be
Note: Requests are set FIFO for every TraceID, and a request is cleared from memory once a request is completed


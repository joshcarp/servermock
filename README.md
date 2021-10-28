# servermock

Dynamic Mocking Tool

## Usage

1. run a servermock server ./servermock

2a. use the go API to set a response for a path

```go
SetResponse("http://localhost:8000", "/path/whatever", "{'Foo':'Bar'}")
SetGRPCResponse("http://localhost:8000", "/foo.service.bar.SomethingAPI/GetWhatever", proto.Message)
SetJsonResponse("http://localhost:8000", "/path/whatever", interface{})

```

2b. OR use the `Set-Data:` Header to post data for a path

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

or whatever the equivalent grpc request would be (not supporting the reflection api)

Note: Requests are set FIFO for every Path, and a request is cleared
from memory once a request is completed


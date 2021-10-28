
<p align="center">
</p>


<h1 align="center">servermock</h1>

<div align="center">

[![Status](https://img.shields.io/badge/status-active-success.svg)]()
[![GitHub Issues](https://img.shields.io/github/issues/joshcarp/servermock)](https://github.com/joshcarp/servermock/issues)
[![GitHub Pull Requests](https://img.shields.io/github/issues-pr/joshcarp/servermock)](https://github.com/joshcarp/servermock/pulls)
[![License](https://img.shields.io/badge/license-apache2-blue.svg)](/LICENSE)

</div>

---


## 📝 Table of Contents
- [About](#about)
- [Getting Started](#getting_started)
- [Usage](#usage)
- [Authors](#authors)
- [Acknowledgments](#acknowledgement)

## 🧐 About <a name = "about"></a>

servermock is a go package that can be used to mock out http or grpc servers simply without any external server implementations.


## 🚀 Usage <a name = "usage"></a>

### Inline in golang
1. Start a server
```golang
err := servermock.Serve(ctx, Printf, ":8000")
```
2. Set the response
```golang
err = servermock.SetResponse("http://localhost:8000", servermock.Request{
		Path:       "/foo.service.bar.SomethingAPI/GetWhatever",
		Body:       []byte(`{"Hello": "true"}`),
		StatusCode: 200,
	})
```
3. Send a request
```golang
resp, err := http.Get("http://localhost:8000/foo.service.bar.SomethingAPI/GetWhatever")
// {"Hello": "true"}
```

### In a docker container
1. Run the docker container
```bash
docker run -p 8000:8000 joshcarp/servermock
```
2. Set the response conforming to the `servermock.Request` type
```bash
curl --header "Content-Type: application/json" --header "SERVERMOCK-MODE: SET" --request POST --data '{"path":"/foo.service.bar.SomethingAPI/GetWhatever","body":"eyJIZWxsbyI6ICJ0cnVlIn0=","status_code":200' http://localhost:8000/foo.service.bar.SomethingAPI/GetWhatever
```
3. Send a request
```bash
curl localhost:8000/foo.service.bar.SomethingAPI/GetWhatever
> {"Hello": "true"}                                                                                                  
```

### gRPC vs REST servers

Setting data always occurs over http 1.0 using the json payload, gRPC servers are, after all, just servers that return some bytes.

see [example/example_test.go](example/example_test.go) for full examples.

## ✍️ Authors <a name = "authors"></a>
- [@joshcarp](https://github.com/joshcarp)

## 🎉 Acknowledgements <a name = "acknowledgement"></a>
- [@emmaCullen](https://github.com/emmaCullen) had the original idea for this package.
- [github.com/dnaeon/go-vcr](https://github.com/dnaeon/go-vcr) is similar but different; whilst any network traffic can be recorded and replayed, servermock tries tosimplify mocking of servers in unit tests/contexts where writing a specific server implementation is a little too much. 

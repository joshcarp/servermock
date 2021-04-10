//go:generate bash -c "docker run --rm -v $(pwd):/example:rw joshcarp/protoc -I./example/ --go-grpc_out=paths=source_relative:example --go_out=paths=source_relative:example example.proto"

package example

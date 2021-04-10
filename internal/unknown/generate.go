//go:generate bash -c "docker run --rm -v $(pwd):/unknown:rw joshcarp/protoc -I./unknown/ --go_out=paths=source_relative:unknown unknown.proto"

package unknown

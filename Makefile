jsonstream: ./cmd/jsonstream/jsonstream.go ./reader.go ./formatter.go ./utils.go
	go build cmd/jsonstream/jsonstream.go

clean:
	rm ./jsonstream

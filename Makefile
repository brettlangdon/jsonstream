jsonstream: ./cmd/jsonstream.go ./reader.go ./formatter.go ./utils.go
	go build cmd/jsonstream.go

clean:
	rm ./jsonstream

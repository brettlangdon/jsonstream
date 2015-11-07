package jsonstream

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"
)

// Reader is used for reading newline delimited JSON objects from an input io.Reader
type Reader struct {
	buffer *bufio.Reader
	keys   map[string]bool
}

// NewReader will create a new Reader from the provided io.Reader and keys
// k, when not an empty slice, will tell the Reader to only include the provided keys from the input JSON object
func NewReader(r io.Reader, k []string) *Reader {
	var keys map[string]bool
	keys = make(map[string]bool, 0)
	for _, key := range k {
		keys[key] = true
	}
	return &Reader{
		buffer: bufio.NewReader(r),
		keys:   keys,
	}
}

func (reader *Reader) processData(data interface{}) (processed map[string]interface{}, err error) {
	var fields map[string]interface{}
	fields, err = getAsMap(data)
	if err != nil {
		return processed, err
	}

	processed = make(map[string]interface{})
	for key, value := range fields {
		if _, ok := reader.keys[key]; ok {
			processed[key] = value
		}
	}
	return processed, err
}

// ReadLine will read the next JSON object from the input io.Reader
func (reader *Reader) ReadLine() (data interface{}, err error) {
	var line []byte
	var isPrefix bool
	line, isPrefix, err = reader.buffer.ReadLine()

	if isPrefix {
		err = errors.New("Line exceeds the length of the buffer")
	}

	if err != nil {
		return data, err
	}

	// skip empty lines, we'll fail at processing them anyways
	if len(line) == 0 {
		return reader.ReadLine()
	}

	err = json.Unmarshal(line, &data)

	if err != nil {
		return data, err
	}

	if len(reader.keys) > 0 {
		data, err = reader.processData(data)
	}

	return data, err
}

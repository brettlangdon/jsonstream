package jsonstream

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"
)

type Reader struct {
	buffer *bufio.Reader
}

func NewReader(r io.Reader) *Reader {
	return &Reader{
		buffer: bufio.NewReader(r),
	}
}

func (reader *Reader) ReadLine() (data interface{}, err error) {
	var line []byte
	var isPrefix bool
	line, isPrefix, err = reader.buffer.ReadLine()

	if isPrefix {
		err = errors.New("Line exceeds the length of the buffer")
	}

	if err == nil {
		err = json.Unmarshal(line, &data)
	}

	return data, err
}

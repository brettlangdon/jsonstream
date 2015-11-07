package jsonstream

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
)

// FormatType denotes the output formatting types supported by jsonstream
type FormatType int

const (
	// FormatJSON informs the Formatter to output as a JSON object
	FormatJSON FormatType = iota
	// FormatTSV informs the Formatter to output as a tab separated JSON encoded values
	FormatTSV
	// FormatTSVKey informs the Formatter to output as a tab separated `<key>=<value>` pairs
	FormatTSVKey
)

// Formatter is used for formatting any data into a byte slice
type Formatter struct {
	format FormatType
}

// NewFormatter will create a new Formatter with the designated output format
func NewFormatter(format FormatType) *Formatter {
	return &Formatter{
		format: format,
	}
}

// Format the provided data into the output format designated for this Formatter
func (formatter *Formatter) Format(data interface{}) (buf []byte, err error) {
	if formatter.format == FormatJSON {
		buf, err = formatter.formatJSON(data)
	} else if formatter.format == FormatTSV {
		buf, err = formatter.formatTSV(data, false)
	} else if formatter.format == FormatTSVKey {
		buf, err = formatter.formatTSV(data, true)
	} else {
		err = fmt.Errorf("Unknown FormatType '%v+'. Options are FormatJSON=0, or FormatTSV=1", formatter.format)
	}
	return buf, err
}

func (formatter *Formatter) formatJSON(data interface{}) (buf []byte, err error) {
	return json.Marshal(data)
}

func (formatter *Formatter) formatTSV(data interface{}, keyed bool) (buf []byte, err error) {
	var fields map[string]interface{}
	fields, err = getAsMap(data)
	if err != nil {
		return buf, err
	}

	var values [][]byte
	var keys []string
	for key := range fields {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		var value []byte
		value, err = formatter.formatJSON(fields[key])
		if err != nil {
			break
		}
		if keyed {
			value = []byte(fmt.Sprintf("%s=%s", key, value))
		}
		values = append(values, value)
	}
	buf = bytes.Join(values, []byte{'\t'})
	return buf, err
}

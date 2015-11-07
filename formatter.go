package jsonstream

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
)

type FormatType int

const (
	FormatJSON FormatType = iota
	FormatTSV
	FormatTSVKey
)

type Formatter struct {
	format FormatType
}

func NewFormatter(format FormatType) *Formatter {
	return &Formatter{
		format: format,
	}
}

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

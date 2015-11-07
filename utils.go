package jsonstream

import (
	"fmt"
	"reflect"
)

func getAsMap(data interface{}) (fields map[string]interface{}, err error) {
	switch value := data.(type) {
	case map[string]interface{}:
		fields = value
	default:
		err = fmt.Errorf("Unexpected data type '%s', expected 'map[string]interface{}'.", reflect.TypeOf(value))
	}
	return fields, err
}

package session

import (
	"errors"
	"reflect"
)

func inputToMap(input any) (map[string]interface{}, error) {
	rv := reflect.ValueOf(input)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

	switch rv.Kind() {
	case reflect.Struct:
		// continue
	case reflect.Map:
		return rv.Interface().(map[string]interface{}), nil
	default:
		return nil, errors.New("input must be a struct or map")
	}

	rt := rv.Type()
	m := make(map[string]interface{}, rt.NumField())
	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)
		if f.PkgPath != "" {
			continue
		}

		m[f.Name] = rv.Field(i).Interface()
	}

	return m, nil
}

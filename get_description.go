package gotoschema

import (
	"encoding"
	"encoding/json"
	"errors"
	"reflect"
	"strings"
)

var (
	jsonMarshalerType = reflect.TypeOf((*json.Marshaler)(nil)).Elem()
	textMarshalerType = reflect.TypeOf((*encoding.TextMarshaler)(nil)).Elem()
)

func (gen *DocGenerator) genSchemaForType(rt reflect.Type, isRefable bool) (typeSchema, error) {

	switch rt.Kind() {
	case reflect.String:
		return typeSchema{
			"type": "string",
		}, nil
	case reflect.Bool:
		return typeSchema{
			"type": "boolean",
		}, nil
	case reflect.Float32:
		return typeSchema{
			"type":   "number",
			"format": "float",
		}, nil
	case reflect.Float64:
		return typeSchema{
			"type":   "number",
			"format": "double",
		}, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return typeSchema{
			"type":   "integer",
			"format": rt.String(),
		}, nil
	case // all number types without specific format
		reflect.Complex64, reflect.Complex128:
		return typeSchema{
			"type":   "number",
			"format": rt.String(),
		}, nil

	case reflect.Array, reflect.Slice:
		desc, err := gen.genSchemaForType(rt.Elem(), true)
		if err != nil {
			return nil, err
		}
		return typeSchema{
			"type":  "array",
			"items": desc,
		}, nil

	case reflect.Struct:

		if rt.Implements(textMarshalerType) {
			return typeSchema{
				"type": "string",
			}, nil
		}

		props := map[string]interface{}{}
		for i := 0; i < rt.NumField(); i++ {
			field := rt.Field(i)
			tags := strings.Split(field.Tag.Get("json"), ",")
			name := tags[0]
			if name == "-" {
				continue
			}
			if name == "" {
				name = field.Name
			}

			fieldDesc, err := gen.genSchemaForType(field.Type, true)
			if err != nil {
				return nil, err
			}
			fieldDesc["title"] = field.Name
			props[name] = fieldDesc
		}

		schema := typeSchema{
			"type":       "object",
			"properties": props,
		}
		gen.addToCache(rt, schema)
		if isRefable {
			return typeSchema{
				"ref": "#/" + rt.Name(),
			}, nil
		} else {
			return schema, nil
		}

	case reflect.Pointer:
		elem := rt.Elem()
		desc, err := gen.genSchemaForType(elem, true)
		if err != nil {
			return nil, err
		}
		desc["nullable"] = true

		return desc, nil
	}

	return nil, errors.New("unsupported type: " + rt.String())
}

func (gen *DocGenerator) addToCache(t reflect.Type, s typeSchema) {
	if e, ok := gen.typeSchemaCache[t]; ok {
		e.RefCount += 1
		gen.typeSchemaCache[t] = e
	} else {
		gen.typeSchemaCache[t] = schemaCacheEntry{
			Name:     t.Name(),
			RefCount: 1,
			Value:    s,
		}
	}
}

func (gen *DocGenerator) getTypeSchemaFromCache(t reflect.Type) (typeSchema, bool) {
	if e, ok := gen.typeSchemaCache[t]; ok {
		e.RefCount += 1
		gen.typeSchemaCache[t] = e
		return e.Value, true
	}
	return nil, false
}

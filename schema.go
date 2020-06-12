package changeset

import (
	"reflect"
	"strings"
	"sync"

	"github.com/azer/snakecase"
)

var (
	fieldsCache       sync.Map
	fieldMappingCache sync.Map
	typesCache        sync.Map
)

type fields interface {
	Fields() map[string]int
}

func inferFields(record interface{}) map[string]int {
	if s, ok := record.(fields); ok {
		return s.Fields()
	}

	rt := reflectInternalType(record)
	// rt := reflectTypePtr(record)

	// check for cache
	if v, cached := fieldsCache.Load((rt)); cached {
		return v.(map[string]int)
	}

	var (
		index  = 0
		fields = make(map[string]int, rt.NumField())
	)

	for i := 0; i < rt.NumField(); i++ {
		var (
			sf   = rt.Field(i)
			name = inferFieldName(sf)
		)

		if name != "" {
			fields[name] = index
			index++
		}
	}

	fieldsCache.Store(rt, fields)

	return fields
}

func inferFieldName(sf reflect.StructField) string {
	if tag := sf.Tag.Get("db"); tag != "" {
		name := strings.Split(tag, ",")[0]

		if name == "-" {
			return ""
		}

		if name != "" {
			return name
		}
	}

	return snakecase.SnakeCase(sf.Name)
}

func inferFieldMapping(record interface{}) map[string]int {
	rt := reflectTypePtr(record)

	// check for cache
	if v, cached := fieldMappingCache.Load((rt)); cached {
		return v.(map[string]int)
	}

	mapping := make(map[string]int, rt.NumField())

	for i := 0; i < rt.NumField(); i++ {
		var (
			sf   = rt.Field(i)
			name = inferFieldName(sf)
		)

		if name != "" {
			mapping[name] = i
		}
	}

	fieldMappingCache.Store(rt, mapping)

	return mapping
}

type types interface {
	fields
	Types() []reflect.Type
}

func inferTypes(record interface{}) []reflect.Type {
	if v, ok := record.(types); ok {
		return v.Types()
	}

	rt := reflectTypePtr(record)

	// check for cache
	if v, cached := typesCache.Load(rt); cached {
		return v.([]reflect.Type)
	}

	var (
		fields  = inferFields(record)
		mapping = inferFieldMapping(record)
		types   = make([]reflect.Type, len(fields))
	)

	for name, index := range fields {
		var (
			structIndex = mapping[name]
			ft          = rt.Field(structIndex).Type
		)

		if ft.Kind() == reflect.Ptr {
			ft = ft.Elem()
		} else if ft.Kind() == reflect.Slice && ft.Elem().Kind() == reflect.Ptr {
			ft = reflect.SliceOf(ft.Elem().Elem())
		}

		types[index] = ft
	}

	typesCache.Store(rt, types)

	return types
}

type values interface {
	fields
	Values() []interface{}
}

func inferValues(record interface{}) []interface{} {
	if v, ok := record.(values); ok {
		return v.Values()
	}

	var (
		rv      = reflectValuePtr(record)
		fields  = inferFields(record)
		mapping = inferFieldMapping(record)
		values  = make([]interface{}, len(fields))
	)

	for name, index := range fields {
		var (
			structIndex = mapping[name]
			fv          = rv.Field(structIndex)
			ft          = fv.Type()
		)

		if ft.Kind() == reflect.Ptr {
			if !fv.IsNil() {
				values[index] = fv.Elem().Interface()
			}
		} else {
			values[index] = fv.Interface()
		}
	}

	return values
}

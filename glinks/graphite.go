package glinks

import (
	"fmt"
	"reflect"
)

type Metric struct {
	Path      string
	Value     interface{}
	Timestamp int64
}

func (m *Metric) String() string {
	return fmt.Sprintf("%s %v %d", m.Path, m.Value, m.Timestamp)
}

func FromStruct(x Data, prefix string) []Metric {

	metrics := make([]Metric, 0)
	value := reflect.ValueOf(x)

	return destruct(metrics, value, prefix, x.SampleTime())
}

func destruct(metrics []Metric, value reflect.Value, prefix string, timestamp int64) []Metric {
	switch value.Kind() {
	case reflect.Struct:
		for i := 0; i < value.NumField(); i++ {
			fieldName := value.Type().Field(i).Name
			// Ignore the Time field
			if fieldName == "Time" {
				continue
			}
			newPath := fmt.Sprintf("%s.%s", prefix, fieldName)
			metrics = destruct(metrics, value.Field(i), newPath, timestamp)
		}
	case reflect.Array, reflect.Slice:
		for i := 0; i < value.Len(); i++ {
			newPath := fmt.Sprintf("%s.%d", prefix, i)
			switch value.Index(i).Kind() {
			case reflect.Struct, reflect.Array, reflect.Slice, reflect.Map:
				metrics = destruct(metrics, value.Index(i), newPath, timestamp)
			default:
				metrics = append(metrics, Metric{Path: newPath,
					Value:     devalue(value.Index(i)),
					Timestamp: timestamp})
			}
		}
	case reflect.Map:
		for _, key := range value.MapKeys() {
			newPath := fmt.Sprintf("%s.%s", prefix, key)
			metrics = append(metrics, Metric{Path: newPath,
				Value:     devalue(value.MapIndex(key)),
				Timestamp: timestamp})
		}
	default:
		metrics = append(metrics, Metric{Path: prefix, Value: devalue(value), Timestamp: timestamp})
	}
	return metrics
}

func devalue(value reflect.Value) interface{} {
	switch value.Kind() {
	case reflect.String:
		return value.String()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return value.Int()
	case reflect.Float32, reflect.Float64:
		return value.Float()
	default:
		panic(fmt.Sprintf("Unhandled type: %s", value.Kind()))
	}
	return 0
}

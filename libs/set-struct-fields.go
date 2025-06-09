package libs

import (
	"reflect"
	"strings"
)

func SetStructFieldsFromMap(s interface{}, updates map[string]interface{}) {
	v := reflect.ValueOf(s).Elem()
	for key, value := range updates {
		field := v.FieldByNameFunc(func(n string) bool {
			return strings.EqualFold(n, key)
		})
		if field.IsValid() && field.CanSet() {
			val := reflect.ValueOf(value)
			if val.Type().ConvertibleTo(field.Type()) {
				field.Set(val.Convert(field.Type()))
			}
		}
	}
}

package utils

import (
	"errors"
	"reflect"
)

func MapFields(source interface{}, dest interface{}) error {
	sourceVal := reflect.ValueOf(source)
	destVal := reflect.ValueOf(dest)

	if sourceVal.Kind() != reflect.Struct || destVal.Kind() != reflect.Ptr || destVal.Elem().Kind() != reflect.Struct {
		return errors.New("source must be a struct, dest must be a pointer to a struct")
	}

	sourceType := sourceVal.Type()
	destType := destVal.Elem().Type()

	for i := 0; i < sourceType.NumField(); i++ {
		sourceField := sourceType.Field(i)
		destField, ok := destType.FieldByName(sourceField.Name)

		if ok {
			destVal.Elem().FieldByName(destField.Name).Set(sourceVal.Field(i))
		}
	}

	return nil
}

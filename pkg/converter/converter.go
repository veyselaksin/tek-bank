package converter

import (
	"errors"
	"reflect"
)

func CopyStruct(src any, dest any) error {
	srcType := reflect.TypeOf(src)
	srcValue := reflect.ValueOf(src)
	destType := reflect.TypeOf(dest)
	destValue := reflect.ValueOf(dest)

	if srcType.Kind() != reflect.Struct || destType.Kind() != reflect.Struct {
		return errors.New("src and dest must be a struct")
	}

	for i := 0; i < srcType.NumField(); i++ {
		srcField := srcType.Field(i)
		srcFieldValue := srcValue.Field(i)

		for j := 0; j < destType.NumField(); j++ {
			destField := destType.Field(j)
			destFieldValue := destValue.Field(j)

			if srcField.Name == destField.Name && srcField.Type == destField.Type {
				destFieldValue.Set(srcFieldValue)
			}
		}
	}

	return nil
}

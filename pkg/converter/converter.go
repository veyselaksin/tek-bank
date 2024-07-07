package converter

import (
	"encoding/json"
	"errors"
	"reflect"
)

func CopyStruct(src any, dest any) error {
	jsonItem, err := json.Marshal(src)
	if err != nil {
		return err
	}

	err = json.Unmarshal(jsonItem, dest)
	if err != nil {
		return err
	}

	return nil
}

// Stos is a function that converts any struct to a string
func Stos(item any) (*string, error) {
	// Check if the item is a struct
	if reflect.TypeOf(item).Kind() != reflect.Struct {
		return nil, errors.New("Item is not a struct")
	}
	jsonItem, err := json.Marshal(item)
	if err != nil {
		return nil, err
	}

	itemString := string(jsonItem)
	return &itemString, nil
}

// Stom is a function that converts a string to a struct
func Stom(item string, dest any) error {
	err := json.Unmarshal([]byte(item), dest)
	if err != nil {
		return err
	}

	return nil
}

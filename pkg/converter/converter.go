package converter

import (
	"encoding/json"
	"errors"
	"reflect"
)

//go:generate mockgen -destination=../../mocks/converter/converter_mock.go -package=converter tek-bank/pkg/converter Converter
type Converter interface {
	CopyStruct(src any, dest any) error
	Stos(item any) (*string, error)
	Stom(item string, dest any) error
}

type converter struct{}

func NewConverter() Converter {
	return &converter{}
}

func (c *converter) CopyStruct(src any, dest any) error {
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
func (c *converter) Stos(item any) (*string, error) {
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
func (c *converter) Stom(item string, dest any) error {
	err := json.Unmarshal([]byte(item), dest)
	if err != nil {
		return err
	}

	return nil
}

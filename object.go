package ioc

import (
	"reflect"
)

type Object struct {
	Name  string
	Value interface{}
	Type  reflect.Type
}

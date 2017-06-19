package configuration

import (
	"fmt"
	"reflect"
)

type value struct {
	typ   reflect.Type
	value interface{}
}

func (v *value) String() string {
	return fmt.Sprintf("%v", v.value)
}

func (v *value) Set(s string) error {
	fmt.Println("SETTING", s)
	return nil
}

func (v *value) Type() string {
	return v.typ.String()
}

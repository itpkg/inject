package ioc

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
)

type Injector struct {
	values []*Object
}

func (p *Injector) GetByName(name string) interface{} {
	for _, o := range p.values {
		if o.Name == name {
			return o.Value
		}
	}
	return nil
}

func (p *Injector) GetByType(rty reflect.Type) interface{} {
	for _, o := range p.values {
		if o.Name == "" && o.Type == rty {
			return o.Value
		}
	}
	return nil
}

func (p *Injector) Provide(objects ...*Object) {
	for _, o := range objects {
		o.Type = reflect.TypeOf(o.Value)
	}
	p.values = append(p.values, objects...)
}

func (p *Injector) String() string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "ID\tNAME\t\tTYPE\t\tVALUE\n")
	for i, o := range p.values {
		fmt.Fprintf(&buf, "%04d: %s\t\t%v\t\t%v\n", i, o.Name, o.Type, o.Value)
	}
	return buf.String()
}

func (p *Injector) Populate() error {
	for _, o := range p.values {
		if o.Type != nil && p.isStruct(o.Type) {
			el := reflect.ValueOf(o.Value).Elem()
			for i := 0; i < el.NumField(); i++ {
				fd := el.Field(i)
				tag := o.Type.Elem().Field(i)

				if tag.Tag == "" {
					continue
				}
				if !fd.CanSet() {
					return errors.New(fmt.Sprintf("inject requested on unexported field %s in type %v", tag.Name, o.Type))
				}
				if !p.isNilOrZero(fd) {
					continue
				}

				var val interface{}
				if name := tag.Tag.Get("inject"); name == "" {
					val = p.GetByType(fd.Type())
				} else {
					val = p.GetByName(name)
				}
				if val == nil {
					return errors.New(fmt.Sprintf("%s is null", tag.Tag))
				} else {
					fd.Set(reflect.ValueOf(val))
				}
			}

		}
	}
	return nil
}

func (*Injector) isStruct(t reflect.Type) bool {
	return t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Struct
}

func (*Injector) isNilOrZero(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	default:
		return reflect.DeepEqual(v.Interface(), reflect.Zero(v.Type()).Interface())
	}
}

func (p *Injector) Run(handler interface{}, args ...interface{}) ([]interface{}, error) {
	rty := reflect.TypeOf(handler)
	if rty.Kind() != reflect.Func {
		return nil, errors.New("Handler must be a callable func.")
	}
	ins := make([]reflect.Value, 0)
	for i := 0; i < rty.NumIn(); i++ {
		ty := rty.In(i)
		val := p.GetByType(ty)
		if val == nil {
			for _, arg := range args {
				if reflect.TypeOf(arg) == ty {
					val = arg
				}
			}
		}
		if val == nil {
			return nil, errors.New(fmt.Sprintf("Null arg []%v]", ty))
		}

		ins = append(ins, reflect.ValueOf(val))
	}

	ret := make([]interface{}, 0)
	for _, vl := range reflect.ValueOf(handler).Call(ins) {
		ret = append(ret, vl.Interface())
	}
	return ret, nil
}

//-----------------------------------------------------------------------------
func New() *Injector {
	return &Injector{values: make([]*Object, 0)}
}

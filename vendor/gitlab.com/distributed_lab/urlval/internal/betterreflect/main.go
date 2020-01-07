package betterreflect

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cast"
)

type Struct struct {
	src reflect.Value
}

func NewStruct(src interface{}) *Struct {
	return &Struct{
		src: reflect.ValueOf(src),
	}
}

func NewStructFromValue(src reflect.Value) *Struct {
	return &Struct{
		src: src,
	}
}

func (s *Struct) srcvalue() reflect.Value {
	return reflect.Indirect(s.src)
}

func (s *Struct) srctype() reflect.Type {
	return s.srcvalue().Type()
}

func (s *Struct) NumField() int {
	return s.srcvalue().Type().NumField()
}

func (s *Struct) Tag(i int, key string) string {
	return s.srctype().Field(i).Tag.Get(key)
}

func (s *Struct) Type(i int) reflect.Type {
	return s.srctype().Field(i).Type
}

func (s *Struct) Value(i int) reflect.Value {
	return s.srcvalue().Field(i)
}

// ErrInvalidType is returned when Set has failed to cast
// value to destination type.
var ErrInvalidType = errors.New("failed to cast value")

func (s *Struct) Set(i int, value interface{}) (err error) {
	t := s.Type(i)
	kind := t.Kind()

	if kind == reflect.Ptr {
		t = t.Elem()
		kind = t.Kind()
	}

	if IsZero(value) {
		return nil
	}

	switch kind {
	case reflect.String:
	case reflect.Bool:
		value, err = cast.ToBoolE(value)
	case reflect.Int:
		value, err = cast.ToIntE(value)
	case reflect.Int8:
		value, err = cast.ToInt8E(value)
	case reflect.Int16:
		value, err = cast.ToInt16E(value)
	case reflect.Int32:
		value, err = cast.ToInt32E(value)
	case reflect.Int64:
		value, err = cast.ToInt64E(value)
	case reflect.Uint:
		value, err = cast.ToUintE(value)
	case reflect.Uint8:
		value, err = cast.ToUint8E(value)
	case reflect.Uint16:
		value, err = cast.ToUint16E(value)
	case reflect.Uint32:
		value, err = cast.ToUint32E(value)
	case reflect.Uint64:
		value, err = cast.ToUint64E(value)
	case reflect.Float32:
		value, err = cast.ToFloat32E(value)
	case reflect.Float64:
		value, err = cast.ToFloat64E(value)
	case reflect.Slice:
		value, err = ConvertSlice(value, t.Elem())
	case reflect.Complex64, reflect.Complex128, reflect.Chan,
		reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr,
		reflect.Array, reflect.Struct, reflect.UnsafePointer:
	default:
		panic(fmt.Sprintf("unknown field kind: %v", kind))
	}
	if err != nil {
		return errors.Wrapf(ErrInvalidType, "expected value to be %s", kind)
	}

	if s.Type(i).Kind() == reflect.Ptr {
		zero := reflect.New(s.Type(i).Elem())
		s.Value(i).Set(zero)
		setValue(s.Value(i).Elem(), value)
	} else {
		setValue(s.Value(i), value)
	}

	return nil
}

// ConvertToUnderlyingType takes a value ot custom type T and converts it to
// it's underlying types. If value is of built-in type, ConvertToUnderlyingType
// just returns value as provided.
func ConvertToUnderlyingType(value interface{}) interface{} {
	var rType reflect.Type
	var rValue = reflect.ValueOf(value)
	var rKind = rValue.Kind()

	if rKind == reflect.Ptr {
		if !rValue.Elem().IsValid() {
			return nil
		}

		rValue = rValue.Elem()
		rKind = rValue.Kind()
	}

	switch rKind {
	case reflect.String:
		rType = reflect.TypeOf("")
	case reflect.Bool:
		rType = reflect.TypeOf(false)
	case reflect.Int:
		rType = reflect.TypeOf(0)
	case reflect.Int8:
		rType = reflect.TypeOf(int8(0))
	case reflect.Int16:
		rType = reflect.TypeOf(int16(0))
	case reflect.Int32:
		rType = reflect.TypeOf(int32(0))
	case reflect.Int64:
		rType = reflect.TypeOf(int64(0))
	case reflect.Uint:
		rType = reflect.TypeOf(uint(0))
	case reflect.Uint8:
		rType = reflect.TypeOf(uint8(0))
	case reflect.Uint16:
		rType = reflect.TypeOf(uint16(0))
	case reflect.Uint32:
		rType = reflect.TypeOf(uint32(0))
	case reflect.Uint64:
		rType = reflect.TypeOf(uint64(0))
	case reflect.Float32:
		rType = reflect.TypeOf(uint32(0))
	case reflect.Float64:
		rType = reflect.TypeOf(uint64(0))
	case reflect.Slice, reflect.Complex64, reflect.Complex128, reflect.Chan,
		reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr,
		reflect.Array, reflect.Struct, reflect.UnsafePointer:
		panic(errors.Errorf("got %s, when ConvertToUnderlyingType works only with primitive types", rValue.Type()))
	default:
		panic(fmt.Sprintf("unknown rKind: %v", rKind))
	}

	return rValue.Convert(rType).Interface()
}

func setValue(dest reflect.Value, v interface{}) {
	// to assign type to their aliases:
	if reflect.TypeOf(v) != dest.Type() && reflect.TypeOf(v).Kind() == dest.Type().Kind() {
		dest.Set(reflect.ValueOf(v).Convert(dest.Type()))
		return
	}

	dest.Set(reflect.ValueOf(v))
}

func IsZero(value interface{}) bool {
	switch reflect.TypeOf(value).Kind() {
	case reflect.Ptr:
		ptr := reflect.ValueOf(value)
		return ptr.IsNil() || ptr.Elem() == reflect.Zero(reflect.TypeOf(value))
	case reflect.Slice:
		slice := reflect.ValueOf(value)
		return slice.Len() == 0
	case reflect.Complex64, reflect.Complex128, reflect.Chan,
		reflect.Func, reflect.Interface, reflect.Map,
		reflect.Array, reflect.Struct, reflect.UnsafePointer:
		return false
	default:
		return value == reflect.Zero(reflect.TypeOf(value)).Interface()
	}
}

// ConvertSlice takes a slice of any value and converts it's items to the "destElemType".
// Returns error, if types are not convertible to each other. Returns error if provided
// rawValue is not slice.
func ConvertSlice(rawValue interface{}, destElemType reflect.Type) (interface{}, error) {
	switch kind := reflect.TypeOf(rawValue).Kind(); kind {
	case reflect.Slice:
	case reflect.String:
		rawValue = strings.Split(cast.ToString(rawValue), ",")
	default:
		return nil, errors.Errorf("expected value to be slice or string, but got %s", kind)
	}

	sourceSlice := reflect.ValueOf(rawValue)

	destSliceType := reflect.SliceOf(destElemType)
	destSlice := reflect.MakeSlice(destSliceType, 0, sourceSlice.Len())

	for j := 0; j < sourceSlice.Len(); j++ {
		sourceElem := sourceSlice.Index(j)
		sourceElemType := sourceElem.Type()

		if !sourceElemType.ConvertibleTo(destElemType) {
			return nil, errors.Errorf("%s is not convertible to %s", sourceElemType.String(), destElemType.String())
		}

		convertedValue := sourceElem.Convert(destElemType)
		destSlice = reflect.Append(destSlice, convertedValue)
	}

	return destSlice.Interface(), nil
}

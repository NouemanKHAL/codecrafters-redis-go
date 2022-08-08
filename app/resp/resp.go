package resp

import (
	"bufio"
	"fmt"
)

type Type rune

const (
	SIMPLE_STRING Type = '+'
	INTEGER       Type = ':'
	BULK_STRING   Type = '$'
	ARRAY         Type = '*'
	ERROR         Type = '-'
)

type Value struct {
	typ   Type
	data  []byte
	array []Value
}

func NewValue(b []byte, v []Value, t Type) (Value, error) {
	if v != nil || t == ARRAY {
		return Value{
			typ:   t,
			array: v,
		}, nil
	}
	return Value{
		typ:  t,
		data: b,
	}, nil
}

func NewSimpleStringValue(s string) Value {
	v, _ := NewValue([]byte(s), nil, SIMPLE_STRING)
	return v
}

func NewBulkStringValue(s string) Value {
	v, _ := NewValue([]byte(s), nil, BULK_STRING)
	return v
}

func NewIntegerValue(i int) Value {
	v, _ := NewValue([]byte(string(i)), nil, INTEGER)
	return v
}

func NewErrorValue(err string) Value {
	v, _ := NewValue([]byte(err), nil, ERROR)
	return v
}

func NewArrayValue(arr []Value) Value {
	v, _ := NewValue(nil, arr, ARRAY)
	return v
}

func (v *Value) Type() Type {
	return v.typ
}

func (v *Value) Data() []byte {
	return v.data
}

func (v *Value) Array() []Value {
	if v.typ == ARRAY {
		return v.array
	}
	return []Value{}
}

func (v *Value) String() string {
	return string(v.data)
}

func (v Value) Encode() ([]byte, error) {
	switch v.typ {
	case INTEGER:
		return encodeInteger(v)
	case SIMPLE_STRING:
		return encodeSimpleString(v), nil
	case BULK_STRING:
		return encodeBulkString(v), nil
	case ERROR:
		return encodeError(v), nil
	case ARRAY:
		return encodeArray(v)
	}
	return []byte{}, fmt.Errorf("Encode was given an unsupported type")
}

func Decode(byteStream *bufio.Reader) (Value, error) {
	b, err := byteStream.ReadByte()
	if err != nil {
		return Value{}, err
	}
	switch Type(b) {
	case SIMPLE_STRING:
		return decodeSimpleString(byteStream)
	case BULK_STRING:
		return decodeBulkString(byteStream)
	case ARRAY:
		return decodeArray(byteStream)
	}
	return Value{}, fmt.Errorf("Decode was given an unsupported data type")
}

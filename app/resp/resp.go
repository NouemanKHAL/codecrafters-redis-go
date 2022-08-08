package resp

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"strconv"
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

func NewValue(b []byte, t Type) (Value, error) {
	if t == ARRAY {
		byteStream := bufio.NewReader(bytes.NewReader(b))
		val, err := decodeArray(byteStream)
		if err != nil {
			return Value{}, err
		}
		return val, nil
	}
	return Value{
		typ:  t,
		data: b,
	}, nil
}

func (v *Value) Array() []Value {
	if v.typ == ARRAY {
		return v.array
	}
	return []Value{}
}

func (v *Value) Bytes() []byte {
	return v.data
}

func (v *Value) String() string {
	return string(v.data)
}
func (v *Value) Type() Type {
	return v.typ
}

func encode(v *Value) ([]byte, error) {
	switch v.typ {
	case INTEGER:
		return []byte(fmt.Sprintf(":%d\r\n", v.data)), nil
	case SIMPLE_STRING:
		return []byte(fmt.Sprintf("+%s\r\n", v.data)), nil
	case BULK_STRING:
		return []byte(fmt.Sprintf("$%d\r\n%s\r\n", len(v.data), v.data)), nil
	case ERROR:
		return []byte(fmt.Sprintf("-%s\r\n", v.data)), nil
	case ARRAY:
		res := []byte{}
		for _, elem := range v.array {
			val, err := encode(&elem)
			if err != nil {
				log.Println(err)
				continue
			}
			res = append(res, val...)
		}
		return []byte(fmt.Sprintf("*%d\r\n%s\r\n", len(res), res)), nil
	}
	return []byte{}, fmt.Errorf("unknown type given to encode")
}

func encodeSimpleString(s string) []byte {
	v := &Value{
		typ:  SIMPLE_STRING,
		data: []byte(s),
	}
	bytes, _ := encode(v)
	return bytes
}

func encodeBulkString(s string) []byte {
	v := &Value{
		typ:  BULK_STRING,
		data: []byte(s),
	}
	bytes, _ := encode(v)
	return bytes
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
	return Value{}, fmt.Errorf("unknown data type")
}

func readToken(byteStream *bufio.Reader) ([]byte, error) {
	bytes, err := byteStream.ReadBytes('\n')
	if err != nil {
		return nil, err
	}
	// discard \r\n
	return bytes[:len(bytes)-2], nil
}

func decodeSimpleString(byteStream *bufio.Reader) (Value, error) {
	t, err := readToken(byteStream)
	if err != nil {
		return Value{}, err
	}
	return Value{
		typ:  SIMPLE_STRING,
		data: t,
	}, nil
}

func decodeBulkString(byteStream *bufio.Reader) (Value, error) {
	t, err := readToken(byteStream)
	if err != nil {
		return Value{}, nil
	}

	size, err := strconv.Atoi(string(t))
	if err != nil {
		return Value{}, err
	}

	str := make([]byte, size+2)

	_, err = io.ReadFull(byteStream, str)
	if err != nil {
		return Value{}, err
	}

	// discard \r\n
	str = str[:size]

	return Value{
		typ:  BULK_STRING,
		data: str,
	}, nil
}

func decodeArray(byteStream *bufio.Reader) (Value, error) {
	t, err := readToken(byteStream)
	if err != nil {
		return Value{}, nil
	}
	length, err := strconv.Atoi(string(t))
	if err != nil {
		return Value{}, err
	}

	arr := make([]Value, length)
	for i := 0; i < len(arr); i++ {
		v, err := Decode(byteStream)
		if err != nil {
			return Value{}, err
		}
		arr[i] = v
	}

	return Value{
		typ:   ARRAY,
		array: arr,
	}, nil
}

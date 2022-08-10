package resp

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

func encodeInteger(v Value) ([]byte, error) {
	int_value, err := strconv.Atoi(string(v.data))
	if err != nil {
		return nil, err
	}
	return []byte(fmt.Sprintf(":%d\r\n", int_value)), nil
}

func encodeSimpleString(v Value) []byte {
	return []byte(fmt.Sprintf("+%s\r\n", v.data))
}

func encodeBulkString(v Value) []byte {
	return []byte(fmt.Sprintf("$%d\r\n%s\r\n", len(v.data), v.data))
}

func encodeError(v Value) []byte {
	return []byte(fmt.Sprintf("-%s\r\n", v.data))
}

func encodeArray(v Value) ([]byte, error) {
	res := []byte{}
	for _, elem := range v.array {
		val, err := elem.Encode()
		if err != nil {
			return nil, err
		}
		res = append(res, val...)
	}
	return []byte(fmt.Sprintf("*%d\r\n%s\r\n", len(res), res)), nil
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
	return NewSimpleStringValue(string(t)), nil
}

func decodeInteger(byteStream *bufio.Reader) (Value, error) {
	t, err := readToken(byteStream)
	if err != nil {
		return Value{}, nil
	}

	num, err := strconv.Atoi(string(t))
	if err != nil {
		return Value{}, err
	}
	return NewIntegerValue(num), nil
}

func decodeError(byteStream *bufio.Reader) (Value, error) {
	t, err := readToken(byteStream)
	if err != nil {
		return Value{}, err
	}
	return NewErrorValue(string(t)), nil
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

	return NewBulkStringValue(string(str)), nil
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

	return NewArrayValue(arr), nil
}

func SendError(err error) []byte {
	e := NewErrorValue("ERR - " + err.Error())
	return encodeError(e)
}

func SendNil() []byte {
	return []byte("$-1\r\n")
}

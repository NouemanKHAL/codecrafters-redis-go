package resp

func Ping() []byte {
	data, err := NewSimpleStringValue("PONG").Encode()
	if err != nil {
		e, _ := NewErrorValue(err.Error()).Encode()
		return e
	}
	return data
}

func Echo(v Value) []byte {
	data, err := v.Encode()
	if err != nil {
		e, _ := NewErrorValue(err.Error()).Encode()
		return e
	}
	return data
}

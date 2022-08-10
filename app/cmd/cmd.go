package cmd

import (
	"build-your-own-redis/app/resp"
	"fmt"
)

var store = make(map[string]resp.Value)

func Ping() []byte {
	data, err := resp.NewSimpleStringValue("PONG").Encode()
	if err != nil {
		return resp.SendError(err)
	}
	return data
}

func Echo(v resp.Value) []byte {
	data, err := v.Encode()
	if err != nil {
		return resp.SendError(err)
	}
	return data
}

func Set(k, v resp.Value) []byte {
	var response []byte
	var err error
	if old, ok := store[k.String()]; ok {
		response, err = resp.NewBulkStringValue(old.String()).Encode()
		if err != nil {
			return resp.SendError(err)
		}
	} else {
		response, err = resp.NewSimpleStringValue("OK").Encode()
		if err != nil {
			return resp.SendError(err)
		}
	}
	if v.Type() != resp.BULK_STRING {
		v = resp.NewBulkStringValue(v.String())
	}
	store[k.String()] = v
	return response
}

func Get(k resp.Value) []byte {
	if data, ok := store[k.String()]; ok {
		bytes, err := data.Encode()
		if err != nil {
			return resp.SendError(err)
		}
		return bytes
	}
	return resp.SendError(fmt.Errorf("key not found"))
}

package cmd

import (
	"build-your-own-redis/app/resp"
	"strings"
	"time"
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

func Set(args []resp.Value) []byte {
	var k, v resp.Value
	var opt string
	var expiry int
	var err error

	if len(args) >= 2 {
		k = args[0]
		v = args[1]
	}

	if len(args) == 4 {
		opt = args[2].String()
		expiry, err = args[3].Integer()
		if err != nil {
			return resp.SendError(err)
		}
	}

	var response []byte

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

	if strings.ToUpper(opt) == "PX" && expiry > 0 {
		go func() {
			ch := time.After(time.Duration(expiry) * time.Millisecond)
			for {
				select {
				case <-ch:
					delete(store, k.String())
				}
			}
		}()
	}
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
	return resp.SendNil()
}

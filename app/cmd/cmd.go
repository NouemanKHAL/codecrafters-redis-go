package cmd

import "build-your-own-redis/app/resp"

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

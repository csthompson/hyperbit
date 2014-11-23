package main

import (
	redis "github.com/xuyu/goredis"
	_ "log"
)

func bitwiseOp(op string, index1 string, index2 string) string {
	//Connect to server (TESTING)
	client, _ := redis.Dial(&redis.DialConfig{Address: "127.0.0.1:6379"})

	client.BitOp(op, index1+op+index2, index1, index2)

	return (index1 + op + index2)
}

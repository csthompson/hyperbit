package main

import (
  redis "github.com/xuyu/goredis"
  "log"
)

// Perform a bitwise AND operation

func main() {
  client, _ := redis.Dial(&redis.DialConfig{Address: "127.0.0.1:6379"})

  client.SetBit("test", 123, 1)

  log.Println(client.GetBit("test", 123))
}

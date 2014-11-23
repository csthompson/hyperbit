package main

import (
	_ "fmt"
	_ "github.com/xuyu/goredis"
)

// Perform a bitwise AND operation
//
func Intersection() {

}

func main() {
	//client, _ := redis.Dial(&redis.DialConfig{Address: "127.0.0.1:6379"})

	//client.SetBit("test", 123, 1)

	//log.Println(client.GetBit("test", 123))

	Load()
	Test()

	//fmt.Println("infix:  ", input)
	//fmt.Println("postfix:", ParseInfix(input))
}

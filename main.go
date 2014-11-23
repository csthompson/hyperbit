package main

import (
	_ "fmt"
  "log"
	redis "github.com/xuyu/goredis"
)

// Perform a bitwise AND operation
//
func Intersection() {

}

func main() {
	client, _ := redis.Dial(&redis.DialConfig{Address: "127.0.0.1:6379"})

	client.SetBit("test", 123, 1)

	//log.Println(client.GetBit("test", 123))

for i :=0; i < 15; i ++ {
  client.SetBit("users:inactive:201401", i, 1)
  client.SetBit("users:Active:201401", i, 1)
  client.SetBit("users:Testing:2014021505", i + 1, 1)
  client.SetBit("users:Testing:2014011210", i , 1)
  client.SetBit("users:Testing:2014021505", i +2, 1)
  client.SetBit("users:Testing:2014011210", i , 1)
  client.SetBit("users:Testing:2014021505", i + 1, 1)
  client.SetBit("users:Testing:2014011210", i, 1)
}


	Load()
	Test()

  rp, _ := client.ExecuteCommand("BITCOUNT", "users:Testing:2014021505ANDusers:Testing:2014011210ORusers:inactive:201401ANDusers:Active:201401")
  log.Println(rp.IntegerValue())
	//fmt.Println("infix:  ", input)
	//fmt.Println("postfix:", ParseInfix(input))
}

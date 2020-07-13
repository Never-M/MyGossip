package main

import (
	"fmt"
	. "github.com/Never-M/MyGossip/pkg/gossiper"
)

func main()  {
	node1 := NewGossiper("node1", "localhost")
	node2 := NewGossiper("node2", "localhost")
	node1.AddPeer(node2)
	go node2.HeartBeatReceiver()
	var input string
	for  {
		fmt.Scanln(&input)
		if input == "1" {
			rc := node1.SendHeartBeat()
			if rc != 0 {
				fmt.Println("Send req error")
			}
		} else {
			break
		}
	}
}
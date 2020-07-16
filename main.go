package main

import (
	"fmt"

	"github/Never-M/MyGossip/pkg/gossiper"
)

func main()  {
	node1 := gossiper.NewGossiper("node1", "localhost")
	node2 := gossiper.NewGossiper("node2", "localhost")
	node1.AddPeer(gossiper.NewPeer("node2","localhost"))
	go node2.HeartBeatReceiver()
	var input string
	for  {
		fmt.Scanln(&input)
		if input == "1" {
			rc, err := node1.SendHeartBeat()
			if err != nil {
				fmt.Println("Send req error, %v",err)
			}
			fmt.Println(rc)
		} else {
			break
		}
	}
}
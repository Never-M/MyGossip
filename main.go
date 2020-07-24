package main

import (
	"fmt"
	gp "github.com/Never-M/MyGossip/pkg/gossiper"
)

func main()  {
	logger := gp.Newlogger()
	logger.Fatal("1")
	node1 := gp.NewGossiper("node1", "localhost")
	node2 := gp.NewGossiper("node2", "localhost")
	node1.AddPeer(gp.NewPeer("node2","localhost"))
	go node2.HeartBeatReceiver()
	var input string
	for  {
		fmt.Scanln(&input)
		if input == "1" {
			_, err := node1.SendHeartBeats()
			if err != nil {
				fmt.Printf("Send req error, %v",err)
			}
		} else {
			break
		}
	}
}
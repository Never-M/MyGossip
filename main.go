package main

import (
	"fmt"
	gp "github.com/Never-M/MyGossip/pkg/gossiper"
)

func main() {
	// logger := gp.Newlogger()
	// logger.Fatal("1")
	// node1 := gp.NewGossiper("node1", "localhost")
	// node2 := gp.NewGossiper("node2", "localhost")
	// node1.AddPeer(gp.NewPeer("node2", "localhost"))
	// go node2.HeartBeatReceiver()
	// var input string
	// for {
	// 	fmt.Scanln(&input)
	// 	if input == "1" {
	// 		_, err := node1.SendHeartBeats()
	// 		if err != nil {start
	// 			logger.Error("Send req err")
	// 		}
	// 	} else {
	// 		break
	// 	}
	// }
	var begin, command, name, ip string
	//commands := []string{"start", "stop", "AddPeer", "RemovePeer", "PrintPeer"}
	for begin != "start" {
		fmt.Println("Enter start to begin")
		fmt.Scan(&begin)
		if begin == "start" {
			fmt.Println("Enter the name and ip")
			fmt.Scan(&name, &ip)
			g := gp.NewGossiper(name, ip)
			g.Start()
			for {
				fmt.Scan(&command)
				switch command {
				case "stop":
					g.Stop()
					break
				case "AddPeer":
					fmt.Println("Enter peer's name and ip to ADD")
					fmt.Scan(&name, &ip)
					g.AddPeer(gp.NewPeer(name, ip))
				case "RemovePeer":
					fmt.Println("Enter peer's name and ip to REMOVE")
					fmt.Scan(&name)
					g.RemovePeer(name)
				case "PrintPeer":
					g.PrintPeerNames()
				case "help":
					fmt.Println("start", "stop", "AddPeer", "RemovePeer", "PrintPeer")
				}
			}
		}
		break
	}
}

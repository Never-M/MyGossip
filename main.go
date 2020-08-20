package main

import (
	"fmt"
	gp "github.com/Never-M/MyGossip/pkg/gossiper"
)

func main() {
	var begin, command, name, ip string
	for begin != "start" {
		fmt.Println("Enter start to begin")
		fmt.Scan(&begin)
		if begin == "start" {
			fmt.Println("Enter node name")
			fmt.Scan(&name)
			fmt.Println("Enter node ip")
			fmt.Scan(&ip)
			g := gp.NewGossiper(name, ip)
			g.Start()
			for {
				fmt.Println()
				fmt.Print(g.GetName() + ": ")
				fmt.Scan(&command)
				switch command {
				case "stop":
					g.Stop()
					break
				case "add":
					fmt.Println("Enter peer name")
					fmt.Scan(&name)
					fmt.Println("Enter peer ip")
					fmt.Scan(&ip)
					g.AddPeer(gp.NewPeer(name, ip))
				case "remove":
					fmt.Println("Enter peer name to remove")
					fmt.Scan(&name)
					g.RemovePeer(name)
				case "show":
					g.PrintPeerNames()
				case "help":
					help()
				default:
					fmt.Println("Invalid input, try again")
					help()
				}
			}
		}
		break
	}
}

func help() {
	fmt.Println("Please use commands below:")
	fmt.Println("stop: Shut down current node")
	fmt.Println("add: Add a peer to current node")
	fmt.Println("remove: remove a peer from current node")
	fmt.Println("show: Print out peers of current node")
}
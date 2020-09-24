package main

import (
	"fmt"
	gp "github.com/Never-M/MyGossip/pkg/gossiper"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"time"
)

func main() {
	var command, name, ip string
	fmt.Println(">> Enter node name")
	fmt.Scan(&name)
	fmt.Println(">> Enter node ip")
	fmt.Scan(&ip)
	g := gp.NewGossiper(name, ip)
	g.Start()
	filename := filepath.Join("/tmp/", name, "/ibltlog")
	f, err := os.OpenFile(filename, os.O_WRONLY | os.O_CREATE, 0755)
	if err != nil {
		fmt.Println("Internel error! try again!")
	}
	logrus.SetOutput(f)
	for {
		fmt.Println()
		fmt.Print("[" + g.GetName() + "]" + ": ")
		fmt.Scan(&command)
		switch command {
		case "exit":
			g.Stop()
			fmt.Println(">> Bye")
			time.Sleep(time.Second)
			return
		case "add":
			fmt.Println(">> Enter peer name")
			fmt.Scan(&name)
			fmt.Println(">> Enter peer ip")
			fmt.Scan(&ip)
			g.AddPeer(gp.NewPeer(name, ip))
		case "remove":
			fmt.Println(">> Enter peer name to remove")
			fmt.Scan(&name)
			g.RemovePeer(name)
		case "show":
			g.PrintPeerNames()
		case "help":
			help()
		case "put":
			var key, value string
			fmt.Println(">> Enter key")
			fmt.Scan(&key)
			fmt.Println(">> Enter value")
			fmt.Scan(&value)
			g.Put(key, value)
		case "get":
			var key string
			fmt.Println(">> Enter key")
			fmt.Scan(&key)
			value := g.Get(key)
			if value == ""{
				fmt.Println(">> Key not found")
			} else {
				fmt.Printf(">> Value of key [%v]: %v\n", key, value)
			}
		case "delete":
			var key string
			fmt.Println(">> Enter key")
			fmt.Scan(&key)
			g.Delete(key)
			fmt.Println(">> Key [" + key + "] deleted!")
		case "list":
			g.GetDB().ListData()
		default:
			fmt.Println(">> Invalid input, try again")
			help()
		}
	}
}

func help() {
	fmt.Println(">> Please use commands below:")
	fmt.Println(">> exit: Shut down and exit current node")
	fmt.Println(">> add: Add a peer to current node")
	fmt.Println(">> remove: remove a peer from current node")
	fmt.Println(">> show: Print out peers of current node")
	fmt.Println(">> put: Put a key & value pair to the database")
	fmt.Println(">> get: get the value of a specific key")
	fmt.Println(">> delete: delete a key from the database")
}

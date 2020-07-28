package gossiper

import (
	"fmt"
	"github.com/Never-M/MyGossip/pkg/types"
	"github.com/stretchr/testify/assert"
	"testing"
	//"time"
)

func TestWriteRead(t *testing.T) {
	g := NewGossiper("node1", "localhost")
	//for running this test, change the names of the peers everytime
	resultCode := g.AddPeer(NewPeer("node2", "localhost"))
	assert.Equal(t, types.SUCCEED, resultCode, "Add peer failed")
	resultCode = g.AddPeer(NewPeer("node3", "localhost"))
	assert.Equal(t, types.SUCCEED, resultCode, "Add peer failed")
	g.WritePeersToFile()
	res := ReadPeersFromFile("node1")
	fmt.Println(res)
	g.db.Close()
}

// func TestAddRemovePeer(t *testing.T) {
// 	g := NewGossiper("node1", "localhost")
// 	resultCode := g.AddPeer(NewPeer("node2", "localhost"))
// 	assert.Equal(t, types.SUCCEED, resultCode, "Add peer failed")

// 	resultCode = g.RemovePeer("node2")
// 	assert.Equal(t, types.SUCCEED, resultCode, "Remove peer failed")
// 	//every test case has a Close, deleting it later
// 	g.db.Close()
// }

// func TestSendSingleHeartBeat(t *testing.T) {
// 	node1 := NewGossiper("node1", "localhost")
// 	node2 := NewGossiper("node2", "localhost")
// 	peer2 := NewPeerFromGossiper(node2)
// 	node1.AddPeer(peer2)
// 	go node2.HeartBeatReceiver()
// 	resultCode, err := node1.SendHeartBeat(peer2)
// 	assert.NoError(t, err)
// 	assert.Equal(t, types.SUCCEED, resultCode, "Send heartbeat failed")
// 	//every test case has a Close, deleting it later
// 	node1.db.Close()
// 	node2.db.Close()
// }

// func TestNodeStartStop(t *testing.T) {
// 	node1 := NewGossiper("node1", "localhost")
// 	node1.Start()
// 	time.Sleep(3000 * time.Millisecond)
// 	node1.Stop()
// 	node1.db.Close()
// }

// func TestGossiperLifeCycle(t *testing.T) {
// 	node1 := NewGossiper("node1", "localhost")
// 	node2 := NewGossiper("node2", "localhost")
// 	node1.Start()
// 	node2.Start()
// 	node1.AddPeer(NewPeerFromGossiper(node2))
// 	time.Sleep(3000 * time.Millisecond)
// 	node1.Stop()
// 	node2.Stop()
// 	node1.db.Close()
// 	node2.db.Close()
// }

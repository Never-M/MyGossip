package gossiper

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/Never-M/MyGossip/pkg/types"
)

func TestAddRemovePeer(t *testing.T)  {
	g := NewGossiper("node1", "localhost")
	resultCode := g.AddPeer(NewPeer("node2", "localhost"))
	assert.Equal(t,types.SUCCEED, resultCode, "Add peer failed")

	resultCode = g.RemovePeer("node2")
	assert.Equal(t,types.SUCCEED, resultCode, "Remove peer failed")
}

func TestSendSingleHeartBeat(t *testing.T)  {
	node1 := NewGossiper("node1", "localhost")
	node2 := NewGossiper("node2", "localhost")
	peer2 := NewPeerFromGossiper(node2)
	node1.AddPeer(peer2)
	go node2.HeartBeatReceiver()
	resultCode,err := node1.SendHeartBeat(peer2)
	assert.NoError(t,err)
	assert.Equal(t, types.SUCCEED, resultCode, "Send heartbeat failed")
}

func TestNodeStartStop(t *testing.T)  {
	node1 := NewGossiper("node1", "localhost")
	node1.Start()
	time.Sleep(3000 * time.Millisecond)
	node1.Stop()
}

func TestGossiperLifeCycle(t *testing.T)  {
	node1 := NewGossiper("node1", "localhost")
	node2 := NewGossiper("node2", "localhost")
	node1.Start()
	node2.Start()
	node1.AddPeer(NewPeerFromGossiper(node2))
	time.Sleep(3000 * time.Millisecond)
	node1.Stop()
	node2.Stop()
}
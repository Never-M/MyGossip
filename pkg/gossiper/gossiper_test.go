package gossiper

import (
	"testing"

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

func TestSendHeartBeat(t *testing.T)  {
	node1 := NewGossiper("node1", "localhost")
	node2 := NewGossiper("node2", "localhost")
	node1.AddPeer(NewPeerFromGossiper(node2))
	go node2.HeartBeatReceiver()
	resultCode,err := node1.SendHeartBeat()
	assert.NoError(t,err)
	assert.Equal(t, types.SUCCEED, resultCode, "Send heartbeat failed")
}

package gossiper

import (
	"github.com/Never-M/MyGossip/pkg/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAddRemovePeer(t *testing.T) {
	g := NewGossiper("node1", "localhost")
	resultCode := g.AddPeer(NewPeer("node2", "localhost"))
	assert.Equal(t, types.SUCCEED, resultCode, "Add peer failed")

	resultCode = g.RemovePeer("node2")
	assert.Equal(t, types.SUCCEED, resultCode, "Remove peer failed")
	//every test case has a Close, deleting it later
	g.db.Close()
}

func TestSendHeartBeat(t *testing.T) {
	node1 := NewGossiper("node1", "localhost")
	node2 := NewGossiper("node2", "localhost")
	node1.AddPeer(NewPeerFromGossiper(node2))
	go node2.HeartBeatReceiver()
	resultCode, err := node1.SendHeartBeat()
	assert.NoError(t, err)
	assert.Equal(t, types.SUCCEED, resultCode, "Send heartbeat failed")
	//every test case has a Close, deleting it later
	node1.db.Close()
	node2.db.Close()
}

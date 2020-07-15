package gossiper

import (
	. "github.com/Never-M/MyGossip/pkg/types"
	"github.com/stretchr/testify/assert"
	"testing"
)


func TestAddRemovePeer(t *testing.T)  {
	assert := assert.New(t)
	g := NewGossiper("node1", "localhost")
	resultCode := g.AddPeer(NewGossiper("node2", "localhost"))
	assert.Equal(SUCCEED, resultCode, "Add peer failed")

	resultCode = g.RemovePeer("node2")
	assert.Equal(SUCCEED, resultCode, "Remove peer failed")
}

func TestSendHeartBeat(t *testing.T)  {
	assert := assert.New(t)
	node1 := NewGossiper("node1", "localhost")
	node2 := NewGossiper("node2", "localhost")
	node1.AddPeer(node2)
	go node2.HeartBeatReceiver()
	resultCode := node1.SendHeartBeat()
	assert.Equal(SUCCEED, resultCode, "Send heartbeat failed")
}

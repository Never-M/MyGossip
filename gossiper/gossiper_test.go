package gossiper

import "testing"

func TestAddRemovePeer(t *testing.T)  {
	g := new(gossiper)
	g.name = "test"
	g.AddPeer(gossiper{"peer", "localhost", make(map[string]gossiper)})
	if _, ok := g.peers["peer"]; !ok {
		t.Errorf("Add Failed")
	}
	g.RemovePeer("peer")
	if _, ok := g.peers["peer"]; ok {
		t.Errorf("Remove Failed")
	}
}
package gossiper

import "fmt"

type gossiper struct {
	name	string
	ip		string
	peers	map[string]gossiper
}

func (g *gossiper) AddPeer(peer gossiper)  {
	g.peers[peer.name] = peer
}

func (g *gossiper) RemovePeer(name string)  {
	delete(g.peers, name)
}

func (g *gossiper) GetPeerNames() {
	fmt.Println(g.peers)
}
package gossiper

import "time"

type peer struct {
	name 	string
	ip 		string
	timer	*time.Timer
}

func NewPeer(name, ip string) *peer {
	return &peer{
		name:name,
		ip:ip,
		timer:time.NewTimer(2 * HEARTBEAT_TIMEOUT * time.Millisecond),
	}
}

func NewPeerFromGossiper(g *gossiper) *peer {
	return &peer{
		name:g.name,
		ip:g.ip,
		timer:time.NewTimer(2 * HEARTBEAT_TIMEOUT * time.Millisecond),
	}
}

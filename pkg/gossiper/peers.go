package gossiper

type peer struct {
	name 	string
	ip 		string
}

func NewPeer(name, ip string) *peer {
	return &peer{
		name:name,
		ip:ip,
	}
}
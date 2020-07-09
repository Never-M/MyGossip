package gossiper

import (
	"context"
	"fmt"
	. "github.com/Never-M/MyGossip/pkg/types"
	"log"
	"net/http"
	"time"
)

// Time formatter
var timeFormat = "2006-01-02 15:04:05"
const HEARTBEAT_PORT = ":9000"
const HEARTBEAT_CTXKEY = "heartbeat"

type gossiper struct {
	name	string
	ip		string
	peers	map[string]*gossiper
	client  *http.Client
}

type heartBeat struct {
	gossiperName	string
	currentTime		time.Time
}

func NewGossiper(name, address string) *gossiper {
	return &gossiper{
		name: name,
		ip:address,
		peers:make(map[string]*gossiper),
		client:&http.Client{},
	}
}

func (g *gossiper) AddPeer(peer *gossiper) int {
	if _, ok := g.peers[peer.name]; ok {
		return OPERATION_ERROR
	}
	g.peers[peer.name] = peer
	return SUCCEED
}

func (g *gossiper) RemovePeer(name string) int {
	delete(g.peers, name)
	return SUCCEED
}

func (g *gossiper) PrintPeerNames() {
	fmt.Println(g.peers)
}

func (g *gossiper) HeartBeatHandler(w http.ResponseWriter, r *http.Request) {
	hb := r.Context().Value(HEARTBEAT_CTXKEY)
	hbStr := fmt.Sprint(hb)
	log.Println(hbStr)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(hbStr))
}

// maybe the host doesn't need context
//func (g *gossiper) HeartBeatContext(next http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		ctx := context.WithValue(r.Context(), HEARTBEAT_CTXKEY, heartBeat{g.name, time.Now()})
//		next.ServeHTTP(w, r.WithContext(ctx))
//	})
//}

func (g *gossiper) SendHeartBeat() int {
	ctx := context.WithValue(context.Background(), "heartbeat", &heartBeat{g.name, time.Now()})
	for _, peer := range g.peers {
		req, err := http.NewRequestWithContext(ctx, "GET", "http://" + peer.ip + HEARTBEAT_PORT +"/" + HEARTBEAT_CTXKEY, nil)
		if err != nil {
			log.Println(err)
			return HEARTBEAT_REQUEST_ERROR
		}
		_, err = g.client.Do(req)
		if err != nil {
			log.Println(err)
			return HEARTBEAT_RESPONSE_ERROR
		}
	}
	return SUCCEED
}

func (g *gossiper) HeartBeatReceiver() {
	mux := http.NewServeMux()
	mux.HandleFunc("/heartbeat", g.HeartBeatHandler)
	// maybe the host doesn't need context
	//contextedMux := g.HeartBeatContext(mux)
	log.Println("Start server on " + g.ip + HEARTBEAT_PORT)
	http.ListenAndServe(HEARTBEAT_PORT, nil)
}

func (hb heartBeat) String() string {
	return "HeartBeat back from: " + hb.gossiperName + ", at " + hb.currentTime.Format(timeFormat)
}

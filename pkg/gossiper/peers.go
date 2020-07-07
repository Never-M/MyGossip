package gossiper

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

// Time formatter
var timeFormat = "2006-01-02 15:04:05"

type gossiper struct {
	name	string
	ip		string
	peers	map[string]gossiper
	client  *http.Client
}

type heartBeat struct {
	gossiperName	string
	currentTime		time.Time
}

func NewGossiper(self string, address string) *gossiper {
	g := &gossiper{
		name: self,
		ip:address,
		peers:make(map[string]gossiper),
		client:&http.Client{},
	}
	return g;
}

func (g *gossiper) AddPeer(peer gossiper)  {
	g.peers[peer.name] = peer
}

func (g *gossiper) RemovePeer(name string)  {
	delete(g.peers, name)
}

func (g *gossiper) PrintPeerNames() {
	fmt.Println(g.peers)
}

func (g *gossiper) HeartBeatHandler(w http.ResponseWriter, r *http.Request) {
	hb := r.Context().Value("heartbeat")
	hbStr := fmt.Sprint(hb)
	log.Println("!!!" + hbStr)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(hbStr))
}

func (g *gossiper) HeartBeatContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "heartbeat", heartBeat{g.name, time.Now()})
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (g *gossiper) SendHeartBeat() {
	ctx := context.WithValue(context.Background(), "heartbeat", &heartBeat{g.name, time.Now()})
	ctx = context.WithValue(ctx, "key", "value")
	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8085/heartbeat", nil)
	if err != nil {
		log.Println(err)
	}
	_, err = g.client.Do(req)
	if err != nil {
		log.Println(err)
	}
}

func (g *gossiper) HeartBeatReceiver() {
	mux := http.NewServeMux()
	mux.HandleFunc("/heartbeat", g.HeartBeatHandler)
	contextedMux := g.HeartBeatContext(mux)
	log.Println("Start server on port :8085")
	http.ListenAndServe(":8085", contextedMux)
}

func (hb heartBeat) String() string {
	return "HeartBeat from: " + hb.gossiperName + ", at " + hb.currentTime.Format(timeFormat);
}

package gossiper

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github/Never-M/MyGossip/pkg/types"
)

// Time formatter
var timeFormat = "2006-01-02 15:04:05"

const HEARTBEAT_PORT = ":8001"
const HEARTBEAT_CTXKEY = "heartbeat"

type gossiper struct {
	name   string
	ip     string
	peers  map[string]*peer
	client *http.Client
}

type heartBeat struct {
	gossiperName string
	currentTime  time.Time
}

func NewGossiper(name, address string) *gossiper {
	var rc int
	if rc != types.SUCCEED {

	}
	return &gossiper{
		name:   name,
		ip:     address,
		peers:  make(map[string]*peer),
		client: &http.Client{},
	}
}

func (g *gossiper) AddPeer(p *peer) int {
	if _, ok := g.peers[p.name]; ok {
		return types.GOSSIPER_PEER_EXIST
	}
	g.peers[p.name] = p
	return types.SUCCEED
}

func (g *gossiper) RemovePeer(name string) int {
	delete(g.peers, name)
	return types.SUCCEED
}

func (g *gossiper) PrintPeerNames() {
	fmt.Println(g.peers)
}

func (g *gossiper) HeartBeatHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Request Context: %+v",r)
	hb := r.Context().Value(HEARTBEAT_CTXKEY)
	hbStr := fmt.Sprint(hb)
	log.Println(hbStr)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(r.Host))

	//incomingPeerIP := r.Host
	//TODO: incoming heartbeat should be recorded. unknown hosts should be added to peer list.
	//TODO
}

// maybe the host doesn't need context
//func (g *gossiper) HeartBeatContext(next http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		ctx := context.WithValue(r.Context(), HEARTBEAT_CTXKEY, heartBeat{g.name, time.Now()})
//		next.ServeHTTP(w, r.WithContext(ctx))
//	})
//}

func (g *gossiper) SendHeartBeat() (int, error) {
	ctx := context.WithValue(context.Background(), HEARTBEAT_CTXKEY, &heartBeat{g.name, time.Now()})
	for _, peer := range g.peers {
		req, err := http.NewRequest("GET", "http://"+peer.ip+HEARTBEAT_PORT+"/"+HEARTBEAT_CTXKEY, nil)
		if err != nil {
			return types.FAIL, err
		}
		newReq := req.WithContext(ctx)
		if err != nil {
			log.Println(err)
			return types.HEARTBEAT_REQUEST_ERROR, err
		}
		_, err = g.client.Do(newReq)
		if err != nil {
			log.Println(err)
			return types.HEARTBEAT_RESPONSE_ERROR, err
		}
	}
	return types.SUCCEED, nil
}

func (g *gossiper) HeartBeatReceiver() {
	mux := http.NewServeMux()
	mux.HandleFunc("/heartbeat", g.HeartBeatHandler)
	// maybe the host doesn't need context
	//contextedMux := g.HeartBeatContext(mux)
	server := &http.Server{
		Addr:         HEARTBEAT_PORT,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
		Handler:      mux,
	}
	log.Println("Start server on " + g.ip + HEARTBEAT_PORT)
	server.ListenAndServe()
}

func (hb heartBeat) String() string {
	return "HeartBeat back from: " + hb.gossiperName + ", at " + hb.currentTime.Format(timeFormat)
}

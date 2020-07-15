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
const DB_BASE_PATH = "/tmp/"

type gossiper struct {
	name   string
	ip     string
	peers  map[string]*gossiper
	client *http.Client
	db     *mydb
}

type heartBeat struct {
	gossiperName string
	currentTime  time.Time
}

func NewGossiper(name, address string) *gossiper {
	newgossiper := &gossiper{}
	newgossiper.name = name
	newgossiper.ip = address
	newgossiper.peers = make(map[string]*gossiper)
	newgossiper.client = &http.Client{}
	var rc int
	rc, newgossiper.db = newdb(DB_BASE_PATH + name)
	if rc != SUCCEED {

	}
	return newgossiper
}

func (g *gossiper) AddPeer(peer *gossiper) int {
	if _, ok := g.peers[peer.name]; ok {
		return GOSSIPER_PEER_EXIST
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
	//hb := r.Context().Value(HEARTBEAT_CTXKEY)
	//hbStr := fmt.Sprint(hb)
	//log.Println(hbStr)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(r.Host))

	incomingPeerIP := r.Host
	//TODO incoming heartbeat should be recorded. unknown hosts should be added to peer list.
	//TODO
}

// maybe the host doesn't need context
//func (g *gossiper) HeartBeatContext(next http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		ctx := context.WithValue(r.Context(), HEARTBEAT_CTXKEY, heartBeat{g.name, time.Now()})
//		next.ServeHTTP(w, r.WithContext(ctx))
//	})
//}

func (g *gossiper) SendHeartBeat() int {
	ctx := context.WithValue(context.Background(), HEARTBEAT_CTXKEY, &heartBeat{g.name, time.Now()})
	for _, peer := range g.peers {
		req, err := http.NewRequest("GET", "http://"+peer.ip+HEARTBEAT_PORT+"/"+HEARTBEAT_CTXKEY, nil)
		req = req.WithContext(ctx)
		if err != nil {
			log.Println(err)
			return HEARTBEAT_REQUEST_ERROR
		}
		resp, err := g.client.Do(req)
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

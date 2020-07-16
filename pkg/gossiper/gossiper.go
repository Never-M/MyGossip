package gossiper

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Never-M/MyGossip/pkg/types"
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
	formData := make(map[string]interface{})
	json.NewDecoder(r.Body).Decode(&formData)
	for key,value := range formData{
		log.Println("key:",key," => value :",value)
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(r.Host))

	//incomingPeerIP := r.Host
	//TODO: incoming heartbeat should be recorded. unknown hosts should be added to peer list.
	//TODO
}

func (g *gossiper) SendHeartBeat() (int, error) {
	for _, peer := range g.peers {
		req, err := http.NewRequest("POST", "http://"+peer.ip+HEARTBEAT_PORT+"/"+HEARTBEAT_CTXKEY, strings.NewReader("name=" + g.name + "&ip=" + g.ip))
		if err != nil {
			return types.HEARTBEAT_RESPONSE_ERROR, err
		}
		_, err = g.client.Do(req)
		if err != nil {
			return types.HEARTBEAT_RESPONSE_ERROR, err
		}
	}
	return types.SUCCEED, nil
}

func (g *gossiper) HeartBeatReceiver() {
	mux := http.NewServeMux()
	mux.HandleFunc("/heartbeat", g.HeartBeatHandler)
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

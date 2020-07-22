package gossiper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Never-M/MyGossip/pkg/types"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// Time formatter
var timeFormat = "2006-01-02 15:04:05"

const HEARTBEAT_PORT = ":8001"
const HEARTBEAT_PATH = "heartbeat"

type gossiper struct {
	name   string
	ip     string
	peers  map[string]*peer
	client *http.Client
	db     *mydb
}

type heartBeat struct {
	Name string `json:"name"`
	Ip   string `json:"ip"`
	Time string `json:"time"`
}

func NewGossiper(name, address string) *gossiper {
	_, db, err := Newdb("/tmp/" + name + "/database")
	if err != nil {
		//TODO: Using logrus later
		fmt.Println(err)
		log.Fatal("database Failed")
	}
	return &gossiper{
		name:  name,
		ip:    address,
		peers: make(map[string]*peer),
		db:    db,
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
	// Decode heartbeat json
	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	hb := &heartBeat{}
	err := json.Unmarshal(body, hb)
	if err != nil {
		panic(err)
	}

	// check node in the peers list or not
	if _, ok := g.peers[hb.Name]; !ok {
		g.peers[hb.Name] = NewPeer(hb.Name, hb.Ip)
	} else {
		//TODO receive heartbeat
	}
}

func (g *gossiper) SendHeartBeat() (int, error) {
	for _, peer := range g.peers {
		// encode heartbeat to json
		hb := &heartBeat{
			g.name,
			g.ip,
			time.Now().Format(timeFormat),
		}
		hbJson, err := json.Marshal(hb)
		if err != nil {
			return types.FAILED, err
		}
		body := bytes.NewBuffer(hbJson)
		// url
		url := "http://" + peer.ip + HEARTBEAT_PORT + "/" + HEARTBEAT_PATH
		//send request
		resp, err := http.Post(url, "application/json;charset=utf-8", body)
		if err != nil {
			return types.HEARTBEAT_RESPONSE_ERROR, err
		}
		resp.Body.Close()
	}
	return types.SUCCEED, nil
}

func (g *gossiper) HeartBeatReceiver() {
	mux := http.NewServeMux()
	mux.HandleFunc("/"+HEARTBEAT_PATH, g.HeartBeatHandler)
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
	return "HeartBeat back from: " + hb.Name + ", at " + hb.Time
}

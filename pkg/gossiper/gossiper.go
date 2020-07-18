package gossiper

import (
	"bytes"
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
	// formData := make(map[string]interface{})
	// json.NewDecoder(r.Body).Decode(&formData)
	// for key,value := range formData{
	// 	log.Println("key:",key," => value :",value)
	// }
	body, _ := ioutil.ReadAll(r.Body)
	fmt.Println("Got the request!")
	fmt.Println("request Body:", string(body))

	//incomingPeerIP := r.Host
	//TODO: incoming heartbeat should be recorded. unknown hosts should be added to peer list.
	//TODO
}

// maybe delete this later
// type HeartBeatBody struct {
// 	name string
// 	ip   string
// }

func (g *gossiper) SendHeartBeat() (int, error) {
	for _, peer := range g.peers {
		//hbb := HeartBeatBody{name: g.name, ip: g.ip}
		var jsonStr = []byte(`{"name":` + g.name + `, "ip" :` + g.ip + `}`)
		url := "http://" + peer.ip + HEARTBEAT_PORT + "/" + HEARTBEAT_PATH
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
		req.Header.Set("X-Custom-Header", "myvalue")
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		fmt.Println("Response Status:", resp.Status)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		// fmt.Println("response Headers:", resp.Header)
		// fmt.Println("send request...")
		// body, _ := ioutil.ReadAll(resp.Body)
		// fmt.Println("response Body:", string(body))
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
	return "HeartBeat back from: " + hb.gossiperName + ", at " + hb.currentTime.Format(timeFormat)
}

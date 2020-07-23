package gossiper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Never-M/MyGossip/pkg/types"
	"io/ioutil"
	"net/http"
	"time"
)

// Time formatter
var timeFormat = "2006-01-02 15:04:05"

const HEARTBEAT_PORT = ":8001"
const HEARTBEAT_PATH = "heartbeat"
const HEARTBEAT_TIMEOUT = 1000

type gossiper struct {
	name   string
	ip     string
	peers  map[string]*peer
	heartbeatTimer *time.Timer
	terminateChan chan int
	db     *mydb
	logger *logger
}

type heartBeat struct {
	Name string `json:"name"`
	Ip   string `json:"ip"`
	Time string `json:"time"`
}

func NewGossiper(name, ip string) *gossiper {
	_, db, err := Newdb("/tmp/" + name + "/database")
	logger := Newlogger()
	logger.SaveToFile("/tmp/" + name + "/log")
	if err != nil {
		fmt.Println(err)
		logger.Fatal("create database Failed", "gossiper", "NewGossiper")
	}
	return &gossiper{
		name:   name,
		ip:     ip,
		peers:  make(map[string]*peer),
		terminateChan:make(chan int),
		db:     db,
		logger: logger,
	}
}

func (g *gossiper) Start()  {
	go g.HeartBeatReceiver()
	g.heartbeatTimer = time.NewTimer(HEARTBEAT_TIMEOUT * time.Millisecond)
	g.logger.Info("%v start", g.name)
	go func() {
		for {
			select {
			case <-g.heartbeatTimer.C:
				g.SendHeartBeats()
				g.heartbeatTimer.Reset(HEARTBEAT_TIMEOUT * time.Millisecond)
			case <-g.terminateChan:
				g.logger.Info("%v stoped", g.name)
				break
			}
		}
	}()
}

func (g *gossiper) Stop()  {
	g.terminateChan <- 1
}

func (g *gossiper) AddPeer(p *peer) int {
	if _, ok := g.peers[p.name]; ok {
		return types.GOSSIPER_PEER_EXIST
	}
	g.peers[p.name] = p
	g.logger.Info("New peer " + p.name + " joined " + g.name)
	go func() {
		for {
			select {
			case <-p.timer.C:
				g.RemovePeer(p.name)
				g.logger.Info("Peer " + p.name + " of " + g.name + " removed")
				break
			}
		}
	}()
	return types.SUCCEED
}

func (g *gossiper) RemovePeer(name string) int {
	delete(g.peers, name)
	g.logger.Info(name + " removed from " + g.name)
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
		g.logger.Error(err.Error(), "gossiper", "HeartBeatHandler")
	}

	g.logger.Info(g.name + " receive heartbeat from " + hb.Name)
	// check node in the peers list or not
	if _, ok := g.peers[hb.Name]; !ok {
		g.AddPeer(NewPeer(hb.Name, hb.Ip))
		g.logger.Info("Receive heartbeat from unknown node " + hb.Name + ", added to peerlist")
	} else {
		g.peers[hb.Name].timer.Reset(2 * HEARTBEAT_TIMEOUT * time.Millisecond)
		g.logger.Info(hb.Name + " timer reset")
	}
}

func (g *gossiper) SendHeartBeats() (int, error) {
	g.logger.Info(g.name + "sending heartbeats...")
	for _, peer := range g.peers {
		go g.SendHeartBeat(peer)
	}
	return types.SUCCEED, nil
}

func (g *gossiper) SendHeartBeat(p *peer) (int, error) {
	// encode heartbeat to json
	hb := &heartBeat{
		g.name,
		g.ip,
		time.Now().Format(timeFormat),
	}
	hbJson, err := json.Marshal(hb)
	if err != nil {
		g.logger.Error(g.name + " json marshal failed", "gossiper", "SendHeartBeat")
		return types.FAILED, err
	}
	body := bytes.NewBuffer(hbJson)
	// url
	url := "http://" + p.ip + HEARTBEAT_PORT + "/" + HEARTBEAT_PATH
	//send request
	resp, err := http.Post(url, "application/json;charset=utf-8", body)
	if err != nil {
		g.logger.Error(g.name + " send heartbeat to " + p.name + " response error", "gossiper", "SendHeartBeat")
		return types.HEARTBEAT_RESPONSE_ERROR, err
	}
	resp.Body.Close()
	g.logger.Info(g.name + " sent heartbeat to " + p.name)
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
	g.logger.Info(g.name + "Start server on " + g.ip + HEARTBEAT_PORT)
	server.ListenAndServe()
}

func (hb heartBeat) String() string {
	return "HeartBeat back from: " + hb.Name + ", at " + hb.Time
}

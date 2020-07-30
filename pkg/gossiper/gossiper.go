package gossiper

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/Never-M/MyGossip/pkg/types"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// Time formatter
var timeFormat = "2006-01-02 15:04:05"

const HEARTBEAT_PORT = ":8001"
const HEARTBEAT_PATH = "heartbeat"
const HEARTBEAT_TIMEOUT = 1000

type Gossiper struct {
	name           string
	ip             string
	peers          map[string]*peer
	heartbeatTimer *time.Timer
	terminateChan  chan int
	db             *mydb
	logger         *logger
}

type heartBeat struct {
	Name string `json:"name"`
	Ip   string `json:"ip"`
	Time string `json:"time"`
}

func NewGossiper(name, ip string) *Gossiper {
	_, db, err := Newdb(filepath.Join("/tmp/", name, "/database"))
	logger := Newlogger()
	logger.SaveToFile(filepath.Join("/tmp/", name, "/log"))
	if err != nil {
		logger.Fatal("create database Failed", "gossiper", "NewGossiper")
	}
	peers := make(map[string]*peer)

	//check if file exsit
	f, err := os.Open(filepath.Join("/tmp/", name, "/peers.csv"))
	f.Close()
	if err == nil {
		//exsit!
		peerPairSlice := ReadPeersFromFile(name)
		if len(peerPairSlice) > 0 {
			for _, item := range peerPairSlice {
				peers[item.Name] = NewPeer(item.Name, item.IP)
			}
		}
	}
	return &Gossiper{
		name:          name,
		ip:            ip,
		peers:         peers,
		terminateChan: make(chan int),
		db:            db,
		logger:        logger,
	}
}

type PeerPair struct {
	Name string
	IP   string
}

func (g *Gossiper) WritePeersToFile() {
	var toSave []PeerPair
	for _, peer := range g.peers {
		toSave = append(toSave, PeerPair{Name: peer.name, IP: peer.ip})
	}

	csvFile, err := os.Create(filepath.Join("/tmp/", g.name, "/peers.csv"))
	if err != nil {
		g.logger.Panic("Create csv file Failed -- <err>: ")
	}
	defer csvFile.Close()
	writer := csv.NewWriter(csvFile)

	for _, pair := range toSave {
		line := []string{pair.Name, pair.IP}
		err = writer.Write(line)
		if err != nil {
			g.logger.Panic("Write Failed -- <err>: ")
		}
	}
	writer.Flush()
}

func ReadPeersFromFile(name string) []PeerPair {
	logger := Newlogger()
	file, err := os.Open(filepath.Join("/tmp/", name, "/peers.csv"))
	if err != nil {
		logger.Panic("Failed to open csv")
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1
	record, err := reader.ReadAll()
	if err != nil {
		logger.Panic("Failed to read csv")
	}
	var peerPairSlice []PeerPair
	for _, item := range record {
		peerPairSlice = append(peerPairSlice, PeerPair{Name: item[0], IP: item[1]})
	}
	return peerPairSlice
}

func (g *Gossiper) Start() {
	go g.HeartBeatReceiver()
	g.heartbeatTimer = time.NewTimer(HEARTBEAT_TIMEOUT * time.Millisecond)
	g.logger.Info(g.name + " started")
	go func() {
		for {
			select {
			case <-g.heartbeatTimer.C:
				g.SendHeartBeats()
				g.heartbeatTimer.Reset(HEARTBEAT_TIMEOUT * time.Millisecond)
			case <-g.terminateChan:
				// save neighbors to a file
				g.WritePeersToFile()
				g.logger.Info(g.name + " stoped")
				break
			}
		}
	}()
}

func (g *Gossiper) Stop() {
	g.terminateChan <- 1
}

func (g *Gossiper) AddPeer(p *peer) int {
	if _, ok := g.peers[p.name]; ok {
		return types.SUCCEED
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

func (g *Gossiper) RemovePeer(name string) int {
	delete(g.peers, name)
	g.logger.Info(name + " removed from " + g.name)
	return types.SUCCEED
}

func (g *Gossiper) PrintPeerNames() {
	fmt.Println(g.peers)
}

func (g *Gossiper) HeartBeatHandler(w http.ResponseWriter, r *http.Request) {
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

func (g *Gossiper) SendHeartBeats() (int, error) {
	g.logger.Info(g.name + "sending heartbeats...")
	for _, peer := range g.peers {
		go g.SendHeartBeat(peer)
	}
	return types.SUCCEED, nil
}

func (g *Gossiper) SendHeartBeat(p *peer) (int, error) {
	// encode heartbeat to json
	hb := &heartBeat{
		g.name,
		g.ip,
		time.Now().Format(timeFormat),
	}
	hbJson, err := json.Marshal(hb)
	if err != nil {
		g.logger.Error(g.name+" json marshal failed", "gossiper", "SendHeartBeat")
		return types.FAILED, err
	}
	body := bytes.NewBuffer(hbJson)
	// url
	url := "http://" + p.ip + HEARTBEAT_PORT + "/" + HEARTBEAT_PATH
	//send request
	resp, err := http.Post(url, "application/json;charset=utf-8", body)
	if err != nil {
		g.logger.Error(g.name+" send heartbeat to "+p.name+" response error", "gossiper", "SendHeartBeat")
		return types.HEARTBEAT_RESPONSE_ERROR, err
	}
	resp.Body.Close()
	g.logger.Info(g.name + " sent heartbeat to " + p.name)
	return types.SUCCEED, nil
}

func (g *Gossiper) HeartBeatReceiver() {
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

func (g *Gossiper) GetName() string {
	return g.name
}

func (g *Gossiper) GetIP() string {
	return g.ip
}

func (g *Gossiper) GetPeers() map[string]*peer {
	return g.peers
}

func (g *Gossiper) GetHeartBeatTimer() *time.Timer {
	return g.heartbeatTimer
}

func (g *Gossiper) GetTerminateChan() chan int {
	return g.terminateChan
}

func (g *Gossiper) GetDB() *mydb {
	return g.db
}
func (g *Gossiper) GetLogger() *logger {
	return g.logger
}

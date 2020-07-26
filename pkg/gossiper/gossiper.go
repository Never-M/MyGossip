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
	"time"
)

// Time formatter
var timeFormat = "2006-01-02 15:04:05"

const HEARTBEAT_PORT = ":8001"
const HEARTBEAT_PATH = "heartbeat"
const HEARTBEAT_TIMEOUT = 1000

type gossiper struct {
	Name           string
	IP             string
	Peers          map[string]*peer
	HeartbeatTimer *time.Timer
	TerminateChan  chan int
	DB             *mydb
	Logger         *logger
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
		logger.Fatal("create database Failed", "gossiper", "NewGossiper")
	}
	peers := make(map[string]*peer)

	//check if file exsit
	_, err = os.Stat("/tmp/" + name + "/peers.csv")
	if err == nil {
		//exsit!
		peerpairslice := Read(name)
		if len(peerpairslice) > 0 {
			for _, item := range peerpairslice {
				peers[item.Name] = NewPeer(item.Name, item.IP)
			}
		}
	}
	return &gossiper{
		Name:          name,
		IP:            ip,
		Peers:         peers,
		TerminateChan: make(chan int),
		DB:            db,
		Logger:        logger,
	}
}

type PeerPair struct {
	Name string
	IP   string
}

func (g *gossiper) Write() {
	var tosave []PeerPair
	for _, peer := range g.Peers {
		tosave = append(tosave, PeerPair{Name: peer.name, IP: peer.ip})
	}

	csvFile, err := os.Create("/tmp/" + g.Name + "/peers.csv")
	if err != nil {
		g.Logger.Panic("Create csv file Failed -- <err>: ")
	}
	defer csvFile.Close()
	writer := csv.NewWriter(csvFile)

	for _, pair := range tosave {
		line := []string{pair.Name, pair.IP}
		err = writer.Write(line)
		if err != nil {
			g.Logger.Panic("Write Failed -- <err>: ")
		}
	}
	writer.Flush()
}

func Read(name string) []PeerPair {
	logger := Newlogger()
	file, err := os.Open("/tmp/" + name + "/peers.csv")
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
	var peerpairslice []PeerPair
	for _, item := range record {
		peerpairslice = append(peerpairslice, PeerPair{Name: item[0], IP: item[1]})
	}
	return peerpairslice
}

func (g *gossiper) Start() {
	go g.HeartBeatReceiver()
	g.HeartbeatTimer = time.NewTimer(HEARTBEAT_TIMEOUT * time.Millisecond)
	g.Logger.Info(g.Name + " started")
	go func() {
		for {
			select {
			case <-g.HeartbeatTimer.C:
				g.SendHeartBeats()
				g.HeartbeatTimer.Reset(HEARTBEAT_TIMEOUT * time.Millisecond)
			case <-g.TerminateChan:
				// save neighbors to a file
				g.Write()
				g.Logger.Info(g.Name + " stoped")
				break
			}
		}
	}()
}

func (g *gossiper) Stop() {
	g.TerminateChan <- 1
}

func (g *gossiper) AddPeer(p *peer) int {
	if _, ok := g.Peers[p.name]; ok {
		return types.SUCCEED
	}
	g.Peers[p.name] = p
	g.Logger.Info("New peer " + p.name + " joined " + g.Name)
	go func() {
		for {
			select {
			case <-p.timer.C:
				g.RemovePeer(p.name)
				g.Logger.Info("Peer " + p.name + " of " + g.Name + " removed")
				break
			}
		}
	}()
	return types.SUCCEED
}

func (g *gossiper) RemovePeer(name string) int {
	delete(g.Peers, name)
	g.Logger.Info(name + " removed from " + g.Name)
	return types.SUCCEED
}

func (g *gossiper) PrintPeerNames() {
	fmt.Println(g.Peers)
}

func (g *gossiper) HeartBeatHandler(w http.ResponseWriter, r *http.Request) {
	// Decode heartbeat json
	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	hb := &heartBeat{}
	err := json.Unmarshal(body, hb)
	if err != nil {
		g.Logger.Error(err.Error(), "gossiper", "HeartBeatHandler")
	}

	g.Logger.Info(g.Name + " receive heartbeat from " + hb.Name)
	// check node in the peers list or not
	if _, ok := g.Peers[hb.Name]; !ok {
		g.AddPeer(NewPeer(hb.Name, hb.Ip))
		g.Logger.Info("Receive heartbeat from unknown node " + hb.Name + ", added to peerlist")
	} else {
		g.Peers[hb.Name].timer.Reset(2 * HEARTBEAT_TIMEOUT * time.Millisecond)
		g.Logger.Info(hb.Name + " timer reset")
	}
}

func (g *gossiper) SendHeartBeats() (int, error) {
	g.Logger.Info(g.Name + "sending heartbeats...")
	for _, peer := range g.Peers {
		go g.SendHeartBeat(peer)
	}
	return types.SUCCEED, nil
}

func (g *gossiper) SendHeartBeat(p *peer) (int, error) {
	// encode heartbeat to json
	hb := &heartBeat{
		g.Name,
		g.IP,
		time.Now().Format(timeFormat),
	}
	hbJson, err := json.Marshal(hb)
	if err != nil {
		g.Logger.Error(g.Name+" json marshal failed", "gossiper", "SendHeartBeat")
		return types.FAILED, err
	}
	body := bytes.NewBuffer(hbJson)
	// url
	url := "http://" + p.ip + HEARTBEAT_PORT + "/" + HEARTBEAT_PATH
	//send request
	resp, err := http.Post(url, "application/json;charset=utf-8", body)
	if err != nil {
		g.Logger.Error(g.Name+" send heartbeat to "+p.name+" response error", "gossiper", "SendHeartBeat")
		return types.HEARTBEAT_RESPONSE_ERROR, err
	}
	resp.Body.Close()
	g.Logger.Info(g.Name + " sent heartbeat to " + p.name)
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
	g.Logger.Info(g.Name + "Start server on " + g.IP + HEARTBEAT_PORT)
	server.ListenAndServe()
}

func (hb heartBeat) String() string {
	return "HeartBeat back from: " + hb.Name + ", at " + hb.Time
}

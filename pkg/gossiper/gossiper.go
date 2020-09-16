package gossiper

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/Never-M/MyGossip/pkg/types"
	"github.com/String-Reconciliation-Ditributed-System/RCDS_GO/pkg/lib/genSync"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
	"github.com/String-Reconciliation-Ditributed-System/RCDS_GO/pkg/lib/algorithm/iblt"
	"github.com/String-Reconciliation-Ditributed-System/RCDS_GO/pkg/set"
)

const (
	PUT = 1
	GET = 2
	DELETE = 3
)
// Time formatter
var timeFormat = "2006-01-02 15:04:05"

const HEARTBEAT_PORT = ":8001"
const SYNC_PORT = 8002
const HEARTBEAT_PATH = "heartbeat"
const HEARTBEAT_TIMEOUT = 1000

const FIXED_DIFF = 4

type Gossiper struct {
	name			string
	ip				string
	peers			map[string]*peer
	heartbeatTimer	*time.Timer
	terminateChan	chan int
	q				queue
	syncServer		genSync.GenSync
	pointer 		int
	db				*mydb
	logger			*logger
	counter 		int
}


type heartBeat struct {
	Name string `json:"name"`
	Ip   string `json:"ip"`
	Time string `json:"time"`
}

type queue []*logEntry

func (g *Gossiper) PopLeft() *logEntry {
	if len(g.q) != 0 {
		res := g.q[0]
		g.q = g.q[1:]
		return res
	} else {
		return nil
	}
}

func (g *Gossiper) Push(l *logEntry) {
	g.q = append(g.q, l)
}

type logEntry struct {
	operation		int
	key 			string
	value			string
	timestamp 		time.Time
}

func NewGossiper(name, ip string) *Gossiper {
	_, db, err := Newdb(filepath.Join("/tmp/", name, "/database"))
	logger := Newlogger()
	logger.SaveToFile(filepath.Join("/tmp/", name, "/log"))
	if err != nil {
		fmt.Println(err)
		logger.Fatal("create database Failed", "gossiper", "NewGossiper")
	}
	peers := make(map[string]*peer)

	//check if file exsit
	f, err := os.Open(filepath.Join("/tmp/", name, "/peers.csv"))
	f.Close()
	//if err == nil {
	//	//exsit!
	//	peerPairSlice := ReadPeersFromFile(name)
	//	if len(peerPairSlice) > 0 {
	//		for _, item := range peerPairSlice {
	//			peers[item.Name] = NewPeer(item.Name, item.IP)
	//		}
	//	}
	//}
	// new Iblt
	syncServer, err := iblt.NewIBLTSetSync(iblt.WithSymmetricSetDiff(FIXED_DIFF))
	return &Gossiper{
		name:          	name,
		ip:            	ip,
		peers:         	peers,
		terminateChan: 	make(chan int),
		db:            	db,
		logger:        	logger,
		syncServer:		syncServer,
		counter: 		0,
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
	go g.SyncServerStart()
	go g.CheckSyncClient()

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
				g.logger.Info(g.name + " stopped")
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
	if len(g.peers) == 0 {
		fmt.Printf("%v has no peer")
		return
	}
	fmt.Printf("Peers of %v are:", g.name)
	for name, _ := range g.peers {
		fmt.Print(name + " ")
	}
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

func (g *Gossiper) SyncClient(logEntryNum int) {
	var logEntriesToCommit []*logEntry

	for i := 0; i < logEntryNum; i++ {
		logEntryToCommit := g.PopLeft()
		// add to commit
		logEntriesToCommit = append(logEntriesToCommit, logEntryToCommit)


		err := g.syncServer.AddElement(encodeLogEntry(logEntryToCommit))
		if err != nil {
			g.logger.Error(g.name + ":Put error " + err.Error())
		}
	}

	var wg sync.WaitGroup
	for peerName, peer := range g.peers {
		wg.Add(1)
		go func() {
			e := g.syncServer.SyncClient(peer.ip, SYNC_PORT)
			if e != nil {
				fmt.Printf("SyncClient error: %v", e)
				g.logger.Error(g.name + ":Sync error with " + peerName)
			}
			wg.Done()
			//s := g.syncServer.GetLocalSet()
			i := g.syncServer.GetTotalBytes()
			fmt.Println(i)
			// TODO 写入queue，并commit到db
		}()
	}
	wg.Wait()
}

func (g *Gossiper) CheckSyncClient() {
	for {
		if len(g.q) > FIXED_DIFF / 2 {
			g.SyncClient(FIXED_DIFF/2)
		} else if len(g.q) < FIXED_DIFF / 2 && len(g.q) > 0 {
			g.SyncClient(len(g.q))
		} else {
			time.Sleep(time.Second)
		}
	}
}

func (g *Gossiper) SyncServerStart() {
	for {
		var wg sync.WaitGroup
		wg.Add(1)
		//go func() {
			err := g.syncServer.SyncServer(SYNC_PORT)
			if err != nil {
				fmt.Println(err)
				g.logger.Error("SyncServer err: " + err.Error())
				return
			}
			wg.Done()
		//}()
		wg.Wait()
		s := set.New()
		s = g.syncServer.GetLocalSet()
		fmt.Println("SyncServer done")
		for k, v := range *s {
			fmt.Println(k)
			fmt.Println(v)
		}
	}
}

func (g *Gossiper) Put(key, value string) {
	entry := &logEntry{
		operation:PUT,
		key:key,
		value:value,
		timestamp:time.Now(),
	}
	g.Push(entry)

	if g.counter++; g.counter == FIXED_DIFF/2 {
		g.SyncClient(FIXED_DIFF/2)
	}
}

func (g *Gossiper) Delete(key string) {
	entry := &logEntry{
		operation:DELETE,
		key:key,
		timestamp:time.Now(),
	}
	g.Push(entry)
	if g.counter++; g.counter == FIXED_DIFF/2 {
		g.SyncClient(FIXED_DIFF/2)
	}
}

func (g *Gossiper) Get(key string) string {
	rc, value, err := g.db.Get(key)
	if rc == types.SUCCEED {
		return value
	} else {
		g.logger.Error(g.name + ":Get error" + string(rc) + " " + err.Error())
		return ""
	}
}

func encodeLogEntry(l *logEntry) []byte {
	buf, err := json.Marshal(l)
	if err != nil {
		fmt.Printf("Encode error: %v", err)
	}
	return buf
}

func decodeLogEntry(b []byte) *logEntry{
	l := &logEntry{}
	if err := json.Unmarshal(b, &l); err != nil {
		fmt.Printf("Decode error: %v", err)
	}
	return l
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

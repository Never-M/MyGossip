package gossiper

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/Never-M/MyGossip/pkg/types"
	"github.com/String-Reconciliation-Ditributed-System/RCDS_GO/pkg/lib/algorithm/iblt"
	"github.com/String-Reconciliation-Ditributed-System/RCDS_GO/pkg/lib/genSync"
)

const (
	PUT    = "PUT"
	DELETE = "DELETE"
)
const HEARTBEAT_PORT = ":8002"
const SYNC_PORT = 8001
const HEARTBEAT_PATH = "heartbeat"
const HEARTBEAT_TIMEOUT = 1000
const FIXED_DIFF = 4

// Time formatter
var timeFormat = "2006-01-02 15:04:05"

type queue []logEntry

type logEntry struct {
	Operation string `json:"operation"`
	Key       string `json:"key"`
	Value     string `json:"value"`
	Timestamp string `json:"timestamp"`
}

type heartBeat struct {
	Name string `json:"name"`
	Ip   string `json:"ip"`
	Time string `json:"time"`
}

type peerTime struct {
	Name string
	t    time.Time
}

func (hb heartBeat) String() string {
	return "HeartBeat back from: " + hb.Name + ", at " + hb.Time
}

type Gossiper struct {
	name           string
	ip             string
	peers          map[string]*peer
	heartbeatTimer *time.Timer
	terminateChan  chan int
	q              queue
	removePeerQ    []peerTime
	SyncServer     genSync.GenSync
	pointer        int
	db             *mydb
	logger         *logger
	counter        int
	logFile        *os.File
	checkLog       map[string]logEntry
	writer         *csv.Writer
}

func (g *Gossiper) PopLeft() logEntry {
	if len(g.q) != 0 {
		res := g.q[0]
		g.q = g.q[1:]
		return res
	} else {
		return logEntry{}
	}
}

func (g *Gossiper) Push(l logEntry) {
	g.q = append(g.q, l)
}

func NewGossiper(name, ip string) *Gossiper {
	logFile, err := os.OpenFile(filepath.Join("/tmp/", name, "/logs.csv"), os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		fmt.Printf("can not create log file, err %+v\n", err)
	}
	w := csv.NewWriter(logFile)
	w.Comma = ','
	w.UseCRLF = true
	_, db, err := Newdb(filepath.Join("/tmp/", name, "/database"))
	logger := Newlogger()
	logger.SaveToFile(filepath.Join("/tmp/", name, "/log"))
	if err != nil {
		logger.Fatal("create database Failed", "gossiper", "NewGossiper")
	}
	peers := make(map[string]*peer)

	// new Iblt
	SyncServer, err := iblt.NewIBLTSetSync(iblt.WithSymmetricSetDiff(FIXED_DIFF))

	g := &Gossiper{
		name:          name,
		ip:            ip,
		peers:         peers,
		terminateChan: make(chan int),
		db:            db,
		logger:        logger,
		SyncServer:    SyncServer,
		counter:       0,
		logFile:       logFile,
		writer:        w,
	}

	//check if file exsit
	f, err := os.Open(filepath.Join("/tmp/", name, "/peers.csv"))
	f.Close()
	if err == nil {
		//exsit!
		peerPairSlice := ReadPeersFromFile(name)
		if len(peerPairSlice) > 0 {
			for _, item := range peerPairSlice {
				g.AddPeer(NewPeer(item.Name, item.IP))
			}
		}
	}

	//read log to fill myc
	for _, logEntry := range g.ReadLogsFromFile() {
		g.checkLog[logEntry.Key] = logEntry
	}

	return g
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
			g.logger.Panic("Write peers Failed -- <err>: ", err.Error())
		}
	}
	writer.Flush()
}

func (g *Gossiper) WriteLogToFile(l logEntry) {
	line := []string{l.Timestamp, l.Operation, l.Key, l.Value}
	err := g.writer.Write(line)
	if err != nil {
		g.logger.Error("Write log Failed -- <err>: ", err.Error())
	}
	g.writer.Flush()
}

func (g *Gossiper) ReadLogsFromFile() []logEntry {
	reader := csv.NewReader(g.logFile)
	reader.FieldsPerRecord = -1
	record, err := reader.ReadAll()
	if err != nil {
		g.logger.Error("Failed to read csv")
	}
	var logs []logEntry
	for _, item := range record {
		logs = append(logs, logEntry{Timestamp: item[0], Operation: item[1], Key: item[2], Value: item[3]})
	}
	return logs
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
		logger.Panic("Failed to read peers")
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
				g.logFile.Close()
				g.logger.Info(g.name + " stopped")
				break
			}
		}
	}()
}

func (g *Gossiper) Stop() {
	// TODO maybe add a waitgroup to wait channel signal to send
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

	// sync with new peer myc
	logs := g.ReadLogsFromFile()
	for _, log := range logs {
		g.Push(log)
	}
	/*
		var localToSyncEntries []*logEntry
		_, pairs, err := g.db.ListData()
		if err != nil {
			g.logger.Error(g.name + "can't get data from db")
		}
		for _, pair := range pairs {
			entry := &logEntry{
				Operation: PUT,
				Key:       pair.key,
				Value:     pair.val,
				Timestamp: time.Now().String(),
			}
			localToSyncEntries = append(localToSyncEntries, entry)
		}
	*/

	return types.SUCCEED
}

func (g *Gossiper) RemovePeer(name string) int {
	delete(g.peers, name)
	g.logger.Info(name + " removed from " + g.name)
	g.removePeerQ = append(g.removePeerQ, peerTime{name, time.Now()})
	return types.SUCCEED
}

func (g *Gossiper) PrintPeerNames() {
	if len(g.peers) == 0 {
		fmt.Printf(">> %v has no peer\n", g.name)
		return
	}
	fmt.Printf(">> Peers of %v:\n", g.name)
	for name, _ := range g.peers {
		fmt.Println(">> " + name)
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
		if len(g.removePeerQ) != 0 {
			for g.removePeerQ[0].t.Add(time.Second).Before(time.Now()) {
				g.removePeerQ = g.removePeerQ[1:]
				if len(g.removePeerQ) == 0 {
					break
				}
			}
		}
		flag := true
		for _, item := range g.removePeerQ {
			if item.Name == hb.Name {
				flag = false
				break
			}
		}
		if flag {
			g.AddPeer(NewPeer(hb.Name, hb.Ip))
			g.logger.Info("Receive heartbeat from unknown node " + hb.Name + ", added to peerlist")
		}
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

func (g *Gossiper) CheckSyncClient() {
	for {
		if len(g.q) >= FIXED_DIFF/2 {
			g.SyncClientStart(FIXED_DIFF / 2)
		} else if len(g.q) < FIXED_DIFF/2 && len(g.q) > 0 {
			g.SyncClientStart(len(g.q))
		} else {
			time.Sleep(time.Second)
		}
	}
}

func (g *Gossiper) SyncClientStart(logEntryNum int) {
	var logEntriesToCommit []logEntry

	for i := 0; i < logEntryNum; i++ {
		logEntryToCommit := g.PopLeft()
		// add to commit
		logEntriesToCommit = append(logEntriesToCommit, logEntryToCommit)

		err := g.SyncServer.AddElement(encodeLogEntry(logEntryToCommit))
		if err != nil {
			g.logger.Error(g.name + ":Put error " + err.Error())
		}
	}

	var wg sync.WaitGroup
	for peerName, peer := range g.peers {
		wg.Add(1)
		go func() {
			e := g.SyncServer.SyncClient(peer.ip, SYNC_PORT)
			if e != nil {
				g.logger.Error(g.name + ":Sync error with " + peerName)
			}
			wg.Done()
			additions := g.SyncServer.GetSetAdditions()
			for k := range *additions {
				//myc
				l := decodeLogEntry([]byte(k.(string)))
				_, ok := g.checkLog[l.Key]
				if ok {
					if l.Timestamp > g.checkLog[l.Key].Timestamp {
						g.Push(l)
					}
				} else {
					g.Push(l)
				}
			}
		}()
	}
	wg.Wait()

	for _, entry := range logEntriesToCommit {
		g.WriteLogToFile(entry)
		if entry.Operation == PUT {
			_, e := g.db.Put(entry.Key, entry.Value)
			if e != nil {
				g.logger.Error(g.name + "can't put (" + entry.Key + "," + entry.Value + ") into database")
			}
		} else if entry.Operation == DELETE { // DELETE
			_, e := g.db.Delete(entry.Key)
			if e != nil {
				g.logger.Error(g.name + "can't delete" + entry.Key + "from database")
			}
		}
	}
}

func (g *Gossiper) SyncServerStart() {
	for {
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			err := g.SyncServer.SyncServer(g.ip, SYNC_PORT)
			if err != nil {
				g.logger.Error(err.Error())
			}
			wg.Done()
		}()
		wg.Wait()
		additions := g.SyncServer.GetSetAdditions()
		for k := range *additions {
			l := decodeLogEntry([]byte(k.(string)))
			//myc
			_, ok := g.checkLog[l.Key]
			if ok {
				if l.Timestamp > g.checkLog[l.Key].Timestamp {
					g.Push(l)
				}
			} else {
				g.Push(l)
			}
		}
	}
}

func (g *Gossiper) Put(key, value string) {
	entry := logEntry{
		Operation: PUT,
		Key:       key,
		Value:     value,
		Timestamp: time.Now().Format(timeFormat),
	}
	g.Push(entry)
}

func (g *Gossiper) Delete(key string) {
	entry := logEntry{
		Operation: DELETE,
		Key:       key,
		Timestamp: time.Now().Format(timeFormat),
	}
	g.Push(entry)
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

func encodeLogEntry(l logEntry) []byte {
	buf, err := json.Marshal(l)
	if err != nil {
		fmt.Printf(">> Encode error: %v\n", err)
	}
	return buf
}

func decodeLogEntry(b []byte) logEntry {
	l := logEntry{}
	if err := json.Unmarshal(b, &l); err != nil {
		fmt.Printf(">> Decode error: %v\n", err)
	}
	return l
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

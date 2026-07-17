package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"
)

type ARecord struct {
	IP  net.IP `json:"ip"`
	TTL uint32 `json:"ttl"`
}

type AAAARecord struct {
	IP  net.IP `json:"ip"`
	TTL uint32 `json:"ttl"`
}

type CNAMERecord struct {
	Target string `json:"target"`
	TTL    uint32 `json:"ttl"`
}

type MXRecord struct {
	Preference uint16 `json:"preference"`
	MX         string `json:"mx"`
	TTL        uint32 `json:"ttl"`
}

type NSRecord struct {
	NS  string `json:"ns"`
	TTL uint32 `json:"ttl"`
}

type TXTRecord struct {
	Txt []string `json:"txt"`
	TTL uint32   `json:"ttl"`
}

type ZoneStore struct {
	mu           sync.RWMutex
	aRecords     map[string][]ARecord
	aaaaRecords  map[string][]AAAARecord
	cnameRecords map[string]CNAMERecord
	mxRecords    map[string][]MXRecord
	nsRecords    map[string][]NSRecord
	txtRecords   map[string][]TXTRecord
	syncPort     int
	isRunning    bool
}

type ZoneData struct {
	ARecords     map[string][]ARecord    `json:"a_records"`
	AAAARecords  map[string][]AAAARecord `json:"aaaa_records"`
	CNAMERecords map[string]CNAMERecord  `json:"cname_records"`
	MXRecords    map[string][]MXRecord   `json:"mx_records"`
	NSRecords    map[string][]NSRecord   `json:"ns_records"`
	TXTRecords   map[string][]TXTRecord  `json:"txt_records"`
}

func NewZoneStore() *ZoneStore {
	return &ZoneStore{
		aRecords:     make(map[string][]ARecord),
		aaaaRecords:  make(map[string][]AAAARecord),
		cnameRecords: make(map[string]CNAMERecord),
		mxRecords:    make(map[string][]MXRecord),
		nsRecords:    make(map[string][]NSRecord),
		txtRecords:   make(map[string][]TXTRecord),
		syncPort:     8081,
	}
}

func (zs *ZoneStore) StartSyncServer(port int) {
	zs.syncPort = port
	http.HandleFunc("/api/v1/zone/sync", zs.handleZoneSync)
	http.HandleFunc("/api/v1/zone/data", zs.handleZoneData)

	go func() {
		log.Printf("Zone sync server starting on :%d", zs.syncPort)
		if err := http.ListenAndServe(fmt.Sprintf(":%d", zs.syncPort), nil); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start zone sync server: %s", err.Error())
		}
	}()
}

func (zs *ZoneStore) handleZoneSync(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var data ZoneData
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	zs.Load(data)
	log.Println("Zone data synced from remote")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func (zs *ZoneStore) handleZoneData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(zs.Export())
}

func (zs *ZoneStore) SyncFromPrimary(primaryAddr string) error {
	resp, err := http.Get(fmt.Sprintf("http://%s/api/v1/zone/data", primaryAddr))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("primary returned status %d", resp.StatusCode)
	}

	var data ZoneData
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return err
	}

	zs.Load(data)
	log.Printf("Zone data synced from primary: %s", primaryAddr)
	return nil
}

func (zs *ZoneStore) SyncToBackup(backupAddr string) error {
	data := zs.Export()
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	resp, err := http.Post(fmt.Sprintf("http://%s/api/v1/zone/sync", backupAddr), "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("backup returned status %d", resp.StatusCode)
	}

	log.Printf("Zone data synced to backup: %s", backupAddr)
	return nil
}

func (zs *ZoneStore) AddA(name string, ip net.IP, ttl uint32) {
	zs.mu.Lock()
	defer zs.mu.Unlock()
	zs.aRecords[name] = append(zs.aRecords[name], ARecord{IP: ip, TTL: ttl})
}

func (zs *ZoneStore) GetA(name string) ([]ARecord, bool) {
	zs.mu.RLock()
	defer zs.mu.RUnlock()
	records, ok := zs.aRecords[name]
	return records, ok
}

func (zs *ZoneStore) AddAAAA(name string, ip net.IP, ttl uint32) {
	zs.mu.Lock()
	defer zs.mu.Unlock()
	zs.aaaaRecords[name] = append(zs.aaaaRecords[name], AAAARecord{IP: ip, TTL: ttl})
}

func (zs *ZoneStore) GetAAAA(name string) ([]AAAARecord, bool) {
	zs.mu.RLock()
	defer zs.mu.RUnlock()
	records, ok := zs.aaaaRecords[name]
	return records, ok
}

func (zs *ZoneStore) AddCNAME(name string, target string, ttl uint32) {
	zs.mu.Lock()
	defer zs.mu.Unlock()
	zs.cnameRecords[name] = CNAMERecord{Target: target, TTL: ttl}
}

func (zs *ZoneStore) GetCNAME(name string) (CNAMERecord, bool) {
	zs.mu.RLock()
	defer zs.mu.RUnlock()
	record, ok := zs.cnameRecords[name]
	return record, ok
}

func (zs *ZoneStore) AddMX(name string, preference uint16, mx string, ttl uint32) {
	zs.mu.Lock()
	defer zs.mu.Unlock()
	zs.mxRecords[name] = append(zs.mxRecords[name], MXRecord{Preference: preference, MX: mx, TTL: ttl})
}

func (zs *ZoneStore) GetMX(name string) ([]MXRecord, bool) {
	zs.mu.RLock()
	defer zs.mu.RUnlock()
	records, ok := zs.mxRecords[name]
	return records, ok
}

func (zs *ZoneStore) AddNS(name string, ns string, ttl uint32) {
	zs.mu.Lock()
	defer zs.mu.Unlock()
	zs.nsRecords[name] = append(zs.nsRecords[name], NSRecord{NS: ns, TTL: ttl})
}

func (zs *ZoneStore) GetNS(name string) ([]NSRecord, bool) {
	zs.mu.RLock()
	defer zs.mu.RUnlock()
	records, ok := zs.nsRecords[name]
	return records, ok
}

func (zs *ZoneStore) AddTXT(name string, txt []string, ttl uint32) {
	zs.mu.Lock()
	defer zs.mu.Unlock()
	zs.txtRecords[name] = append(zs.txtRecords[name], TXTRecord{Txt: txt, TTL: ttl})
}

func (zs *ZoneStore) GetTXT(name string) ([]TXTRecord, bool) {
	zs.mu.RLock()
	defer zs.mu.RUnlock()
	records, ok := zs.txtRecords[name]
	return records, ok
}

func (zs *ZoneStore) Export() ZoneData {
	zs.mu.RLock()
	defer zs.mu.RUnlock()
	return ZoneData{
		ARecords:     zs.aRecords,
		AAAARecords:  zs.aaaaRecords,
		CNAMERecords: zs.cnameRecords,
		MXRecords:    zs.mxRecords,
		NSRecords:    zs.nsRecords,
		TXTRecords:   zs.txtRecords,
	}
}

func (zs *ZoneStore) Load(data ZoneData) {
	zs.mu.Lock()
	defer zs.mu.Unlock()
	if data.ARecords != nil {
		zs.aRecords = data.ARecords
	}
	if data.AAAARecords != nil {
		zs.aaaaRecords = data.AAAARecords
	}
	if data.CNAMERecords != nil {
		zs.cnameRecords = data.CNAMERecords
	}
	if data.MXRecords != nil {
		zs.mxRecords = data.MXRecords
	}
	if data.NSRecords != nil {
		zs.nsRecords = data.NSRecords
	}
	if data.TXTRecords != nil {
		zs.txtRecords = data.TXTRecords
	}
}

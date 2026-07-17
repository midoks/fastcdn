package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/miekg/dns"
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

type SOARecord struct {
	Ns      string `json:"ns"`
	Mbox    string `json:"mbox"`
	Serial  uint32 `json:"serial"`
	Refresh uint32 `json:"refresh"`
	Retry   uint32 `json:"retry"`
	Expire  uint32 `json:"expire"`
	Minttl  uint32 `json:"minttl"`
	TTL     uint32 `json:"ttl"`
}

type ZoneStore struct {
	mu           sync.RWMutex
	aRecords     map[string][]ARecord
	aaaaRecords  map[string][]AAAARecord
	cnameRecords map[string]CNAMERecord
	mxRecords    map[string][]MXRecord
	nsRecords    map[string][]NSRecord
	txtRecords   map[string][]TXTRecord
	soaRecords   map[string]SOARecord
	syncPort     int
	isRunning    bool
	lastSerial   uint32
}

type ZoneData struct {
	ARecords     map[string][]ARecord    `json:"a_records"`
	AAAARecords  map[string][]AAAARecord `json:"aaaa_records"`
	CNAMERecords map[string]CNAMERecord  `json:"cname_records"`
	MXRecords    map[string][]MXRecord   `json:"mx_records"`
	NSRecords    map[string][]NSRecord   `json:"ns_records"`
	TXTRecords   map[string][]TXTRecord  `json:"txt_records"`
	SOARecords   map[string]SOARecord    `json:"soa_records"`
}

func NewZoneStore() *ZoneStore {
	return &ZoneStore{
		aRecords:     make(map[string][]ARecord),
		aaaaRecords:  make(map[string][]AAAARecord),
		cnameRecords: make(map[string]CNAMERecord),
		mxRecords:    make(map[string][]MXRecord),
		nsRecords:    make(map[string][]NSRecord),
		txtRecords:   make(map[string][]TXTRecord),
		soaRecords:   make(map[string]SOARecord),
		syncPort:     8081,
		lastSerial:   uint32(time.Now().Unix()),
	}
}

func (zs *ZoneStore) incrementSerial(zone string) {
	zs.mu.Lock()
	defer zs.mu.Unlock()

	if soa, ok := zs.soaRecords[zone]; ok {
		soa.Serial++
		zs.soaRecords[zone] = soa
	} else {
		zs.soaRecords[zone] = SOARecord{
			Ns:      "ns1." + zone,
			Mbox:    "admin." + zone,
			Serial:  zs.lastSerial + 1,
			Refresh: 3600,
			Retry:   900,
			Expire:  604800,
			Minttl:  86400,
			TTL:     3600,
		}
		zs.lastSerial++
	}
}

func (zs *ZoneStore) GetSOA(zone string) (SOARecord, bool) {
	zs.mu.RLock()
	defer zs.mu.RUnlock()
	record, ok := zs.soaRecords[zone]
	return record, ok
}

func (zs *ZoneStore) AddSOA(zone string, soa SOARecord) {
	zs.mu.Lock()
	defer zs.mu.Unlock()
	zs.soaRecords[zone] = soa
	if soa.Serial > zs.lastSerial {
		zs.lastSerial = soa.Serial
	}
}

func (zs *ZoneStore) AddA(name string, ip net.IP, ttl uint32) {
	zs.mu.Lock()
	defer zs.mu.Unlock()
	zs.aRecords[name] = append(zs.aRecords[name], ARecord{IP: ip, TTL: ttl})
	zs.incrementSerialInternal(name)
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
	zs.incrementSerialInternal(name)
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
	zs.incrementSerialInternal(name)
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
	zs.incrementSerialInternal(name)
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
	zs.incrementSerialInternal(name)
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
	zs.incrementSerialInternal(name)
}

func (zs *ZoneStore) GetTXT(name string) ([]TXTRecord, bool) {
	zs.mu.RLock()
	defer zs.mu.RUnlock()
	records, ok := zs.txtRecords[name]
	return records, ok
}

func (zs *ZoneStore) incrementSerialInternal(name string) {
	zone := extractZone(name)
	if zone != "" {
		if soa, ok := zs.soaRecords[zone]; ok {
			soa.Serial++
			zs.soaRecords[zone] = soa
			if soa.Serial > zs.lastSerial {
				zs.lastSerial = soa.Serial
			}
		}
	}
}

func extractZone(name string) string {
	parts := dns.SplitDomainName(name)
	if len(parts) >= 2 {
		return dns.Fqdn(parts[len(parts)-2] + "." + parts[len(parts)-1])
	}
	return ""
}

func (zs *ZoneStore) ToRRs(zone string) []dns.RR {
	zs.mu.RLock()
	defer zs.mu.RUnlock()

	var rrs []dns.RR

	if soa, ok := zs.soaRecords[zone]; ok {
		rrs = append(rrs, &dns.SOA{
			Hdr: dns.RR_Header{
				Name:   zone,
				Rrtype: dns.TypeSOA,
				Class:  dns.ClassINET,
				Ttl:    soa.TTL,
			},
			Ns:      soa.Ns,
			Mbox:    soa.Mbox,
			Serial:  soa.Serial,
			Refresh: soa.Refresh,
			Retry:   soa.Retry,
			Expire:  soa.Expire,
			Minttl:  soa.Minttl,
		})
	}

	for name, records := range zs.nsRecords {
		if dns.IsSubDomain(zone, name) || name == zone {
			for _, record := range records {
				rrs = append(rrs, &dns.NS{
					Hdr: dns.RR_Header{
						Name:   name,
						Rrtype: dns.TypeNS,
						Class:  dns.ClassINET,
						Ttl:    record.TTL,
					},
					Ns: record.NS,
				})
			}
		}
	}

	for name, records := range zs.aRecords {
		if dns.IsSubDomain(zone, name) || name == zone {
			for _, record := range records {
				rrs = append(rrs, &dns.A{
					Hdr: dns.RR_Header{
						Name:   name,
						Rrtype: dns.TypeA,
						Class:  dns.ClassINET,
						Ttl:    record.TTL,
					},
					A: record.IP,
				})
			}
		}
	}

	for name, records := range zs.aaaaRecords {
		if dns.IsSubDomain(zone, name) || name == zone {
			for _, record := range records {
				rrs = append(rrs, &dns.AAAA{
					Hdr: dns.RR_Header{
						Name:   name,
						Rrtype: dns.TypeAAAA,
						Class:  dns.ClassINET,
						Ttl:    record.TTL,
					},
					AAAA: record.IP,
				})
			}
		}
	}

	for name, record := range zs.cnameRecords {
		if dns.IsSubDomain(zone, name) || name == zone {
			rrs = append(rrs, &dns.CNAME{
				Hdr: dns.RR_Header{
					Name:   name,
					Rrtype: dns.TypeCNAME,
					Class:  dns.ClassINET,
					Ttl:    record.TTL,
				},
				Target: record.Target,
			})
		}
	}

	for name, records := range zs.mxRecords {
		if dns.IsSubDomain(zone, name) || name == zone {
			for _, record := range records {
				rrs = append(rrs, &dns.MX{
					Hdr: dns.RR_Header{
						Name:   name,
						Rrtype: dns.TypeMX,
						Class:  dns.ClassINET,
						Ttl:    record.TTL,
					},
					Preference: record.Preference,
					Mx:         record.MX,
				})
			}
		}
	}

	for name, records := range zs.txtRecords {
		if dns.IsSubDomain(zone, name) || name == zone {
			for _, record := range records {
				rrs = append(rrs, &dns.TXT{
					Hdr: dns.RR_Header{
						Name:   name,
						Rrtype: dns.TypeTXT,
						Class:  dns.ClassINET,
						Ttl:    record.TTL,
					},
					Txt: record.Txt,
				})
			}
		}
	}

	if soa, ok := zs.soaRecords[zone]; ok {
		rrs = append(rrs, &dns.SOA{
			Hdr: dns.RR_Header{
				Name:   zone,
				Rrtype: dns.TypeSOA,
				Class:  dns.ClassINET,
				Ttl:    soa.TTL,
			},
			Ns:      soa.Ns,
			Mbox:    soa.Mbox,
			Serial:  soa.Serial,
			Refresh: soa.Refresh,
			Retry:   soa.Retry,
			Expire:  soa.Expire,
			Minttl:  soa.Minttl,
		})
	}

	return rrs
}

func (zs *ZoneStore) FromRRs(rrs []dns.RR) {
	zs.mu.Lock()
	defer zs.mu.Unlock()

	for _, rr := range rrs {
		hdr := rr.Header()
		name := hdr.Name
		ttl := hdr.Ttl

		switch rr := rr.(type) {
		case *dns.SOA:
			zs.soaRecords[name] = SOARecord{
				Ns:      rr.Ns,
				Mbox:    rr.Mbox,
				Serial:  rr.Serial,
				Refresh: rr.Refresh,
				Retry:   rr.Retry,
				Expire:  rr.Expire,
				Minttl:  rr.Minttl,
				TTL:     ttl,
			}
			if rr.Serial > zs.lastSerial {
				zs.lastSerial = rr.Serial
			}
		case *dns.NS:
			zs.nsRecords[name] = append(zs.nsRecords[name], NSRecord{
				NS:  rr.Ns,
				TTL: ttl,
			})
		case *dns.A:
			zs.aRecords[name] = append(zs.aRecords[name], ARecord{
				IP:  rr.A,
				TTL: ttl,
			})
		case *dns.AAAA:
			zs.aaaaRecords[name] = append(zs.aaaaRecords[name], AAAARecord{
				IP:  rr.AAAA,
				TTL: ttl,
			})
		case *dns.CNAME:
			zs.cnameRecords[name] = CNAMERecord{
				Target: rr.Target,
				TTL:    ttl,
			}
		case *dns.MX:
			zs.mxRecords[name] = append(zs.mxRecords[name], MXRecord{
				Preference: rr.Preference,
				MX:         rr.Mx,
				TTL:        ttl,
			})
		case *dns.TXT:
			zs.txtRecords[name] = append(zs.txtRecords[name], TXTRecord{
				Txt: rr.Txt,
				TTL: ttl,
			})
		}
	}
}

func (zs *ZoneStore) GetZones() []string {
	zs.mu.RLock()
	defer zs.mu.RUnlock()

	zones := make(map[string]bool)
	for name := range zs.soaRecords {
		zones[name] = true
	}
	for name := range zs.nsRecords {
		zones[extractZone(name)] = true
	}

	result := make([]string, 0, len(zones))
	for zone := range zones {
		if zone != "" {
			result = append(result, zone)
		}
	}
	return result
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
		SOARecords:   zs.soaRecords,
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
	if data.SOARecords != nil {
		zs.soaRecords = data.SOARecords
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

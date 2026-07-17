package main

import (
	"log"
	"time"
)

type ZoneSync struct {
	dnsServer     *DnsServer
	primaryAddr   string
	syncInterval  time.Duration
	isRunning     bool
	stopChan      chan struct{}
	syncType      string
}

func NewZoneSync(dnsServer *DnsServer, primaryAddr string, interval int, syncType string) *ZoneSync {
	return &ZoneSync{
		dnsServer:    dnsServer,
		primaryAddr:  primaryAddr,
		syncInterval: time.Duration(interval) * time.Second,
		syncType:     syncType,
		stopChan:     make(chan struct{}),
	}
}

func (zs *ZoneSync) Start() {
	if zs.primaryAddr == "" {
		log.Println("ZoneSync: No primary address configured, skipping")
		return
	}

	zs.isRunning = true
	log.Printf("ZoneSync starting with interval %ds, sync type: %s", zs.syncInterval/time.Second, zs.syncType)

	go zs.syncLoop()
}

func (zs *ZoneSync) Stop() {
	if !zs.isRunning {
		return
	}

	zs.isRunning = false
	close(zs.stopChan)
	log.Println("ZoneSync stopped")
}

func (zs *ZoneSync) syncLoop() {
	zs.syncOnce()

	ticker := time.NewTicker(zs.syncInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			zs.syncOnce()
		case <-zs.stopChan:
			return
		}
	}
}

func (zs *ZoneSync) syncOnce() {
	zones := zs.dnsServer.zoneStore.GetZones()

	if len(zones) == 0 {
		zones = []string{"example.com."}
	}

	for _, zone := range zones {
		switch zs.syncType {
		case "ixfr":
			if err := zs.dnsServer.IXFRFromPrimary(zs.primaryAddr, zone); err != nil {
				log.Printf("[ZoneSync] IXFR failed for zone %s: %s, falling back to AXFR", zone, err.Error())
				if err := zs.dnsServer.AXFRFromPrimary(zs.primaryAddr, zone); err != nil {
					log.Printf("[ZoneSync] AXFR also failed for zone %s: %s", zone, err.Error())
				}
			} else {
				log.Printf("[ZoneSync] IXFR completed for zone %s", zone)
			}
		case "axfr":
			if err := zs.dnsServer.AXFRFromPrimary(zs.primaryAddr, zone); err != nil {
				log.Printf("[ZoneSync] AXFR failed for zone %s: %s", zone, err.Error())
			} else {
				log.Printf("[ZoneSync] AXFR completed for zone %s", zone)
			}
		default:
			if err := zs.dnsServer.IXFRFromPrimary(zs.primaryAddr, zone); err != nil {
				log.Printf("[ZoneSync] IXFR failed for zone %s: %s", zone, err.Error())
			} else {
				log.Printf("[ZoneSync] IXFR completed for zone %s", zone)
			}
		}
	}
}

func (zs *ZoneSync) SyncNow() {
	log.Println("[ZoneSync] Manual sync triggered")
	zs.syncOnce()
}
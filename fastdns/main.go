package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	config := ParseConfig()

	log.Printf("Starting fastdns with role: %s", config.Role)
	log.Printf("Bind address: %s:%d", config.BindAddr, config.Port)

	dnsServer := NewDnsServer(config)
	if err := dnsServer.Start(); err != nil {
		log.Fatalf("Failed to start DNS server: %s", err.Error())
	}

	roleManager := NewRoleManager(config)

	dnsServer.zoneStore.StartSyncServer(config.SyncPort)

	var zoneSync *ZoneSync
	if config.Role == "backup" && config.PrimaryAddr != "" {
		zoneSync = NewZoneSync(dnsServer, config.PrimaryAddr, config.SyncInterval, config.SyncType)
		zoneSync.Start()
	}

	registerDefaultRecords(dnsServer.zoneStore)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for role := range roleManager.Subscribe() {
			log.Printf("Role changed: %s", role)
		}
	}()

	log.Println("fastdns started successfully")

	<-sigChan

	log.Println("Shutting down fastdns...")
	if zoneSync != nil {
		zoneSync.Stop()
	}
	dnsServer.Stop()
	log.Println("fastdns stopped")
}

func registerDefaultRecords(zs *ZoneStore) {
	zs.AddSOA("example.com.", SOARecord{
		Ns:      "ns1.example.com.",
		Mbox:    "admin.example.com.",
		Serial:  2024010100,
		Refresh: 3600,
		Retry:   900,
		Expire:  604800,
		Minttl:  86400,
		TTL:     3600,
	})

	zs.AddNS("example.com.", "ns1.example.com.", 3600)
	zs.AddNS("example.com.", "ns2.example.com.", 3600)
	zs.AddA("example.com.", net.ParseIP("192.168.1.100"), 3600)
	zs.AddA("example.com.", net.ParseIP("192.168.1.101"), 3600)
	zs.AddAAAA("example.com.", net.ParseIP("::1"), 3600)
	zs.AddCNAME("www.example.com.", "example.com.", 3600)
	zs.AddMX("example.com.", 10, "mail.example.com.", 3600)
	zs.AddTXT("example.com.", []string{"v=spf1 include:_spf.example.com ~all"}, 3600)

	log.Println("Default DNS records registered")
}
package main

import (
	"flag"
)

type Config struct {
	Role               string
	BindAddr           string
	Port               int
	HttpPort           int
	SyncPort           int
	PrimaryAddr        string
	BackupAddr         string
	HealthCheckInterval int
}

func ParseConfig() *Config {
	config := &Config{}

	flag.StringVar(&config.Role, "role", "primary", "DNS role: primary or backup")
	flag.StringVar(&config.BindAddr, "bind-addr", "0.0.0.0", "Bind address")
	flag.IntVar(&config.Port, "port", 53, "DNS port")
	flag.IntVar(&config.HttpPort, "http-port", 8080, "HTTP management port")
	flag.IntVar(&config.SyncPort, "sync-port", 8081, "Zone sync port")
	flag.StringVar(&config.PrimaryAddr, "primary-addr", "", "Primary DNS address (for backup)")
	flag.StringVar(&config.BackupAddr, "backup-addr", "", "Backup DNS address (for primary)")
	flag.IntVar(&config.HealthCheckInterval, "health-check-interval", 10, "Health check interval in seconds")

	flag.Parse()

	return config
}
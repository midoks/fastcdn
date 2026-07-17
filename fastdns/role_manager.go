package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

type Role string

const (
	RolePrimary   Role = "primary"
	RoleBackup    Role = "backup"
	RoleUnknown   Role = "unknown"
)

type RoleManager struct {
	mu           sync.RWMutex
	role         Role
	primaryAddr  string
	backupAddr   string
	healthCheck  time.Duration
	statusChan   chan Role
	httpPort     int
	isRunning    bool
}

type HealthResponse struct {
	Role     Role   `json:"role"`
	Status   string `json:"status"`
	ZoneCount int   `json:"zone_count"`
}

func NewRoleManager(config *Config) *RoleManager {
	rm := &RoleManager{
		role:        Role(config.Role),
		primaryAddr: config.PrimaryAddr,
		backupAddr:  config.BackupAddr,
		healthCheck: time.Duration(config.HealthCheckInterval) * time.Second,
		statusChan:  make(chan Role, 1),
		httpPort:    config.HttpPort,
	}

	if rm.role == RoleBackup && rm.primaryAddr != "" {
		rm.startHealthCheck()
	}

	rm.startHTTPServer()

	return rm
}

func (rm *RoleManager) startHTTPServer() {
	http.HandleFunc("/api/v1/health", rm.handleHealth)
	http.HandleFunc("/api/v1/role", rm.handleRole)
	http.HandleFunc("/api/v1/promote", rm.handlePromote)

	go func() {
		log.Printf("Role manager HTTP server starting on :%d", rm.httpPort)
		if err := http.ListenAndServe(fmt.Sprintf(":%d", rm.httpPort), nil); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start role manager HTTP server: %s", err.Error())
		}
	}()
}

func (rm *RoleManager) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(HealthResponse{
		Role:     rm.GetRole(),
		Status:   "healthy",
		ZoneCount: 0,
	})
}

func (rm *RoleManager) handleRole(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"role": string(rm.GetRole()),
	})
}

func (rm *RoleManager) handlePromote(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	rm.SetRole(RolePrimary)
	log.Println("Promoted to primary role")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "success",
		"role":   string(rm.GetRole()),
	})
}

func (rm *RoleManager) startHealthCheck() {
	ticker := time.NewTicker(rm.healthCheck)
	go func() {
		for range ticker.C {
			if rm.isPrimaryHealthy() {
				log.Printf("Primary %s is healthy, staying as backup", rm.primaryAddr)
				rm.SetRole(RoleBackup)
			} else {
				log.Printf("Primary %s is unhealthy, promoting to primary", rm.primaryAddr)
				rm.SetRole(RolePrimary)
			}
		}
	}()
}

func (rm *RoleManager) isPrimaryHealthy() bool {
	if rm.primaryAddr == "" {
		return false
	}

	resp, err := http.Get(fmt.Sprintf("http://%s/api/v1/health", rm.primaryAddr))
	if err != nil {
		log.Printf("Health check failed: %s", err.Error())
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}

func (rm *RoleManager) GetRole() Role {
	rm.mu.RLock()
	defer rm.mu.RUnlock()
	return rm.role
}

func (rm *RoleManager) SetRole(role Role) {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	
	if rm.role != role {
		rm.role = role
		select {
		case rm.statusChan <- role:
		default:
		}
		log.Printf("Role changed to: %s", role)
	}
}

func (rm *RoleManager) Subscribe() <-chan Role {
	return rm.statusChan
}

func (rm *RoleManager) IsPrimary() bool {
	return rm.GetRole() == RolePrimary
}

func (rm *RoleManager) IsBackup() bool {
	return rm.GetRole() == RoleBackup
}
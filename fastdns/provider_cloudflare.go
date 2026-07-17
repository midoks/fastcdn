package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type CloudflareProvider struct {
	apiToken string
	zoneID   string
	httpClient *http.Client
}

type CloudflareRecord struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Type       string `json:"type"`
	Content    string `json:"content"`
	TTL        int    `json:"ttl"`
	Priority   *int   `json:"priority,omitempty"`
	Proxied    bool   `json:"proxied"`
}

type CloudflareListResponse struct {
	Result []CloudflareRecord `json:"result"`
}

type CloudflareResponse struct {
	Result CloudflareRecord `json:"result"`
}

func NewCloudflareProvider(apiToken, zoneID string) *CloudflareProvider {
	return &CloudflareProvider{
		apiToken: apiToken,
		zoneID:   zoneID,
		httpClient: &http.Client{
			Timeout: 30,
		},
	}
}

func (p *CloudflareProvider) Name() string {
	return "cloudflare"
}

func (p *CloudflareProvider) ListRecords(zone string) ([]DnsRecord, error) {
	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records", p.zoneID)
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	
	req.Header.Set("Authorization", "Bearer "+p.apiToken)
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Cloudflare API returned status %d", resp.StatusCode)
	}
	
	var cfResp CloudflareListResponse
	if err := json.NewDecoder(resp.Body).Decode(&cfResp); err != nil {
		return nil, err
	}
	
	var records []DnsRecord
	for _, r := range cfResp.Result {
		priority := uint16(0)
		if r.Priority != nil {
			priority = uint16(*r.Priority)
		}
		
		records = append(records, DnsRecord{
			ID:         r.ID,
			Name:       r.Name,
			Type:       RecordType(r.Type),
			Content:    r.Content,
			TTL:        uint32(r.TTL),
			Priority:   priority,
			Proxied:    r.Proxied,
		})
	}
	
	return records, nil
}

func (p *CloudflareProvider) CreateRecord(zone string, record DnsRecord) (string, error) {
	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records", p.zoneID)
	
	priority := int(record.Priority)
	var priorityPtr *int
	if record.Type == TypeMX {
		priorityPtr = &priority
	}
	
	data := map[string]interface{}{
		"name":    record.Name,
		"type":    string(record.Type),
		"content": record.Content,
		"ttl":     record.TTL,
		"proxied": record.Proxied,
	}
	
	if priorityPtr != nil {
		data["priority"] = *priorityPtr
	}
	
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	
	req.Header.Set("Authorization", "Bearer "+p.apiToken)
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := p.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("Cloudflare API returned status %d", resp.StatusCode)
	}
	
	var cfResp CloudflareResponse
	if err := json.NewDecoder(resp.Body).Decode(&cfResp); err != nil {
		return "", err
	}
	
	return cfResp.Result.ID, nil
}

func (p *CloudflareProvider) UpdateRecord(zone string, record DnsRecord) error {
	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records/%s", p.zoneID, record.ID)
	
	priority := int(record.Priority)
	var priorityPtr *int
	if record.Type == TypeMX {
		priorityPtr = &priority
	}
	
	data := map[string]interface{}{
		"name":    record.Name,
		"type":    string(record.Type),
		"content": record.Content,
		"ttl":     record.TTL,
		"proxied": record.Proxied,
	}
	
	if priorityPtr != nil {
		data["priority"] = *priorityPtr
	}
	
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	
	req.Header.Set("Authorization", "Bearer "+p.apiToken)
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := p.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Cloudflare API returned status %d", resp.StatusCode)
	}
	
	return nil
}

func (p *CloudflareProvider) DeleteRecord(zone string, recordID string) error {
	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records/%s", p.zoneID, recordID)
	
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}
	
	req.Header.Set("Authorization", "Bearer "+p.apiToken)
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := p.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Cloudflare API returned status %d", resp.StatusCode)
	}
	
	return nil
}

func (p *CloudflareProvider) PullRecords(zone string) ([]DnsRecord, error) {
	return p.ListRecords(zone)
}

func (p *CloudflareProvider) PushRecords(zone string, records []DnsRecord) error {
	remoteRecords, err := p.ListRecords(zone)
	if err != nil {
		return err
	}
	
	remoteMap := make(map[string]DnsRecord)
	for _, r := range remoteRecords {
		key := r.Name + "_" + string(r.Type) + "_" + r.Content
		remoteMap[key] = r
	}
	
	for _, local := range records {
		key := local.Name + "_" + string(local.Type) + "_" + local.Content
		
		if remote, ok := remoteMap[key]; ok {
			delete(remoteMap, key)
			
			if remote.TTL != local.TTL || remote.Proxied != local.Proxied || remote.Priority != local.Priority {
				local.ID = remote.ID
				if err := p.UpdateRecord(zone, local); err != nil {
					log.Printf("Failed to update record %s: %s", key, err.Error())
				} else {
					log.Printf("Updated record: %s", key)
				}
			}
		} else {
			if _, err := p.CreateRecord(zone, local); err != nil {
				log.Printf("Failed to create record %s: %s", key, err.Error())
			} else {
				log.Printf("Created record: %s", key)
			}
		}
	}
	
	for _, remote := range remoteMap {
		if err := p.DeleteRecord(zone, remote.ID); err != nil {
			log.Printf("Failed to delete record %s: %s", remote.Name, err.Error())
		} else {
			log.Printf("Deleted record: %s_%s", remote.Name, remote.Type)
		}
	}
	
	return nil
}

func (p *CloudflareProvider) GetZoneID(domain string) (string, error) {
	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones?name=%s", domain)
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	
	req.Header.Set("Authorization", "Bearer "+p.apiToken)
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := p.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Cloudflare API returned status %d", resp.StatusCode)
	}
	
	var result struct {
		Result []struct {
			ID string `json:"id"`
		} `json:"result"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	
	if len(result.Result) == 0 {
		return "", fmt.Errorf("zone not found for domain %s", domain)
	}
	
	return result.Result[0].ID, nil
}

func (p *CloudflareProvider) ParseTTL(ttlStr string) uint32 {
	ttl, err := strconv.Atoi(ttlStr)
	if err != nil {
		return 3600
	}
	return uint32(ttl)
}

func (p *CloudflareProvider) FormatName(name, zone string) string {
	if strings.HasSuffix(name, "."+zone) {
		return name
	}
	if strings.HasSuffix(name, zone) {
		return name + "."
	}
	return name + "." + zone + "."
}
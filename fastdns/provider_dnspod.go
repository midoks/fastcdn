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

type DNSPodProvider struct {
	apiToken string
	domain   string
	httpClient *http.Client
}

type DNSPodRecord struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	Content  string `json:"value"`
	TTL      int    `json:"ttl"`
	Priority int    `json:"mx"`
}

type DNSPodListResponse struct {
	Status struct {
		Code string `json:"code"`
		Message string `json:"message"`
	} `json:"status"`
	Records []DNSPodRecord `json:"records"`
}

type DNSPodResponse struct {
	Status struct {
		Code string `json:"code"`
		Message string `json:"message"`
	} `json:"status"`
	Record DNSPodRecord `json:"record"`
}

func NewDNSPodProvider(apiToken, domain string) *DNSPodProvider {
	return &DNSPodProvider{
		apiToken: apiToken,
		domain:   domain,
		httpClient: &http.Client{
			Timeout: 30,
		},
	}
}

func (p *DNSPodProvider) Name() string {
	return "dnspod"
}

func (p *DNSPodProvider) ListRecords(zone string) ([]DnsRecord, error) {
	url := fmt.Sprintf("https://dnsapi.cn/Record.List")
	
	domain := p.domain
	if zone != "" {
		domain = zone
	}
	
	data := fmt.Sprintf("login_token=%s&format=json&domain=%s", p.apiToken, domain)
	
	req, err := http.NewRequest("POST", url, bytes.NewBufferString(data))
	if err != nil {
		return nil, err
	}
	
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	
	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	var dpResp DNSPodListResponse
	if err := json.NewDecoder(resp.Body).Decode(&dpResp); err != nil {
		return nil, err
	}
	
	if dpResp.Status.Code != "1" {
		return nil, fmt.Errorf("DNSPod API returned error: %s", dpResp.Status.Message)
	}
	
	var records []DnsRecord
	for _, r := range dpResp.Records {
		records = append(records, DnsRecord{
			ID:         strconv.Itoa(r.ID),
			Name:       r.Name + "." + domain + ".",
			Type:       RecordType(r.Type),
			Content:    r.Content,
			TTL:        uint32(r.TTL),
			Priority:   uint16(r.Priority),
			Proxied:    false,
		})
	}
	
	return records, nil
}

func (p *DNSPodProvider) CreateRecord(zone string, record DnsRecord) (string, error) {
	url := fmt.Sprintf("https://dnsapi.cn/Record.Create")
	
	domain := p.domain
	if zone != "" {
		domain = zone
	}
	
	name := record.Name
	if strings.HasSuffix(name, "."+domain+".") {
		name = strings.TrimSuffix(name, "."+domain+".")
	} else if strings.HasSuffix(name, "."+domain) {
		name = strings.TrimSuffix(name, "."+domain)
	} else if strings.HasSuffix(name, ".") {
		name = strings.TrimSuffix(name, ".")
	}
	
	data := fmt.Sprintf(
		"login_token=%s&format=json&domain=%s&sub_domain=%s&record_type=%s&record_line=默认&value=%s&ttl=%d",
		p.apiToken, domain, name, string(record.Type), record.Content, record.TTL,
	)
	
	if record.Type == TypeMX {
		data += fmt.Sprintf("&mx=%d", record.Priority)
	}
	
	req, err := http.NewRequest("POST", url, bytes.NewBufferString(data))
	if err != nil {
		return "", err
	}
	
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	
	resp, err := p.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	
	var dpResp DNSPodResponse
	if err := json.NewDecoder(resp.Body).Decode(&dpResp); err != nil {
		return "", err
	}
	
	if dpResp.Status.Code != "1" {
		return "", fmt.Errorf("DNSPod API returned error: %s", dpResp.Status.Message)
	}
	
	return strconv.Itoa(dpResp.Record.ID), nil
}

func (p *DNSPodProvider) UpdateRecord(zone string, record DnsRecord) error {
	url := fmt.Sprintf("https://dnsapi.cn/Record.Modify")
	
	domain := p.domain
	if zone != "" {
		domain = zone
	}
	
	name := record.Name
	if strings.HasSuffix(name, "."+domain+".") {
		name = strings.TrimSuffix(name, "."+domain+".")
	} else if strings.HasSuffix(name, "."+domain) {
		name = strings.TrimSuffix(name, "."+domain)
	} else if strings.HasSuffix(name, ".") {
		name = strings.TrimSuffix(name, ".")
	}
	
	data := fmt.Sprintf(
		"login_token=%s&format=json&domain=%s&record_id=%s&sub_domain=%s&record_type=%s&record_line=默认&value=%s&ttl=%d",
		p.apiToken, domain, record.ID, name, string(record.Type), record.Content, record.TTL,
	)
	
	if record.Type == TypeMX {
		data += fmt.Sprintf("&mx=%d", record.Priority)
	}
	
	req, err := http.NewRequest("POST", url, bytes.NewBufferString(data))
	if err != nil {
		return err
	}
	
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	
	resp, err := p.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	var dpResp DNSPodResponse
	if err := json.NewDecoder(resp.Body).Decode(&dpResp); err != nil {
		return err
	}
	
	if dpResp.Status.Code != "1" {
		return fmt.Errorf("DNSPod API returned error: %s", dpResp.Status.Message)
	}
	
	return nil
}

func (p *DNSPodProvider) DeleteRecord(zone string, recordID string) error {
	url := fmt.Sprintf("https://dnsapi.cn/Record.Remove")
	
	domain := p.domain
	if zone != "" {
		domain = zone
	}
	
	data := fmt.Sprintf("login_token=%s&format=json&domain=%s&record_id=%s", p.apiToken, domain, recordID)
	
	req, err := http.NewRequest("POST", url, bytes.NewBufferString(data))
	if err != nil {
		return err
	}
	
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	
	resp, err := p.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	var dpResp DNSPodResponse
	if err := json.NewDecoder(resp.Body).Decode(&dpResp); err != nil {
		return err
	}
	
	if dpResp.Status.Code != "1" {
		return fmt.Errorf("DNSPod API returned error: %s", dpResp.Status.Message)
	}
	
	return nil
}

func (p *DNSPodProvider) PullRecords(zone string) ([]DnsRecord, error) {
	return p.ListRecords(zone)
}

func (p *DNSPodProvider) PushRecords(zone string, records []DnsRecord) error {
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
			
			if remote.TTL != local.TTL || remote.Priority != local.Priority {
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

func (p *DNSPodProvider) GetDomainID(domain string) (string, error) {
	url := fmt.Sprintf("https://dnsapi.cn/Domain.Info")
	
	data := fmt.Sprintf("login_token=%s&format=json&domain=%s", p.apiToken, domain)
	
	req, err := http.NewRequest("POST", url, bytes.NewBufferString(data))
	if err != nil {
		return "", err
	}
	
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	
	resp, err := p.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	
	var result struct {
		Status struct {
			Code string `json:"code"`
		} `json:"status"`
		Domain struct {
			ID int `json:"id"`
		} `json:"domain"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	
	if result.Status.Code != "1" {
		return "", fmt.Errorf("failed to get domain ID")
	}
	
	return strconv.Itoa(result.Domain.ID), nil
}
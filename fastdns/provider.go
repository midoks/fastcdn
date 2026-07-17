package main

import (
	"net"
)

type RecordType string

const (
	TypeA     RecordType = "A"
	TypeAAAA  RecordType = "AAAA"
	TypeCNAME RecordType = "CNAME"
	TypeMX    RecordType = "MX"
	TypeNS    RecordType = "NS"
	TypeTXT   RecordType = "TXT"
)

type DnsRecord struct {
	ID         string
	Name       string
	Type       RecordType
	Content    string
	TTL        uint32
	Priority   uint16
	Proxied    bool
}

type DnsProvider interface {
	Name() string
	ListRecords(zone string) ([]DnsRecord, error)
	CreateRecord(zone string, record DnsRecord) (string, error)
	UpdateRecord(zone string, record DnsRecord) error
	DeleteRecord(zone string, recordID string) error
	PullRecords(zone string) ([]DnsRecord, error)
	PushRecords(zone string, records []DnsRecord) error
}

func IsPublicIP(ip net.IP) bool {
	if ip == nil {
		return false
	}
	
	if ip.IsLoopback() {
		return false
	}
	
	if ip.IsLinkLocalUnicast() {
		return false
	}
	
	if ip.IsLinkLocalMulticast() {
		return false
	}
	
	if ip.IsInterfaceLocalMulticast() {
		return false
	}
	
	if ip.IsPrivate() {
		return false
	}
	
	if !ip.IsGlobalUnicast() {
		return false
	}
	
	return true
}
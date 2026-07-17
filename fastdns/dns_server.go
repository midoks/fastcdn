package main

import (
	"fmt"
	"log"
	"net"

	"github.com/miekg/dns"
)

type DnsServer struct {
	udpAddr   *net.UDPAddr
	tcpAddr   *net.TCPAddr
	udpServer *dns.Server
	tcpServer *dns.Server
	zoneStore *ZoneStore
	role      string
}

func NewDnsServer(config *Config) *DnsServer {
	udpAddr, _ := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", config.BindAddr, config.Port))
	tcpAddr, _ := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", config.BindAddr, config.Port))

	return &DnsServer{
		udpAddr:   udpAddr,
		tcpAddr:   tcpAddr,
		zoneStore: NewZoneStore(),
		role:      config.Role,
	}
}

func (s *DnsServer) Start() error {
	mux := dns.NewServeMux()
	mux.HandleFunc(".", s.handleRequest)

	s.udpServer = &dns.Server{
		Addr:    s.udpAddr.String(),
		Net:     "udp",
		Handler: mux,
		UDPSize: 65535,
	}

	s.tcpServer = &dns.Server{
		Addr:    s.tcpAddr.String(),
		Net:     "tcp",
		Handler: mux,
	}

	go func() {
		log.Printf("DNS server starting on %s (UDP)", s.udpAddr.String())
		if err := s.udpServer.ListenAndServe(); err != nil {
			log.Fatalf("Failed to start UDP server: %s\n", err.Error())
		}
	}()

	go func() {
		log.Printf("DNS server starting on %s (TCP)", s.tcpAddr.String())
		if err := s.tcpServer.ListenAndServe(); err != nil {
			log.Fatalf("Failed to start TCP server: %s\n", err.Error())
		}
	}()

	return nil
}

func (s *DnsServer) handleRequest(w dns.ResponseWriter, r *dns.Msg) {
	log.Printf("[%s] Query: %s %s from %s", s.role, dns.TypeToString[r.Question[0].Qtype], r.Question[0].Name, w.RemoteAddr())

	m := new(dns.Msg)
	m.SetReply(r)
	m.Authoritative = true

	for _, q := range r.Question {
		switch q.Qtype {
		case dns.TypeA:
			if records, ok := s.zoneStore.GetA(q.Name); ok {
				for _, record := range records {
					rr := &dns.A{
						Hdr: dns.RR_Header{
							Name:   q.Name,
							Rrtype: dns.TypeA,
							Class:  dns.ClassINET,
							Ttl:    record.TTL,
						},
						A: record.IP,
					}
					m.Answer = append(m.Answer, rr)
				}
			}
		case dns.TypeAAAA:
			if records, ok := s.zoneStore.GetAAAA(q.Name); ok {
				for _, record := range records {
					rr := &dns.AAAA{
						Hdr: dns.RR_Header{
							Name:   q.Name,
							Rrtype: dns.TypeAAAA,
							Class:  dns.ClassINET,
							Ttl:    record.TTL,
						},
						AAAA: record.IP,
					}
					m.Answer = append(m.Answer, rr)
				}
			}
		case dns.TypeCNAME:
			if record, ok := s.zoneStore.GetCNAME(q.Name); ok {
				rr := &dns.CNAME{
					Hdr: dns.RR_Header{
						Name:   q.Name,
						Rrtype: dns.TypeCNAME,
						Class:  dns.ClassINET,
						Ttl:    record.TTL,
					},
					Target: record.Target,
				}
				m.Answer = append(m.Answer, rr)
			}
		case dns.TypeMX:
			if records, ok := s.zoneStore.GetMX(q.Name); ok {
				for _, record := range records {
					rr := &dns.MX{
						Hdr: dns.RR_Header{
							Name:   q.Name,
							Rrtype: dns.TypeMX,
							Class:  dns.ClassINET,
							Ttl:    record.TTL,
						},
						Preference: record.Preference,
						Mx:         record.MX,
					}
					m.Answer = append(m.Answer, rr)
				}
			}
		case dns.TypeNS:
			if records, ok := s.zoneStore.GetNS(q.Name); ok {
				for _, record := range records {
					rr := &dns.NS{
						Hdr: dns.RR_Header{
							Name:   q.Name,
							Rrtype: dns.TypeNS,
							Class:  dns.ClassINET,
							Ttl:    record.TTL,
						},
						Ns: record.NS,
					}
					m.Answer = append(m.Answer, rr)
				}
			}
		case dns.TypeTXT:
			if records, ok := s.zoneStore.GetTXT(q.Name); ok {
				for _, record := range records {
					rr := &dns.TXT{
						Hdr: dns.RR_Header{
							Name:   q.Name,
							Rrtype: dns.TypeTXT,
							Class:  dns.ClassINET,
							Ttl:    record.TTL,
						},
						Txt: record.Txt,
					}
					m.Answer = append(m.Answer, rr)
				}
			}
		}
	}

	if len(m.Answer) == 0 {
		m.SetRcode(r, dns.RcodeNameError)
	}

	w.WriteMsg(m)
}

func (s *DnsServer) Stop() {
	if s.udpServer != nil {
		s.udpServer.Shutdown()
	}
	if s.tcpServer != nil {
		s.tcpServer.Shutdown()
	}
}

func (s *DnsServer) SendQuery(server string, domain string, qtype uint16) (*dns.Msg, error) {
	m := new(dns.Msg)
	m.SetQuestion(domain, qtype)

	c := new(dns.Client)
	reply, _, err := c.Exchange(m, server+":53")
	if err != nil {
		return nil, err
	}

	return reply, nil
}

func (s *DnsServer) SendUpdate(server string, zone string, records []dns.RR) error {
	m := new(dns.Msg)
	m.SetUpdate(zone)

	for _, rr := range records {
		m.Ns = append(m.Ns, rr)
	}

	c := new(dns.Client)
	_, _, err := c.Exchange(m, server+":53")
	return err
}
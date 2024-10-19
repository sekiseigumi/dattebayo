// Fully functional DNS server that can be used to resolve DNS queries on localhost.
package dns

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/miekg/dns"
	"github.com/sekiseigumi/dattebayo/internal/logger"
	"github.com/sekiseigumi/dattebayo/shared"
)

type RecordType int

const (
	A RecordType = iota
	CNAME
	MX
	TXT
)

const LOGGER_SOURCE = "DNS SERVER"

type DomainRecord struct {
	Type     RecordType
	Value    string
	Priority int // Only used for MX records
}

type Domain struct {
	Name           string
	Records        map[string]DomainRecord // Key is subdomain, @ for root
	IsSystemDomain bool
}

type DNSServer struct {
	PrimaryPort  int
	FallbackPort int
	running      bool
	currentPort  int
	domains      map[string]*Domain
	TLDs         map[string]bool
	mu           sync.RWMutex
	logger       *logger.Logger
}

func NewDNSServer(config shared.Config, log *logger.Logger) *DNSServer {
	primaryPort := func() int {
		if config.DNS.PrimaryPort == 0 {
			return 53
		}

		return config.DNS.PrimaryPort
	}()

	fallbackPort := func() int {
		if config.DNS.FallbackPort == 0 {
			return 53535
		}

		return config.DNS.FallbackPort
	}()

	s := &DNSServer{
		PrimaryPort:  primaryPort,
		FallbackPort: fallbackPort,
		domains:      make(map[string]*Domain),
		TLDs:         make(map[string]bool),
		logger:       log,
	}

	// supportedTLDs := []string{
	// 	"ao", "ara", "epic", "fcuk", "internal", "ki", "local", "localhost",
	// 	"lore", "mail", "mi", "myth", "neko", "os", "pwn", "root", "test",
	// 	"thc", "waifu",
	// }

	// for _, tld := range supportedTLDs {
	// 	s.TLDs[tld] = true
	// }

	// // System domains
	// s.AddDomain("domains.internal", true)
	// s.AddSubdomain("domains.internal", "@", []DomainRecord{{Type: A, Value: "127.0.0.1"}})
	// s.AddSubdomain("domains.internal", "api", []DomainRecord{{Type: A, Value: "127.0.0.1"}})

	// s.AddDomain("admin.mail", true)
	// s.AddSubdomain("admin.mail", "@", []DomainRecord{{Type: A, Value: "127.0.0.1"}})

	// s.AddDomain("user.mail", true)
	// s.AddSubdomain("user.mail", "@", []DomainRecord{{Type: A, Value: "127.0.0.1"}})

	s.initializeTLDs()
	s.initializeSystemDomains()

	return s
}

func (s *DNSServer) initializeTLDs() {
	supportedTLDs := []string{
		"ao", "ara", "epic", "fuck", "internal", "ki", "local", "localhost",
		"lore", "mail", "mi", "myth", "neko", "os", "pwn", "root", "test",
		"thc", "waifu",
	}
	for _, tld := range supportedTLDs {
		s.TLDs[tld] = true
	}
}

func (s *DNSServer) initializeSystemDomains() {
	systemDomains := map[string][]string{
		"domains.internal": {"@", "api"},
		"admin.mail":       {"@"},
		"user.mail":        {"@"},
	}
	for domain, subdomains := range systemDomains {
		s.AddDomain(domain, true)
		for _, subdomain := range subdomains {
			s.AddSubdomain(domain, subdomain, []DomainRecord{{Type: A, Value: "127.0.0.1"}})
		}
	}
}

func (s *DNSServer) Start() error {
	dns.HandleFunc(".", s.handleDNS)

	if err := s.startOnPort(s.PrimaryPort); err == nil {
		return nil // Successfully started on primary port
	}

	// If primary port failed, try fallback port
	s.logger.Log(LOGGER_SOURCE, fmt.Sprintf("Falling back to port %d", s.FallbackPort))
	if err := s.startOnPort(s.FallbackPort); err == nil {
		return nil // Successfully started on fallback port
	}

	return fmt.Errorf("failed to start DNS server on both primary (%d) and fallback (%d) ports", s.PrimaryPort, s.FallbackPort)
}

func (s *DNSServer) startOnPort(port int) error {
	s.logger.Log(LOGGER_SOURCE, fmt.Sprintf("Attempting to start DNS Server on port %d", port))
	server := &dns.Server{Addr: fmt.Sprintf(":%d", port), Net: "udp"}

	errChan := make(chan error, 1)
	go func() {
		errChan <- server.ListenAndServe()
	}()

	select {
	case err := <-errChan:
		s.logger.Log(LOGGER_SOURCE, fmt.Sprintf("Failed to start DNS Server on port %d: %v", port, err))
		return err
	case <-time.After(500 * time.Millisecond):
		s.running = true
		s.currentPort = port
		s.logger.Log(LOGGER_SOURCE, fmt.Sprintf("DNS Server is now running on port %d", port))
		return nil
	}
}

func (s *DNSServer) handleDNS(w dns.ResponseWriter, r *dns.Msg) {
	msg := new(dns.Msg)
	msg.SetReply(r)
	msg.Compress = false

	switch r.Opcode {
	case dns.OpcodeQuery:
		s.handleQuery(msg)
	}

	w.WriteMsg(msg)
}

func (s *DNSServer) handleQuery(msg *dns.Msg) {
	for _, q := range msg.Question {
		s.logger.Log(LOGGER_SOURCE, fmt.Sprintf("Received DNS Query for %s", q.Name))
		name := strings.TrimSuffix(q.Name, ".")
		parts := strings.Split(name, ".")
		if len(parts) < 2 {
			s.logger.Log(LOGGER_SOURCE, fmt.Sprintf("Invalid domain name: %s", q.Name))
			continue
		}

		tld := parts[len(parts)-1]
		if !s.TLDs[tld] {
			s.logger.Log(LOGGER_SOURCE, fmt.Sprintf("Unsupported TLD: %s", tld))
			continue
		}

		domain := strings.Join(parts[len(parts)-2:], ".")
		subdomain := strings.Join(parts[:len(parts)-2], ".")
		if subdomain == "" {
			subdomain = "@"
		}

		s.mu.RLock()
		d, ok := s.domains[domain]
		s.mu.RUnlock()

		if !ok {
			s.logger.Log(LOGGER_SOURCE, fmt.Sprintf("Domain not found: %s", domain))
			continue
		}

		record, ok := d.Records[subdomain]
		if !ok {
			s.logger.Log(LOGGER_SOURCE, fmt.Sprintf("Subdomain not found: %s", subdomain))
			continue
		}

		if record.Type == A && q.Qtype == dns.TypeA {
			rr, err := dns.NewRR(fmt.Sprintf("%s %d IN A %s", q.Name, q.Qclass, record.Value))
			if err != nil {
				s.logger.Log(LOGGER_SOURCE, fmt.Sprintf("Failed to create A record: %v", err))
			} else {
				msg.Answer = append(msg.Answer, rr)
			}
		}

		s.logger.Log(LOGGER_SOURCE, fmt.Sprintf("Responded to DNS Query for %s", q.Name))
	}
}

func (s *DNSServer) AddDomain(name string, isSystemDomain bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	parts := strings.Split(name, ".")
	if len(parts) < 2 {
		s.logger.Log(LOGGER_SOURCE, fmt.Sprintf("Invalid domain name: %s", name))
		return fmt.Errorf("invalid domain name: %s", name)
	}

	tld := parts[len(parts)-1]
	if !s.TLDs[tld] {
		s.logger.Log(LOGGER_SOURCE, fmt.Sprintf("Unsupported TLD: %s", tld))
		return fmt.Errorf("unsupported TLD: %s", tld)
	}

	if _, exists := s.domains[name]; exists {
		s.logger.Log(LOGGER_SOURCE, fmt.Sprintf("Domain already exists: %s", name))
		return fmt.Errorf("domain already exists: %s", name)
	}

	s.domains[name] = &Domain{
		Name:           name,
		Records:        make(map[string]DomainRecord),
		IsSystemDomain: isSystemDomain,
	}

	s.logger.Log(LOGGER_SOURCE, fmt.Sprintf("Added Domain: %s", name))
	return nil
}

func (s *DNSServer) AddSubdomain(domain, subdomain string, records []DomainRecord) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	d, ok := s.domains[domain]
	if !ok {
		s.logger.Log(LOGGER_SOURCE, fmt.Sprintf("Domain not found: %s", domain))
		return fmt.Errorf("domain not found: %s", domain)
	}

	for _, record := range records {
		d.Records[subdomain] = record
	}

	s.logger.Log(LOGGER_SOURCE, fmt.Sprintf("Added Subdomain %s to Domain %s", subdomain, domain))
	return nil
}

func (s *DNSServer) RemoveDomain(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	d, exists := s.domains[name]
	if !exists {
		s.logger.Log(LOGGER_SOURCE, fmt.Sprintf("Domain not found: %s", name))
		return fmt.Errorf("domain not found: %s", name)
	}

	if d.IsSystemDomain {
		s.logger.Log(LOGGER_SOURCE, fmt.Sprintf("Cannot remove system domain: %s", name))
		return fmt.Errorf("cannot remove system domain: %s", name)
	}

	delete(s.domains, name)
	s.logger.Log(LOGGER_SOURCE, fmt.Sprintf("Removed Domain: %s", name))
	return nil
}

func (s *DNSServer) RemoveSubdomain(domain, subdomain string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	d, exists := s.domains[domain]
	if !exists {
		s.logger.Log(LOGGER_SOURCE, fmt.Sprintf("Domain not found: %s", domain))
		return fmt.Errorf("domain not found: %s", domain)
	}

	if _, exists := d.Records[subdomain]; !exists {
		s.logger.Log(LOGGER_SOURCE, fmt.Sprintf("Subdomain not found: %s", subdomain))
		return fmt.Errorf("subdomain not found: %s", subdomain)
	}

	delete(d.Records, subdomain)
	s.logger.Log(LOGGER_SOURCE, fmt.Sprintf("Removed Subdomain %s from Domain %s", subdomain, domain))
	return nil
}

func (s *DNSServer) ListDomains() map[string]*Domain {
	s.mu.RLock()
	defer s.mu.RUnlock()

	domains := make(map[string]*Domain)
	for k, v := range s.domains {
		domains[k] = v
	}

	return domains
}

func (s *DNSServer) Status() string {
	if s.running {
		return fmt.Sprintf("DNS Server is running on Port %d", s.currentPort)
	}

	return "DNS Server is not running"
}

func (s *DNSServer) Stop() {
	s.running = false
	s.logger.Log(LOGGER_SOURCE, "Stopping DNS Server")
}

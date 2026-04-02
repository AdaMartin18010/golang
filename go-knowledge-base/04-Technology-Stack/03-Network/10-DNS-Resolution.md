# TS-NET-010: DNS Resolution in Go

> **维度**: Technology Stack > Network
> **级别**: S (16+ KB)
> **标签**: #dns #resolution #go #net #service-discovery
> **权威来源**:
>
> - [Go net Package](https://golang.org/pkg/net/) - Go standard library
> - [DNS RFC 1035](https://tools.ietf.org/html/rfc1035) - IETF

---

## 1. DNS Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                       DNS Resolution Architecture                            │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌─────────────┐                                                             │
│  │ Application │                                                             │
│  └──────┬──────┘                                                             │
│         │ Resolve "api.example.com"                                         │
│         ▼                                                                    │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                      Local DNS Resolver (Go net)                     │   │
│  │  - Check /etc/hosts                                                  │   │
│  │  - Check cache                                                       │   │
│  │  - Query DNS servers                                                 │   │
│  └───────────────────────────────┬─────────────────────────────────────┘   │
│                                  │                                           │
│                                  ▼                                           │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                      DNS Resolution Flow                             │   │
│  │                                                                      │   │
│  │  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐             │   │
│  │  │   Root      │───►│    TLD      │───►│  Authoritative│            │   │
│  │  │   Server    │    │   Server    │    │    Server     │            │   │
│  │  │   (.)       │    │  (.com)     │    │ (example.com) │            │   │
│  │  └─────────────┘    └─────────────┘    └─────────────┘             │   │
│  │        │                  │                  │                      │   │
│  │        │  NS for .com     │ NS for example.com                     │   │
│  │        │  198.41.0.4      │ 192.0.2.1                              │   │
│  │        ▼                  ▼                  ▼                      │   │
│  │  "I don't know,          "I don't know,        "api.example.com      │   │
│  │   ask root server"       ask .com server"     is 203.0.113.5"       │   │
│  │                                                                      │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
│  Record Types:                                                               │
│  - A: IPv4 address                                                           │
│  - AAAA: IPv6 address                                                        │
│  - CNAME: Canonical name (alias)                                             │
│  - MX: Mail exchange                                                         │
│  - NS: Name server                                                           │
│  - TXT: Text record                                                          │
│  - SRV: Service record                                                       │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 2. DNS Resolution in Go

```go
package main

import (
    "context"
    "fmt"
    "net"
    "time"
)

// Basic DNS lookup
func dnsLookup() {
    // Lookup IP addresses
    ips, err := net.LookupIP("example.com")
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    for _, ip := range ips {
        fmt.Println("IP:", ip)
    }

    // Lookup hostname
    names, err := net.LookupAddr("203.0.113.5")
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    for _, name := range names {
        fmt.Println("Hostname:", name)
    }

    // Lookup CNAME
    cname, err := net.LookupCNAME("www.example.com")
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println("CNAME:", cname)

    // Lookup MX records
    mxRecords, err := net.LookupMX("example.com")
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    for _, mx := range mxRecords {
        fmt.Printf("MX: %s (priority: %d)\n", mx.Host, mx.Pref)
    }

    // Lookup NS records
    nsRecords, err := net.LookupNS("example.com")
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    for _, ns := range nsRecords {
        fmt.Println("NS:", ns.Host)
    }

    // Lookup TXT records
    txtRecords, err := net.LookupTXT("example.com")
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    for _, txt := range txtRecords {
        fmt.Println("TXT:", txt)
    }

    // Lookup SRV records
    _, srvRecords, err := net.LookupSRV("http", "tcp", "example.com")
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    for _, srv := range srvRecords {
        fmt.Printf("SRV: %s:%d (priority: %d, weight: %d)\n",
            srv.Target, srv.Port, srv.Priority, srv.Weight)
    }
}

// Resolver with custom settings
func customResolver() {
    resolver := &net.Resolver{
        PreferGo: true, // Use Go's built-in resolver instead of system
        Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
            d := net.Dialer{
                Timeout: time.Second * 3,
            }
            return d.DialContext(ctx, network, address)
        },
    }

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    ips, err := resolver.LookupIPAddr(ctx, "example.com")
    if err != nil {
        fmt.Println("Error:", err)
        return
    }

    for _, ip := range ips {
        fmt.Println("IP:", ip.IP)
    }
}

// DNS caching
var dnsCache = make(map[string][]net.IP)
var dnsCacheMu sync.RWMutex
var dnsCacheTTL = 5 * time.Minute

type cacheEntry struct {
    ips       []net.IP
    timestamp time.Time
}

var cache = make(map[string]cacheEntry)

func cachedLookup(hostname string) ([]net.IP, error) {
    // Check cache
    dnsCacheMu.RLock()
    entry, found := cache[hostname]
    dnsCacheMu.RUnlock()

    if found && time.Since(entry.timestamp) < dnsCacheTTL {
        return entry.ips, nil
    }

    // Perform lookup
    ips, err := net.LookupIP(hostname)
    if err != nil {
        return nil, err
    }

    // Update cache
    dnsCacheMu.Lock()
    cache[hostname] = cacheEntry{
        ips:       ips,
        timestamp: time.Now(),
    }
    dnsCacheMu.Unlock()

    return ips, nil
}
```

---

## 3. Service Discovery

```go
// Service discovery using DNS SRV records
type ServiceDiscovery struct {
    domain string
}

func NewServiceDiscovery(domain string) *ServiceDiscovery {
    return &ServiceDiscovery{domain: domain}
}

func (sd *ServiceDiscovery) DiscoverService(service, proto string) ([]net.SRV, error) {
    _, srvs, err := net.LookupSRV(service, proto, sd.domain)
    if err != nil {
        return nil, err
    }
    return srvs, nil
}

// Load balanced HTTP client with DNS-based discovery
func (sd *ServiceDiscovery) CreateHTTPClient(service string) (*http.Client, error) {
    srvs, err := sd.DiscoverService(service, "tcp")
    if err != nil {
        return nil, err
    }

    if len(srvs) == 0 {
        return nil, errors.New("no services found")
    }

    // Create transport with custom dialer
    transport := &http.Transport{
        DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
            // Pick a service based on priority and weight
            srv := sd.selectSRV(srvs)
            address := fmt.Sprintf("%s:%d", srv.Target, srv.Port)
            return net.Dial(network, address)
        },
    }

    return &http.Client{Transport: transport}, nil
}

func (sd *ServiceDiscovery) selectSRV(srvs []net.SRV) *net.SRV {
    // Sort by priority
    sort.Slice(srvs, func(i, j int) bool {
        return srvs[i].Priority < srvs[j].Priority
    })

    // Select based on weight within same priority
    // Simplified: just return first for now
    return &srvs[0]
}
```

---

## 4. Checklist

```
DNS Resolution Checklist:
□ Use context for timeout control
□ Implement DNS caching
□ Handle DNS failures gracefully
□ Use SRV records for service discovery
□ Monitor DNS resolution time
□ Configure appropriate DNS servers
□ Handle both IPv4 and IPv6
□ Implement retry logic
```

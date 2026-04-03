# Week 4: System Design & Architecture

## Module Overview

**Duration:** 40 hours (5 days)
**Prerequisites:** Week 3 completion (Cloud-Native Patterns)
**Learning Goal:** Design scalable, reliable distributed systems and communicate architecture decisions

---

## Learning Objectives

By the end of this week, you will be able to:

1. **System Design Methodology**
   - Gather and analyze requirements using structured frameworks
   - Break down complex problems into manageable components
   - Make appropriate technology choices based on trade-off analysis
   - Document architecture decisions using ADR format

2. **Scalability Patterns**
   - Design for horizontal scaling and stateless architectures
   - Implement sharding strategies for data distribution
   - Apply caching at multiple layers (CDN, application, database)
   - Use load balancing effectively with health checks

3. **Data Consistency**
   - Apply CAP theorem trade-offs appropriately
   - Choose consistency models based on business requirements
   - Design distributed transactions using Saga pattern
   - Implement eventual consistency patterns

4. **Security Architecture**
   - Design authentication and authorization systems (OAuth2, OIDC)
   - Implement zero-trust principles
   - Secure inter-service communication with mTLS
   - Handle secrets and credentials with Vault

5. **Communication**
   - Present architecture to technical and non-technical stakeholders
   - Write Architecture Decision Records (ADRs)
   - Create clear architecture diagrams using C4 model
   - Participate effectively in design reviews

---

## Reading Assignments

### Required Reading (Complete by Day 3)

1. **[System Design Interview Guide](../05-Application-Domains/AD-010-System-Design-Interview.md)**
   - Study: The 4S framework (Scenario, Sketch, Study, Scale)
   - Learn: Back-of-envelope calculations for capacity planning
   - Master: Trade-off analysis and decision documentation

2. **[CAP Theorem Formal](../01-Formal-Theory/FT-003-CAP-Theorem-Formal.md)**
   - Understand: Consistency, Availability, Partition tolerance trade-offs
   - Learn: PACELC theorem extensions
   - Study: Real-world system classifications (CP vs AP systems)

3. **[Distributed Systems Foundations](../01-Formal-Theory/FT-001-Distributed-Systems-Foundation-Formal.md)**
   - Master: Time and ordering concepts
   - Learn: Failure models (fail-stop, Byzantine, network partitions)
   - Study: Consensus requirements and limitations

4. **[DDD Strategic Patterns](../05-Application-Domains/AD-001-DDD-Strategic-Patterns-Formal.md)**
   - Understand: Bounded contexts and ubiquitous language
   - Learn: Context mapping patterns (partnership, customer-supplier)
   - Study: Anti-corruption layer implementation

5. **[Security Architecture](../05-Application-Domains/AD-013-Security-Architecture.md)**
   - Learn: Defense in depth principles
   - Study: Zero-trust network architecture
   - Understand: Threat modeling methodology (STRIDE)

### Supplementary Reading (Complete by Day 5)

1. **[High Availability Design](../05-Application-Domains/AD-012-High-Availability-Design.md)**
   - Learn: Availability calculations (99.9%, 99.99%, 99.999%)
   - Study: Redundancy patterns (active-active, active-passive)
   - Understand: Failover strategies and detection

2. **[Event-Driven Architecture](../05-Application-Domains/AD-004-Event-Driven-Architecture-Formal.md)**
   - Master: Event patterns (event notification, event-carried state)
   - Learn: Event sourcing implementation patterns
   - Study: CQRS read model optimization

3. **[Capacity Planning](../05-Application-Domains/AD-009-Capacity-Planning.md)**
   - Understand: Load forecasting methods
   - Learn: Resource estimation techniques
   - Study: Scaling triggers and auto-scaling policies

---

## Hands-on Exercises

### Day 1: Design Methodology

#### Exercise 1.1: Requirements Gathering (2 hours)

Practice extracting requirements from problem statements using structured analysis:

**Problem:** Design a URL shortening service like bit.ly

**Requirements Analysis Template:**

```
FUNCTIONAL REQUIREMENTS:
- Core Functions:
  * Shorten long URLs to short codes
  * Redirect short codes to original URLs
  * Support custom aliases

- Extended Functions:
  * Analytics (click tracking, referrer info)
  * QR code generation
  * URL expiration
  * API rate limiting per user

NON-FUNCTIONAL REQUIREMENTS:
- Availability: 99.99% uptime (52 minutes downtime/year)
- Latency: P99 < 50ms for redirects, P99 < 200ms for API
- Scale: 10M new URLs/day, 100M redirects/day
- Durability: Zero data loss for created URLs
- Security: HTTPS only, protection against abuse

CONSTRAINTS:
- Budget: Cost-efficient at scale
- Regulatory: GDPR compliance for analytics
- Technical: Existing infrastructure integration
```

**Capacity Estimation:**

```go
package estimation

import "fmt"

// SystemRequirements captures all capacity needs
type SystemRequirements struct {
    DailyActiveUsers    int64
    RequestsPerUserDay  int64
    AverageResponseSize int64 // bytes
    ReadWriteRatio      float64

    // Derived metrics
    QPS           int64
    WriteQPS      int64
    ReadQPS       int64
    StoragePerDay int64
    StoragePerYear int64
    BandwidthQPS  int64
}

func CalculateURLShortenerRequirements() SystemRequirements {
    sr := SystemRequirements{
        DailyActiveUsers:    100_000_000, // 100M DAU
        RequestsPerUserDay:  10,
        AverageResponseSize: 500, // bytes
        ReadWriteRatio:      100, // 100:1 read:write
    }

    // Peak QPS (assume 2x average for peak)
    sr.QPS = (sr.DailyActiveUsers * sr.RequestsPerUserDay) / 86400 * 2
    sr.WriteQPS = int64(float64(sr.QPS) / (sr.ReadWriteRatio + 1))
    sr.ReadQPS = sr.QPS - sr.WriteQPS

    // Storage calculations
    sr.StoragePerDay = sr.WriteQPS * 86400 * sr.AverageResponseSize
    sr.StoragePerYear = sr.StoragePerDay * 365

    sr.BandwidthQPS = sr.QPS * sr.AverageResponseSize

    return sr
}

func PrintRequirements(sr SystemRequirements) {
    fmt.Println("=== URL Shortener Capacity Requirements ===")
    fmt.Printf("Daily Active Users: %d\n", sr.DailyActiveUsers)
    fmt.Printf("Peak QPS: %d\n", sr.QPS)
    fmt.Printf("  - Write QPS: %d\n", sr.WriteQPS)
    fmt.Printf("  - Read QPS: %d\n", sr.ReadQPS)
    fmt.Printf("Storage per year: %.2f TB\n", float64(sr.StoragePerYear)/1e12)
    fmt.Printf("Bandwidth required: %.2f Gbps\n", float64(sr.BandwidthQPS*8)/1e9)
}
```

**Tasks:**

1. Estimate storage requirements for 5-year horizon
2. Calculate read/write ratios impact on database choice
3. Identify potential bottlenecks at each layer
4. Define SLIs (Service Level Indicators) and SLOs (Service Level Objectives)

**Deliverable:** Requirements document with detailed capacity calculations

#### Exercise 1.2: Back-of-Envelope Calculation Practice (2 hours)

Practice estimation for different scenarios:

```go
package estimation

// Twitter-like feed system
func CalculateFeedRequirements() {
    // DAU: 500M
    // Avg following: 200 users
    // Tweets per user per day: 5
    // Timeline views per day: 10
}

// Video streaming platform
func CalculateStreamingRequirements() {
    // Concurrent viewers: 10M
    // Avg bitrate: 5 Mbps
    // Video catalog: 1M videos
    // Upload rate: 1000 videos/day
}

// E-commerce checkout system
func CalculateCheckoutRequirements() {
    // Daily orders: 1M
    // Peak: 10x average
    // Cart abandonment rate: 70%
    // Payment methods: 10+
}
```

**Deliverable:** Estimation spreadsheet for 3 different systems with justifications

---

### Day 2: Architecture Design

#### Exercise 2.1: Component Design (3 hours)

Design system components with clear interfaces:

**URL Shortener Architecture:**

```
┌─────────────────────────────────────────────────────────────┐
│                        DNS / CDN                            │
│              (Route 53, CloudFlare)                         │
└──────────────────────┬──────────────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────────────┐
│                   Load Balancer                             │
│              (Nginx / AWS ALB)                              │
└──────────────────────┬──────────────────────────────────────┘
                       │
        ┌──────────────┼──────────────┐
        │              │              │
┌───────▼──────┐ ┌─────▼──────┐ ┌────▼─────┐
│  API Server  │ │  Redirect  │ │  Worker  │
│   (REST)     │ │  Service   │ │  (Async) │
└───────┬──────┘ └─────┬──────┘ └────┬─────┘
        │              │             │
        └──────────────┼─────────────┘
                       │
        ┌──────────────┼──────────────┐
        │              │              │
┌───────▼──────┐ ┌─────▼──────┐ ┌────▼─────┐
│ PostgreSQL   │ │   Redis    │ │  Kafka   │
│ (Primary DB) │ │  (Cache)   │ │ (Events) │
└──────────────┘ └────────────┘ └──────────┘
```

**API Design:**

```go
package design

// URL Shortener Service Interface
type URLShortener interface {
    // Shorten creates a short URL
    Shorten(ctx context.Context, req ShortenRequest) (*ShortenResponse, error)

    // Resolve returns original URL from short code
    Resolve(ctx context.Context, shortCode string) (*ResolveResponse, error)

    // GetAnalytics returns click analytics
    GetAnalytics(ctx context.Context, shortCode string) (*AnalyticsResponse, error)

    // Delete removes a short URL
    Delete(ctx context.Context, shortCode string) error

    // Update allows modifying URL properties
    Update(ctx context.Context, shortCode string, req UpdateRequest) error
}

type ShortenRequest struct {
    URL          string            `json:"url" validate:"required,url"`
    CustomAlias  string            `json:"custom_alias,omitempty" validate:"omitempty,alphanum,min=4,max=20"`
    ExpiresAt    *time.Time        `json:"expires_at,omitempty"`
    Metadata     map[string]string `json:"metadata,omitempty"`
    UserID       string            `json:"user_id,omitempty"`
}

type ShortenResponse struct {
    ShortCode    string    `json:"short_code"`
    ShortURL     string    `json:"short_url"`
    OriginalURL  string    `json:"original_url"`
    CreatedAt    time.Time `json:"created_at"`
    ExpiresAt    *time.Time `json:"expires_at,omitempty"`
    QRCodeURL    string    `json:"qr_code_url,omitempty"`
}

type ResolveResponse struct {
    OriginalURL  string            `json:"original_url"`
    Metadata     map[string]string `json:"metadata,omitempty"`
    ClickCount   int64             `json:"click_count"`
    ExpiresAt    *time.Time        `json:"expires_at,omitempty"`
}

type AnalyticsResponse struct {
    ShortCode     string       `json:"short_code"`
    TotalClicks   int64        `json:"total_clicks"`
    UniqueClicks  int64        `json:"unique_clicks"`
    TopReferrers  []Referrer   `json:"top_referrers"`
    ClicksByDay   []DailyStats `json:"clicks_by_day"`
    ClicksByCountry []CountryStats `json:"clicks_by_country"`
    Devices       []DeviceStats  `json:"devices"`
}

type Referrer struct {
    URL   string `json:"url"`
    Count int64  `json:"count"`
}

type DailyStats struct {
    Date  string `json:"date"`
    Count int64  `json:"count"`
}
```

**Database Schema:**

```sql
-- URLs table (normalized)
CREATE TABLE urls (
    id BIGSERIAL PRIMARY KEY,
    short_code VARCHAR(10) UNIQUE NOT NULL,
    original_url TEXT NOT NULL,
    custom_alias BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    expires_at TIMESTAMP WITH TIME ZONE,
    user_id BIGINT REFERENCES users(id),
    metadata JSONB,
    click_count BIGINT DEFAULT 0,
    last_accessed_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_urls_short_code ON urls(short_code);
CREATE INDEX idx_urls_expires_at ON urls(expires_at) WHERE expires_at IS NOT NULL;
CREATE INDEX idx_urls_user_id ON urls(user_id);

-- Analytics table (time-series, partitioned)
CREATE TABLE url_analytics (
    id BIGSERIAL,
    url_id BIGINT REFERENCES urls(id),
    clicked_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    ip_address INET,
    user_agent TEXT,
    referrer TEXT,
    country_code VARCHAR(2),
    device_type VARCHAR(20),
    browser VARCHAR(50)
) PARTITION BY RANGE (clicked_at);

-- Create monthly partitions
CREATE TABLE url_analytics_2024_01 PARTITION OF url_analytics
    FOR VALUES FROM ('2024-01-01') TO ('2024-02-01');

CREATE INDEX idx_analytics_url_id ON url_analytics(url_id);
CREATE INDEX idx_analytics_clicked_at ON url_analytics(clicked_at);
CREATE INDEX idx_analytics_country ON url_analytics(country_code);
```

**Deliverable:** Complete design document with architecture diagram and API specs

#### Exercise 2.2: Scalability Design (2 hours)

Design for horizontal scaling:

```
SCALING STRATEGY:

1. Read Scaling:
   - L1 Cache: In-memory LRU per instance (10K hot URLs)
   - L2 Cache: Redis Cluster with 1hr TTL (99% hit rate target)
   - L3 Cache: Read replicas for analytics queries
   - CDN: Edge caching for redirects (CloudFlare)

2. Write Scaling:
   - Sharding by short_code prefix (16 shards: 0-9, a-f)
   - Async analytics ingestion via Kafka
   - Batch database writes for analytics
   - Write-behind caching for hot URLs

3. Database Scaling:
   - Primary-Replica setup for OLTP
   - Separate analytics warehouse (ClickHouse)
   - Automated partitioning by date
   - Archive old data to S3 after 1 year

4. Cache Strategy:
   - Cache-aside pattern for reads
   - Write-through for URL creation
   - TTL based on URL popularity
   - Hot key protection with local caching
```

**Capacity Numbers:**

```
Current Scale:
- 10M URLs/day = 116 URLs/second write
- 100M redirects/day = 1,157/second read
- Storage: ~50GB/year for metadata

At 10x Scale:
- 100M URLs/day = 1,157 URLs/second write
- 1B redirects/day = 11,574/second read
- Storage: ~500GB/year

Sharding Plan:
- 16 shards initially (0-9, a-f prefixes)
- Each shard handles ~72 writes/second
- Can expand to 256 shards if needed
```

**Deliverable:** Scaling plan with capacity numbers and growth projections

---

### Day 3: Trade-off Analysis

#### Exercise 3.1: CAP Analysis (2 hours)

Analyze CAP trade-offs for different scenarios:

```
URL SHORTENER CAP ANALYSIS:

Scenario 1: URL Creation (Write Operation)
- Preference: CONSISTENCY over Availability
- Reason: Cannot have duplicate short codes; data loss is unacceptable
- Implementation: Synchronous replication to majority before acknowledgment
- Trade-off: 50-100ms additional latency for stronger consistency

Scenario 2: URL Resolution (Read Operation)
- Preference: AVAILABILITY over Consistency
- Reason: Serving stale data briefly is acceptable; redirects should work
- Implementation: Read from cache first, tolerate stale data
- Trade-off: Potential for 1-second inconsistency during updates

Scenario 3: Analytics Collection (Write-heavy, Read-light)
- Preference: EVENTUAL CONSISTENCY
- Reason: Real-time not critical; batch processing acceptable
- Implementation: Async queue (Kafka), batch insert to analytics DB
- Trade-off: 5-minute delay in analytics visibility

Partition Handling:
- During partition: Continue serving from cache, queue writes
- After partition: Reconcile using vector clocks if needed
- Decision: AP for reads, CP for writes during partition
```

**Deliverable:** CAP analysis document for 3 different systems

#### Exercise 3.2: Technology Selection (2 hours)

Practice technology choices with decision matrix:

```
DATABASE SELECTION FOR URL SHORTENER:

┌────────────────────────────────────────────────────────────────┐
│ Criteria          │ Weight │ PostgreSQL │ Redis │ Cassandra │
├───────────────────┼────────┼────────────┼───────┼───────────┤
│ Consistency       │  30%   │     10     │   3   │     5     │
│ Availability      │  25%   │     7      │  10   │     9     │
│ Scalability       │  20%   │     6      │   7   │    10     │
│ Query Flexibility │  15%   │     10     │   4   │     6     │
│ Operational Cost  │  10%   │     8      │   9   │     5     │
├───────────────────┼────────┼────────────┼───────┼───────────┤
│ WEIGHTED SCORE    │  100%  │    8.35    │ 6.85  │    7.15   │
└────────────────────────────────────────────────────────────────┘

Decision: PostgreSQL for primary storage, Redis for cache, ClickHouse for analytics

Justification:
- PostgreSQL provides ACID guarantees needed for URL uniqueness
- JSONB support for flexible metadata
- Mature operational tooling
- Can scale vertically + read replicas initially
- Consider Citus or Spanner if horizontal scaling needed
```

**Alternative Analysis:**

```
SHORT CODE GENERATION OPTIONS:

Option 1: Auto-increment ID + Base62
- Pros: Simple, guaranteed unique, sequential
- Cons: Predictable, exposes system scale, easy to scrape
- Verdict: REJECTED - Security concerns

Option 2: Hash of URL (MD5/SHA256)
- Pros: Deterministic, distributed generation possible
- Cons: Collision risk, same URL always same code, irreversible
- Verdict: REJECTED - No custom aliases possible

Option 3: Twitter Snowflake / Sonyflake
- Pros: Unique, distributed, roughly time-ordered, 64-bit
- Cons: Clock synchronization required, slightly longer codes
- Verdict: ACCEPTED - Best balance of features

Option 4: Random Alphanumeric + Uniqueness Check
- Pros: Unpredictable, shorter codes possible
- Cons: Need uniqueness check, storage overhead for attempts
- Verdict: BACKUP - Use for custom aliases only
```

**Deliverable:** Technology decision document with quantitative analysis

---

### Day 4: Architecture Documentation

#### Exercise 4.1: ADR Writing (2 hours)

Write Architecture Decision Records:

```markdown
# ADR-001: URL Shortening Algorithm Selection

## Status
Accepted

## Context
We need to generate short, unique codes for URLs that satisfy:
- Short length (6-8 characters) for shareability
- Guaranteed uniqueness across the system
- URL-safe characters only
- Non-sequential to prevent scraping/prediction

## Options Considered

### Option 1: Base62 Encoding of Auto-increment ID
**Description:** Use database auto-increment and encode as Base62

**Pros:**
- Simplest implementation
- Guaranteed uniqueness via database constraint
- No coordination needed

**Cons:**
- Sequential and predictable
- Exposes system scale (easy to count total URLs)
- Vulnerable to scraping attacks
- Cannot support custom aliases

**Verdict:** Rejected due to security concerns

### Option 2: Hash of URL (MD5/SHA256 truncated)
**Description:** Hash the original URL and use first N characters

**Pros:**
- Deterministic - same URL always produces same code
- Distributed generation possible without coordination
- Simple implementation

**Cons:**
- Collision risk (birthday paradox)
- Cannot have multiple shorts for same URL
- Cannot support custom aliases
- Irreversible if original URL needed

**Verdict:** Rejected due to collision risk and lack of flexibility

### Option 3: Snowflake ID + Base62 Encoding
**Description:** Use Twitter Snowflake algorithm (or Sonyflake) for unique IDs

**Pros:**
- Globally unique without coordination
- Distributed generation (no single point of contention)
- Time-ordered (roughly)
- 64-bit provides sufficient space (9.2 quintillion)

**Cons:**
- Requires clock synchronization (NTP)
- Slightly longer codes (11 characters vs 6-7)
- Clock skew can cause out-of-order IDs

**Verdict:** Accepted - best balance of features

### Option 4: Random Alphanumeric with Uniqueness Check
**Description:** Generate random strings and check uniqueness

**Pros:**
- Truly unpredictable
- Can generate shorter codes
- Supports custom aliases naturally

**Cons:**
- Requires uniqueness check (database round-trip)
- Storage overhead for collision tracking
- Retry logic needed
- Birthday paradox requires more space than expected

**Verdict:** Backup option for custom aliases only

## Decision
Use **Sonyflake** (Go implementation of Snowflake) for auto-generated codes,
with random alphanumeric for custom aliases.

Sonyflake configuration:
- 41-bit timestamp (milliseconds)
- 10-bit machine ID (up to 1024 nodes)
- 12-bit sequence (4096 IDs per millisecond per node)

Base62 encoding for human-friendly representation.

## Consequences

### Positive
- No single point of failure for ID generation
- 4096 IDs per millisecond per machine (4 million/second total)
- Can generate IDs independently on any node
- Codes are roughly chronological

### Negative
- Must maintain unique machine IDs across cluster
- NTP synchronization required (within 10ms)
- 11-character codes (vs 6-7 with sequential)
- Clock skew can cause duplicate IDs if NTP fails

## Mitigations
- Automated NTP monitoring and alerting
- Machine ID assignment via Kubernetes or configuration
- Database unique constraint as safety net
- Fallback to random generation if clock skew detected

## References
- Twitter Snowflake: https://github.com/twitter-archive/snowflake
- Sonyflake: https://github.com/sony/sonyflake
- Base62 encoding: https://en.wikipedia.org/wiki/Base62
```

**Deliverable:** 3 ADRs covering key architecture decisions

#### Exercise 4.2: Architecture Diagrams (2 hours)

Create comprehensive C4 model diagrams:

```
C4 MODEL DIAGRAMS:

┌─────────────────────────────────────────────────────────────────┐
│ LEVEL 1: SYSTEM CONTEXT                                          │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  ┌──────────┐     Shorten URLs    ┌──────────────────┐          │
│  │          │ ───────────────────>│                  │          │
│  │   User   │                     │ URL Shortener    │          │
│  │          │ <───────────────────│ System           │          │
│  └──────────┘     Redirects       │                  │          │
│                                   └────────┬─────────┘          │
│                                            │                    │
│                                            │ Analytics          │
│                                            ▼                    │
│                                   ┌──────────────────┐          │
│                                   │ Analytics Store  │          │
│                                   └──────────────────┘          │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────┐
│ LEVEL 2: CONTAINER DIAGRAM                                       │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  ┌──────────┐                                                    │
│  │  User    │                                                    │
│  └────┬─────┘                                                    │
│       │ HTTPS                                                    │
│       ▼                                                          │
│  ┌────────────────────────────────────────────────────────┐     │
│  │              [Load Balancer: Nginx]                     │     │
│  └──────────────┬────────────────────────┬─────────────────┘     │
│                 │                        │                       │
│       ┌─────────┴────────┐    ┌──────────┴─────────┐            │
│       │                  │    │                    │            │
│  ┌────▼─────┐      ┌────▼────▼───┐        ┌───────▼────────┐   │
│  │  [Web    │      │  [API        │        │  [Analytics    │   │
│  │   App]   │      │   Service]   │        │   Worker]      │   │
│  │ React    │      │ Go/Echo      │        │ Go             │   │
│  └────┬─────┘      └──────┬───────┘        └───────┬────────┘   │
│       │                   │                        │            │
│       │            ┌──────┴──────┐                 │            │
│       │            │             │                 │            │
│       │       ┌────▼────┐   ┌────▼────┐      ┌─────▼────┐      │
│       │       │[Redis]  │   │[PostgreSQL]    │[ClickHouse]     │
│       │       │ Cache   │   │ Primary       │ Analytics       │
│       │       └─────────┘   └─────────┘      └──────────┘      │
│       │                                                          │
│  ┌────┴─────────────────────────────────────────────────────┐   │
│  │                    [Kafka]                                │   │
│  │              Event Streaming                              │   │
│  └───────────────────────────────────────────────────────────┘   │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────┐
│ LEVEL 3: COMPONENT DIAGRAM (API Service)                         │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  ┌────────────────────────────────────────────────────────┐     │
│  │                  [API Service Container]                │     │
│  │                                                         │     │
│  │  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐  │     │
│  │  │   Handler    │  │   Handler    │  │   Handler    │  │     │
│  │  │  (REST API)  │  │  (Metrics)   │  │   (Health)   │  │     │
│  │  └──────┬───────┘  └──────────────┘  └──────────────┘  │     │
│  │         │                                              │     │
│  │         ▼                                              │     │
│  │  ┌─────────────────────────────────────┐               │     │
│  │  │         Service Layer               │               │     │
│  │  │  ┌───────────┐    ┌───────────┐    │               │     │
│  │  │  │  URL      │    │ Analytics │    │               │     │
│  │  │  │ Service   │    │ Service   │    │               │     │
│  │  │  └─────┬─────┘    └─────┬─────┘    │               │     │
│  │  │        │                │          │               │     │
│  │  │        ▼                ▼          │               │     │
│  │  │  ┌─────────────────────────────┐   │               │     │
│  │  │  │    Repository Layer         │   │               │     │
│  │  │  │  ┌───────┐    ┌─────────┐   │   │               │     │
│  │  │  │  │  URL  │    │ Analytics│   │   │               │     │
│  │  │  │  │ Repo  │    │  Repo    │   │   │               │     │
│  │  │  │  └───┬───┘    └────┬────┘   │   │               │     │
│  │  │  └──────┼─────────────┼────────┘   │               │     │
│  │  └─────────┼─────────────┼────────────┘               │     │
│  │            │             │                            │     │
│  │            ▼             ▼                            │     │
│  │       [PostgreSQL]   [Redis]                          │     │
│  │                                                        │     │
│  └────────────────────────────────────────────────────────┘     │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

**Deliverable:** Complete C4 model diagrams (Level 1-3)

---

### Day 5: Design Review Preparation

#### Exercise 5.1: Presentation Preparation (3 hours)

Prepare architecture presentation:

```
PRESENTATION STRUCTURE (30 MINUTES):

SLIDE 1: Title (1 min)
- System name and version
- Presenter name and date

SLIDE 2: Problem Statement (2 min)
- What business problem are we solving?
- Why existing solutions don't work?
- Success criteria

SLIDE 3-4: Requirements (3 min)
- Functional requirements (must-haves vs nice-to-haves)
- Non-functional requirements with numbers
- Constraints and assumptions

SLIDE 5-6: High-Level Design (5 min)
- Architecture diagram
- Component overview
- Technology stack
- Data flow walkthrough

SLIDE 7-10: Deep Dive (10 min)
- Critical component details
- Database schema
- Caching strategy
- Scaling approach
- Security considerations

SLIDE 11-12: Trade-offs (5 min)
- Alternatives considered
- Why chosen approach
- What would change at 10x scale

SLIDE 13: Operations (3 min)
- Deployment strategy
- Monitoring plan
- Incident response

SLIDE 14: Roadmap (1 min)
- Phase 1: MVP
- Phase 2: Scale
- Phase 3: Advanced features

Q&A: (5 min)
```

**Common Questions to Prepare:**

```
SCALABILITY QUESTIONS:
- How would you handle 10x traffic?
- Where are the bottlenecks?
- What's your caching strategy and invalidation approach?
- How do you shard the data?

RELIABILITY QUESTIONS:
- What happens if the database goes down?
- How do you handle network partitions?
- What's your disaster recovery plan?
- How do you detect and handle failures?

SECURITY QUESTIONS:
- How do you prevent abuse/rate limiting?
- How do you authenticate users?
- How do you protect sensitive data?
- What's your threat model?

OPERATIONAL QUESTIONS:
- How do you deploy without downtime?
- How do you monitor this system?
- How do you debug production issues?
- What's your rollback strategy?
```

**Deliverable:** Complete presentation deck with speaker notes

#### Exercise 5.2: Mock Interview (2 hours)

Practice system design interview:

```
PRACTICE QUESTIONS:

1. Design a distributed cache (like Redis Cluster)
   - Partitioning strategy
   - Replication for availability
   - Cache eviction policies
   - Handling node failures

2. Design a message queue (like Kafka)
   - Topic partitioning
   - Consumer groups
   - Message retention
   - Ordering guarantees

3. Design a rate limiter for a global API
   - Single node vs distributed
   - Token bucket vs sliding window
   - Handling clock skew
   - Fairness between users

4. Design a real-time chat system
   - Message delivery guarantees
   - Presence detection
   - Message history
   - Push notifications

5. Design a video streaming platform
   - Content delivery
   - Adaptive bitrate
   - DRM protection
   - Live streaming
```

**Deliverable:** Written responses to 20 common interview questions

---

## Code Review Checklist for System Design

### Architecture

- [ ] Clear separation of concerns between layers
- [ ] Appropriate abstraction levels (not too high or low)
- [ ] Consistent patterns across all components
- [ ] Proper error handling strategy defined
- [ ] Well-defined interfaces between components
- [ ] Dependency direction is correct (dependencies point inward)

### Scalability

- [ ] Stateless design where possible
- [ ] Horizontal scaling support documented
- [ ] Efficient data structures chosen
- [ ] Caching strategy documented with invalidation
- [ ] Database sharding plan with criteria
- [ ] Load balancing approach specified
- [ ] CDN usage considered for static assets

### Reliability

- [ ] Failure scenarios identified and documented
- [ ] Retry strategies with backoff defined
- [ ] Circuit breaker patterns where appropriate
- [ ] Graceful degradation plan
- [ ] Recovery procedures documented
- [ ] Backup and restore procedures
- [ ] Disaster recovery plan with RTO/RPO

### Security

- [ ] Authentication mechanism specified
- [ ] Authorization strategy (RBAC/ABAC) defined
- [ ] Data encryption at rest and in transit
- [ ] Secret management approach (Vault, etc.)
- [ ] Audit logging requirements
- [ ] Input validation strategy
- [ ] Rate limiting and DDoS protection
- [ ] Threat model documented

### Data

- [ ] Schema normalization level justified
- [ ] Indexing strategy for queries
- [ ] Data retention policy defined
- [ ] Archiving strategy for old data
- [ ] GDPR/privacy compliance addressed
- [ ] Data consistency model documented

### Operations

- [ ] Deployment strategy (blue-green, canary)
- [ ] Monitoring and alerting plan
- [ ] Health checks defined
- [ ] Metrics collection strategy
- [ ] Log aggregation approach
- [ ] Tracing strategy for distributed systems
- [ ] Runbooks for common issues

---

## Assessment Criteria

### Knowledge Assessment (30%)

**Written Quiz Topics:**

1. System design methodology (4S framework)
2. Scalability patterns and trade-offs
3. Consistency models (strong, eventual, causal)
4. Security principles and threat modeling
5. CAP theorem and PACELC
6. Back-of-envelope calculations
7. Caching strategies and invalidation
8. Database selection criteria
9. Microservices vs monolith trade-offs
10. Event-driven architecture patterns

**Passing Score:** 75%

### Design Challenge (50%)

**Problem:** Design a distributed task scheduler

**Requirements:**

- Schedule millions of tasks with cron expressions
- Support task dependencies (DAG)
- Handle task failures and retries
- Distribute across worker nodes
- Provide observability and alerting
- Scale horizontally

**Evaluation Criteria:**

- Requirements gathering: 10%
- Architecture design: 20%
- Technology choices: 10%
- Scalability plan: 10%

### Presentation (20%)

**Requirements:**

- 30-minute presentation to panel
- Clear communication of complex concepts
- Handle Q&A effectively
- Professional delivery

**Evaluation Criteria:**

- Structure and organization: 5%
- Technical accuracy: 5%
- Communication clarity: 5%
- Q&A handling: 5%

---

## Resources and References

### Books

- "Designing Data-Intensive Applications" by Martin Kleppmann
- "Building Microservices" by Sam Newman
- "Release It!" by Michael Nygard
- "The Architecture of Open Source Applications"
- "Site Reliability Engineering" by Google

### Online Resources

- System Design Primer (GitHub)
- High Scalability blog
- AWS Architecture Center
- Azure Architecture Center
- Google Cloud Architecture Center

### Tools

- C4 Model (c4model.com)
- PlantUML for diagrams
- Mermaid for markdown diagrams
- Excalidraw for sketches

---

*Next: [Advanced Distributed Systems](advanced-distributed.md)*

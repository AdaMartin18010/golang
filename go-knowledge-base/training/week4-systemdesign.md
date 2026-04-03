# Week 4: System Design & Architecture

## Module Overview

**Duration:** 40 hours (5 days)
**Prerequisites:** Week 3 completion (Cloud-Native Patterns)
**Learning Goal:** Design scalable, reliable distributed systems and communicate architecture decisions

---

## Learning Objectives

By the end of this week, you will be able to:

1. **System Design Methodology**
   - Gather and analyze requirements
   - Break down problems into components
   - Make appropriate technology choices
   - Document architecture decisions

2. **Scalability Patterns**
   - Design for horizontal scaling
   - Implement sharding strategies
   - Apply caching at multiple layers
   - Use load balancing effectively

3. **Data Consistency**
   - Apply CAP theorem trade-offs
   - Choose consistency models appropriately
   - Design distributed transactions
   - Implement eventual consistency patterns

4. **Security Architecture**
   - Design authentication and authorization systems
   - Implement zero-trust principles
   - Secure inter-service communication
   - Handle secrets and credentials

5. **Communication**
   - Present architecture to stakeholders
   - Write Architecture Decision Records (ADRs)
   - Create clear architecture diagrams
   - Participate in design reviews

---

## Reading Assignments

### Required Reading (Complete by Day 3)

1. **[System Design Interview Guide](../05-Application-Domains/AD-010-System-Design-Interview.md)**
   - Study: The 4S framework (Scenario, Sketch, Study, Scale)
   - Learn: Back-of-envelope calculations
   - Master: Trade-off analysis

2. **[CAP Theorem Formal](../01-Formal-Theory/FT-003-CAP-Theorem-Formal.md)**
   - Understand: Consistency, Availability, Partition tolerance
   - Learn: PACELC theorem
   - Study: Real-world system classifications

3. **[Distributed Systems Foundations](../01-Formal-Theory/FT-001-Distributed-Systems-Foundation-Formal.md)**
   - Master: Time and ordering
   - Learn: Failure models
   - Study: Consensus requirements

4. **[DDD Strategic Patterns](../05-Application-Domains/AD-001-DDD-Strategic-Patterns-Formal.md)**
   - Understand: Bounded contexts
   - Learn: Context mapping
   - Study: Ubiquitous language

5. **[Security Architecture](../05-Application-Domains/AD-013-Security-Architecture.md)**
   - Learn: Defense in depth
   - Study: Zero-trust principles
   - Understand: Threat modeling

---

## Hands-on Exercises

### Day 1: Design Methodology

#### Exercise 1.1: Requirements Gathering (2 hours)

Practice extracting requirements from problem statements:

**Problem:** Design a URL shortening service like bit.ly

**Requirements Analysis Template:**

```
Functional Requirements:
- Core: Shorten URLs, redirect to original
- Extended: Custom aliases, analytics, expiration

Non-Functional Requirements:
- Availability: 99.99% uptime
- Latency: <50ms for redirection
- Scale: 10M new URLs/day, 100M reads/day
- Durability: No URL loss

Constraints:
- Budget considerations
- Regulatory (GDPR for analytics)
- Technical (existing infrastructure)
```

**Deliverable:** Requirements document with calculations

#### Exercise 1.2: Back-of-Envelope Calculation (2 hours)

Practice estimation:

```go
package estimation

import "fmt"

type SystemRequirements struct {
    DailyActiveUsers    int64
    RequestsPerUserDay  int64
    AverageResponseSize int64
    ReadWriteRatio      float64
    QPS                 int64
    StoragePerYear      int64
}

func CalculateRequirements(dau, rpu, respSize int64, rwRatio float64) SystemRequirements {
    qps := (dau * rpu) / 86400 * 2
    storagePerDay := (qps / int64(rwRatio+1)) * 86400 * respSize

    return SystemRequirements{
        DailyActiveUsers:    dau,
        RequestsPerUserDay:  rpu,
        QPS:                 qps,
        StoragePerYear:      storagePerDay * 365,
    }
}
```

**Deliverable:** Estimation spreadsheet for 3 systems

---

### Day 2: Architecture Design

#### Exercise 2.1: Component Design (3 hours)

Design system components with API specifications.

**Deliverable:** Complete design document with architecture diagram

#### Exercise 2.2: Scalability Design (2 hours)

Design for horizontal scaling with caching strategy.

**Deliverable:** Scaling plan with capacity numbers

---

### Day 3: Trade-off Analysis

#### Exercise 3.1: CAP Analysis (2 hours)

Analyze CAP trade-offs for different scenarios.

**Deliverable:** CAP analysis document

#### Exercise 3.2: Technology Selection (2 hours)

Practice technology choices with decision matrix.

**Deliverable:** Technology decision document

---

### Day 4: Architecture Documentation

#### Exercise 4.1: ADR Writing (2 hours)

Write Architecture Decision Records following template.

**Deliverable:** 3 ADRs for your system design

#### Exercise 4.2: Architecture Diagrams (2 hours)

Create C4 model diagrams for your system.

**Deliverable:** Complete C4 model diagrams

---

### Day 5: Design Review Preparation

#### Exercise 5.1: Presentation Preparation (3 hours)

Prepare architecture presentation following structure.

**Deliverable:** Complete presentation deck

#### Exercise 5.2: Mock Interview (2 hours)

Practice system design interview questions.

**Deliverable:** Responses to 20 common questions

---

## Code Review Checklist for System Design

### Architecture

- [ ] Clear separation of concerns
- [ ] Appropriate abstraction levels
- [ ] Consistent patterns across components
- [ ] Proper error handling strategy
- [ ] Defined interfaces between components

### Scalability

- [ ] Stateless design where possible
- [ ] Horizontal scaling support
- [ ] Efficient data structures
- [ ] Caching strategy documented
- [ ] Database sharding plan

### Reliability

- [ ] Failure scenarios identified
- [ ] Retry strategies defined
- [ ] Circuit breaker patterns
- [ ] Graceful degradation plan
- [ ] Recovery procedures documented

### Security

- [ ] Authentication mechanism
- [ ] Authorization strategy
- [ ] Data encryption approach
- [ ] Secret management plan
- [ ] Audit logging

---

## Assessment Criteria

### Knowledge Assessment (30%)

- System design methodology
- Scalability patterns
- Consistency models
- Security principles

**Passing Score:** 75%

### Design Challenge (50%)

**Problem:** Design a distributed task scheduler

**Requirements:**

- Schedule millions of tasks
- Support cron expressions
- Handle task failures and retries
- Provide observability

### Presentation (20%)

**Requirements:**

- 30-minute presentation
- Clear communication
- Handle Q&A effectively

---

*Next: [Advanced Distributed Systems](advanced-distributed.md)*

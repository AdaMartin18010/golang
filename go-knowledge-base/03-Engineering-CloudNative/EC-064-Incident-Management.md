# Incident Management

> **分类**: 工程与云原生
> **标签**: #incident-management #sre #response #command #communication
> **参考**: Google SRE, PagerDuty, NIST SP 800-61

---

## 1. Formal Definition

### 1.1 What is Incident Management?

Incident Management is the systematic approach to identifying, analyzing, and resolving incidents that disrupt normal service operations. It encompasses the processes, tools, and roles required to restore service quickly while minimizing impact to business operations.

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      Incident Management Lifecycle                           │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│   DETECTION          RESPONSE           MITIGATION         RESOLUTION       │
│      │                  │                   │                  │            │
│      ▼                  ▼                   ▼                  ▼            │
│   ┌───────┐         ┌───────┐          ┌───────┐         ┌───────┐         │
│   │Alert  │────────►│Triage │─────────►│Contain│────────►│Restore│         │
│   │Trigger│         │Assess │          │Isolate│         │Service│         │
│   └───────┘         └───┬───┘          └───┬───┘         └───┬───┘         │
│       │                 │                  │                 │             │
│       │                 ▼                  ▼                 │             │
│       │            ┌─────────────────────────────┐           │             │
│       │            │      COMMAND STRUCTURE      │           │             │
│       │            ├─────────────────────────────┤           │             │
│       │            │                             │           │             │
│       │            │   Incident Commander (IC)   │           │             │
│       │            │   • Overall coordination    │           │             │
│       │            │   • External communication  │           │             │
│       │            │   • Decision authority      │           │             │
│       │            │                             │           │             │
│       │            │   Communications Lead (CL)  │           │             │
│       │            │   • Status updates          │           │             │
│       │            │   • Stakeholder updates     │           │             │
│       │            │   • Status page updates     │           │             │
│       │            │                             │           │             │
│       │            │   Operations Lead (OL)      │           │             │
│       │            │   • Technical resolution    │           │             │
│       │            │   • Resource coordination   │           │             │
│       │            │   • Mitigation execution    │           │             │
│       │            │                             │           │             │
│       │            └─────────────────────────────┘           │             │
│       │                                                      │             │
│       ▼                                                      ▼             │
│   ┌─────────────────────────────────────────────────────────────────────┐  │
│   │                        POST-INCIDENT                                 │  │
│   │  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐              │  │
│   │  │   Review    │───►│  Postmortem │───►│   Action    │              │  │
│   │  │  Timeline   │    │  Document   │    │   Items     │              │  │
│   │  └─────────────┘    └─────────────┘    └─────────────┘              │  │
│   └─────────────────────────────────────────────────────────────────────┘  │
│                                                                             │
│   Metrics: MTTD (Mean Time To Detect) │ MTTR (Mean Time To Resolve)        │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 1.2 Incident Severity Classification

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Incident Severity Classification                          │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  SEV 1 - CRITICAL                     SEV 2 - HIGH                          │
│  ━━━━━━━━━━━━━━━━━━━━                 ━━━━━━━━━━━━━━━                       │
│                                                                             │
│  ✓ Complete service outage            ✓ Major feature unavailable           │
│  ✓ Data loss or corruption            ✓ Significant performance degradation│
│  ✓ Security breach (active)           ✓ Partial data loss                  │
│  ✓ Revenue impact > $100K/hr          ✓ Revenue impact $10K-100K/hr        │
│  ✓ All users affected                 ✓ Many users affected                 │
│                                                                             │
│  Response: Immediate (24/7)           Response: Within 30 min (24/7)        │
│  Commander: Senior IC required        Commander: Experienced IC             │
│  Communication: 15-min updates        Communication: 30-min updates         │
│                                                                             │
│  ─────────────────────────────────────────────────────────────────────────  │
│                                                                             │
│  SEV 3 - MEDIUM                       SEV 4 - LOW                           │
│  ━━━━━━━━━━━━━━━━━━━━                 ━━━━━━━━━━━━━━━                       │
│                                                                             │
│  ✓ Minor feature issues               ✓ Cosmetic issues                     │
│  ✓ Non-critical bug                   ✓ Single user issues                  │
│  ✓ Degraded experience                ✓ Workaround available                │
│  ✓ Low revenue impact                 ✓ No revenue impact                   │
│  ✓ Limited user impact                ✓ Minimal user impact                 │
│                                                                             │
│  Response: Business hours             Response: Next business day           │
│  Commander: On-call engineer          Commander: Not required               │
│  Communication: As needed             Communication: Ticket updates         │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 2. Implementation Patterns in Go

### 2.1 Incident Command System

```go
package incident

import (
    "context"
    "errors"
    "fmt"
    "sync"
    "time"
)

// Severity levels
type Severity string

const (
    SeverityCritical Severity = "critical"
    SeverityHigh     Severity = "high"
    SeverityMedium   Severity = "medium"
    SeverityLow      Severity = "low"
)

// IncidentState represents the state of an incident
type IncidentState string

const (
    StateDetected      IncidentState = "detected"
    StateTriaged       IncidentState = "triaged"
    StateMitigating    IncidentState = "mitigating"
    StateMitigated     IncidentState = "mitigated"
    StateMonitoring    IncidentState = "monitoring"
    StateResolved      IncidentState = "resolved"
    StatePostmortem    IncidentState = "postmortem"
    StateClosed        IncidentState = "closed"
)

// Role represents an incident response role
type Role string

const (
    RoleIncidentCommander Role = "incident_commander"
    RoleOpsLead           Role = "ops_lead"
    RoleCommunications    Role = "communications"
    RoleScribe            Role = "scribe"
    RoleSubjectMatterExpert Role = "sme"
)

// Participant represents an incident participant
type Participant struct {
    ID       string `json:"id"`
    Name     string `json:"name"`
    Email    string `json:"email"`
    Role     Role   `json:"role"`
    JoinedAt time.Time `json:"joined_at"`
}

// Incident represents a managed incident
type Incident struct {
    ID          string    `json:"id"`
    Title       string    `json:"title"`
    Description string    `json:"description"`
    Severity    Severity  `json:"severity"`
    State       IncidentState `json:"state"`

    // Timestamps
    DetectedAt   time.Time  `json:"detected_at"`
    AcknowledgedAt *time.Time `json:"acknowledged_at,omitempty"`
    MitigatedAt  *time.Time `json:"mitigated_at,omitempty"`
    ResolvedAt   *time.Time `json:"resolved_at,omitempty"`
    ClosedAt     *time.Time `json:"closed_at,omitempty"`

    // Command structure
    CommanderID  string        `json:"commander_id"`
    Participants []Participant `json:"participants"`

    // Impact
    AffectedServices []string `json:"affected_services"`
    AffectedRegions  []string `json:"affected_regions"`
    CustomerImpact   *CustomerImpact `json:"customer_impact,omitempty"`

    // Communication
    SlackChannel  string `json:"slack_channel,omitempty"`
    ZoomLink      string `json:"zoom_link,omitempty"`
    StatusPageURL string `json:"status_page_url,omitempty"`

    // Tracking
    Timeline  []TimelineEvent `json:"timeline"`
    Updates   []StatusUpdate  `json:"updates"`
    ActionItems []ActionItem  `json:"action_items"`

    mu sync.RWMutex
}

// CustomerImpact describes customer impact
type CustomerImpact struct {
    UsersAffected   int     `json:"users_affected"`
    RegionsAffected int     `json:"regions_affected"`
    RevenueImpact   float64 `json:"revenue_impact_per_hour"`
    Description     string  `json:"description"`
}

// TimelineEvent represents an incident timeline entry
type TimelineEvent struct {
    Timestamp   time.Time `json:"timestamp"`
    AuthorID    string    `json:"author_id"`
    Description string    `json:"description"`
    State       IncidentState `json:"state,omitempty"`
}

// StatusUpdate represents a status update
type StatusUpdate struct {
    Timestamp   time.Time `json:"timestamp"`
    AuthorID    string    `json:"author_id"`
    Message     string    `json:"message"`
    Audience    string    `json:"audience"` // internal, external, executive
}

// ActionItem represents a follow-up action
type ActionItem struct {
    ID          string    `json:"id"`
    Description string    `json:"description"`
    AssigneeID  string    `json:"assignee_id"`
    DueDate     time.Time `json:"due_date"`
    Status      string    `json:"status"` // open, in_progress, closed
    Priority    string    `json:"priority"`
}

// IncidentManager manages the incident lifecycle
type IncidentManager struct {
    incidents map[string]*Incident
    mu        sync.RWMutex

    notifier  IncidentNotifier
    store     IncidentStore
    idCounter int
}

// IncidentNotifier sends incident notifications
type IncidentNotifier interface {
    NotifyIncidentDeclared(ctx context.Context, incident *Incident) error
    NotifyStatusUpdate(ctx context.Context, incident *Incident, update StatusUpdate) error
    NotifyParticipantJoined(ctx context.Context, incident *Incident, participant Participant) error
    NotifyStateChange(ctx context.Context, incident *Incident, oldState, newState IncidentState) error
    NotifyIncidentResolved(ctx context.Context, incident *Incident) error
}

// IncidentStore persists incident data
type IncidentStore interface {
    SaveIncident(ctx context.Context, incident *Incident) error
    LoadIncident(ctx context.Context, id string) (*Incident, error)
    ListIncidents(ctx context.Context, filter IncidentFilter) ([]*Incident, error)
}

// IncidentFilter filters incidents
type IncidentFilter struct {
    State      *IncidentState
    Severity   *Severity
    StartTime  *time.Time
    EndTime    *time.Time
    Commander  string
}

// NewIncidentManager creates a new incident manager
func NewIncidentManager(notifier IncidentNotifier, store IncidentStore) *IncidentManager {
    return &IncidentManager{
        incidents: make(map[string]*Incident),
        notifier:  notifier,
        store:     store,
    }
}

// DeclareIncident declares a new incident
func (m *IncidentManager) DeclareIncident(ctx context.Context, title, description string, severity Severity, commanderID string, affectedServices []string) (*Incident, error) {
    m.mu.Lock()
    m.idCounter++
    id := fmt.Sprintf("INC-%d-%d", time.Now().Year(), m.idCounter)
    m.mu.Unlock()

    now := time.Now().UTC()
    incident := &Incident{
        ID:               id,
        Title:            title,
        Description:      description,
        Severity:         severity,
        State:            StateDetected,
        DetectedAt:       now,
        CommanderID:      commanderID,
        AffectedServices: affectedServices,
        Participants: []Participant{
            {
                ID:       commanderID,
                Role:     RoleIncidentCommander,
                JoinedAt: now,
            },
        },
        Timeline: []TimelineEvent{
            {
                Timestamp:   now,
                AuthorID:    commanderID,
                Description: "Incident declared",
                State:       StateDetected,
            },
        },
        Updates:     make([]StatusUpdate, 0),
        ActionItems: make([]ActionItem, 0),
    }

    m.mu.Lock()
    m.incidents[id] = incident
    m.mu.Unlock()

    if err := m.store.SaveIncident(ctx, incident); err != nil {
        return nil, fmt.Errorf("failed to save incident: %w", err)
    }

    if err := m.notifier.NotifyIncidentDeclared(ctx, incident); err != nil {
        // Log but don't fail
        fmt.Printf("Failed to notify incident declaration: %v\n", err)
    }

    return incident, nil
}

// UpdateState updates the incident state
func (m *IncidentManager) UpdateState(ctx context.Context, incidentID, authorID string, newState IncidentState, description string) error {
    incident, err := m.getIncident(incidentID)
    if err != nil {
        return err
    }

    incident.mu.Lock()
    defer incident.mu.Unlock()

    oldState := incident.State
    if oldState == newState {
        return nil
    }

    // Validate state transition
    if !isValidStateTransition(oldState, newState) {
        return fmt.Errorf("invalid state transition from %s to %s", oldState, newState)
    }

    now := time.Now().UTC()
    incident.State = newState

    // Update timestamps
    switch newState {
    case StateMitigated:
        incident.MitigatedAt = &now
    case StateResolved:
        incident.ResolvedAt = &now
    case StateClosed:
        incident.ClosedAt = &now
    }

    // Add timeline event
    event := TimelineEvent{
        Timestamp:   now,
        AuthorID:    authorID,
        Description: description,
        State:       newState,
    }
    incident.Timeline = append(incident.Timeline, event)

    if err := m.store.SaveIncident(ctx, incident); err != nil {
        return fmt.Errorf("failed to save incident: %w", err)
    }

    if err := m.notifier.NotifyStateChange(ctx, incident, oldState, newState); err != nil {
        fmt.Printf("Failed to notify state change: %v\n", err)
    }

    return nil
}

// AddParticipant adds a participant to the incident
func (m *IncidentManager) AddParticipant(ctx context.Context, incidentID string, participant Participant) error {
    incident, err := m.getIncident(incidentID)
    if err != nil {
        return err
    }

    incident.mu.Lock()
    defer incident.mu.Unlock()

    // Check if already participating
    for _, p := range incident.Participants {
        if p.ID == participant.ID {
            return nil
        }
    }

    participant.JoinedAt = time.Now().UTC()
    incident.Participants = append(incident.Participants, participant)

    if err := m.store.SaveIncident(ctx, incident); err != nil {
        return err
    }

    return m.notifier.NotifyParticipantJoined(ctx, incident, participant)
}

// AddStatusUpdate adds a status update
func (m *IncidentManager) AddStatusUpdate(ctx context.Context, incidentID, authorID, message, audience string) error {
    incident, err := m.getIncident(incidentID)
    if err != nil {
        return err
    }

    incident.mu.Lock()
    defer incident.mu.Unlock()

    update := StatusUpdate{
        Timestamp: time.Now().UTC(),
        AuthorID:  authorID,
        Message:   message,
        Audience:  audience,
    }
    incident.Updates = append(incident.Updates, update)

    if err := m.store.SaveIncident(ctx, incident); err != nil {
        return err
    }

    return m.notifier.NotifyStatusUpdate(ctx, incident, update)
}

// AddActionItem adds an action item
func (m *IncidentManager) AddActionItem(ctx context.Context, incidentID string, item ActionItem) error {
    incident, err := m.getIncident(incidentID)
    if err != nil {
        return err
    }

    incident.mu.Lock()
    defer incident.mu.Unlock()

    item.ID = fmt.Sprintf("AI-%d", len(incident.ActionItems)+1)
    item.Status = "open"
    incident.ActionItems = append(incident.ActionItems, item)

    return m.store.SaveIncident(ctx, incident)
}

// UpdateCustomerImpact updates customer impact information
func (m *IncidentManager) UpdateCustomerImpact(ctx context.Context, incidentID string, impact CustomerImpact) error {
    incident, err := m.getIncident(incidentID)
    if err != nil {
        return err
    }

    incident.mu.Lock()
    defer incident.mu.Unlock()

    incident.CustomerImpact = &impact

    return m.store.SaveIncident(ctx, incident)
}

// getIncident retrieves an incident by ID
func (m *IncidentManager) getIncident(id string) (*Incident, error) {
    m.mu.RLock()
    incident, exists := m.incidents[id]
    m.mu.RUnlock()

    if !exists {
        return nil, fmt.Errorf("incident not found: %s", id)
    }

    return incident, nil
}

// GetActiveIncidents returns all active (non-closed) incidents
func (m *IncidentManager) GetActiveIncidents() []*Incident {
    m.mu.RLock()
    defer m.mu.RUnlock()

    active := make([]*Incident, 0)
    for _, incident := range m.incidents {
        if incident.State != StateClosed {
            active = append(active, incident)
        }
    }

    return active
}

// CalculateMetrics calculates incident metrics
func (m *IncidentManager) CalculateMetrics(incidentID string) (*IncidentMetrics, error) {
    incident, err := m.getIncident(incidentID)
    if err != nil {
        return nil, err
    }

    incident.mu.RLock()
    defer incident.mu.RUnlock()

    metrics := &IncidentMetrics{}

    // MTTD
    if incident.AcknowledgedAt != nil {
        metrics.MTTD = incident.AcknowledgedAt.Sub(incident.DetectedAt)
    }

    // MTTR
    if incident.ResolvedAt != nil {
        metrics.MTTR = incident.ResolvedAt.Sub(incident.DetectedAt)
    }

    // Time to Mitigate
    if incident.MitigatedAt != nil {
        metrics.TimeToMitigate = incident.MitigatedAt.Sub(incident.DetectedAt)
    }

    return metrics, nil
}

// IncidentMetrics contains incident metrics
type IncidentMetrics struct {
    MTTD           time.Duration `json:"mttd"`
    MTTR           time.Duration `json:"mttr"`
    TimeToMitigate time.Duration `json:"time_to_mitigate"`
}

// isValidStateTransition checks if a state transition is valid
func isValidStateTransition(from, to IncidentState) bool {
    transitions := map[IncidentState][]IncidentState{
        StateDetected:   {StateTriaged, StateMitigating, StateResolved},
        StateTriaged:    {StateMitigating, StateResolved},
        StateMitigating: {StateMitigated, StateResolved},
        StateMitigated:  {StateMonitoring, StateResolved},
        StateMonitoring: {StateResolved},
        StateResolved:   {StatePostmortem},
        StatePostmortem: {StateClosed},
    }

    valid, exists := transitions[from]
    if !exists {
        return false
    }

    for _, v := range valid {
        if v == to {
            return true
        }
    }

    return false
}
```

### 2.2 Status Page Integration

```go
package incident

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "time"
)

// StatusPageClient integrates with status page services
type StatusPageClient struct {
    apiKey     string
    pageID     string
    apiURL     string
    httpClient *http.Client
}

// StatusPageIncident represents a status page incident
type StatusPageIncident struct {
    ID          string    `json:"id,omitempty"`
    Name        string    `json:"name"`
    Status      string    `json:"status"` // investigating, identified, monitoring, resolved
    Impact      string    `json:"impact"` // none, minor, major, critical
    CreatedAt   time.Time `json:"created_at,omitempty"`
    UpdatedAt   time.Time `json:"updated_at,omitempty"`
    ResolvedAt  *time.Time `json:"resolved_at,omitempty"`
    Shortlink   string    `json:"shortlink,omitempty"`
    IncidentUpdates []StatusPageUpdate `json:"incident_updates,omitempty"`
    Components  []Component `json:"components,omitempty"`
}

// StatusPageUpdate represents a status page update
type StatusPageUpdate struct {
    ID        string    `json:"id,omitempty"`
    Status    string    `json:"status"`
    Body      string    `json:"body"`
    CreatedAt time.Time `json:"created_at,omitempty"`
}

// Component represents a status page component
type Component struct {
    ID     string `json:"id"`
    Name   string `json:"name"`
    Status string `json:"status"` // operational, degraded, partial_outage, major_outage
}

// NewStatusPageClient creates a new StatusPage client
func NewStatusPageClient(apiKey, pageID string) *StatusPageClient {
    return &StatusPageClient{
        apiKey: apiKey,
        pageID: pageID,
        apiURL: "https://api.statuspage.io/v1",
        httpClient: &http.Client{
            Timeout: 30 * time.Second,
        },
    }
}

// CreateIncident creates a new incident on the status page
func (c *StatusPageClient) CreateIncident(ctx context.Context, name, impact string, componentIDs []string) (*StatusPageIncident, error) {
    components := make(map[string]string)
    for _, id := range componentIDs {
        components[id] = "major_outage"
    }

    payload := map[string]interface{}{
        "incident": map[string]interface{}{
            "name":       name,
            "status":     "investigating",
            "impact":     impact,
            "components": components,
        },
    }

    body, err := json.Marshal(payload)
    if err != nil {
        return nil, err
    }

    url := fmt.Sprintf("%s/pages/%s/incidents", c.apiURL, c.pageID)
    req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
    if err != nil {
        return nil, err
    }

    req.Header.Set("Authorization", "OAuth "+c.apiKey)
    req.Header.Set("Content-Type", "application/json")

    resp, err := c.httpClient.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusCreated {
        return nil, fmt.Errorf("status page API returned %d", resp.StatusCode)
    }

    var incident StatusPageIncident
    if err := json.NewDecoder(resp.Body).Decode(&incident); err != nil {
        return nil, err
    }

    return &incident, nil
}

// UpdateIncident updates an existing incident
func (c *StatusPageClient) UpdateIncident(ctx context.Context, incidentID, status, message string) (*StatusPageIncident, error) {
    payload := map[string]interface{}{
        "incident": map[string]interface{}{
            "status": status,
        },
    }

    if message != "" {
        payload["incident"].(map[string]interface{})["message"] = message
    }

    if status == "resolved" {
        now := time.Now().UTC()
        payload["incident"].(map[string]interface{})["resolved_at"] = now.Format(time.RFC3339)
    }

    body, err := json.Marshal(payload)
    if err != nil {
        return nil, err
    }

    url := fmt.Sprintf("%s/pages/%s/incidents/%s", c.apiURL, c.pageID, incidentID)
    req, err := http.NewRequestWithContext(ctx, http.MethodPatch, url, bytes.NewReader(body))
    if err != nil {
        return nil, err
    }

    req.Header.Set("Authorization", "OAuth "+c.apiKey)
    req.Header.Set("Content-Type", "application/json")

    resp, err := c.httpClient.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("status page API returned %d", resp.StatusCode)
    }

    var incident StatusPageIncident
    if err := json.NewDecoder(resp.Body).Decode(&incident); err != nil {
        return nil, err
    }

    return &incident, nil
}

// PostUpdate posts an update to an incident
func (c *StatusPageClient) PostUpdate(ctx context.Context, incidentID, status, body string) error {
    payload := map[string]interface{}{
        "incident_update": map[string]interface{}{
            "status": status,
            "body":   body,
        },
    }

    payloadBytes, err := json.Marshal(payload)
    if err != nil {
        return err
    }

    url := fmt.Sprintf("%s/pages/%s/incidents/%s/incident_updates", c.apiURL, c.pageID, incidentID)
    req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(payloadBytes))
    if err != nil {
        return err
    }

    req.Header.Set("Authorization", "OAuth "+c.apiKey)
    req.Header.Set("Content-Type", "application/json")

    resp, err := c.httpClient.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusCreated {
        return fmt.Errorf("status page API returned %d", resp.StatusCode)
    }

    return nil
}

// ComponentStatus maps severity to component status
func ComponentStatus(severity Severity) string {
    switch severity {
    case SeverityCritical:
        return "major_outage"
    case SeverityHigh:
        return "partial_outage"
    case SeverityMedium:
        return "degraded_performance"
    default:
        return "operational"
    }
}

// ImpactLevel maps severity to status page impact
func ImpactLevel(severity Severity) string {
    switch severity {
    case SeverityCritical:
        return "critical"
    case SeverityHigh:
        return "major"
    case SeverityMedium:
        return "minor"
    default:
        return "none"
    }
}
```

### 2.3 War Room Coordination

```go
package incident

import (
    "context"
    "fmt"
    "sync"
    "time"
)

// WarRoom manages incident war room coordination
type WarRoom struct {
    ID          string    `json:"id"`
    IncidentID  string    `json:"incident_id"`
    Name        string    `json:"name"`
    CreatedAt   time.Time `json:"created_at"`
    Active      bool      `json:"active"`

    // Communication channels
    SlackChannel string `json:"slack_channel"`
    ZoomLink     string `json:"zoom_link"`
    BridgeNumber string `json:"bridge_number,omitempty"`

    // Participants
    Participants []WarRoomParticipant `json:"participants"`

    // Status
    CurrentFocus string `json:"current_focus"`
    NextUpdateDue time.Time `json:"next_update_due"`

    mu sync.RWMutex
}

// WarRoomParticipant represents a war room participant
type WarRoomParticipant struct {
    ID         string    `json:"id"`
    Name       string    `json:"name"`
    Role       Role      `json:"role"`
    JoinedAt   time.Time `json:"joined_at"`
    LeftAt     *time.Time `json:"left_at,omitempty"`
    IsActive   bool      `json:"is_active"`
}

// WarRoomManager manages war rooms
type WarRoomManager struct {
    rooms map[string]*WarRoom
    mu    sync.RWMutex

    roomProvider RoomProvider
}

// RoomProvider creates communication channels
type RoomProvider interface {
    CreateSlackChannel(ctx context.Context, name string) (string, error)
    CreateZoomMeeting(ctx context.Context, topic string) (string, error)
    InviteParticipants(ctx context.Context, channelID string, userIDs []string) error
}

// NewWarRoomManager creates a new war room manager
func NewWarRoomManager(provider RoomProvider) *WarRoomManager {
    return &WarRoomManager{
        rooms:        make(map[string]*WarRoom),
        roomProvider: provider,
    }
}

// CreateWarRoom creates a new war room for an incident
func (m *WarRoomManager) CreateWarRoom(ctx context.Context, incident *Incident) (*WarRoom, error) {
    roomName := fmt.Sprintf("war-room-%s", incident.ID)

    // Create Slack channel
    slackChannel, err := m.roomProvider.CreateSlackChannel(ctx, roomName)
    if err != nil {
        return nil, fmt.Errorf("failed to create slack channel: %w", err)
    }

    // Create Zoom meeting
    zoomLink, err := m.roomProvider.CreateZoomMeeting(ctx, fmt.Sprintf("Incident %s: %s", incident.ID, incident.Title))
    if err != nil {
        return nil, fmt.Errorf("failed to create zoom meeting: %w", err)
    }

    room := &WarRoom{
        ID:            generateID(),
        IncidentID:    incident.ID,
        Name:          roomName,
        CreatedAt:     time.Now().UTC(),
        Active:        true,
        SlackChannel:  slackChannel,
        ZoomLink:      zoomLink,
        Participants:  make([]WarRoomParticipant, 0),
        NextUpdateDue: time.Now().UTC().Add(15 * time.Minute),
    }

    // Add initial participants
    for _, p := range incident.Participants {
        room.Participants = append(room.Participants, WarRoomParticipant{
            ID:       p.ID,
            Name:     p.Name,
            Role:     p.Role,
            JoinedAt: time.Now().UTC(),
            IsActive: true,
        })
    }

    m.mu.Lock()
    m.rooms[room.ID] = room
    m.mu.Unlock()

    return room, nil
}

// AddParticipant adds a participant to the war room
func (m *WarRoomManager) AddParticipant(roomID, userID, name string, role Role) error {
    m.mu.Lock()
    defer m.mu.Unlock()

    room, exists := m.rooms[roomID]
    if !exists {
        return fmt.Errorf("war room not found")
    }

    room.mu.Lock()
    defer room.mu.Unlock()

    participant := WarRoomParticipant{
        ID:       userID,
        Name:     name,
        Role:     role,
        JoinedAt: time.Now().UTC(),
        IsActive: true,
    }

    room.Participants = append(room.Participants, participant)

    return nil
}

// UpdateFocus updates the current focus of the war room
func (m *WarRoomManager) UpdateFocus(roomID, focus string) error {
    m.mu.Lock()
    defer m.mu.Unlock()

    room, exists := m.rooms[roomID]
    if !exists {
        return fmt.Errorf("war room not found")
    }

    room.mu.Lock()
    defer room.mu.Unlock()

    room.CurrentFocus = focus

    return nil
}

// ScheduleUpdate schedules the next status update
func (m *WarRoomManager) ScheduleUpdate(roomID string, dueTime time.Time) error {
    m.mu.Lock()
    defer m.mu.Unlock()

    room, exists := m.rooms[roomID]
    if !exists {
        return fmt.Errorf("war room not found")
    }

    room.mu.Lock()
    defer room.mu.Unlock()

    room.NextUpdateDue = dueTime

    return nil
}

// CloseWarRoom closes a war room
func (m *WarRoomManager) CloseWarRoom(roomID string) error {
    m.mu.Lock()
    defer m.mu.Unlock()

    room, exists := m.rooms[roomID]
    if !exists {
        return fmt.Errorf("war room not found")
    }

    room.mu.Lock()
    defer room.mu.Unlock()

    room.Active = false
    now := time.Now().UTC()
    for i := range room.Participants {
        if room.Participants[i].IsActive {
            room.Participants[i].IsActive = false
            room.Participants[i].LeftAt = &now
        }
    }

    return nil
}

// GetActiveWarRooms returns all active war rooms
func (m *WarRoomManager) GetActiveWarRooms() []*WarRoom {
    m.mu.RLock()
    defer m.mu.RUnlock()

    active := make([]*WarRoom, 0)
    for _, room := range m.rooms {
        if room.Active {
            active = append(active, room)
        }
    }

    return active
}
```

---

## 3. Production-Ready Configurations

### 3.1 Incident Response Automation

```yaml
# incident-response-automation.yaml
apiVersion: automation.io/v1
kind: IncidentResponseWorkflow
metadata:
  name: automated-incident-response
spec:
  # Trigger conditions
  triggers:
    - type: alert
      condition: severity == 'critical'
    - type: alert
      condition: service == 'payment-gateway' and error_rate > 0.01

  # Automated actions
  workflow:
    # Step 1: Create incident
    - name: create-incident
      action: incident.create
      params:
        title: "{{ alert.title }}"
        severity: "{{ alert.severity }}"
        services: "{{ alert.affected_services }}"
      output: incident_id

    # Step 2: Create war room
    - name: create-war-room
      action: warroom.create
      params:
        incident_id: "{{ steps.create-incident.output.incident_id }}"
      output: war_room

    # Step 3: Page on-call
    - name: page-oncall
      action: pagerduty.page
      params:
        service: "{{ alert.service }}"
        incident_id: "{{ steps.create-incident.output.incident_id }}"
        priority: "P1"

    # Step 4: Update status page (if customer-facing)
    - name: update-status-page
      action: statuspage.create
      condition: alert.customer_facing == true
      params:
        name: "{{ alert.title }}"
        impact: "{{ alert.severity | to_impact }}"
        components: "{{ alert.affected_services }}"
      output: status_page_url

    # Step 5: Create Slack channel
    - name: create-slack-channel
      action: slack.create_channel
      params:
        name: "incident-{{ steps.create-incident.output.incident_id }}"
        topic: "{{ alert.title }} - {{ steps.create-incident.output.incident_id }}"
      output: slack_channel

    # Step 6: Send initial communication
    - name: notify-stakeholders
      action: slack.post_message
      params:
        channel: "#incidents"
        message: |
          🚨 **Incident Declared: {{ steps.create-incident.output.incident_id }}**

          **Title:** {{ alert.title }}
          **Severity:** {{ alert.severity }}
          **Status Page:** {{ steps.update-status-page.output.status_page_url }}
          **War Room:** {{ steps.create-war-room.output.war_room.zoom_link }}

          @channel - Please join the war room for coordination.

---
# Auto-Remediation Rules
apiVersion: automation.io/v1
kind: AutoRemediation
metadata:
  name: common-issues
spec:
  rules:
    # Auto-restart service on OOM
    - name: oom-restart
      condition: alert.reason == 'OOMKilled'
      actions:
        - type: kubernetes.restart
          params:
            deployment: "{{ alert.labels.deployment }}"
            namespace: "{{ alert.labels.namespace }}"
        - type: incident.add_note
          params:
            text: "Auto-restarted pod due to OOMKill"

    # Auto-scale on high CPU
    - name: cpu-scale-up
      condition: alert.metric == 'cpu_utilization' and alert.value > 80
      throttle: 1 per 10m
      actions:
        - type: kubernetes.scale
          params:
            deployment: "{{ alert.labels.deployment }}"
            namespace: "{{ alert.labels.namespace }}"
            replicas: "+2"
        - type: slack.notify
          params:
            channel: "#scaling-events"
            message: "Auto-scaled {{ alert.labels.deployment }} due to high CPU"

    # Circuit breaker for downstream failures
    - name: circuit-breaker
      condition: alert.type == 'dependency_timeout' and alert.consecutive_failures > 5
      actions:
        - type: config.update
          params:
            key: "circuit_breaker.{{ alert.dependency }}.enabled"
            value: "true"
        - type: incident.add_note
          params:
            text: "Circuit breaker opened for {{ alert.dependency }}"
```

---

## 4. Security Considerations

### 4.1 Incident Security Matrix

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Incident Management Security Matrix                       │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  Phase              │  Risk                │  Control                       │
├─────────────────────┼──────────────────────┼────────────────────────────────│
│  Detection          │  False positives     │  Corroboration requirements    │
│                     │  Alert suppression   │  Tamper-proof alerting         │
│                     │                      │  Multi-source confirmation     │
│  ───────────────────┼──────────────────────┼────────────────────────────────│
│  Communication      │  Information leak    │  Private channels for details  │
│                     │  Social engineering  │  Identity verification         │
│                     │  Unauthorized access │  Channel access controls       │
│  ───────────────────┼──────────────────────┼────────────────────────────────│
│  Response           │  Insider threat      │  Dual authorization            │
│                     │  Privilege escalation│  Just-in-time access           │
│                     │  Evidence tampering  │  Immutable audit logs          │
│  ───────────────────┼──────────────────────┼────────────────────────────────│
│  Resolution         │  Backdoor insertion  │  Change review requirements    │
│                     │  Incomplete fix      │  Validation checklists         │
│                     │  Configuration drift │  Infrastructure as Code        │
│  ───────────────────┼──────────────────────┼────────────────────────────────│
│  Post-Incident      │  Data retention      │  Encrypted storage             │
│                     │  Knowledge loss      │  Documented procedures         │
│                     │  Repeat incidents    │  Action item tracking          │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 5. Compliance Requirements

### 5.1 Regulatory Compliance Mapping

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                 Incident Management Compliance Requirements                  │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  SOC 2 TYPE II                    ISO 27001                                 │
│  ━━━━━━━━━━━━━━━━                 ━━━━━━━━━━━                               │
│                                                                             │
│  CC6.1 - Logical access           A.16.1 - Responsibilities                 │
│  CC6.2 - Prior to access          A.16.2 - Reporting                        │
│  CC7.3 - Incident detection       A.16.3 - Assessment                       │
│  CC7.4 - Incident response        A.16.4 - Evidence                         │
│  CC7.5 - Incident recovery        A.16.5 - Learn lessons                    │
│                                                                             │
│  HIPAA                            PCI DSS                                   │
│  ━━━━━━━━                         ━━━━━━━━━                                 │
│                                                                             │
│  §164.312(a)(1) - Access control  Req 12.10.1 - IR plan                     │
│  §164.312(a)(2)(ii) - Emergency   Req 12.10.2 - IR roles                    │
│  §164.312(b) - Audit controls     Req 12.10.3 - IR testing                  │
│  §164.312(c)(1) - Integrity       Req 12.10.4 - 24/7 response               │
│  §164.312(d) - Person auth        Req 12.10.5 - IR escalation               │
│                                                                             │
│  GDPR                             NIST SP 800-61                            │
│  ━━━━━━                           ━━━━━━━━━━━━━━                            │
│                                                                             │
│  Art. 33 - Breach notification    Detection & Analysis                      │
│  Art. 34 - Data subject notif     Containment/Eradication                   │
│  Art. 35 - DPIA                   Post-Incident Activity                    │
│  Art. 5(1)(f) - Security          Communication                             │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 6. Decision Matrices

### 6.1 Severity Assignment Matrix

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                       Severity Assignment Matrix                             │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  Impact →        │  Complete  │  Major     │  Minor     │  None            │
│  Likelihood ↓    │  Outage    │  Degraded  │  Issue     │  Impact          │
├──────────────────┼────────────┼────────────┼────────────┼──────────────────│
│  Certain         │  SEV 1     │  SEV 1     │  SEV 2     │  SEV 3           │
│                  │  Critical  │  Critical  │  High      │  Medium          │
├──────────────────┼────────────┼────────────┼────────────┼──────────────────│
│  Likely          │  SEV 1     │  SEV 2     │  SEV 3     │  SEV 4           │
│                  │  Critical  │  High      │  Medium    │  Low             │
├──────────────────┼────────────┼────────────┼────────────┼──────────────────│
│  Possible        │  SEV 2     │  SEV 3     │  SEV 4     │  SEV 4           │
│                  │  High      │  Medium    │  Low       │  Low             │
├──────────────────┼────────────┼────────────┼────────────┼──────────────────│
│  Unlikely        │  SEV 3     │  SEV 4     │  SEV 4     │  Track           │
│                  │  Medium    │  Low       │  Low       │                  │
│                                                                             │
│  Additional Escalators:                                                    │
│  • Data breach → Always SEV 1 minimum                                      │
│  • Regulatory requirement → +1 severity level                              │
│  • Public media attention → +1 severity level                              │
│  • Financial impact > $100K/hr → SEV 1                                     │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 6.2 Communication Cadence Matrix

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      Communication Cadence Matrix                            │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  Audience        │  SEV 1      │  SEV 2      │  SEV 3      │  SEV 4        │
├──────────────────┼─────────────┼─────────────┼─────────────┼───────────────│
│  War Room        │  Real-time  │  Real-time  │  As needed  │  N/A          │
│  (Responders)    │  Chat       │  Chat       │             │               │
├──────────────────┼─────────────┼─────────────┼─────────────┼───────────────│
│  Internal        │  15 min     │  30 min     │  1 hour     │  4 hours      │
│  (#incidents)    │             │             │             │               │
├──────────────────┼─────────────┼─────────────┼─────────────┼───────────────│
│  Status Page     │  30 min     │  1 hour     │  4 hours    │  Next day     │
│  (External)      │             │             │             │               │
├──────────────────┼─────────────┼─────────────┼─────────────┼───────────────│
│  Executive       │  1 hour     │  4 hours    │  24 hours   │  Weekly       │
│  (Leadership)    │             │             │             │  summary      │
├──────────────────┼─────────────┼─────────────┼─────────────┼───────────────│
│  Customers       │  Proactive  │  On inquiry │  N/A        │  N/A          │
│  (Direct)        │  outreach   │             │             │               │
├──────────────────┼─────────────┼─────────────┼─────────────┼───────────────│
│  Regulators      │  Immediate  │  Per policy │  Per policy │  Per policy   │
│                  │  (if req)   │             │             │               │
│                                                                             │
│  Template Updates:                                                         │
│  • Ongoing: "We are investigating [issue] affecting [scope]"               │
│  • Identified: "We have identified [cause] and are working on [action]"    │
│  • Monitoring: "We have implemented [fix] and are monitoring"              │
│  • Resolved: "Service is restored. [Summary of impact]"                    │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 6.3 Resource Escalation Matrix

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                       Resource Escalation Matrix                             │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  Time from      │  Response       │  Additional    │  External            │
│  Detection      │  Team           │  Resources     │  Support             │
├─────────────────┼─────────────────┼────────────────┼──────────────────────│
│  0-15 min       │  On-call        │  IC assigned   │                      │
│                 │  engineer       │  War room      │                      │
│                 │                 │  opened        │                      │
├─────────────────┼─────────────────┼────────────────┼──────────────────────│
│  15-30 min      │  Secondary      │  Team lead     │  Vendor support      │
│                 │  on-call        │  joined        │  (if applicable)     │
├─────────────────┼─────────────────┼────────────────┼──────────────────────│
│  30 min - 1 hr  │  Full team      │  Engineering   │  External            │
│                 │  assembled      │  manager       │  consultants         │
├─────────────────┼─────────────────┼────────────────┼──────────────────────│
│  1-4 hours      │  Cross-team     │  Director      │  Executive           │
│                 │  support        │  engaged       │  briefing            │
├─────────────────┼─────────────────┼────────────────┼──────────────────────│
│  > 4 hours      │  All hands      │  C-level       │  Legal/PR/           │
│                 │                 │  engaged       │  Regulatory          │
│                                                                             │
│  Trigger Conditions:                                                       │
│  • SEV 1: All resources within 30 minutes                                  │
│  • Customer data involved: Legal immediately                               │
│  • Media attention: PR immediately                                         │
│  • Regulatory requirement: Compliance immediately                          │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 7. Best Practices Summary

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Incident Management Best Practices                        │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  PREPARATION                                                                │
│  ✓ Document and test incident response procedures                           │
│  ✓ Maintain updated service dependency maps                                 │
│  ✓ Pre-provision war room tools (Slack, Zoom, etc.)                         │
│  ✓ Create incident response runbooks                                        │
│  ✓ Conduct regular incident response drills                                 │
│  ✓ Define clear severity classification                                     │
│  ✓ Establish communication templates                                        │
│                                                                             │
│  DETECTION & RESPONSE                                                       │
│  ✓ Automate incident creation for critical alerts                           │
│  ✓ Assign Incident Commander early                                          │
│  ✓ Focus on mitigation before root cause                                    │
│  ✓ Communicate early and often                                              │
│  ✓ Document timeline as events occur                                        │
│  ✓ Preserve evidence before remediation                                     │
│  ✓ Use feature flags for rapid rollback                                     │
│                                                                             │
│  COMMUNICATION                                                              │
│  ✓ Separate technical and customer communication channels                   │
│  ✓ Status page updates for customer-facing issues                           │
│  ✓ Executive briefings at appropriate cadence                               │
│  ✓ Be honest about impact and timeline                                      │
│  ✓ Never promise specific resolution times                                  │
│  ✓ Include next update time in all communications                           │
│                                                                             │
│  POST-INCIDENT                                                              │
│  ✓ Hold blameless postmortem within 1 week                                  │
│  ✓ Focus on systemic issues, not individual mistakes                        │
│  ✓ Share postmortems company-wide                                           │
│  ✓ Track action items to completion                                         │
│  ✓ Update runbooks with lessons learned                                     │
│  ✓ Measure and improve MTTD/MTTR over time                                  │
│                                                                             │
│  CONTINUOUS IMPROVEMENT                                                     │
│  ✓ Regular review of incident trends                                        │
│  ✓ Update automation based on common issues                                 │
│  ✓ Refine severity classifications based on experience                      │
│  ✓ Cross-train team members in IC role                                      │
│  ✓ Celebrate successful incident responses                                  │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## References

1. Google SRE Book - Managing Incidents
2. NIST SP 800-61 - Computer Security Incident Handling Guide
3. PagerDuty Incident Response Guide
4. Atlassian Incident Management Handbook
5. Site Reliability Workbook - Incident Response

# On-Call Procedures

> **分类**: 工程与云原生
> **标签**: #oncall #sre #incident-response #operations #rotations
> **参考**: Google SRE, PagerDuty, Incident Management Best Practices

---

## 1. Formal Definition

### 1.1 What is On-Call?

On-call is an operational responsibility model where engineers are designated to respond to alerts, incidents, and operational issues outside of normal business hours. It is a critical component of Site Reliability Engineering (SRE) that ensures continuous service availability and rapid incident response.

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                        On-Call Ecosystem                                     │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│   ┌─────────────┐    ┌─────────────┐    ┌─────────────┐    ┌─────────────┐ │
│   │   Primary   │    │  Secondary  │    │   Shadow    │    │  Incident   │ │
│   │   On-Call   │◄──►│   On-Call   │    │   Engineer  │    │  Commander  │ │
│   │             │    │             │    │             │    │             │ │
│   │ • Responds  │    │ • Backup    │    │ • Learns    │    │ • Major     │ │
│   │ • Triage    │    │ • Escalate  │    │ • Observes  │    │   incidents │ │
│   │ • Fixes     │    │ • Support   │    │ • Assists   │    │ • Coordinates│ │
│   │ • Pages     │    │             │    │             │    │             │ │
│   └──────┬──────┘    └─────────────┘    └─────────────┘    └─────────────┘ │
│          │                                                                  │
│          ▼                                                                  │
│   ┌─────────────────────────────────────────────────────────────────┐      │
│   │                      Response Flow                               │      │
│   ├─────────────────────────────────────────────────────────────────┤      │
│   │                                                                 │      │
│   │   Page Received ──► Acknowledge ──► Triage ──► Resolve/escalate│      │
│   │        │                │             │              │          │      │
│   │        │                │             │              ▼          │      │
│   │        │                │             │         ┌──────────┐     │      │
│   │        │                │             │         │ Escalate │─────┼──────┼──►
│   │        │                │             │         │ if needed│     │      │
│   │        │                │             │         └──────────┘     │      │
│   │        │                │             ▼                          │      │
│   │        │                │      ┌──────────┐                      │      │
│   │        │                │      │  Resolve │◄─────────────────────┼──────┘
│   │        │                │      └────┬─────┘                      │
│   │        │                │           │                            │
│   │        │                │           ▼                            │
│   │        │                │      ┌──────────┐                      │
│   │        │                └─────►│ Post-mortem/Runbook update     │
│   │        │                       └──────────┘                      │
│   │        ▼                                                          │
│   │   No ack in 5 min ──► Escalate to secondary                       │
│   │                                                                 │      │
│   └─────────────────────────────────────────────────────────────────┘      │
│                                                                             │
│   SLA Targets:  Ack < 5 min │ Triage < 15 min │ Resolution varies          │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 1.2 On-Call Rotation Models

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                     On-Call Rotation Models                                  │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  FOLLOW-THE-SUN              WEEKLY ROTATION              DAILY ROTATION    │
│  (Global teams)              (Standard)                   (High-intensity)  │
│                                                                             │
│  ┌─────────┐                ┌─────────┐                  ┌─────────┐       │
│  │  APAC   │──►             │  Week 1 │  Engineer A       │  Mon    │  A   │
│  │ Primary │                ├─────────┤                  ├─────────┤       │
│  └─────────┘                │  Week 2 │  Engineer B       │  Tue    │  B   │
│       │                     ├─────────┤                  ├─────────┤       │
│       ▼                     │  Week 3 │  Engineer C       │  Wed    │  C   │
│  ┌─────────┐                ├─────────┤                  ├─────────┤       │
│  │  EMEA   │──►             │  Week 4 │  Engineer D       │  Thu    │  D   │
│  │ Primary │                └─────────┘                  ├─────────┤       │
│  └─────────┘                                              │  Fri    │  E   │
│       │                                                   ├─────────┤       │
│       ▼                                                   │  Sat    │  F   │
│  ┌─────────┐                                              ├─────────┤       │
│  │ Americas│──► (cycle)                                   │  Sun    │  G   │
│  │ Primary │                                              └─────────┘       │
│  └─────────┘                                                               │
│                                                                             │
│  PROS: 24/7 coverage        PROS: Less context switching  PROS: Fair burden│
│        No night shifts            Deep investigation       Consistent rest │
│        Regional expertise         Time for complex fixes   Fresh eyes      │
│  CONS: Handoff complexity   CONS: Long on-call weeks    CONS: Context loss│
│        Coordination needed        Fatigue accumulation     Daily stress     │
│                                                                             │
│  HYBRID MODEL (Recommended for most teams)                                  │
│  ┌─────────────────────────────────────────────────────────────────┐       │
│  │  • 1-week primary rotation                                      │       │
│  │  • Secondary on-call (flexible, no fixed schedule)              │       │
│  │  • Follow-the-sun for P0 only (global team escalation)          │       │
│  │  • Business hours coverage by dedicated team                    │       │
│  └─────────────────────────────────────────────────────────────────┘       │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 2. Implementation Patterns in Go

### 2.1 On-Call Schedule Manager

```go
package oncall

import (
    "context"
    "errors"
    "fmt"
    "sync"
    "time"
)

// RotationType defines the type of rotation
type RotationType string

const (
    RotationTypeWeekly  RotationType = "weekly"
    RotationTypeDaily   RotationType = "daily"
    RotationTypeCustom  RotationType = "custom"
)

// Role defines the on-call role
type Role string

const (
    RolePrimary   Role = "primary"
    RoleSecondary Role = "secondary"
    RoleShadow    Role = "shadow"
    RoleManager   Role = "manager"
)

// Engineer represents an on-call engineer
type Engineer struct {
    ID        string   `json:"id"`
    Name      string   `json:"name"`
    Email     string   `json:"email"`
    Phone     string   `json:"phone"`
    SlackID   string   `json:"slack_id"`
    Timezone  string   `json:"timezone"`
    Skills    []string `json:"skills"`
    EscalationLevel int `json:"escalation_level"`

    // Preferences
    BlackoutDates []time.Time `json:"blackout_dates,omitempty"`
    PreferredRole Role        `json:"preferred_role,omitempty"`
}

// Shift represents a single on-call shift
type Shift struct {
    ID          string    `json:"id"`
    StartTime   time.Time `json:"start_time"`
    EndTime     time.Time `json:"end_time"`
    EngineerID  string    `json:"engineer_id"`
    Role        Role      `json:"role"`
    Status      string    `json:"status"` // scheduled, active, completed, overridden
    OverrideBy  *string   `json:"override_by,omitempty"`
    OverrideReason string `json:"override_reason,omitempty"`
}

// Schedule represents the on-call schedule for a team
type Schedule struct {
    ID              string       `json:"id"`
    Name            string       `json:"name"`
    Team            string       `json:"team"`
    RotationType    RotationType `json:"rotation_type"`
    ShiftDuration   time.Duration `json:"shift_duration"`
    HandoffTime     time.Time    `json:"handoff_time"`

    Engineers       []Engineer   `json:"engineers"`
    Shifts          []Shift      `json:"shifts"`

    // Configuration
    EscalationPolicy EscalationPolicy `json:"escalation_policy"`

    mu              sync.RWMutex
}

// EscalationPolicy defines how to escalate alerts
type EscalationPolicy struct {
    Levels []EscalationLevel `json:"levels"`
}

// EscalationLevel defines a single escalation level
type EscalationLevel struct {
    Level          int           `json:"level"`
    NotifyAfter    time.Duration `json:"notify_after"`
    Targets        []string      `json:"targets"` // Engineer IDs
    ContactMethods []string      `json:"contact_methods"` // push, sms, phone
}

// ScheduleManager manages on-call schedules
type ScheduleManager struct {
    schedules map[string]*Schedule
    mu        sync.RWMutex

    notifier  Notifier
    store     ScheduleStore
}

// Notifier sends notifications
type Notifier interface {
    NotifyShiftStart(ctx context.Context, shift *Shift, engineer *Engineer) error
    NotifyShiftEnd(ctx context.Context, shift *Shift, engineer *Engineer) error
    NotifyEscalation(ctx context.Context, level EscalationLevel, alert *Alert) error
}

// ScheduleStore persists schedule data
type ScheduleStore interface {
    SaveSchedule(ctx context.Context, schedule *Schedule) error
    LoadSchedule(ctx context.Context, id string) (*Schedule, error)
    SaveShift(ctx context.Context, shift *Shift) error
    GetCurrentShift(ctx context.Context, scheduleID string, role Role) (*Shift, error)
}

// NewScheduleManager creates a new schedule manager
func NewScheduleManager(notifier Notifier, store ScheduleStore) *ScheduleManager {
    return &ScheduleManager{
        schedules: make(map[string]*Schedule),
        notifier:  notifier,
        store:     store,
    }
}

// CreateSchedule creates a new on-call schedule
func (m *ScheduleManager) CreateSchedule(name, team string, rotationType RotationType, shiftDuration time.Duration) (*Schedule, error) {
    if shiftDuration == 0 {
        shiftDuration = 7 * 24 * time.Hour // Default to weekly
    }

    schedule := &Schedule{
        ID:            generateID(),
        Name:          name,
        Team:          team,
        RotationType:  rotationType,
        ShiftDuration: shiftDuration,
        HandoffTime:   time.Date(0, 1, 1, 9, 0, 0, 0, time.UTC), // 9 AM UTC
        Engineers:     make([]Engineer, 0),
        Shifts:        make([]Shift, 0),
    }

    m.mu.Lock()
    m.schedules[schedule.ID] = schedule
    m.mu.Unlock()

    if err := m.store.SaveSchedule(context.Background(), schedule); err != nil {
        return nil, fmt.Errorf("failed to save schedule: %w", err)
    }

    return schedule, nil
}

// AddEngineer adds an engineer to the schedule
func (m *ScheduleManager) AddEngineer(scheduleID string, engineer Engineer) error {
    m.mu.Lock()
    defer m.mu.Unlock()

    schedule, exists := m.schedules[scheduleID]
    if !exists {
        return errors.New("schedule not found")
    }

    schedule.Engineers = append(schedule.Engineers, engineer)

    return m.store.SaveSchedule(context.Background(), schedule)
}

// GenerateSchedule generates shifts for the given period
func (m *ScheduleManager) GenerateSchedule(scheduleID string, startDate, endDate time.Time) error {
    m.mu.Lock()
    defer m.mu.Unlock()

    schedule, exists := m.schedules[scheduleID]
    if !exists {
        return errors.New("schedule not found")
    }

    if len(schedule.Engineers) == 0 {
        return errors.New("no engineers in schedule")
    }

    // Generate shifts
    current := startDate
    engineerIndex := 0

    for current.Before(endDate) {
        // Calculate shift times
        shiftStart := time.Date(
            current.Year(), current.Month(), current.Day(),
            schedule.HandoffTime.Hour(), schedule.HandoffTime.Minute(), 0, 0,
            schedule.HandoffTime.Location(),
        )
        shiftEnd := shiftStart.Add(schedule.ShiftDuration)

        // Create primary shift
        primaryShift := Shift{
            ID:         generateID(),
            StartTime:  shiftStart,
            EndTime:    shiftEnd,
            EngineerID: schedule.Engineers[engineerIndex].ID,
            Role:       RolePrimary,
            Status:     "scheduled",
        }
        schedule.Shifts = append(schedule.Shifts, primaryShift)

        // Create secondary shift (next engineer)
        secondaryIndex := (engineerIndex + 1) % len(schedule.Engineers)
        secondaryShift := Shift{
            ID:         generateID(),
            StartTime:  shiftStart,
            EndTime:    shiftEnd,
            EngineerID: schedule.Engineers[secondaryIndex].ID,
            Role:       RoleSecondary,
            Status:     "scheduled",
        }
        schedule.Shifts = append(schedule.Shifts, secondaryShift)

        // Move to next shift
        current = shiftEnd
        engineerIndex = (engineerIndex + 1) % len(schedule.Engineers)
    }

    return m.store.SaveSchedule(context.Background(), schedule)
}

// GetCurrentOnCall returns the current on-call engineer(s)
func (m *ScheduleManager) GetCurrentOnCall(scheduleID string, role Role) (*Shift, *Engineer, error) {
    now := time.Now()

    m.mu.RLock()
    schedule, exists := m.schedules[scheduleID]
    m.mu.RUnlock()

    if !exists {
        return nil, nil, errors.New("schedule not found")
    }

    // Find current shift
    for i := range schedule.Shifts {
        shift := &schedule.Shifts[i]
        if shift.Role == role && now.After(shift.StartTime) && now.Before(shift.EndTime) {
            // Find engineer
            for _, eng := range schedule.Engineers {
                if eng.ID == shift.EngineerID {
                    return shift, &eng, nil
                }
            }
        }
    }

    return nil, nil, errors.New("no active shift found")
}

// OverrideShift allows manual override of a shift
func (m *ScheduleManager) OverrideShift(shiftID, newEngineerID, reason string) error {
    m.mu.Lock()
    defer m.mu.Unlock()

    for _, schedule := range m.schedules {
        for i := range schedule.Shifts {
            if schedule.Shifts[i].ID == shiftID {
                schedule.Shifts[i].OverrideBy = &newEngineerID
                schedule.Shifts[i].OverrideReason = reason
                schedule.Shifts[i].Status = "overridden"

                return m.store.SaveShift(context.Background(), &schedule.Shifts[i])
            }
        }
    }

    return errors.New("shift not found")
}

// StartShiftMonitoring begins monitoring for shift changes
func (m *ScheduleManager) StartShiftMonitoring(ctx context.Context) {
    ticker := time.NewTicker(1 * time.Minute)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            m.checkShiftTransitions(ctx)
        case <-ctx.Done():
            return
        }
    }
}

// checkShiftTransitions checks for and handles shift transitions
func (m *ScheduleManager) checkShiftTransitions(ctx context.Context) {
    m.mu.RLock()
    schedules := make([]*Schedule, 0, len(m.schedules))
    for _, s := range m.schedules {
        schedules = append(schedules, s)
    }
    m.mu.RUnlock()

    now := time.Now()

    for _, schedule := range schedules {
        for i := range schedule.Shifts {
            shift := &schedule.Shifts[i]

            // Check if shift just started
            if now.After(shift.StartTime) && now.Before(shift.StartTime.Add(1*time.Minute)) {
                if shift.Status == "scheduled" {
                    shift.Status = "active"

                    // Notify engineer
                    for _, eng := range schedule.Engineers {
                        if eng.ID == shift.EngineerID {
                            m.notifier.NotifyShiftStart(ctx, shift, &eng)
                            break
                        }
                    }

                    m.store.SaveShift(ctx, shift)
                }
            }

            // Check if shift just ended
            if now.After(shift.EndTime) && now.Before(shift.EndTime.Add(1*time.Minute)) {
                if shift.Status == "active" {
                    shift.Status = "completed"

                    // Notify engineer
                    for _, eng := range schedule.Engineers {
                        if eng.ID == shift.EngineerID {
                            m.notifier.NotifyShiftEnd(ctx, shift, &eng)
                            break
                        }
                    }

                    m.store.SaveShift(ctx, shift)
                }
            }
        }
    }
}

// generateID generates a unique ID (simplified)
func generateID() string {
    return fmt.Sprintf("sch_%d", time.Now().UnixNano())
}
```

### 2.2 Incident Response Workflow

```go
package oncall

import (
    "context"
    "fmt"
    "sync"
    "time"
)

// IncidentState represents the state of an incident
type IncidentState string

const (
    IncidentStateNew        IncidentState = "new"
    IncidentStateAcknowledged IncidentState = "acknowledged"
    IncidentStateInvestigating IncidentState = "investigating"
    IncidentStateMitigated  IncidentState = "mitigated"
    IncidentStateResolved   IncidentState = "resolved"
    IncidentStatePostmortem IncidentState = "postmortem"
    IncidentStateClosed     IncidentState = "closed"
)

// IncidentSeverity represents incident severity
type IncidentSeverity string

const (
    IncidentSeverityCritical IncidentSeverity = "critical"
    IncidentSeverityHigh     IncidentSeverity = "high"
    IncidentSeverityMedium   IncidentSeverity = "medium"
    IncidentSeverityLow      IncidentSeverity = "low"
)

// Incident represents an operational incident
type Incident struct {
    ID          string           `json:"id"`
    Title       string           `json:"title"`
    Description string           `json:"description"`
    Severity    IncidentSeverity `json:"severity"`
    State       IncidentState    `json:"state"`

    // Timestamps
    CreatedAt       time.Time  `json:"created_at"`
    AcknowledgedAt  *time.Time `json:"acknowledged_at,omitempty"`
    AcknowledgedBy  *string    `json:"acknowledged_by,omitempty"`
    MitigatedAt     *time.Time `json:"mitigated_at,omitempty"`
    ResolvedAt      *time.Time `json:"resolved_at,omitempty"`
    ClosedAt        *time.Time `json:"closed_at,omitempty"`

    // Assignments
    CommanderID    string   `json:"commander_id"`
    Responders     []string `json:"responders"`

    // Details
    Alerts      []string          `json:"alerts"`
    AffectedServices []string     `json:"affected_services"`
    Labels      map[string]string `json:"labels"`

    // Communication
    SlackChannel   string `json:"slack_channel,omitempty"`
    ConferenceURL  string `json:"conference_url,omitempty"`

    // Updates
    Updates []IncidentUpdate `json:"updates"`

    mu sync.RWMutex
}

// IncidentUpdate represents an update to an incident
type IncidentUpdate struct {
    ID        string    `json:"id"`
    Timestamp time.Time `json:"timestamp"`
    AuthorID  string    `json:"author_id"`
    Message   string    `json:"message"`
    State     IncidentState `json:"state,omitempty"`
}

// IncidentManager manages incident lifecycle
type IncidentManager struct {
    incidents    map[string]*Incident
    incidentsMu  sync.RWMutex

    currentID    int
    idMu         sync.Mutex

    notifier     IncidentNotifier
    store        IncidentStore
}

// IncidentNotifier sends incident notifications
type IncidentNotifier interface {
    NotifyNewIncident(ctx context.Context, incident *Incident) error
    NotifyIncidentUpdate(ctx context.Context, incident *Incident, update *IncidentUpdate) error
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
    Severity   *IncidentSeverity
    Commander  string
    StartTime  *time.Time
    EndTime    *time.Time
}

// NewIncidentManager creates a new incident manager
func NewIncidentManager(notifier IncidentNotifier, store IncidentStore) *IncidentManager {
    return &IncidentManager{
        incidents: make(map[string]*Incident),
        notifier:  notifier,
        store:     store,
    }
}

// CreateIncident creates a new incident
func (m *IncidentManager) CreateIncident(ctx context.Context, title, description string, severity IncidentSeverity, commanderID string) (*Incident, error) {
    m.idMu.Lock()
    m.currentID++
    id := fmt.Sprintf("INC-%d", m.currentID)
    m.idMu.Unlock()

    incident := &Incident{
        ID:           id,
        Title:        title,
        Description:  description,
        Severity:     severity,
        State:        IncidentStateNew,
        CreatedAt:    time.Now().UTC(),
        CommanderID:  commanderID,
        Responders:   []string{commanderID},
        Alerts:       make([]string, 0),
        Labels:       make(map[string]string),
        Updates:      make([]IncidentUpdate, 0),
    }

    m.incidentsMu.Lock()
    m.incidents[id] = incident
    m.incidentsMu.Unlock()

    if err := m.store.SaveIncident(ctx, incident); err != nil {
        return nil, fmt.Errorf("failed to save incident: %w", err)
    }

    // Notify
    if err := m.notifier.NotifyNewIncident(ctx, incident); err != nil {
        // Log but don't fail
        fmt.Printf("Failed to notify new incident: %v\n", err)
    }

    return incident, nil
}

// AcknowledgeIncident acknowledges an incident
func (m *IncidentManager) AcknowledgeIncident(ctx context.Context, incidentID, engineerID string) error {
    incident, err := m.getIncident(incidentID)
    if err != nil {
        return err
    }

    incident.mu.Lock()
    defer incident.mu.Unlock()

    if incident.State != IncidentStateNew {
        return fmt.Errorf("incident already acknowledged")
    }

    now := time.Now().UTC()
    incident.State = IncidentStateAcknowledged
    incident.AcknowledgedAt = &now
    incident.AcknowledgedBy = &engineerID

    update := IncidentUpdate{
        ID:        generateID(),
        Timestamp: now,
        AuthorID:  engineerID,
        Message:   "Incident acknowledged",
        State:     IncidentStateAcknowledged,
    }
    incident.Updates = append(incident.Updates, update)

    if err := m.store.SaveIncident(ctx, incident); err != nil {
        return fmt.Errorf("failed to save incident: %w", err)
    }

    return nil
}

// UpdateIncidentState updates the incident state
func (m *IncidentManager) UpdateIncidentState(ctx context.Context, incidentID, engineerID string, newState IncidentState, message string) error {
    incident, err := m.getIncident(incidentID)
    if err != nil {
        return err
    }

    incident.mu.Lock()
    defer incident.mu.Unlock()

    // Validate state transition
    if !isValidStateTransition(incident.State, newState) {
        return fmt.Errorf("invalid state transition from %s to %s", incident.State, newState)
    }

    now := time.Now().UTC()
    incident.State = newState

    // Update timestamps based on state
    switch newState {
    case IncidentStateMitigated:
        incident.MitigatedAt = &now
    case IncidentStateResolved:
        incident.ResolvedAt = &now
    case IncidentStateClosed:
        incident.ClosedAt = &now
    }

    update := IncidentUpdate{
        ID:        generateID(),
        Timestamp: now,
        AuthorID:  engineerID,
        Message:   message,
        State:     newState,
    }
    incident.Updates = append(incident.Updates, update)

    if err := m.store.SaveIncident(ctx, incident); err != nil {
        return fmt.Errorf("failed to save incident: %w", err)
    }

    // Notify
    if newState == IncidentStateResolved {
        if err := m.notifier.NotifyIncidentResolved(ctx, incident); err != nil {
            fmt.Printf("Failed to notify incident resolved: %v\n", err)
        }
    } else {
        if err := m.notifier.NotifyIncidentUpdate(ctx, incident, &update); err != nil {
            fmt.Printf("Failed to notify incident update: %v\n", err)
        }
    }

    return nil
}

// AddResponder adds a responder to the incident
func (m *IncidentManager) AddResponder(incidentID, engineerID string) error {
    incident, err := m.getIncident(incidentID)
    if err != nil {
        return err
    }

    incident.mu.Lock()
    defer incident.mu.Unlock()

    for _, r := range incident.Responders {
        if r == engineerID {
            return nil // Already added
        }
    }

    incident.Responders = append(incident.Responders, engineerID)
    return m.store.SaveIncident(context.Background(), incident)
}

// AddUpdate adds an update to the incident
func (m *IncidentManager) AddUpdate(ctx context.Context, incidentID, engineerID, message string) error {
    incident, err := m.getIncident(incidentID)
    if err != nil {
        return err
    }

    incident.mu.Lock()
    defer incident.mu.Unlock()

    update := IncidentUpdate{
        ID:        generateID(),
        Timestamp: time.Now().UTC(),
        AuthorID:  engineerID,
        Message:   message,
    }
    incident.Updates = append(incident.Updates, update)

    if err := m.store.SaveIncident(ctx, incident); err != nil {
        return fmt.Errorf("failed to save incident: %w", err)
    }

    return m.notifier.NotifyIncidentUpdate(ctx, incident, &update)
}

// getIncident retrieves an incident by ID
func (m *IncidentManager) getIncident(id string) (*Incident, error) {
    m.incidentsMu.RLock()
    incident, exists := m.incidents[id]
    m.incidentsMu.RUnlock()

    if !exists {
        return nil, fmt.Errorf("incident not found: %s", id)
    }

    return incident, nil
}

// GetOpenIncidents returns all open incidents
func (m *IncidentManager) GetOpenIncidents() []*Incident {
    m.incidentsMu.RLock()
    defer m.incidentsMu.RUnlock()

    open := make([]*Incident, 0)
    for _, incident := range m.incidents {
        if incident.State != IncidentStateResolved && incident.State != IncidentStateClosed {
            open = append(open, incident)
        }
    }

    return open
}

// isValidStateTransition checks if a state transition is valid
func isValidStateTransition(from, to IncidentState) bool {
    validTransitions := map[IncidentState][]IncidentState{
        IncidentStateNew:           {IncidentStateAcknowledged},
        IncidentStateAcknowledged:  {IncidentStateInvestigating},
        IncidentStateInvestigating: {IncidentStateMitigated},
        IncidentStateMitigated:     {IncidentStateResolved},
        IncidentStateResolved:      {IncidentStatePostmortem},
        IncidentStatePostmortem:    {IncidentStateClosed},
    }

    valid, exists := validTransitions[from]
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

// SLA Breach calculation
func (m *IncidentManager) CheckSLABreach(incident *Incident) (bool, time.Duration) {
    var sla time.Duration

    switch incident.Severity {
    case IncidentSeverityCritical:
        sla = 5 * time.Minute
    case IncidentSeverityHigh:
        sla = 15 * time.Minute
    case IncidentSeverityMedium:
        sla = 1 * time.Hour
    case IncidentSeverityLow:
        sla = 4 * time.Hour
    default:
        sla = 1 * time.Hour
    }

    elapsed := time.Since(incident.CreatedAt)
    breached := elapsed > sla

    return breached, sla - elapsed
}
```

### 2.3 Handoff Report Generator

```go
package oncall

import (
    "bytes"
    "context"
    "fmt"
    "html/template"
    "time"
)

// HandoffReport represents a shift handoff report
type HandoffReport struct {
    Shift          Shift         `json:"shift"`
    Engineer       Engineer      `json:"engineer"`
    GeneratedAt    time.Time     `json:"generated_at"`

    // Activity Summary
    Incidents      []IncidentSummary `json:"incidents"`
    AlertsReceived int               `json:"alerts_received"`
    AlertsAcked    int               `json:"alerts_acked"`

    // Ongoing Issues
    OngoingIssues  []OngoingIssue    `json:"ongoing_issues"`

    // System Health
    SystemHealth   map[string]HealthStatus `json:"system_health"`

    // Action Items
    ActionItems    []ActionItem      `json:"action_items"`

    // Notes
    Notes          string            `json:"notes"`
}

// IncidentSummary summarizes an incident
type IncidentSummary struct {
    ID           string           `json:"id"`
    Title        string           `json:"title"`
    Severity     IncidentSeverity `json:"severity"`
    State        IncidentState    `json:"state"`
    Duration     time.Duration    `json:"duration"`
    RootCause    string           `json:"root_cause,omitempty"`
}

// OngoingIssue represents an unresolved issue
type OngoingIssue struct {
    ID          string    `json:"id"`
    Title       string    `json:"title"`
    Severity    string    `json:"severity"`
    Status      string    `json:"status"`
    NextAction  string    `json:"next_action"`
}

// HealthStatus represents system health
type HealthStatus struct {
    Status      string  `json:"status"`
    LastChecked time.Time `json:"last_checked"`
    Details     string  `json:"details,omitempty"`
}

// ActionItem represents a follow-up action
type ActionItem struct {
    ID          string    `json:"id"`
    Description string    `json:"description"`
    Assignee    string    `json:"assignee"`
    DueDate     time.Time `json:"due_date"`
    Priority    string    `json:"priority"`
}

// HandoffReportGenerator generates handoff reports
type HandoffReportGenerator struct {
    incidentManager *IncidentManager
    scheduleManager *ScheduleManager
}

// NewHandoffReportGenerator creates a new report generator
func NewHandoffReportGenerator(im *IncidentManager, sm *ScheduleManager) *HandoffReportGenerator {
    return &HandoffReportGenerator{
        incidentManager: im,
        scheduleManager: sm,
    }
}

// GenerateReport generates a handoff report for a shift
func (g *HandoffReportGenerator) GenerateReport(ctx context.Context, shiftID string) (*HandoffReport, error) {
    // Get shift details
    shift, engineer, err := g.getShiftDetails(shiftID)
    if err != nil {
        return nil, err
    }

    report := &HandoffReport{
        Shift:       *shift,
        Engineer:    *engineer,
        GeneratedAt: time.Now().UTC(),
        Incidents:   make([]IncidentSummary, 0),
        OngoingIssues: make([]OngoingIssue, 0),
        SystemHealth: make(map[string]HealthStatus),
        ActionItems: make([]ActionItem, 0),
    }

    // Query incidents during shift
    incidents := g.getIncidentsDuringShift(shift)
    for _, inc := range incidents {
        duration := time.Duration(0)
        if inc.ResolvedAt != nil {
            duration = inc.ResolvedAt.Sub(inc.CreatedAt)
        }

        summary := IncidentSummary{
            ID:       inc.ID,
            Title:    inc.Title,
            Severity: inc.Severity,
            State:    inc.State,
            Duration: duration,
        }

        if inc.State == IncidentStateResolved {
            summary.RootCause = "See incident details" // Would be populated from postmortem
        }

        report.Incidents = append(report.Incidents, summary)
    }

    // Get ongoing issues (incidents not resolved)
    openIncidents := g.incidentManager.GetOpenIncidents()
    for _, inc := range openIncidents {
        if inc.CreatedAt.Before(shift.EndTime) {
            issue := OngoingIssue{
                ID:         inc.ID,
                Title:      inc.Title,
                Severity:   string(inc.Severity),
                Status:     string(inc.State),
                NextAction: "Continue investigation",
            }
            report.OngoingIssues = append(report.OngoingIssues, issue)
        }
    }

    return report, nil
}

// GenerateHTML generates HTML report
func (g *HandoffReportGenerator) GenerateHTML(report *HandoffReport) (string, error) {
    const tmpl = `
<!DOCTYPE html>
<html>
<head>
    <title>Handoff Report - {{ .Shift.ID }}</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; }
        h1 { color: #333; }
        h2 { color: #666; border-bottom: 2px solid #ddd; padding-bottom: 5px; }
        .summary { background: #f5f5f5; padding: 15px; border-radius: 5px; }
        .incident { margin: 10px 0; padding: 10px; border-left: 4px solid #999; }
        .incident.critical { border-color: #d32f2f; }
        .incident.high { border-color: #f57c00; }
        .incident.medium { border-color: #fbc02d; }
        .ongoing { background: #fff3cd; padding: 10px; margin: 5px 0; border-radius: 3px; }
        table { width: 100%; border-collapse: collapse; margin: 10px 0; }
        th, td { padding: 8px; text-align: left; border-bottom: 1px solid #ddd; }
        th { background: #f5f5f5; }
    </style>
</head>
<body>
    <h1>🔄 On-Call Handoff Report</h1>

    <div class="summary">
        <h2>Shift Summary</h2>
        <p><strong>Engineer:</strong> {{ .Engineer.Name }} ({{ .Engineer.Email }})</p>
        <p><strong>Shift Period:</strong> {{ .Shift.StartTime.Format "2006-01-02 15:04" }} - {{ .Shift.EndTime.Format "2006-01-02 15:04" }}</p>
        <p><strong>Total Incidents:</strong> {{ len .Incidents }}</p>
        <p><strong>Ongoing Issues:</strong> {{ len .OngoingIssues }}</p>
    </div>

    <h2>📋 Incidents During Shift</h2>
    {{ range .Incidents }}
    <div class="incident {{ .Severity }}">
        <strong>{{ .ID }}:</strong> {{ .Title }}<br>
        <small>Severity: {{ .Severity }} | State: {{ .State }} | Duration: {{ .Duration }}</small>
    </div>
    {{ end }}

    <h2>⚠️ Ongoing Issues</h2>
    {{ range .OngoingIssues }}
    <div class="ongoing">
        <strong>{{ .ID }}:</strong> {{ .Title }} ({{ .Severity }})<br>
        <small>Status: {{ .Status }} | Next Action: {{ .NextAction }}</small>
    </div>
    {{ end }}

    <h2>📝 Notes</h2>
    <p>{{ .Notes }}</p>

    <hr>
    <p><small>Generated at {{ .GeneratedAt.Format "2006-01-02 15:04:05 MST" }}</small></p>
</body>
</html>
`

    t := template.Must(template.New("handoff").Parse(tmpl))
    var buf bytes.Buffer
    if err := t.Execute(&buf, report); err != nil {
        return "", err
    }

    return buf.String(), nil
}

// getShiftDetails retrieves shift and engineer details
func (g *HandoffReportGenerator) getShiftDetails(shiftID string) (*Shift, *Engineer, error) {
    // Implementation would retrieve from schedule manager
    return nil, nil, fmt.Errorf("not implemented")
}

// getIncidentsDuringShift retrieves incidents during the shift
func (g *HandoffReportGenerator) getIncidentsDuringShift(shift *Shift) []*Incident {
    // Query all incidents and filter by time
    all := g.incidentManager.GetOpenIncidents()
    var filtered []*Incident

    // Would also include resolved incidents from store
    _ = all

    return filtered
}
```

---

## 3. Production-Ready Configurations

### 3.1 PagerDuty Integration

```yaml
# pagerduty-service.yaml
apiVersion: pagerduty.com/v1
kind: Service
metadata:
  name: production-api
  team: platform-engineering
spec:
  description: "Production API Service - Customer Facing"
  escalation_policy: platform-critical

  # Alert grouping
  alert_grouping:
    enabled: true
    type: intelligent  # or time, content
    timeout: 300  # seconds

  # Auto-resolve
  auto_resolve_timeout: 14400  # 4 hours

  # Acknowledge timeout
  acknowledge_timeout: 1800  # 30 minutes

  # Incident urgency rules
  incident_urgency_rules:
    type: use_support_hours
    during_support_hours:
      type: constant
      urgency: high
    outside_support_hours:
      type: constant
      urgency: low

  # Support hours
  support_hours:
    type: fixed_time_per_day
    start_time: "09:00:00"
    end_time: "17:00:00"
    time_zone: "America/New_York"
    days_of_week:
      - Monday
      - Tuesday
      - Wednesday
      - Thursday
      - Friday

  # Integrations
  integrations:
    - type: events_api_v2
      name: prometheus-alerts
    - type: events_api_v2
      name: datadog-alerts

---
# Escalation Policy
apiVersion: pagerduty.com/v1
kind: EscalationPolicy
metadata:
  name: platform-critical
spec:
  description: "Critical escalation for platform issues"

  escalation_rules:
    # Level 1: Primary on-call
    - escalation_delay_in_minutes: 5
      targets:
        - type: schedule_reference
          id: platform-primary

    # Level 2: Secondary on-call + Manager
    - escalation_delay_in_minutes: 10
      targets:
        - type: schedule_reference
          id: platform-secondary
        - type: user_reference
          id: platform-manager

    # Level 3: Director + Additional engineers
    - escalation_delay_in_minutes: 15
      targets:
        - type: user_reference
          id: platform-director
        - type: schedule_reference
          id: platform-emergency

---
# Schedule Configuration
apiVersion: pagerduty.com/v1
kind: Schedule
metadata:
  name: platform-primary
spec:
  description: "Primary on-call rotation for platform team"
  time_zone: "America/New_York"

  schedule_layers:
    - name: "Weekday Rotation"
      start: "2024-01-01T00:00:00-05:00"
      rotation_virtual_start: "2024-01-01T00:00:00-05:00"
      rotation_turn_length_in_seconds: 604800  # 1 week
      users:
        - id: user1@example.com
        - id: user2@example.com
        - id: user3@example.com
        - id: user4@example.com

      # Restrictions (e.g., business hours only for some)
      restrictions:
        - type: daily_restriction
          start_time_of_day: "00:00:00"
          duration_seconds: 86400  # All day

---
# Event Orchestration
apiVersion: pagerduty.com/v1
kind: EventOrchestration
metadata:
  name: platform-routing
spec:
  description: "Route events based on severity and service"

  sets:
    - id: start
      rules:
        # Critical alerts - Page immediately
        - label: "Critical Infrastructure Alert"
          conditions:
            - expression: "event.severity == 'critical'"
          actions:
            - route_to: platform-critical
            - priority: P1
            - urgency: high

        # Security alerts - Special handling
        - label: "Security Alert"
          conditions:
            - expression: "event.tags contains 'security'"
          actions:
            - route_to: security-team
            - priority: P1
            - annotate: "Security incident - follow security playbook"

        # Non-prod alerts - Lower priority
        - label: "Non-Production"
          conditions:
            - expression: "event.tags contains 'env:staging' or event.tags contains 'env:dev'"
          actions:
            - suppress: true  # Don't page for non-prod
            - automation_action:
                name: "Auto-restart staging service"

        # Default - Route to primary
        - label: "Default"
          actions:
            - route_to: platform-primary
            - priority: P3
```

---

## 4. Security Considerations

### 4.1 On-Call Security Matrix

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    On-Call Security Considerations                           │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  CATEGORY           │  RISK                │  MITIGATION                    │
├─────────────────────┼──────────────────────┼────────────────────────────────│
│  Access Control     │                      │                                │
│  ├─ On-call access  │  Unauthorized access │  • Just-in-time access         │
│  │  to production  │                      │  • Time-limited credentials    │
│  ├─ Shared accounts │  Account compromise  │  • Individual accounts only    │
│  └─ Offboarding     │  Orphaned access     │  • Auto-revoke on rotation     │
│  ───────────────────┼──────────────────────┼────────────────────────────────│
│  Communication      │                      │                                │
│  ├─ Unencrypted     │  Data interception   │  • End-to-end encryption       │
│  │  channels       │                      │  • Secure incident rooms       │
│  ├─ Information     │  Data leak           │  • Need-to-know sharing        │
│  │  oversharing    │                      │  • Sanitized alerts            │
│  └─ Social          │  Impersonation       │  • Identity verification       │
│     engineering     │                      │  • Multi-factor auth           │
│  ───────────────────┼──────────────────────┼────────────────────────────────│
│  Incident Response  │                      │                                │
│  ├─ Insider threat  │  Malicious insider   │  • Dual authorization          │
│  │                   │                      │  • Audit logging               │
│  ├─ Credential      │  Privilege escalation│  • Break-glass procedures      │
│  │  compromise     │                      │  • Session recording           │
│  └─ Data exposure   │  Sensitive data leak │  • Data masking                │
│                     │                      │  • Secure log storage          │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 4.2 Break-Glass Access Pattern

```go
package oncall

import (
    "context"
    "crypto/rand"
    "encoding/base64"
    "fmt"
    "time"
)

// BreakGlassRequest represents a break-glass access request
type BreakGlassRequest struct {
    ID          string    `json:"id"`
    EngineerID  string    `json:"engineer_id"`
    Reason      string    `json:"reason"`
    IncidentID  string    `json:"incident_id,omitempty"`

    RequestedAt time.Time `json:"requested_at"`
    ExpiresAt   time.Time `json:"expires_at"`

    // Approval
    ApprovedBy  *string   `json:"approved_by,omitempty"`
    ApprovedAt  *time.Time `json:"approved_at,omitempty"`

    // Access details
    Credentials *BreakGlassCredentials `json:"credentials,omitempty"`
    Status      string    `json:"status"` // pending, approved, denied, expired, revoked
}

// BreakGlassCredentials contains temporary access credentials
type BreakGlassCredentials struct {
    Username    string    `json:"username"`
    Token       string    `json:"token"`
    ExpiresAt   time.Time `json:"expires_at"`
    Permissions []string  `json:"permissions"`
}

// BreakGlassManager manages emergency access
type BreakGlassManager struct {
    requests    map[string]*BreakGlassRequest
    approvers   []string  // List of approved approver IDs
    maxDuration time.Duration

    credentialProvider CredentialProvider
    auditLogger        AuditLogger
}

// CredentialProvider issues temporary credentials
type CredentialProvider interface {
    IssueCredentials(ctx context.Context, engineerID string, permissions []string, duration time.Duration) (*BreakGlassCredentials, error)
    RevokeCredentials(ctx context.Context, username string) error
}

// AuditLogger logs access events
type AuditLogger interface {
    LogAccess(ctx context.Context, event AccessEvent) error
}

// AccessEvent represents an access audit event
type AccessEvent struct {
    Timestamp   time.Time
    EngineerID  string
    Action      string
    Resource    string
    Success     bool
    Details     map[string]string
}

// NewBreakGlassManager creates a break-glass manager
func NewBreakGlassManager(approvers []string, maxDuration time.Duration, provider CredentialProvider, logger AuditLogger) *BreakGlassManager {
    if maxDuration == 0 {
        maxDuration = 4 * time.Hour
    }

    return &BreakGlassManager{
        requests:           make(map[string]*BreakGlassRequest),
        approvers:          approvers,
        maxDuration:        maxDuration,
        credentialProvider: provider,
        auditLogger:        logger,
    }
}

// RequestAccess requests break-glass access
func (m *BreakGlassManager) RequestAccess(ctx context.Context, engineerID, reason, incidentID string, duration time.Duration) (*BreakGlassRequest, error) {
    if duration > m.maxDuration {
        duration = m.maxDuration
    }

    request := &BreakGlassRequest{
        ID:          generateID(),
        EngineerID:  engineerID,
        Reason:      reason,
        IncidentID:  incidentID,
        RequestedAt: time.Now().UTC(),
        ExpiresAt:   time.Now().UTC().Add(duration),
        Status:      "pending",
    }

    m.requests[request.ID] = request

    // Notify approvers
    m.notifyApprovers(request)

    return request, nil
}

// ApproveRequest approves a break-glass request
func (m *BreakGlassManager) ApproveRequest(ctx context.Context, requestID, approverID string) error {
    request, exists := m.requests[requestID]
    if !exists {
        return fmt.Errorf("request not found")
    }

    // Verify approver
    isApprover := false
    for _, a := range m.approvers {
        if a == approverID {
            isApprover = true
            break
        }
    }

    if !isApprover {
        return fmt.Errorf("not an authorized approver")
    }

    // Self-approval not allowed
    if request.EngineerID == approverID {
        return fmt.Errorf("cannot approve own request")
    }

    now := time.Now().UTC()
    request.ApprovedBy = &approverID
    request.ApprovedAt = &now
    request.Status = "approved"

    // Issue credentials
    credentials, err := m.credentialProvider.IssueCredentials(ctx, request.EngineerID, []string{"break-glass-emergency"}, time.Until(request.ExpiresAt))
    if err != nil {
        return fmt.Errorf("failed to issue credentials: %w", err)
    }

    request.Credentials = credentials

    // Log approval
    m.auditLogger.LogAccess(ctx, AccessEvent{
        Timestamp:  now,
        EngineerID: request.EngineerID,
        Action:     "break_glass_approved",
        Resource:   requestID,
        Success:    true,
        Details: map[string]string{
            "approver": approverID,
            "reason":   request.Reason,
        },
    })

    return nil
}

// RevokeAccess revokes break-glass access
func (m *BreakGlassManager) RevokeAccess(ctx context.Context, requestID string) error {
    request, exists := m.requests[requestID]
    if !exists {
        return fmt.Errorf("request not found")
    }

    if request.Credentials != nil {
        if err := m.credentialProvider.RevokeCredentials(ctx, request.Credentials.Username); err != nil {
            return err
        }
    }

    request.Status = "revoked"

    // Log revocation
    m.auditLogger.LogAccess(ctx, AccessEvent{
        Timestamp:  time.Now().UTC(),
        EngineerID: request.EngineerID,
        Action:     "break_glass_revoked",
        Resource:   requestID,
        Success:    true,
    })

    return nil
}

// CleanupExpired removes expired requests and revokes credentials
func (m *BreakGlassManager) CleanupExpired(ctx context.Context) {
    now := time.Now().UTC()

    for _, request := range m.requests {
        if now.After(request.ExpiresAt) && request.Status == "approved" {
            m.RevokeAccess(ctx, request.ID)
            request.Status = "expired"
        }
    }
}

// notifyApprovers sends approval requests to approvers
func (m *BreakGlassManager) notifyApprovers(request *BreakGlassRequest) {
    // Implementation would send notifications to approvers
}

// generateSecureToken generates a secure random token
func generateSecureToken() (string, error) {
    b := make([]byte, 32)
    if _, err := rand.Read(b); err != nil {
        return "", err
    }
    return base64.URLEncoding.EncodeToString(b), nil
}
```

---

## 5. Compliance Requirements

### 5.1 Compliance Mapping

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    On-Call Compliance Requirements                           │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  SOC 2 TYPE II                                                              │
│  ├─ CC6.1: Logical access - Access during on-call is logged and reviewed   │
│  ├─ CC6.2: Prior to access - Emergency access requires approval            │
│  ├─ CC6.3: Access removal - Automatic revocation after shift               │
│  ├─ CC7.2: System monitoring - 24/7 coverage for critical systems          │
│  ├─ CC7.3: Incident detection - Defined escalation procedures              │
│  └─ CC7.4: Incident response - Documented response procedures              │
│                                                                             │
│  ISO 27001                                                                  │
│  ├─ A.6.1.3: Contact with authorities - On-call for regulatory contact     │
│  ├─ A.12.4.1: Event logging - All on-call activities logged                │
│  ├─ A.16.1: Incident management - Defined roles and procedures             │
│  └─ A.16.2: Incident reporting - Reporting to relevant authorities         │
│                                                                             │
│  HIPAA                                                                      │
│  ├─ §164.312(a)(2)(ii): Emergency access - Break-glass procedures          │
│  ├─ §164.312(b): Audit controls - Access to PHI during on-call             │
│  └─ §164.312(c)(1): Integrity - Data integrity during incident response    │
│                                                                             │
│  PCI DSS                                                                    │
│  ├─ Req 7.1.1: Access restrictions - Need-to-know during incidents         │
│  ├─ Req 10.1: Audit trails - All on-call access to CDE logged              │
│  ├─ Req 10.2.5: Use of IDs - Individual accounts for on-call               │
│  └─ Req 12.10.4: Incident response - 24/7 incident response capability     │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 6. Decision Matrices

### 6.1 Escalation Decision Matrix

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      Escalation Decision Matrix                              │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  Scenario                              │  Action          │  Timeline       │
├────────────────────────────────────────┼──────────────────┼─────────────────│
│  Alert received, no acknowledgment     │  Page on-call    │  +5 minutes     │
│                                                                             │
│  On-call acknowledged, no progress     │  Escalate to     │  +15 minutes    │
│  after initial triage                  │  secondary       │                 │
│                                                                             │
│  Multiple related alerts firing        │  Declare         │  Immediately    │
│                                        │  incident        │                 │
│                                                                             │
│  Customer-facing impact confirmed      │  Page manager    │  Immediately    │
│  (revenue affecting)                   │  + notify comms  │                 │
│                                                                             │
│  Security-related alert                │  Page security   │  Immediately    │
│                                        │  team + manager  │                 │
│                                                                             │
│  Data loss suspected                   │  Page legal +    │  Immediately    │
│                                        │  compliance      │                 │
│                                                                             │
│  Vendor/dependency down                │  Page vendor     │  After 10 min   │
│                                        │  liaison         │  internal work  │
│                                                                             │
│  On-call cannot be reached             │  Skip to next    │  After 5 min    │
│                                        │  escalation      │                 │
│                                                                             │
│  Weekend/after-hours critical          │  Page director   │  After 30 min   │
│                                        │                  │                 │
│                                                                             │
│  Recovery requires major change        │  Require CAB     │  Before change  │
│  (outside runbook)                     │  approval        │                 │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 6.2 Response Action Matrix

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                       Response Action Decision Matrix                        │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  Alert Type              │  Initial Action        │  Next Steps             │
├──────────────────────────┼────────────────────────┼─────────────────────────│
│  High error rate         │  Check recent          │  If code change:        │
│  (> 5%)                  │  deployments           │  Consider rollback      │
│                          │                        │  If dependency:         │
│                          │                        │  Failover or degrade    │
│  ────────────────────────┼────────────────────────┼─────────────────────────│
│  High latency            │  Check resource        │  If capacity:           │
│  (P99 > 2s)              │  utilization           │  Scale up               │
│                          │                        │  If dependency:         │
│                          │                        │  Circuit break          │
│  ────────────────────────┼────────────────────────┼─────────────────────────│
│  Resource exhaustion     │  Check trending        │  If predictable:        │
│                          │  alerts                │  Scale preemptively     │
│                          │                        │  If sudden spike:       │
│                          │                        │  Investigate cause      │
│  ────────────────────────┼────────────────────────┼─────────────────────────│
│  Database issue          │  Check connection      │  If primary down:       │
│                          │  pools, slow queries   │  Failover to replica    │
│                          │                        │  If lock/contention:    │
│                          │                        │  Kill blocking queries  │
│  ────────────────────────┼────────────────────────┼─────────────────────────│
│  Security alert          │  Follow security       │  If breach confirmed:   │
│                          │  runbook               │  Initiate IR plan       │
│                          │                        │  If false positive:     │
│                          │                        │  Tune alert rules       │
│  ────────────────────────┼────────────────────────┼─────────────────────────│
│  Network issue           │  Check connectivity    │  If internal:           │
│                          │  tests                 │  Check routing/firewall  │
│                          │                        │  If external:           │
│                          │                        │  Contact provider       │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 6.3 Shift Override Decision Matrix

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      Shift Override Decision Matrix                          │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  Reason for Override         │  Approver Required  │  Documentation        │
├──────────────────────────────┼─────────────────────┼───────────────────────│
│  Medical emergency           │  Manager            │  Reason documented    │
│                              │                     │  (private)            │
│  ────────────────────────────┼─────────────────────┼───────────────────────│
│  Vacation/travel conflict    │  Team lead          │  Advance notice       │
│  (planned)                   │                     │  (> 2 weeks)          │
│  ────────────────────────────┼─────────────────────┼───────────────────────│
│  Family emergency            │  Manager            │  Reason documented    │
│                              │                     │  (private)            │
│  ────────────────────────────┼─────────────────────┼───────────────────────│
│  Professional conflict       │  Team lead          │  Meeting details      │
│  (conference, interview)     │                     │                       │
│  ────────────────────────────┼─────────────────────┼───────────────────────│
│  Shift swap (peer to peer)   │  Self-service       │  Both parties agree   │
│                              │                     │  (system tracked)     │
│  ────────────────────────────┼─────────────────────┼───────────────────────│
│  Burnout/fatigue             │  Manager + HR       │  Wellness check       │
│                              │                     │  triggered            │
│  ────────────────────────────┼─────────────────────┼───────────────────────│
│  Skill mismatch              │  Team lead          │  Training scheduled   │
│  (unable to handle service)  │                     │                       │
│                                                                             │
│  Override Limits per Quarter:                                              │
│  • Self-initiated: 2 (requires manager awareness)                          │
│  • Manager-approved: No limit (with valid reason)                          │
│  • Unexcused no-show: Performance review trigger                           │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 7. Best Practices Summary

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      On-Call Best Practices Summary                          │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  SCHEDULING                                                                 │
│  ✓ Rotate at least every week (daily for high-intensity)                    │
│  ✓ Include secondary/backup on-call                                         │
│  ✓ Account for timezones in global teams                                    │
│  ✓ Honor blackout dates and preferences                                     │
│  ✓ Plan rotations at least 1 month in advance                               │
│  ✓ Ensure coverage during holidays                                          │
│  ✓ Shadow program for new team members                                      │
│                                                                             │
│  PREPAREDNESS                                                               │
│  ✓ Laptop with VPN access always available                                  │
│  ✓ Secondary internet (mobile hotspot)                                      │
│  ✓ Access credentials tested before shift                                   │
│  ✓ Runbooks bookmarked and accessible offline                               │
│  ✓ Team contact list up to date                                             │
│  ✓ Test alert delivery before first shift                                   │
│                                                                             │
│  RESPONSE                                                                   │
│  ✓ Acknowledge alerts within 5 minutes                                      │
│  ✓ Update status page if customer-impacting                                 │
│  ✓ Communicate early and often in incident channel                          │
│  ✓ Escalate early if stuck (don't struggle alone)                           │
│  ✓ Document actions taken during incident                                   │
│  ✓ Hand off active incidents properly                                       │
│                                                                             │
│  POST-INCIDENT                                                              │
│  ✓ Complete post-mortem within 1 week                                       │
│  ✓ Update runbooks with lessons learned                                     │
│  ✓ File tickets for follow-up action items                                  │
│  ✓ Review alert effectiveness                                               │
│  ✓ Share knowledge with team                                                │
│                                                                             │
│  WELLNESS                                                                   │
│  ✓ Compensate for off-hours work (time off or pay)                          │
│  ✓ Limit consecutive on-call shifts                                         │
│  ✓ Monitor for alert fatigue                                                │
│  ✓ Provide wellness resources                                               │
│  ✓ No-blame culture for incidents                                           │
│  ✓ Regular team retrospectives                                              │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## References

1. Google SRE Book - Being On-Call
2. PagerDuty Incident Response Guide
3. Atlassian Incident Management Handbook
4. The Phoenix Project - Gene Kim
5. Site Reliability Workbook - On-Call Practices

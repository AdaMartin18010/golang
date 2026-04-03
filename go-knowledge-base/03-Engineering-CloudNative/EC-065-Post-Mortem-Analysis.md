# Post-Mortem Analysis

> **分类**: 工程与云原生
> **标签**: #postmortem #blameless #sre #learning #continuous-improvement
> **参考**: Google SRE, Etsy Blameless Post-Mortems, Etsy Morgue

---

## 1. Formal Definition

### 1.1 What is a Post-Mortem?

A post-mortem (or post-incident review) is a structured process for documenting the root causes, timeline, impact, and lessons learned from an incident. The primary goal is organizational learning and prevention of recurrence, conducted in a blameless manner that focuses on systemic issues rather than individual fault.

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      Post-Mortem Process Flow                                │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│   INCIDENT RESOLVED                                                         │
│        │                                                                    │
│        ▼                                                                    │
│   ┌─────────────────┐                                                       │
│   │  Within 24-72h  │                                                       │
│   │  Schedule       │                                                       │
│   │  Post-Mortem    │                                                       │
│   │  Meeting        │                                                       │
│   └────────┬────────┘                                                       │
│            │                                                                │
│            ▼                                                                │
│   ┌─────────────────────────────────────────────────────────────────────┐   │
│   │                     PREPARATION PHASE                                │   │
│   │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  │   │
│   │  │  Timeline   │  │   Impact    │  │   Metrics   │  │   Evidence  │  │   │
│   │  │  Collection │  │  Assessment │  │  Collection │  │  Gathering  │  │   │
│   │  └─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘  │   │
│   └─────────────────────────────────────────────────────────────────────┘   │
│            │                                                                │
│            ▼                                                                │
│   ┌─────────────────────────────────────────────────────────────────────┐   │
│   │                      MEETING PHASE                                   │   │
│   │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  │   │
│   │  │  Timeline   │  │  Root Cause │  │  Contributing│  │  Lessons    │  │   │
│   │  │   Review    │  │  Analysis   │  │   Factors   │  │   Learned   │  │   │
│   │  └─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘  │   │
│   └─────────────────────────────────────────────────────────────────────┘   │
│            │                                                                │
│            ▼                                                                │
│   ┌─────────────────────────────────────────────────────────────────────┐   │
│   │                    DOCUMENTATION PHASE                               │   │
│   │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  │   │
│   │  │  Write      │  │  Review &   │  │  Publish    │  │  Share      │  │   │
│   │  │  Document   │  │  Approve    │  │  Internally │  │  Externally │  │   │
│   │  └─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘  │   │
│   └─────────────────────────────────────────────────────────────────────┘   │
│            │                                                                │
│            ▼                                                                │
│   ┌─────────────────────────────────────────────────────────────────────┐   │
│   │                     FOLLOW-UP PHASE                                  │   │
│   │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  │   │
│   │  │  Create     │  │  Assign     │  │  Track      │  │  Verify     │  │   │
│   │  │  Action     │  │  Owners     │  │  Progress   │  │  Completion │  │   │
│   │  │  Items      │  │             │  │             │  │             │  │   │
│   │  └─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘  │   │
│   └─────────────────────────────────────────────────────────────────────┘   │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 1.2 Blameless Culture Principles

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      Blameless Culture Principles                            │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│   ❌ TRADITIONAL APPROACH              ✅ BLAMELESS APPROACH                │
│   ─────────────────────────            ─────────────────────                │
│                                                                             │
│   "Who caused this issue?"             "What about our system allowed       │
│                                        this issue to occur?"                │
│                                                                             │
│   "Why did you make that change?"      "What safeguards were missing that   │
│                                        would have caught this?"             │
│                                                                             │
│   "You should have known better"       "How can we improve documentation    │
│                                        to prevent similar issues?"          │
│                                                                             │
│   Individual punishment                System improvement                   │
│   Fear and hiding                      Transparency and learning            │
│   Single point of failure              Resilient systems                    │
│                                                                             │
│   ────────────────────────────────────────────────────────────────────────  │
│                                                                             │
│   BLAMELESS DOES NOT MEAN:                                                  │
│   • No accountability                                                       │
│   • No consequences for negligence                                          │
│   • Ignoring willful misconduct                                             │
│   • Avoiding root cause analysis                                            │
│                                                                             │
│   BLAMELESS DOES MEAN:                                                      │
│   • Focusing on system factors                                              │
│   • Understanding human error context                                       │
│   • Multiple contributing factors                                           │
│   • Process and tooling improvements                                        │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 2. Implementation Patterns in Go

### 2.1 Post-Mortem Document Model

```go
package postmortem

import (
    "context"
    "fmt"
    "html/template"
    "strings"
    "time"
)

// PostMortem represents a post-mortem document
type PostMortem struct {
    ID          string    `json:"id"`
    IncidentID  string    `json:"incident_id"`
    Title       string    `json:"title"`
    Status      string    `json:"status"` // draft, review, published, closed

    // Metadata
    CreatedAt   time.Time `json:"created_at"`
    CreatedBy   string    `json:"created_by"`
    UpdatedAt   time.Time `json:"updated_at"`
    PublishedAt *time.Time `json:"published_at,omitempty"`

    // Authors and Reviewers
    Authors     []string  `json:"authors"`
    Reviewers   []string  `json:"reviewers"`
    ApprovedBy  []string  `json:"approved_by"`

    // Incident Summary
    Summary     IncidentSummary   `json:"summary"`

    // Detailed Sections
    Timeline    []TimelineEntry   `json:"timeline"`
    RootCause   RootCauseAnalysis `json:"root_cause"`
    Impact      ImpactAssessment  `json:"impact"`
    Lessons     []LessonLearned   `json:"lessons"`
    ActionItems []ActionItem      `json:"action_items"`

    // Metrics
    Metrics     PostMortemMetrics `json:"metrics"`
}

// IncidentSummary provides high-level incident information
type IncidentSummary struct {
    Description     string     `json:"description"`
    Severity        string     `json:"severity"`
    DetectionMethod string     `json:"detection_method"`
    ResolvedBy      string     `json:"resolved_by"`
    ResolutionTime  time.Duration `json:"resolution_time"`
}

// TimelineEntry represents a single timeline event
type TimelineEntry struct {
    Timestamp   time.Time `json:"timestamp"`
    Description string    `json:"description"`
    Category    string    `json:"category"` // detection, response, mitigation, resolution
    Actor       string    `json:"actor"`
    Evidence    []string  `json:"evidence,omitempty"`
}

// RootCauseAnalysis contains the 5 Whys analysis
type RootCauseAnalysis struct {
    ProblemStatement string   `json:"problem_statement"`
    Whys             []Why    `json:"whys"`
    RootCause        string   `json:"root_cause"`
    ContributingFactors []string `json:"contributing_factors"`
}

// Why represents one level of the 5 Whys
type Why struct {
    Level       int    `json:"level"`
    Question    string `json:"question"`
    Answer      string `json:"answer"`
    Evidence    string `json:"evidence,omitempty"`
}

// ImpactAssessment details the incident impact
type ImpactAssessment struct {
    Duration        time.Duration `json:"duration"`
    ServicesAffected []string     `json:"services_affected"`
    UsersAffected   int           `json:"users_affected"`
    RegionsAffected []string      `json:"regions_affected"`
    FinancialImpact *FinancialImpact `json:"financial_impact,omitempty"`
    DataImpact      *DataImpact      `json:"data_impact,omitempty"`
    ReputationImpact string        `json:"reputation_impact,omitempty"`
}

// FinancialImpact represents monetary impact
type FinancialImpact struct {
    RevenueLost    float64 `json:"revenue_lost"`
    SLA Penalties float64 `json:"sla_penalties"`
    RecoveryCost   float64 `json:"recovery_cost"`
    Currency       string  `json:"currency"`
}

// DataImpact represents data-related impact
type DataImpact struct {
    RecordsAffected int    `json:"records_affected"`
    DataTypes       []string `json:"data_types"`
    ExposureDuration time.Duration `json:"exposure_duration,omitempty"`
}

// LessonLearned captures a lesson from the incident
type LessonLearned struct {
    ID          string `json:"id"`
    Category    string `json:"category"` // what_went_well, what_went_wrong, where_we_got_lucky
    Description string `json:"description"`
    Context     string `json:"context"`
}

// ActionItem represents a follow-up action
type ActionItem struct {
    ID          string    `json:"id"`
    Description string    `json:"description"`
    Owner       string    `json:"owner"`
    Priority    string    `json:"priority"` // critical, high, medium, low
    DueDate     time.Time `json:"due_date"`
    Status      string    `json:"status"` // open, in_progress, closed
    RelatedLesson string  `json:"related_lesson,omitempty"`
}

// PostMortemMetrics tracks post-mortem metrics
type PostMortemMetrics struct {
    TimeToPostMortem   time.Duration `json:"time_to_post_mortem"`
    MeetingDuration    time.Duration `json:"meeting_duration"`
    Attendees          int           `json:"attendees"`
    ActionItemsCreated int           `json:"action_items_created"`
}

// PostMortemGenerator generates post-mortem documents
type PostMortemGenerator struct {
    template *template.Template
    store    PostMortemStore
}

// PostMortemStore persists post-mortem data
type PostMortemStore interface {
    Save(ctx context.Context, pm *PostMortem) error
    Load(ctx context.Context, id string) (*PostMortem, error)
    List(ctx context.Context, filter PostMortemFilter) ([]*PostMortem, error)
}

// PostMortemFilter filters post-mortems
type PostMortemFilter struct {
    Status     string
    Severity   string
    StartTime  *time.Time
    EndTime    *time.Time
    Service    string
}

// NewPostMortemGenerator creates a new generator
func NewPostMortemGenerator(store PostMortemStore) (*PostMortemGenerator, error) {
    tmpl, err := template.New("postmortem").Funcs(template.FuncMap{
        "formatDuration": formatDuration,
        "formatTime":     formatTime,
        "toUpper":        strings.ToUpper,
    }).Parse(postMortemTemplate)

    if err != nil {
        return nil, err
    }

    return &PostMortemGenerator{
        template: tmpl,
        store:    store,
    }, nil
}

// CreatePostMortem creates a new post-mortem from incident data
func (g *PostMortemGenerator) CreatePostMortem(incidentID, title, createdBy string) (*PostMortem, error) {
    pm := &PostMortem{
        ID:         generateID(),
        IncidentID: incidentID,
        Title:      title,
        Status:     "draft",
        CreatedAt:  time.Now().UTC(),
        CreatedBy:  createdBy,
        UpdatedAt:  time.Now().UTC(),
        Authors:    []string{createdBy},
        Timeline:   make([]TimelineEntry, 0),
        Lessons:    make([]LessonLearned, 0),
        ActionItems: make([]ActionItem, 0),
    }

    return pm, nil
}

// AddTimelineEntry adds a timeline entry
func (pm *PostMortem) AddTimelineEntry(timestamp time.Time, description, category, actor string) {
    entry := TimelineEntry{
        Timestamp:   timestamp,
        Description: description,
        Category:    category,
        Actor:       actor,
    }
    pm.Timeline = append(pm.Timeline, entry)
    pm.UpdatedAt = time.Now().UTC()
}

// AddLesson adds a lesson learned
func (pm *PostMortem) AddLesson(category, description, context string) {
    lesson := LessonLearned{
        ID:          fmt.Sprintf("LL-%d", len(pm.Lessons)+1),
        Category:    category,
        Description: description,
        Context:     context,
    }
    pm.Lessons = append(pm.Lessons, lesson)
    pm.UpdatedAt = time.Now().UTC()
}

// AddActionItem adds an action item
func (pm *PostMortem) AddActionItem(description, owner, priority string, dueDate time.Time, relatedLesson string) {
    item := ActionItem{
        ID:            fmt.Sprintf("AI-%d", len(pm.ActionItems)+1),
        Description:   description,
        Owner:         owner,
        Priority:      priority,
        DueDate:       dueDate,
        Status:        "open",
        RelatedLesson: relatedLesson,
    }
    pm.ActionItems = append(pm.ActionItems, item)
    pm.UpdatedAt = time.Now().UTC()
}

// GenerateMarkdown generates Markdown document
func (g *PostMortemGenerator) GenerateMarkdown(pm *PostMortem) (string, error) {
    var buf strings.Builder
    if err := g.template.Execute(&buf, pm); err != nil {
        return "", err
    }
    return buf.String(), nil
}

// Helper functions
func formatDuration(d time.Duration) string {
    if d < time.Minute {
        return fmt.Sprintf("%ds", int(d.Seconds()))
    }
    if d < time.Hour {
        return fmt.Sprintf("%dm", int(d.Minutes()))
    }
    return fmt.Sprintf("%dh", int(d.Hours()))
}

func formatTime(t time.Time) string {
    return t.Format("2006-01-02 15:04:05 MST")
}

func generateID() string {
    return fmt.Sprintf("pm-%d", time.Now().UnixNano())
}

const postMortemTemplate = `# Post-Mortem: {{ .Title }}

**Incident ID:** {{ .IncidentID }}
**Post-Mortem ID:** {{ .ID }}
**Status:** {{ .Status | toUpper }}
**Created:** {{ .CreatedAt | formatTime }}
**Authors:** {{ .Authors }}

---

## Executive Summary

{{ .Summary.Description }}

**Severity:** {{ .Summary.Severity }}
**Detection Method:** {{ .Summary.DetectionMethod }}
**Resolution Time:** {{ .Summary.ResolutionTime | formatDuration }}

---

## Timeline

| Time (UTC) | Event | Category | Actor |
|------------|-------|----------|-------|
{{ range .Timeline -}}
| {{ .Timestamp | formatTime }} | {{ .Description }} | {{ .Category }} | {{ .Actor }} |
{{ end }}

---

## Root Cause Analysis

### Problem Statement
{{ .RootCause.ProblemStatement }}

### 5 Whys Analysis
{{ range .RootCause.Whys }}
**{{ .Level }}. Why?** {{ .Question }}

**Answer:** {{ .Answer }}
{{ if .Evidence }}_Evidence: {{ .Evidence }}_{{ end }}

{{ end }}

### Root Cause
{{ .RootCause.RootCause }}

### Contributing Factors
{{ range .RootCause.ContributingFactors -}}
- {{ . }}
{{ end }}

---

## Impact Assessment

**Duration:** {{ .Impact.Duration | formatDuration }}
**Services Affected:** {{ .Impact.ServicesAffected }}
**Users Affected:** {{ .Impact.UsersAffected }}
**Regions Affected:** {{ .Impact.RegionsAffected }}

{{ if .Impact.FinancialImpact }}
### Financial Impact
- Revenue Lost: {{ .Impact.FinancialImpact.Currency }} {{ .Impact.FinancialImpact.RevenueLost }}
- SLA Penalties: {{ .Impact.FinancialImpact.Currency }} {{ .Impact.FinancialImpact.SLAPenalties }}
- Recovery Cost: {{ .Impact.FinancialImpact.Currency }} {{ .Impact.FinancialImpact.RecoveryCost }}
{{ end }}

---

## Lessons Learned

### What Went Well
{{ range .Lessons }}{{ if eq .Category "what_went_well" }}
#### {{ .ID }}: {{ .Description }}
{{ .Context }}
{{ end }}{{ end }}

### What Went Wrong
{{ range .Lessons }}{{ if eq .Category "what_went_wrong" }}
#### {{ .ID }}: {{ .Description }}
{{ .Context }}
{{ end }}{{ end }}

### Where We Got Lucky
{{ range .Lessons }}{{ if eq .Category "where_we_got_lucky" }}
#### {{ .ID }}: {{ .Description }}
{{ .Context }}
{{ end }}{{ end }}

---

## Action Items

| ID | Description | Owner | Priority | Due Date | Status |
|----|-------------|-------|----------|----------|--------|
{{ range .ActionItems -}}
| {{ .ID }} | {{ .Description }} | {{ .Owner }} | {{ .Priority }} | {{ .DueDate | formatTime }} | {{ .Status }} |
{{ end }}

---

## Supporting Documentation

- [Incident Timeline](link)
- [Monitoring Dashboards](link)
- [Communication Log](link)

---

*This post-mortem was conducted following our blameless post-mortem process.*
`
```

### 2.2 5 Whys Analysis Engine

```go
package postmortem

import (
    "fmt"
    "strings"
)

// FiveWhysEngine performs automated 5 Whys analysis
type FiveWhysEngine struct {
    maxDepth int
}

// NewFiveWhysEngine creates a new 5 Whys engine
func NewFiveWhysEngine() *FiveWhysEngine {
    return &FiveWhysEngine{
        maxDepth: 5,
    }
}

// AnalysisResult contains the analysis result
type AnalysisResult struct {
    Whys         []WhyAnalysis
    RootCause    string
    IsComplete   bool
    Confidence   float64
}

// WhyAnalysis represents a single why analysis
type WhyAnalysis struct {
    Level       int
    Observation string
    Question    string
    Answer      string
    Evidence    []string
    SuggestedNextQuestions []string
}

// Analyze performs 5 Whys analysis on an incident description
func (e *FiveWhysEngine) Analyze(problem string, timeline []TimelineEntry) *AnalysisResult {
    result := &AnalysisResult{
        Whys: make([]WhyAnalysis, 0),
    }

    currentProblem := problem

    for level := 1; level <= e.maxDepth; level++ {
        analysis := e.analyzeLevel(level, currentProblem, timeline)
        result.Whys = append(result.Whys, analysis)

        currentProblem = analysis.Answer

        // Check if we've reached a root cause
        if e.isRootCause(analysis.Answer) {
            result.RootCause = analysis.Answer
            result.IsComplete = true
            result.Confidence = e.calculateConfidence(result.Whys)
            break
        }
    }

    if result.RootCause == "" && len(result.Whys) > 0 {
        result.RootCause = result.Whys[len(result.Whys)-1].Answer
        result.Confidence = e.calculateConfidence(result.Whys)
    }

    return result
}

// analyzeLevel analyzes a single level
func (e *FiveWhysEngine) analyzeLevel(level int, problem string, timeline []TimelineEntry) WhyAnalysis {
    question := fmt.Sprintf("Why did %s?", strings.ToLower(problem))

    // Find relevant timeline entries
    evidence := e.extractEvidence(problem, timeline)

    // Generate answer based on evidence
    answer := e.inferCause(problem, evidence)

    // Suggest next questions
    nextQuestions := e.suggestNextQuestions(answer)

    return WhyAnalysis{
        Level:                  level,
        Observation:            problem,
        Question:               question,
        Answer:                 answer,
        Evidence:               evidence,
        SuggestedNextQuestions: nextQuestions,
    }
}

// extractEvidence extracts relevant evidence from timeline
func (e *FiveWhysEngine) extractEvidence(problem string, timeline []TimelineEntry) []string {
    evidence := make([]string, 0)

    keywords := e.extractKeywords(problem)

    for _, entry := range timeline {
        entryText := fmt.Sprintf("%s %s %s", entry.Description, entry.Category, entry.Actor)
        entryLower := strings.ToLower(entryText)

        for _, keyword := range keywords {
            if strings.Contains(entryLower, strings.ToLower(keyword)) {
                evidence = append(evidence, fmt.Sprintf("%s: %s", entry.Timestamp.Format("15:04:05"), entry.Description))
                break
            }
        }
    }

    return evidence
}

// extractKeywords extracts keywords from text
func (e *FiveWhysEngine) extractKeywords(text string) []string {
    // Simple keyword extraction - in production use NLP
    words := strings.Fields(text)
    keywords := make([]string, 0)

    stopWords := map[string]bool{
        "the": true, "a": true, "an": true, "is": true, "was": true,
        "were": true, "be": true, "been": true, "being": true,
        "have": true, "has": true, "had": true, "do": true,
        "does": true, "did": true, "will": true, "would": true,
    }

    for _, word := range words {
        word = strings.ToLower(strings.Trim(word, ",.!?;:"))
        if len(word) > 3 && !stopWords[word] {
            keywords = append(keywords, word)
        }
    }

    return keywords
}

// inferCause infers a cause from evidence
func (e *FiveWhysEngine) inferCause(problem string, evidence []string) string {
    // This is a simplified implementation
    // In production, this would use ML/NLP to analyze patterns

    if len(evidence) == 0 {
        return "[Root cause not identifiable from available data - investigation needed]"
    }

    // Pattern matching for common issues
    for _, ev := range evidence {
        evLower := strings.ToLower(ev)

        if strings.Contains(evLower, "oom") || strings.Contains(evLower, "out of memory") {
            return "the application ran out of memory"
        }
        if strings.Contains(evLower, "cpu") && strings.Contains(evLower, "high") {
            return "CPU utilization exceeded available resources"
        }
        if strings.Contains(evLower, "timeout") {
            return "the operation exceeded the configured timeout"
        }
        if strings.Contains(evLower, "deployment") || strings.Contains(evLower, "deploy") {
            return "a recent deployment introduced the issue"
        }
        if strings.Contains(evLower, "database") || strings.Contains(evLower, "db") {
            return "the database was unavailable or overloaded"
        }
    }

    return "[Systemic factor contributing to the issue - further analysis required]"
}

// suggestNextQuestions suggests follow-up questions
func (e *FiveWhysEngine) suggestNextQuestions(answer string) []string {
    suggestions := make([]string, 0)
    answerLower := strings.ToLower(answer)

    if strings.Contains(answerLower, "memory") {
        suggestions = append(suggestions, "Why did memory usage increase?")
        suggestions = append(suggestions, "Why wasn't memory usage monitored?")
    }
    if strings.Contains(answerLower, "timeout") {
        suggestions = append(suggestions, "Why did the operation take longer than expected?")
        suggestions = append(suggestions, "Why was the timeout set to this value?")
    }
    if strings.Contains(answerLower, "deployment") {
        suggestions = append(suggestions, "Why wasn't the issue caught in testing?")
        suggestions = append(suggestions, "Why didn't the canary detect the issue?")
    }

    return suggestions
}

// isRootCause checks if an answer indicates a root cause
func (e *FiveWhysEngine) isRootCause(answer string) bool {
    // Check for systemic factors
    rootCauseIndicators := []string{
        "process", "procedure", "automation", "monitoring",
        "testing", "documentation", "training", "tooling",
        "architecture", "design", "policy",
    }

    answerLower := strings.ToLower(answer)
    for _, indicator := range rootCauseIndicators {
        if strings.Contains(answerLower, indicator) {
            return true
        }
    }

    return false
}

// calculateConfidence calculates confidence score
func (e *FiveWhysEngine) calculateConfidence(whys []WhyAnalysis) float64 {
    if len(whys) == 0 {
        return 0
    }

    // Base confidence on depth and evidence
    depthScore := float64(len(whys)) / float64(e.maxDepth)

    evidenceScore := 0.0
    for _, why := range whys {
        if len(why.Evidence) > 0 {
            evidenceScore += 0.2
        }
    }
    evidenceScore = min(evidenceScore, 1.0)

    return (depthScore + evidenceScore) / 2
}

func min(a, b float64) float64 {
    if a < b {
        return a
    }
    return b
}
```

### 2.3 Action Item Tracker

```go
package postmortem

import (
    "context"
    "fmt"
    "sort"
    "sync"
    "time"
)

// ActionItemTracker tracks post-mortem action items
type ActionItemTracker struct {
    items map[string]*TrackedActionItem
    mu    sync.RWMutex

    store ActionItemStore
}

// TrackedActionItem extends ActionItem with tracking
type TrackedActionItem struct {
    ActionItem

    PostMortemID  string    `json:"post_mortem_id"`
    IncidentID    string    `json:"incident_id"`
    CreatedAt     time.Time `json:"created_at"`
    UpdatedAt     time.Time `json:"updated_at"`
    CompletedAt   *time.Time `json:"completed_at,omitempty"`

    Progress      int       `json:"progress"` // 0-100
    Notes         []ProgressNote `json:"notes"`
    Blockers      []string  `json:"blockers"`

    RemindersSent int       `json:"reminders_sent"`
    LastReminder  *time.Time `json:"last_reminder,omitempty"`
}

// ProgressNote tracks progress updates
type ProgressNote struct {
    Timestamp time.Time `json:"timestamp"`
    Author    string    `json:"author"`
    Note      string    `json:"note"`
    Progress  int       `json:"progress"`
}

// ActionItemStore persists action items
type ActionItemStore interface {
    Save(ctx context.Context, item *TrackedActionItem) error
    Load(ctx context.Context, id string) (*TrackedActionItem, error)
    List(ctx context.Context, filter ActionItemFilter) ([]*TrackedActionItem, error)
    ListOverdue(ctx context.Context) ([]*TrackedActionItem, error)
}

// ActionItemFilter filters action items
type ActionItemFilter struct {
    Owner      string
    Status     string
    Priority   string
    PostMortemID string
    DueBefore  *time.Time
}

// NewActionItemTracker creates a new tracker
func NewActionItemTracker(store ActionItemStore) *ActionItemTracker {
    return &ActionItemTracker{
        items: make(map[string]*TrackedActionItem),
        store: store,
    }
}

// CreateActionItem creates a new action item
func (t *ActionItemTracker) CreateActionItem(ctx context.Context, pmID, incidentID string, item ActionItem) (*TrackedActionItem, error) {
    tracked := &TrackedActionItem{
        ActionItem:   item,
        PostMortemID: pmID,
        IncidentID:   incidentID,
        CreatedAt:    time.Now().UTC(),
        UpdatedAt:    time.Now().UTC(),
        Progress:     0,
        Notes:        make([]ProgressNote, 0),
        Blockers:     make([]string, 0),
    }

    t.mu.Lock()
    t.items[tracked.ID] = tracked
    t.mu.Unlock()

    if err := t.store.Save(ctx, tracked); err != nil {
        return nil, err
    }

    return tracked, nil
}

// UpdateProgress updates action item progress
func (t *ActionItemTracker) UpdateProgress(ctx context.Context, itemID, author, note string, progress int) error {
    t.mu.Lock()
    item, exists := t.items[itemID]
    if !exists {
        t.mu.Unlock()
        return fmt.Errorf("action item not found")
    }

    item.Progress = progress
    item.UpdatedAt = time.Now().UTC()

    if progress == 100 {
        now := time.Now().UTC()
        item.CompletedAt = &now
        item.Status = "closed"
    } else if progress > 0 {
        item.Status = "in_progress"
    }

    progressNote := ProgressNote{
        Timestamp: time.Now().UTC(),
        Author:    author,
        Note:      note,
        Progress:  progress,
    }
    item.Notes = append(item.Notes, progressNote)

    t.mu.Unlock()

    return t.store.Save(ctx, item)
}

// AddBlocker adds a blocker to an action item
func (t *ActionItemTracker) AddBlocker(ctx context.Context, itemID, blocker string) error {
    t.mu.Lock()
    item, exists := t.items[itemID]
    if !exists {
        t.mu.Unlock()
        return fmt.Errorf("action item not found")
    }

    item.Blockers = append(item.Blockers, blocker)
    item.UpdatedAt = time.Now().UTC()
    t.mu.Unlock()

    return t.store.Save(ctx, item)
}

// RemoveBlocker removes a blocker
func (t *ActionItemTracker) RemoveBlocker(ctx context.Context, itemID, blocker string) error {
    t.mu.Lock()
    item, exists := t.items[itemID]
    if !exists {
        t.mu.Unlock()
        return fmt.Errorf("action item not found")
    }

    newBlockers := make([]string, 0)
    for _, b := range item.Blockers {
        if b != blocker {
            newBlockers = append(newBlockers, b)
        }
    }
    item.Blockers = newBlockers
    item.UpdatedAt = time.Now().UTC()
    t.mu.Unlock()

    return t.store.Save(ctx, item)
}

// GetOverdueItems returns overdue action items
func (t *ActionItemTracker) GetOverdueItems() []*TrackedActionItem {
    t.mu.RLock()
    defer t.mu.RUnlock()

    now := time.Now().UTC()
    overdue := make([]*TrackedActionItem, 0)

    for _, item := range t.items {
        if item.Status != "closed" && now.After(item.DueDate) {
            overdue = append(overdue, item)
        }
    }

    // Sort by due date
    sort.Slice(overdue, func(i, j int) bool {
        return overdue[i].DueDate.Before(overdue[j].DueDate)
    })

    return overdue
}

// GetCompletionStats returns completion statistics
func (t *ActionItemTracker) GetCompletionStats(postMortemID string) CompletionStats {
    t.mu.RLock()
    defer t.mu.RUnlock()

    stats := CompletionStats{}

    for _, item := range t.items {
        if postMortemID != "" && item.PostMortemID != postMortemID {
            continue
        }

        stats.Total++

        switch item.Status {
        case "closed":
            stats.Closed++
        case "in_progress":
            stats.InProgress++
        default:
            stats.Open++
        }

        if item.Priority == "critical" {
            stats.CriticalTotal++
            if item.Status == "closed" {
                stats.CriticalClosed++
            }
        }

        if item.Status != "closed" && time.Now().UTC().After(item.DueDate) {
            stats.Overdue++
        }
    }

    if stats.Total > 0 {
        stats.CompletionRate = float64(stats.Closed) / float64(stats.Total) * 100
    }

    return stats
}

// CompletionStats contains completion statistics
type CompletionStats struct {
    Total          int
    Open           int
    InProgress     int
    Closed         int
    Overdue        int
    CriticalTotal  int
    CriticalClosed int
    CompletionRate float64
}

// StartReminderScheduler starts the reminder scheduler
func (t *ActionItemTracker) StartReminderScheduler(ctx context.Context, interval time.Duration, notifier ReminderNotifier) {
    ticker := time.NewTicker(interval)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            t.sendReminders(ctx, notifier)
        case <-ctx.Done():
            return
        }
    }
}

// ReminderNotifier sends reminders
type ReminderNotifier interface {
    SendReminder(ctx context.Context, item *TrackedActionItem) error
}

// sendReminders sends reminders for approaching/overdue items
func (t *ActionItemTracker) sendReminders(ctx context.Context, notifier ReminderNotifier) {
    now := time.Now().UTC()

    t.mu.RLock()
    items := make([]*TrackedActionItem, 0)
    for _, item := range t.items {
        if item.Status == "closed" {
            continue
        }

        // Check if reminder needed
        daysUntilDue := item.DueDate.Sub(now).Hours() / 24
        shouldRemind := false

        if daysUntilDue < 0 {
            // Overdue - remind daily
            shouldRemind = item.LastReminder == nil || now.Sub(*item.LastReminder) > 24*time.Hour
        } else if daysUntilDue <= 3 {
            // Due soon - remind every 2 days
            shouldRemind = item.LastReminder == nil || now.Sub(*item.LastReminder) > 48*time.Hour
        } else if daysUntilDue <= 7 {
            // Due in a week - remind once
            shouldRemind = item.LastReminder == nil
        }

        if shouldRemind {
            items = append(items, item)
        }
    }
    t.mu.RUnlock()

    for _, item := range items {
        if err := notifier.SendReminder(ctx, item); err == nil {
            t.mu.Lock()
            item.RemindersSent++
            now := time.Now().UTC()
            item.LastReminder = &now
            t.mu.Unlock()

            t.store.Save(ctx, item)
        }
    }
}
```

---

## 3. Production-Ready Configurations

### 3.1 Post-Mortem Automation Configuration

```yaml
# postmortem-automation.yaml
apiVersion: automation.io/v1
kind: PostMortemAutomation
metadata:
  name: postmortem-workflow
spec:
  # Trigger conditions
  triggers:
    - type: incident_resolved
      condition: severity in ['critical', 'high']
    - type: manual

  # Timeline configuration
  timeline:
    auto_collect: true
    sources:
      - type: alertmanager
        fields:
          - alert_name
          - fired_at
          - resolved_at
          - labels
      - type: pagerduty
        fields:
          - log_entries
          - notes
      - type: slack
        channel_pattern: "incident-*"
        extract_mentions: true
      - type: git
        extract_deployments: true
        lookback: 24h

  # Analysis configuration
  analysis:
    five_whys:
      enabled: true
      suggest_questions: true

    impact_calculation:
      enabled: true
      metrics:
        - error_rate
        - latency_p99
        - availability
      financial:
        enabled: true
        revenue_per_minute: 1000

  # Meeting scheduling
  scheduling:
    auto_schedule: true
    delay_after_incident: 48h
    duration: 60m
    required_attendees:
      - incident_commander
      - ops_lead
    optional_attendees:
      - affected_service_owners

    # Find first available slot
    scheduling_window:
      start: "09:00"
      end: "17:00"
      timezone: "America/New_York"

  # Documentation
  documentation:
    template: blameless-postmortem
    auto_populate: true
    fields:
      - timeline
      - metrics
      - root_cause_suggestions
      - action_item_templates

    # Review workflow
    review:
      enabled: true
      reviewers:
        - incident_commander
        - engineering_manager
      approval_required: true
      auto_publish_after: 7d

  # Action item management
  action_items:
    auto_create: true
    templates:
      - name: "Add monitoring"
        condition: "detection_time > 5m"
        description: "Add early detection for {{ incident.symptom }}"
        priority: high
        default_due: 14d

      - name: "Improve runbook"
        condition: "resolution_time > 30m"
        description: "Update runbook for {{ incident.service }} incidents"
        priority: medium
        default_due: 7d

      - name: "Add automated remediation"
        condition: "type == 'known_issue'"
        description: "Automate recovery for {{ incident.root_cause }}"
        priority: medium
        default_due: 30d

    tracking:
      enabled: true
      reminders:
        - at: 3d_before_due
        - at: 1d_before_due
        - at: overdue_daily

      escalations:
        - after: 7d_overdue
          notify: engineering_manager
        - after: 14d_overdue
          notify: director

  # Metrics and reporting
  metrics:
    track:
      - mttr
      - mttd
      - action_item_completion_rate
      - time_to_postmortem
      - recurrence_rate

    dashboards:
      - name: "Post-Mortem Metrics"
        widgets:
          - type: completion_rate
          - type: overdue_items
          - type: mttr_trend
          - type: category_breakdown
```

---

## 4. Security Considerations

### 4.1 Data Handling in Post-Mortems

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Post-Mortem Data Handling                                 │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  DATA TYPE          │  HANDLING              │  RETENTION                   │
├─────────────────────┼────────────────────────┼──────────────────────────────│
│  Timeline events    │  Include in document   │  7 years (compliance)        │
│  (logs, metrics)    │  Anonymized if needed  │                              │
│  ───────────────────┼────────────────────────┼──────────────────────────────│
│  Chat transcripts   │  Summarize key points  │  1 year                      │
│  (Slack, etc.)      │  Don't include PII     │                              │
│  ───────────────────┼────────────────────────┼──────────────────────────────│
│  Customer data      │  NEVER include         │  Per data retention policy   │
│  (names, emails)    │  Use aggregates only   │                              │
│  ───────────────────┼────────────────────────┼──────────────────────────────│
│  Internal IPs       │  Redact or use ranges  │  7 years                     │
│  System details     │                        │                              │
│  ───────────────────┼────────────────────────┼──────────────────────────────│
│  Security incidents │  Restricted access     │  7 years + legal hold        │
│                     │  Legal review required │                              │
│  ───────────────────┼────────────────────────┼──────────────────────────────│
│  Personal           │  Anonymize             │  Remove after review         │
│  performance        │                        │                              │
│                                                                             │
│  ACCESS CONTROLS:                                                           │
│  • Public post-mortems: Remove internal details                             │
│  • Internal post-mortems: Full employee access                              │
│  • Security post-mortems: Need-to-know only                                 │
│  • Legal hold: Restricted, legal team approval                              │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 5. Compliance Requirements

### 5.1 Regulatory Requirements

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                  Post-Mortem Compliance Requirements                         │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  SOC 2 TYPE II                    ISO 27001                                 │
│  ━━━━━━━━━━━━━━━━                 ━━━━━━━━━━━                               │
│                                                                             │
│  CC7.3 - Incident detection       A.16.1 - Management responsibilities      │
│  CC7.4 - Incident response        A.16.2 - Reporting information security   │
│  CC7.5 - Incident recovery        A.16.3 - Assessment and decision          │
│                                   A.16.4 - Evidence collection              │
│                                   A.16.5 - Lessons learned                  │
│                                                                             │
│  HIPAA                            PCI DSS                                   │
│  ━━━━━━━━                         ━━━━━━━━━                                 │
│                                                                             │
│  §164.312(a)(2)(ii) - Emergency   Req 12.10.1 - Incident response plan      │
│  §164.308(a)(6) - Security        Req 12.10.2 - Incident response           │
│        incidents                  procedures                                │
│  §164.312(b) - Audit controls     Req 12.10.4 - 24/7 response               │
│                                                                             │
│  GDPR                                                                     │
│  ━━━━━━                                                                     │
│                                                                             │
│  Art. 33 - Breach notification procedures                                   │
│  Art. 34 - Communication to data subjects                                   │
│  Art. 35 - Data Protection Impact Assessment                                │
│                                                                             │
│  COMMON REQUIREMENTS:                                                       │
│  ✓ Document all security incidents                                          │
│  ✓ Conduct root cause analysis                                              │
│  ✓ Implement corrective actions                                             │
│  ✓ Retain records for specified periods                                     │
│  ✓ Regular review of incident trends                                        │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 6. Decision Matrices

### 6.1 When to Conduct Post-Mortem Matrix

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Post-Mortem Trigger Decision Matrix                       │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  Condition                         │  Post-Mortem Required │  Timeframe     │
├────────────────────────────────────┼───────────────────────┼────────────────│
│  SEV 1 (Critical) incident         │  Mandatory            │  Within 48h    │
│  ──────────────────────────────────┼───────────────────────┼────────────────│
│  SEV 2 (High) incident             │  Mandatory            │  Within 72h    │
│  ──────────────────────────────────┼───────────────────────┼────────────────│
│  SEV 3 (Medium) incident           │  Optional             │  Within 1 week │
│  ──────────────────────────────────┼───────────────────────┼────────────────│
│  Data breach (any size)            │  Mandatory + Legal    │  Within 24h    │
│  ──────────────────────────────────┼───────────────────────┼────────────────│
│  Security incident                 │  Mandatory + Security │  Within 48h    │
│  ──────────────────────────────────┼───────────────────────┼────────────────│
│  Customer-impacting bug            │  Optional             │  Within 1 week │
│  ──────────────────────────────────┼───────────────────────┼────────────────│
│  Recurring issue (3+ times)        │  Mandatory            │  Within 1 week │
│  ──────────────────────────────────┼───────────────────────┼────────────────│
│  Near-miss (almost SEV 1)          │  Optional             │  Within 2 weeks│
│  ──────────────────────────────────┼───────────────────────┼────────────────│
│  On-call pain > 2 hours            │  Optional             │  Within 1 week │
│                                                                             │
│  EXEMPTIONS (with manager approval):                                        │
│  • Same root cause as recent post-mortem with open action items             │
│  • Duplicate incident (same ongoing issue)                                  │
│  • Planned maintenance incident                                             │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 6.2 Severity Classification Decision Matrix

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Action Item Priority Matrix                               │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  Impact →        │  Critical  │  High      │  Medium    │  Low             │
│  Effort ↓        │            │            │            │                  │
├──────────────────┼────────────┼────────────┼────────────┼──────────────────│
│  Low (hours)     │  P0 - Do   │  P1 - Do   │  P1 - Do   │  P2 - Schedule   │
│                  │  immediate │  immediate │  soon      │                  │
├──────────────────┼────────────┼────────────┼────────────┼──────────────────│
│  Medium (days)   │  P1 - Do   │  P1 - Do   │  P2 - Plan │  P3 - Backlog    │
│                  │  soon      │  soon      │            │                  │
├──────────────────┼────────────┼────────────┼────────────┼──────────────────│
│  High (weeks)    │  P1 - Plan │  P2 - Plan │  P2 - Plan │  P3 - Evaluate   │
│                  │  carefully │            │            │                  │
├──────────────────┼────────────┼────────────┼────────────┼──────────────────│
│  Very High       │  P2 - May  │  P3 -      │  P3 - May  │  P4 - Don't do   │
│  (months+)       │  need exec │  Evaluate  │  not do    │                  │
│                  │  approval  │            │            │                  │
│                                                                             │
│  Priority Definitions:                                                      │
│  • P0: Interrupt current work, complete immediately                         │
│  • P1: Complete within 1 sprint                                             │
│  • P2: Complete within 1 quarter                                            │
│  • P3: Next quarter planning                                                │
│  • P4: Not now, re-evaluate later                                           │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 6.3 Meeting Structure Decision Matrix

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Post-Mortem Meeting Structure Matrix                      │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  Incident Severity │  Duration  │  Required Attendees    │  Format          │
├────────────────────┼────────────┼────────────────────────┼──────────────────│
│  SEV 1 (Critical)  │  90 min    │  IC, Ops Lead, EM,     │  Synchronous,    │
│                    │            │  All responders,       │  Video required  │
│                    │            │  Service owners,       │                  │
│                    │            │  Engineering Director  │                  │
├────────────────────┼────────────┼────────────────────────┼──────────────────│
│  SEV 2 (High)      │  60 min    │  IC, Ops Lead, EM,     │  Synchronous,    │
│                    │            │  Service owners        │  Video preferred │
├────────────────────┼────────────┼────────────────────────┼──────────────────│
│  SEV 3 (Medium)    │  30 min    │  IC, EM                │  Async document  │
│                    │  or async  │                        │  review OK       │
├────────────────────┼────────────┼────────────────────────┼──────────────────│
│  Security Incident │  90 min    │  + Security team,      │  Synchronous,    │
│                    │            │  + Legal (if needed)   │  Restricted      │
├────────────────────┼────────────┼────────────────────────┼──────────────────│
│  Recurring Issue   │  60 min    │  IC, EM, Previous PM   │  Synchronous,    │
│                    │            │  authors               │  Focus on        │
│                    │            │                        │  patterns        │
│                                                                             │
│  Meeting Agenda Template (60 min):                                         │
│  • 0-5 min:   Introduction and context                                     │
│  • 5-20 min:  Timeline review                                              │
│  • 20-40 min: 5 Whys analysis                                              │
│  • 40-50 min: Impact assessment                                            │
│  • 50-60 min: Action items and owners                                      │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 7. Best Practices Summary

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Post-Mortem Best Practices Summary                        │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  BEFORE THE POST-MORTEM                                                     │
│  ✓ Schedule within 48 hours of incident resolution                          │
│  ✓ Collect timeline from multiple sources (logs, chat, alerts)              │
│  ✓ Calculate impact metrics (MTTR, affected users)                          │
│  ✓ Prepare preliminary 5 Whys analysis                                      │
│  ✓ Invite all relevant participants                                         │
│  ✓ Send pre-read with timeline and data                                     │
│                                                                             │
│  DURING THE POST-MORTEM                                                     │
│  ✓ Start with blameless culture reminder                                    │
│  ✓ Review timeline chronologically                                          │
│  ✓ Ask "why" until reaching systemic factors                                │
│  ✓ Focus on process, not people                                             │
│  ✓ Identify multiple contributing factors                                   │
│  ✓ Note what went well (not just what went wrong)                           │
│  ✓ Define specific, measurable action items                                 │
│  ✓ Assign owners and due dates                                              │
│                                                                             │
│  DOCUMENTATION                                                              │
│  ✓ Write clearly for future readers                                         │
│  ✓ Include complete timeline with evidence                                  │
│  ✓ Document root cause and contributing factors                             │
│  ✓ Quantify impact where possible                                           │
│  ✓ Include lessons learned from all categories                              │
│  ✓ List all action items with owners                                        │
│  ✓ Get review and approval before publishing                                │
│                                                                             │
│  FOLLOW-UP                                                                  │
│  ✓ Track action items to completion                                         │
│  ✓ Send regular reminders for overdue items                                 │
│  ✓ Measure completion rates                                                 │
│  ✓ Review trends across multiple post-mortems                               │
│  ✓ Update runbooks and procedures                                           │
│  ✓ Share learnings with broader organization                                │
│                                                                             │
│  CULTURE                                                                    │
│  ✓ Make post-mortems blameless by policy                                    │
│  ✓ Share post-mortems company-wide (when appropriate)                       │
│  ✓ Celebrate good catches and effective responses                           │
│  ✓ Recognize action item completion                                         │
│  ✓ Continuously improve the process                                         │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## References

1. Google SRE Book - Postmortem Culture
2. Etsy Engineering - Blameless PostMortems and a Just Culture
3. PagerDuty - Post-Incident Reviews
4. Site Reliability Workbook - Postmortems
5. John Allspaw - Blameless PostMortems and a Just Culture

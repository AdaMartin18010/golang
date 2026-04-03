# EC-057: Blue-Green Deployment Pattern

> **Dimension**: Engineering-CloudNative
> **Level**: S (>15KB)
> **Tags**: #blue-green-deployment #zero-downtime #deployment-strategy #traffic-switching #rollback
> **Authoritative Sources**:
>
> - [Continuous Delivery](https://continuousdelivery.com/) - Humble & Farley (2010)
> - [The Release It!](https://pragprog.com/titles/mnee2/release-it-second-edition/) - Michael Nygard (2018)

---

## 1. Problem Formalization

### 1.1 System Context and Constraints

**Definition 1.1 (Blue-Green Deployment Domain)**
Let $\mathcal{E} = \{Blue, Green\}$ be two identical production environments where:

- At any time $t$, exactly one environment $E_{active}(t)$ serves 100% traffic
- The other environment $E_{standby}(t)$ is either idle or being prepared
- Switching $swap(E_{active}, E_{standby})$ must be atomic

**System Constraints:**

| Constraint | Formal Definition | Impact |
|------------|-------------------|--------|
| **Zero Downtime** | $downtime(swap) = 0$ | Must maintain availability |
| **Quick Rollback** | $T_{rollback} < 1min$ | Instant switch required |
| **Resource Cost** | $resources(\mathcal{E}) = 2 \times resources(single)$ | Double infrastructure |
| **Data Consistency** | $\forall d \in data: consistent(d, E_{Blue}) \Leftrightarrow consistent(d, E_{Green})$ | Shared data or migration |

### 1.2 Problem Statement

**Problem 1.1 (Instant Deployment)**
Given new version $V_{new}$, deploy to $E_{standby}$ and switch traffic such that:

$$\forall user: downtime(user) = 0 \land detect(failure) \Rightarrow swap(E_{active}, E_{standby})$$

**Key Challenges:**

1. **Database Migration**: Schema changes between versions
2. **Session Persistence**: Maintaining user sessions across switch
3. **Resource Cost**: Running double infrastructure
4. **Warm-up Time**: Pre-warming standby environment
5. **Smoke Testing**: Validating before switching

---

## 2. Solution Architecture

### 2.1 Blue-Green States

| State | Blue | Green | Router |
|-------|------|-------|--------|
| **Initial** | Active (v1) | Idle | Blue |
| **Preparing** | Active (v1) | Deploying v2 | Blue |
| **Ready** | Active (v1) | Ready (v2) | Blue |
| **Switching** | Standby (v1) | Active (v2) | Green |
| **Cleanup** | Idle (v1) | Active (v2) | Green |

---

## 3. Visual Representations

### 3.1 Blue-Green Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    BLUE-GREEN DEPLOYMENT ARCHITECTURE                        │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  INITIAL STATE:                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                         LOAD BALANCER / ROUTER                       │   │
│  │                                                                      │   │
│  │  ┌─────────────────────────┐                                         │   │
│  │  │  Traffic: 100%          │                                         │   │
│  │  │  Target: BLUE           │                                         │   │
│  │  │                         │                                         │   │
│  │  │  ┌─────┐                │                                         │   │
│  │  │  │BLUE │◄───────────────┼──── All traffic                         │   │
│  │  │  │     │                │                                         │   │
│  │  │  └─────┘                │                                         │   │
│  │  │       ┌─────┐           │                                         │   │
│  │  │       │GREEN│ (Idle)    │                                         │   │
│  │  │       │     │           │                                         │   │
│  │  │       └─────┘           │                                         │   │
│  │  └─────────────────────────┘                                         │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                             │
│  DEPLOYMENT STATE:                                                          │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                                                                      │   │
│  │  ┌─────────────────────────┐                                         │   │
│  │  │  Traffic: 100%          │                                         │   │
│  │  │  Target: BLUE           │                                         │   │
│  │  │                         │                                         │   │
│  │  │  ┌─────┐                │                                         │   │
│  │  │  │BLUE │◄───────────────┼──── Production traffic                  │   │
│  │  │  │v1   │                │                                         │   │
│  │  │  └─────┘                │                                         │   │
│  │  │       ┌─────┐           │                                         │   │
│  │  │       │GREEN│◄──┐       │                                         │   │
│  │  │       │v2   │   │       │                                         │   │
│  │  │       └─────┘   │       │                                         │   │
│  │  │                 │       │                                         │   │
│  │  │  ┌──────────────┘       │                                         │   │
│  │  │  │  Smoke tests         │                                         │   │
│  │  │  │  (validation)        │                                         │   │
│  │  │  └──────────────────────┘                                         │   │
│  │  └─────────────────────────┘                                         │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                             │
│  SWITCHED STATE:                                                            │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                                                                      │   │
│  │  ┌─────────────────────────┐                                         │   │
│  │  │  Traffic: 100%          │                                         │   │
│  │  │  Target: GREEN          │                                         │   │
│  │  │                         │                                         │   │
│  │  │  ┌─────┐                │                                         │   │
│  │  │  │BLUE │ (Standby)      │                                         │   │
│  │  │  │v1   │                │                                         │   │
│  │  │  └─────┘                │                                         │   │
│  │  │       ┌─────┐           │                                         │   │
│  │  │       │GREEN│◄───────────┼──── All traffic                         │   │
│  │  │       │v2   │            │                                         │   │
│  │  │       └─────┘            │                                         │   │
│  │  └─────────────────────────┘                                         │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                             │
│  ROLLBACK PATH (Instant):                                                   │
│                                                                             │
│  Detect issue ──► Switch router back to BLUE ──► GREEN idle               │
│  (< 1 second)                                                                 │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 4. Production Go Implementation

### 4.1 Blue-Green Controller

```go
package bluegreen

import (
 "context"
 "fmt"
 "time"

 "go.uber.org/zap"
)

// Environment represents a deployment environment
type Environment string

const (
 EnvironmentBlue  Environment = "blue"
 EnvironmentGreen Environment = "green"
)

// State represents the deployment state
type State int

const (
 StateIdle State = iota
 StateDeploying
 StateTesting
 StateSwitching
 StateActive
 StateRollingBack
)

// Deployment represents a deployment
type Deployment struct {
 ID            string
 Version       string
 Environment   Environment
 State         State
 CreatedAt     time.Time
 ActivatedAt   *time.Time
 DeactivatedAt *time.Time
}

// Router manages traffic routing
type Router interface {
 RouteTraffic(ctx context.Context, target Environment) error
 GetCurrentTarget(ctx context.Context) (Environment, error)
}

// Deployer manages deployments
type Deployer interface {
 Deploy(ctx context.Context, env Environment, version string) error
 HealthCheck(ctx context.Context, env Environment) error
 Stop(ctx context.Context, env Environment) error
}

// Controller manages blue-green deployments
type Controller struct {
 router   Router
 deployer Deployer
 logger   *zap.Logger

 blue  *Deployment
 green *Deployment

 mu sync.RWMutex
}

// NewController creates a new blue-green controller
func NewController(router Router, deployer Deployer, logger *zap.Logger) *Controller {
 return &Controller{
  router:   router,
  deployer: deployer,
  logger:   logger,
  blue: &Deployment{
   Environment: EnvironmentBlue,
   State:       StateIdle,
  },
  green: &Deployment{
   Environment: EnvironmentGreen,
   State:       StateIdle,
  },
 }
}

// Deploy performs a blue-green deployment
func (c *Controller) Deploy(ctx context.Context, version string) error {
 // Determine standby environment
 currentTarget, err := c.router.GetCurrentTarget(ctx)
 if err != nil {
  return fmt.Errorf("failed to get current target: %w", err)
 }

 standby := c.getStandby(currentTarget)

 c.logger.Info("Starting blue-green deployment",
  zap.String("version", version),
  zap.String("standby", string(standby)))

 // Deploy to standby environment
 deployment := c.getDeployment(standby)
 deployment.Version = version
 deployment.State = StateDeploying
 deployment.CreatedAt = time.Now()

 if err := c.deployer.Deploy(ctx, standby, version); err != nil {
  deployment.State = StateIdle
  return fmt.Errorf("failed to deploy to %s: %w", standby, err)
 }

 // Health check
 deployment.State = StateTesting
 if err := c.deployer.HealthCheck(ctx, standby); err != nil {
  c.deployer.Stop(ctx, standby)
  deployment.State = StateIdle
  return fmt.Errorf("health check failed for %s: %w", standby, err)
 }

 // Switch traffic
 deployment.State = StateSwitching
 if err := c.router.RouteTraffic(ctx, standby); err != nil {
  // Rollback on switch failure
  c.router.RouteTraffic(ctx, currentTarget)
  c.deployer.Stop(ctx, standby)
  deployment.State = StateIdle
  return fmt.Errorf("failed to switch traffic: %w", err)
 }

 now := time.Now()
 deployment.State = StateActive
 deployment.ActivatedAt = &now

 // Deactivate previous environment
 previous := c.getDeployment(currentTarget)
 previous.State = StateIdle
 previous.DeactivatedAt = &now

 c.logger.Info("Blue-green deployment complete",
  zap.String("version", version),
  zap.String("environment", string(standby)))

 return nil
}

// Rollback switches traffic back to the previous environment
func (c *Controller) Rollback(ctx context.Context) error {
 currentTarget, err := c.router.GetCurrentTarget(ctx)
 if err != nil {
  return fmt.Errorf("failed to get current target: %w", err)
 }

 standby := c.getStandby(currentTarget)

 c.logger.Info("Rolling back deployment",
  zap.String("from", string(currentTarget)),
  zap.String("to", string(standby)))

 // Switch traffic to standby
 if err := c.router.RouteTraffic(ctx, standby); err != nil {
  return fmt.Errorf("failed to rollback traffic: %w", err)
 }

 // Update states
 current := c.getDeployment(currentTarget)
 current.State = StateRollingBack

 previous := c.getDeployment(standby)
 now := time.Now()
 previous.State = StateActive
 previous.ActivatedAt = &now

 c.logger.Info("Rollback complete",
  zap.String("active_environment", string(standby)))

 return nil
}

func (c *Controller) getStandby(current Environment) Environment {
 if current == EnvironmentBlue {
  return EnvironmentGreen
 }
 return EnvironmentBlue
}

func (c *Controller) getDeployment(env Environment) *Deployment {
 if env == EnvironmentBlue {
  return c.blue
 }
 return c.green
}
```

---

## 5. Failure Scenarios and Mitigations

| Scenario | Impact | Detection | Mitigation |
|----------|--------|-----------|------------|
| **Switch Failure** | Split traffic | Router error | Automatic retry |
| **Green Crash** | 100% errors | Health check | Instant rollback |
| **Database Incompatibility** | Data corruption | Migration error | Schema versioning |
| **Session Loss** | Users logged out | Session errors | Shared session store |

---

## 6. Semantic Trade-off Analysis

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    BLUE-GREEN vs CANARY COMPARISON                           │
├─────────────────────┬─────────────────┬─────────────────────────────────────┤
│     Dimension       │   Blue-Green    │              Canary                 │
├─────────────────────┼─────────────────┼─────────────────────────────────────┤
│ Resource Cost       │ 2x (double)     │ 1.25x-1.5x (gradual)                │
│ Rollback Speed      │ Instant         │ Gradual (minutes)                   │
│ Risk Level          │ All or nothing  │ Gradual exposure                    │
│ User Impact         │ Zero (if ok)    │ Some affected (if issues)           │
│ Testing Window      │ Limited time    │ Extended (days possible)            │
│ Database Complexity │ High (migrations)│ Medium                              │
│ Session Handling    │ Complex         │ Simple                              │
└─────────────────────┴─────────────────┴─────────────────────────────────────┘
```

---

## 7. References

1. Humble, J., & Farley, D. (2010). *Continuous Delivery*. Addison-Wesley.
2. Nygard, M. (2018). *Release It! (2nd Edition)*. Pragmatic Programmers.

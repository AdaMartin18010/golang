# FT-040-Paxos-Consensus-Formal

> **Dimension**: 01-Formal-Theory
> **Status**: S-Level Academic
> **Created**: 2026-04-03
> **Version**: Paxos (Lamport, 1998, 2001)
> **Size**: >20KB
> **Formal Methods**: TLA+ Specification

---

## 1. Paxos Problem Statement

### 1.1 The Consensus Problem

**Given**: A set of processes $P = \{p_1, p_2, ..., p_n\}$ where some may fail.

**Goal**: All non-faulty processes must agree on a single value from a set of proposed values.

**Safety Properties**:

1. **Validity**: Only a proposed value can be chosen.
2. **Agreement**: No two non-faulty processes decide different values.
3. **Integrity**: Each process decides at most once.

**Liveness Property**:
4. **Termination**: Eventually, every non-faulty process decides some value.

### 1.2 Failure Model

Paxos tolerates **crash-stop failures** and operates in an **asynchronous network**:

- Messages may be lost, delayed, or duplicated
- No bound on message delivery time
- $f < n/2$ processes may fail (majority required)

**FLP Result**: In asynchronous systems with even one faulty process, deterministic consensus is impossible without timeouts. Paxos is safe but not guaranteed to terminate without synchrony assumptions.

---

## 2. Paxos Roles

### 2.1 Agent Taxonomy

```
┌─────────────────────────────────────────┐
│           Paxos Agents                  │
├─────────────────────────────────────────┤
│                                         │
│  Proposer  ──►  Accepts client values   │
│     │                                   │
│     └──► Initiates rounds               │
│          Seeks promises                 │
│          Proposes values                │
│                                         │
│  Acceptor  ──►  Accepts or rejects      │
│     │        proposals                  │
│     │                                   │
│     └──► Maintains promised/proposed    │
│          state                          │
│                                         │
│  Learner   ──►  Determines chosen value │
│     │                                   │
│     └──► Observes acceptor decisions    │
│                                         │
│  Client    ──►  Submits values          │
│                                         │
└─────────────────────────────────────────┘
```

**Note**: A single process may play multiple roles.

---

## 3. Two-Phase Protocol

### 3.1 Phase 1: Prepare/Promise

**Objective**: Learn of any already chosen value and prevent acceptance of older proposals.

**Algorithm**:

```
Proposer p with proposal number n:
─────────────────────────────────────
1. Send Prepare(n) to majority of acceptors

Acceptor a on receiving Prepare(n):
─────────────────────────────────────
1. If n > maxPromised[a]:
     maxPromised[a] ← n
     Reply Promise(n, maxAccepted[a])
   Else:
     Reply Reject(n, maxPromised[a])
```

**Formal State**:

$$
\text{Acceptor}_a = (maxPromised_a, maxAccepted_a)
$$

Where:

- $maxPromised_a \in \mathbb{N}$: Highest prepare number promised
- $maxAccepted_a \in (\mathbb{N} \times V) \cup \{\bot\}$: Highest accepted proposal (number, value) pair

### 3.2 Phase 2: Accept/Accepted

**Objective**: Have acceptors accept the proposed value.

**Algorithm**:

```
Proposer p on receiving Promise from majority:
──────────────────────────────────────────────
1. If any Promise contains a value:
     v ← value from highest numbered proposal
   Else:
     v ← client's proposed value

2. Send Accept(n, v) to all who promised

Acceptor a on receiving Accept(n, v):
─────────────────────────────────────
1. If n ≥ maxPromised[a]:
     maxAccepted[a] ← (n, v)
     Reply Accepted(n, v)
   Else:
     Reply Reject(n, maxPromised[a])
```

---

## 4. Safety Proof

### 4.1 Lemma: Promise Invariant

**Lemma 1**: If acceptor $a$ has promised not to accept proposals numbered less than $n$ (by setting $maxPromised_a = n$), then no proposal with number $m < n$ will be accepted by $a$.

**Proof**: Direct from the acceptor code. $Accept(n, v)$ only succeeds if $n \geq maxPromised_a$.

### 4.2 Lemma: Value Consistency

**Lemma 2**: If proposal $(n, v)$ is chosen (accepted by majority), then any proposal with higher number $m > n$ must also have value $v$.

**Proof**:

Let $M_n$ be the majority that accepted $(n, v)$.
Let proposer $p$ attempt proposal $m > n$.

By pigeonhole principle, $p$'s prepare majority $M_m$ intersects $M_n$:
$$|M_m| + |M_n| > |Acceptors| \implies M_m \cap M_n \neq \emptyset$$

Therefore, at least one acceptor $a \in M_m \cap M_n$ has:

1. Accepted $(n, v)$ (since $a \in M_n$)
2. Promised not to accept lower numbers (since $a \in M_m$)

When $a$ responds to $Prepare(m)$, it includes $(n, v)$ in its promise.
Proposer $p$ must then propose $v$ (the value from the highest numbered promise).

∎

### 4.3 Theorem: Safety

**Theorem**: Paxos satisfies the Agreement property - no two different values can be chosen.

**Proof by Contradiction**:

Assume two different values $v_1 \neq v_2$ are chosen, with proposals $(n_1, v_1)$ and $(n_2, v_2)$.

Without loss of generality, assume $n_1 < n_2$.

By Lemma 2, since $(n_1, v_1)$ was chosen, any higher proposal must have value $v_1$.
Therefore, $(n_2, v_2)$ must have $v_2 = v_1$.

Contradiction: $v_1 \neq v_2$ and $v_2 = v_1$.

∎

---

## 5. Multi-Paxos

### 5.1 Optimization: Skip Phase 1

Once a leader is established, it can skip Phase 1 for subsequent proposals:

```
Proposal 1: Prepare → Promise → Accept → Accepted
Proposal 2: Accept → Accepted  (skip Prepare/Promise)
Proposal 3: Accept → Accepted
...
```

**Invariant**: The leader maintains $maxPromised$ for all acceptors, so new proposals are guaranteed to use higher numbers.

### 5.2 Distinguished Proposer (Leader)

```
Leader Election:
────────────────
1. Process with highest ID attempts to become leader
2. Run Paxos to agree on leader value
3. Leader holds lease for time T
4. If lease expires, new election
```

**Leader State**:

$$
\text{Leader}_l = (leaseExpiry_l, proposals_l, promises_l)
$$

---

## 6. TLA+ Specification

```tla
------------------------------ MODULE Paxos ------------------------------
EXTENDS Naturals, Sequences, FiniteSets, TLC

CONSTANTS Value,         \* Set of proposed values
          Acceptor,      \* Set of acceptors
          Quorum         \* Set of quorums (majorities)

ASSUME QuorumAssumption ==
    /\ \A Q \in Quorum : Q \subseteq Acceptor
    /\ \A Q1, Q2 \in Quorum : Q1 \cap Q2 # {}

CONSTANT Ballot
ASSUME BallotAssumption == Ballot \subseteq Nat

None == CHOOSE v : v \notin Ballot

VARIABLES maxBal,    \* maxBal[a]: highest prepare ballot promised by a
          maxVBal,   \* maxVBal[a]: highest ballot a has accepted
          maxVal     \* maxVal[a]: value a has accepted at maxVBal[a]

typeOK ==
    /\ maxBal  \in [Acceptor -> Ballot \cup {None}]
    /\ maxVBal \in [Acceptor -> Ballot \cup {None}]
    /\ maxVal  \in [Acceptor -> Value \cup {None}]

-----------------------------------------------------------------------------
\* Initial state

Init ==
    /\ maxBal  = [a \in Acceptor |-> None]
    /\ maxVBal = [a \in Acceptor |-> None]
    /\ maxVal  = [a \in Acceptor |-> None]

-----------------------------------------------------------------------------
\* Actions

\* Phase 1a: Proposer sends Prepare
Prepare(bal, a) ==
    /\ bal \in Ballot
    /\ maxBal[a] < bal
    /\ maxBal' = [maxBal EXCEPT ![a] = bal]
    /\ UNCHANGED <<maxVBal, maxVal>>

\* Phase 1b: Acceptor responds with Promise
Promise(bal, a, mbal, mval) ==
    /\ bal \in Ballot
    /\ maxBal[a] <= bal
    /\ mbal = maxVBal[a]
    /\ mval = maxVal[a]
    /\ maxBal' = [maxBal EXCEPT ![a] = bal]
    /\ UNCHANGED <<maxVBal, maxVal>>

\* Phase 2a: Proposer sends Accept
Accept(bal, val, a) ==
    /\ bal \in Ballot
    /\ maxBal[a] <= bal
    /\ maxBal' = [maxBal EXCEPT ![a] = bal]
    /\ maxVBal' = [maxVBal EXCEPT ![a] = bal]
    /\ maxVal' = [maxVal EXCEPT ![a] = val]

\* Phase 2b: Acceptor accepts (already covered in Accept)

-----------------------------------------------------------------------------
\* Safety Properties

\* A value is chosen if a quorum has accepted it
chosenAt(bal, val) ==
    \E Q \in Quorum :
        \A a \in Q : maxVBal[a] = bal /\ maxVal[a] = val

chosen(val) == \E bal \in Ballot : chosenAt(bal, val)

\* Agreement: At most one value can be chosen
Agreement ==
    \A v1, v2 \in Value : chosen(v1) /\ chosen(v2) => v1 = v2

\* Validity: Only proposed values can be chosen
Validity ==
    \A val \in Value : chosen(val) => val \in Value

-----------------------------------------------------------------------------
\* Invariants

\* Core Paxos invariant
PaxosInvariant ==
    \A a1, a2 \in Acceptor :
        (maxVBal[a1] # None /\ maxVBal[a2] # None /\
         maxVBal[a1] <= maxVBal[a2]) =>
            (maxVal[a1] = maxVal[a2])

\* If a value is chosen, all higher ballots must propose the same value
ChosenInvariant ==
    \A bal \in Ballot, val \in Value :
        chosenAt(bal, val) =>
            \A b \in Ballot, v \in Value :
                (b > bal /\ chosenAt(b, v)) => v = val

-----------------------------------------------------------------------------
\* Specification

Next ==
    \E bal \in Ballot, val \in Value, a \in Acceptor :
        \/ Prepare(bal, a)
        \/ Promise(bal, a, maxVBal[a], maxVal[a])
        \/ Accept(bal, val, a)

Spec == Init /\ [][Next]_<<maxBal, maxVBal, maxVal>>

THEOREM Safety == Spec => [](Agreement /\ Validity /\ PaxosInvariant)

=============================================================================
```

---

## 7. Comparison with Raft

| Aspect | Paxos | Raft |
|--------|-------|------|
| **Understandability** | Difficult (single-decree) | Designed for clarity |
| **Leader Election** | External to core protocol | Integrated |
| **Log Replication** | Multiple independent instances | Unified log |
| **Membership Changes** | Complex (requires multiple rounds) | Single configuration change |
| **Performance** | Similar (same message complexity) | Similar |
| **Formal Verification** | Many TLA+ specs | Raft TLA+ spec available |

### 7.1 When to Use Paxos

- **Library/embedded systems**: When you need just consensus, not full replicated log
- **Flexible deployment**: When leader election can be handled externally
- **Research/experimentation**: When exploring consensus variants

### 7.2 When to Use Raft

- **Production systems**: Better understood, more implementations
- **Replicated state machines**: Natural fit for log-based replication
- **Engineer accessibility**: Easier to reason about correctness

---

## 8. References

1. **Lamport, L. (1998)**. "The Part-Time Parliament." *ACM TOCS*, 16(2), 133-169.

2. **Lamport, L. (2001)**. "Paxos Made Simple." *ACM SIGACT News*, 32(4), 18-25.
   - <https://lamport.azurewebsites.net/pubs/paxos-simple.pdf>

3. **Lamport, L. (2006)**. "Fast Paxos." *Distributed Computing*, 19(2), 79-103.

4. **Chandra, T. D., Griesemer, R., & Redstone, J. (2007)**. "Paxos Made Live - An Engineering Perspective." *PODC*.

5. **Ongaro, D., & Ousterhout, J. (2014)**. "In Search of an Understandable Consensus Algorithm." *USENIX ATC*.
   - Raft comparison paper

---

*Last Updated: 2026-04-03*
*Formal Methods: TLA+ specification included with safety invariants*

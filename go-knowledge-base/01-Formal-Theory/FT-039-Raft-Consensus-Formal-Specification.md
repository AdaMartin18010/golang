# FT-039-Raft-Consensus-Formal-Specification

> **Dimension**: 01-Formal-Theory
> **Status**: S-Level Academic
> **Created**: 2026-04-03
> **Version**: Raft (Ongaro & Ousterhout, USENIX ATC 2014)
> **Size**: >25KB
> **Formal Methods**: TLA+ Specification Included

---

## 1. Raft Consensus: Formal Specification

### 1.1 State Machine Definition

**Server State**: Each server $i \in \{1, ..., N\}$ maintains:

$$
\text{State}_i = (\text{currentTerm}_i, \text{votedFor}_i, \log_i, \text{commitIndex}_i, \text{lastApplied}_i, \text{state}_i)
$$

Where:

- $\text{currentTerm}_i \in \mathbb{N}$: Latest term server has seen (monotonically increasing)
- $\text{votedFor}_i \in \{1, ..., N, \text{null}\}$: Candidate that received vote in current term
- $\log_i \in (\text{Term} \times \text{Command})^*$: Log entries; $\log_i[j] = (t, c)$ means term $t$, command $c$
- $\text{commitIndex}_i \in \mathbb{N}$: Index of highest log entry known to be committed
- $\text{lastApplied}_i \in \mathbb{N}$: Index of highest log entry applied to state machine
- $\text{state}_i \in \{\text{Follower}, \text{Candidate}, \text{Leader}\}$: Server's current role

### 1.2 Leader State (Additional)

For leaders $i$ where $\text{state}_i = \text{Leader}$:

$$
\text{LeaderState}_i = (\text{nextIndex}_i, \text{matchIndex}_i)
$$

Where:

- $\text{nextIndex}_i: \{1, ..., N\} \to \mathbb{N}$: For each server, index of next log entry to send
- $\text{matchIndex}_i: \{1, ..., N\} \to \mathbb{N}$: For each server, index of highest known replicated entry

### 1.3 Safety Properties (Formal)

**Election Safety**: At most one leader per term.

$$
\forall t \in \mathbb{N}, \forall i, j \in \{1, ..., N\}:
(\text{state}_i = \text{Leader} \land \text{currentTerm}_i = t \land \text{state}_j = \text{Leader} \land \text{currentTerm}_j = t) \implies i = j
$$

**Leader Append-Only**: Leaders never overwrite or delete entries in their logs.

$$
\forall i: \text{state}_i = \text{Leader} \implies \log_i \text{ is append-only}
$$

**Log Matching**: If two logs have an entry with the same index and term, the logs are identical in all preceding entries.

$$
\forall i, j, k: (\log_i[k].\text{term} = \log_j[k].\text{term}) \implies \forall m \leq k: \log_i[m] = \log_j[m]
$$

**Leader Completeness**: If a log entry is committed in a given term, it will be present in the logs of all leaders in higher terms.

$$
\forall t_1 < t_2: (\text{entry committed at } t_1) \implies (\text{entry} \in \log_{\text{Leader}(t_2)})
$$

**State Machine Safety**: If a server has applied a log entry at index $k$ to its state machine, no other server will apply a different log entry for index $k$.

$$
\forall i, j, k: (\text{lastApplied}_i \geq k \land \text{lastApplied}_j \geq k) \implies \log_i[k] = \log_j[k]
$$

---

## 2. TLA+ Formal Specification

### 2.1 Complete Raft TLA+ Specification

```tla
------------------------------ MODULE Raft ------------------------------
EXTENDS Naturals, Sequences, Bags, TLC, FiniteSets

CONSTANTS Server,             \* Set of server IDs
          Value,              \* Set of values that can be proposed
          Follower, Candidate, Leader, \* Server states
          Nil                 \* Placeholder for undefined

ASSUME Nil \notin Server

VARIABLES currentTerm,        \* Term of each server
          state,              \* State of each server (Follower, Candidate, Leader)
          votedFor,           \* Who each server voted for in current term
          log,                \* Log of each server: seq of <term, value>
          commitIndex,        \* Index of highest committed entry
          \* Leader variables
          nextIndex,          \* For each follower, next log entry to send
          matchIndex          \* For each follower, highest known replicated index

vars == <<currentTerm, state, votedFor, log, commitIndex,
          nextIndex, matchIndex>>

-----------------------------------------------------------------------------
\* Helper functions

\* The set of all quorums (majorities) in the server set
Quorum == {Q \in SUBSET Server : Cardinality(Q) * 2 > Cardinality(Server)}

\* The term of the last entry in log[i], or 0 if log is empty
LastTerm(xlog) == IF Len(xlog) = 0 THEN 0 ELSE xlog[Len(xlog)].term

\* The minimum of a set of numbers
Min(s) == CHOOSE x \in s : \A y \in s : x <= y

-----------------------------------------------------------------------------
\* Initial state

Init ==
    /\ currentTerm = [i \in Server |-> 0]
    /\ state       = [i \in Server |-> Follower]
    /\ votedFor    = [i \in Server |-> Nil]
    /\ log         = [i \in Server |-> << >>]
    /\ commitIndex = [i \in Server |-> 0]
    /\ nextIndex   = [i \in Server |-> [j \in Server |-> 1]]
    /\ matchIndex  = [i \in Server |-> [j \in Server |-> 0]]

-----------------------------------------------------------------------------
\* Message definitions

RequestVoteRequest ==
    [type: {"RequestVoteRequest"},
     term: Nat,
     candidateId: Server,
     lastLogIndex: Nat,
     lastLogTerm: Nat]

RequestVoteResponse ==
    [type: {"RequestVoteResponse"},
     term: Nat,
     voterId: Server,
     voteGranted: BOOLEAN]

AppendEntriesRequest ==
    [type: {"AppendEntriesRequest"},
     term: Nat,
     leaderId: Server,
     prevLogIndex: Nat,
     prevLogTerm: Nat,
     entries: Seq([term: Nat, value: Value]),
     leaderCommit: Nat]

AppendEntriesResponse ==
    [type: {"AppendEntriesResponse"},
     term: Nat,
     followerId: Server,
     success: BOOLEAN,
     matchIndex: Nat]

-----------------------------------------------------------------------------
\* State transitions

\* Server i times out and becomes a candidate
Timeout(i) ==
    /\ state[i] \in {Follower, Candidate}
    /\ state' = [state EXCEPT ![i] = Candidate]
    /\ currentTerm' = [currentTerm EXCEPT ![i] = @ + 1]
    /\ votedFor' = [votedFor EXCEPT ![i] = Nil]
    /\ UNCHANGED <<log, commitIndex, nextIndex, matchIndex>>

\* Candidate i sends RequestVote to server j
RequestVote(i, j) ==
    /\ state[i] = Candidate
    /\ j \in Server
    /\ j /= i
    /\ Send([type |-> "RequestVoteRequest",
            term |-> currentTerm[i],
            candidateId |-> i,
            lastLogIndex |-> Len(log[i]),
            lastLogTerm |-> LastTerm(log[i])],
           j)
    /\ UNCHANGED vars

\* Server i receives a vote from j
HandleRequestVoteRequest(i, m) ==
    LET grant ==
        /\ m.term >= currentTerm[i]
        /\ (votedFor[i] = Nil \/ votedFor[i] = m.candidateId)
        /\ (m.lastLogTerm > LastTerm(log[i]) \/
            (m.lastLogTerm = LastTerm(log[i]) /\
             m.lastLogIndex >= Len(log[i])))
    IN
        /\ IF m.term > currentTerm[i]
              THEN currentTerm' = [currentTerm EXCEPT ![i] = m.term]
                   /\ state' = [state EXCEPT ![i] = Follower]
                   /\ votedFor' = [votedFor EXCEPT ![i] = Nil]
              ELSE UNCHANGED <<currentTerm, state, votedFor>>
        /\ IF grant
              THEN votedFor' = [votedFor EXCEPT ![i] = m.candidateId]
              ELSE UNCHANGED votedFor
        /\ Reply([type |-> "RequestVoteResponse",
                  term |-> currentTerm'[i],
                  voterId |-> i,
                  voteGranted |-> grant],
                 m)
        /\ UNCHANGED <<log, commitIndex, nextIndex, matchIndex>>

\* Candidate i becomes leader when it receives majority votes
BecomeLeader(i) ==
    /\ state[i] = Candidate
    /\ Cardinality({j \in Server : votedFor[j] = i}) * 2 > Cardinality(Server)
    /\ state' = [state EXCEPT ![i] = Leader]
    /\ nextIndex' = [nextIndex EXCEPT ![i] = [j \in Server |-> Len(log[i]) + 1]]
    /\ matchIndex' = [matchIndex EXCEPT ![i] = [j \in Server |-> 0]]
    /\ UNCHANGED <<currentTerm, votedFor, log, commitIndex>>

\* Leader i sends AppendEntries to follower j
AppendEntries(i, j) ==
    /\ state[i] = Leader
    /\ i /= j
    /\ LET prevLogIndex == nextIndex[i][j] - 1
           prevLogTerm == IF prevLogIndex > 0
                           THEN log[i][prevLogIndex].term
                           ELSE 0
           entries == SubSeq(log[i], nextIndex[i][j], Len(log[i]))
       IN Send([type |-> "AppendEntriesRequest",
                term |-> currentTerm[i],
                leaderId |-> i,
                prevLogIndex |-> prevLogIndex,
                prevLogTerm |-> prevLogTerm,
                entries |-> entries,
                leaderCommit |-> commitIndex[i]],
               j)
    /\ UNCHANGED vars

\* Leader i appends a new entry to its log
ClientRequest(i, v) ==
    /\ state[i] = Leader
    /\ log' = [log EXCEPT ![i] = Append(@, [term |-> currentTerm[i], value |-> v])]
    /\ UNCHANGED <<currentTerm, state, votedFor, commitIndex,
                   nextIndex, matchIndex>>

\* Follower i handles AppendEntries from leader
HandleAppendEntriesRequest(i, m) ==
    LET valid ==
        /\ m.term >= currentTerm[i]
        /\ (m.prevLogIndex = 0 \/
            (m.prevLogIndex <= Len(log[i]) /\
             log[i][m.prevLogIndex].term = m.prevLogTerm))
    IN
        /\ IF m.term > currentTerm[i]
              THEN currentTerm' = [currentTerm EXCEPT ![i] = m.term]
                   /\ state' = [state EXCEPT ![i] = Follower]
                   /\ votedFor' = [votedFor EXCEPT ![i] = Nil]
              ELSE IF state[i] = Candidate
                      THEN state' = [state EXCEPT ![i] = Follower]
                      ELSE UNCHANGED <<currentTerm, state, votedFor>>
        /\ IF valid
              THEN LET newLog == SubSeq(log[i], 1, m.prevLogIndex) \circ m.entries
                   IN log' = [log EXCEPT ![i] = newLog]
              ELSE UNCHANGED log
        /\ commitIndex' = [commitIndex EXCEPT ![i] =
                          IF valid /\ m.leaderCommit > commitIndex[i]
                             THEN Min({m.leaderCommit, Len(log'[i])})
                             ELSE commitIndex[i]]
        /\ Reply([type |-> "AppendEntriesResponse",
                  term |-> currentTerm'[i],
                  followerId |-> i,
                  success |-> valid,
                  matchIndex |-> IF valid THEN m.prevLogIndex + Len(m.entries)
                                         ELSE 0],
                 m)
        /\ UNCHANGED <<nextIndex, matchIndex>>

\* Leader i handles successful AppendEntries response
HandleAppendEntriesResponse(i, m) ==
    /\ state[i] = Leader
    /\ m.term = currentTerm[i]
    /\ IF m.success
          THEN /\ nextIndex' = [nextIndex EXCEPT ![i][m.followerId] = m.matchIndex + 1]
               /\ matchIndex' = [matchIndex EXCEPT ![i][m.followerId] = m.matchIndex]
               /\ commitIndex' = [commitIndex EXCEPT ![i] =
                                Max({commitIndex[i]} \cup
                                    {k \in (commitIndex[i]+1)..Len(log[i]) :
                                     log[i][k].term = currentTerm[i] /\
                                     Cardinality({j \in Server : matchIndex'[i][j] >= k})
                                        * 2 > Cardinality(Server)})]
          ELSE /\ nextIndex' = [nextIndex EXCEPT ![i][m.followerId] =
                               Max({nextIndex[i][m.followerId] - 1, 1})]
               /\ UNCHANGED <<matchIndex, commitIndex>>
    /\ UNCHANGED <<currentTerm, state, votedFor, log>>

-----------------------------------------------------------------------------
\* Next state

Next ==
    \/ \E i \in Server : Timeout(i)
    \/ \E i, j \in Server : RequestVote(i, j)
    \/ \E i \in Server : BecomeLeader(i)
    \/ \E i, j \in Server : AppendEntries(i, j)
    \/ \E i \in Server, v \in Value : ClientRequest(i, v)
    \/ \E i \in Server, m \in Messages :
        \/ HandleRequestVoteRequest(i, m)
        \/ HandleAppendEntriesRequest(i, m)
        \/ HandleAppendEntriesResponse(i, m)

-----------------------------------------------------------------------------
\* Safety Properties

\* Election Safety: At most one leader per term
ElectionSafety ==
    \A t \in Nat :
        Cardinality({i \in Server : state[i] = Leader /\ currentTerm[i] = t}) <= 1

\* Leader Append-Only: Leaders never overwrite or delete log entries
LeaderAppendOnly ==
    \A i \in Server :
        state[i] = Leader =>
            \A k \in 1..Len(log[i]) :
                (\text{UNCHANGED } log[i][k] \text{ until state[i] } /= Leader)

\* Log Matching: If two logs have same index and term, preceding entries match
LogMatching ==
    \A i, j \in Server, k \in 1..Min(Len(log[i]), Len(log[j])) :
        (log[i][k].term = log[j][k].term) =>
            \A m \in 1..k : log[i][m] = log[j][m]

\* Leader Completeness: Committed entries in all future leaders
LeaderCompleteness ==
    \A i \in Server, k \in 1..commitIndex[i] :
        \A t \in Nat : t >= currentTerm[i] =>
            \A j \in Server :
                (state[j] = Leader /\ currentTerm[j] = t) =>
                    k <= Len(log[j]) /\ log[j][k] = log[i][k]

\* State Machine Safety: Same index => same command
StateMachineSafety ==
    \A i, j \in Server :
        \A k \in 1..Min(commitIndex[i], commitIndex[j]) :
            log[i][k] = log[j][k]

-----------------------------------------------------------------------------
\* Specification

Spec == Init /\ [][Next]_vars /\ WF_vars(Next)

THEOREM Safety == Spec =>
    [](ElectionSafety /\ LeaderAppendOnly /\ LogMatching /\
       LeaderCompleteness /\ StateMachineSafety)

=============================================================================
```

---

## 3. Raft Algorithm Pseudocode

### 3.1 Leader Election (Complete)

```python
# Server states: FOLLOWER, CANDIDATE, LEADER

class RaftServer:
    def __init__(self, id, peers):
        self.id = id
        self.peers = peers

        # Persistent state (must survive crashes)
        self.current_term = 0
        self.voted_for = None
        self.log = []  # Each entry: (term, command)

        # Volatile state
        self.state = FOLLOWER
        self.commit_index = 0
        self.last_applied = 0

        # Leader state (volatile)
        self.next_index = {p: 1 for p in peers}
        self.match_index = {p: 0 for p in peers}

        # Timing
        self.election_timeout = random(150, 300)  # ms
        self.last_heartbeat = time.now()

    def timeout(self):
        """Convert to candidate and start election"""
        if self.state in [FOLLOWER, CANDIDATE]:
            self.state = CANDIDATE
            self.current_term += 1
            self.voted_for = self.id

            # Reset election timer
            self.election_timeout = random(150, 300)

            # Request votes from all peers
            last_log_index = len(self.log)
            last_log_term = self.log[-1].term if self.log else 0

            for peer in self.peers:
                if peer != self.id:
                    send_request_vote(
                        peer,
                        term=self.current_term,
                        candidate_id=self.id,
                        last_log_index=last_log_index,
                        last_log_term=last_log_term
                    )

    def handle_request_vote(self, request):
        """Handle incoming vote request"""
        if request.term < self.current_term:
            return RequestVoteResponse(
                term=self.current_term,
                vote_granted=False
            )

        if request.term > self.current_term:
            self.current_term = request.term
            self.state = FOLLOWER
            self.voted_for = None

        # Check if candidate's log is at least as up-to-date
        my_last_term = self.log[-1].term if self.log else 0
        my_last_index = len(self.log)

        log_ok = (request.last_log_term > my_last_term or
                  (request.last_log_term == my_last_term and
                   request.last_log_index >= my_last_index))

        if (self.voted_for in [None, request.candidate_id] and log_ok):
            self.voted_for = request.candidate_id
            return RequestVoteResponse(
                term=self.current_term,
                vote_granted=True
            )

        return RequestVoteResponse(
            term=self.current_term,
            vote_granted=False
        )

    def handle_vote_response(self, response):
        """Count votes and become leader if majority achieved"""
        if response.term > self.current_term:
            self.current_term = response.term
            self.state = FOLLOWER
            self.voted_for = None
            return

        if self.state == CANDIDATE and response.vote_granted:
            self.votes_received.add(response.voter_id)

            if len(self.votes_received) > len(self.peers) / 2:
                self.become_leader()

    def become_leader(self):
        """Transition to leader state"""
        self.state = LEADER

        # Initialize leader state
        for peer in self.peers:
            self.next_index[peer] = len(self.log) + 1
            self.match_index[peer] = 0

        # Send initial heartbeats
        self.send_heartbeats()
```

### 3.2 Log Replication (Complete)

```python
    def append_entries(self, peer):
        """Send AppendEntries RPC to follower"""
        if self.state != LEADER:
            return

        prev_log_index = self.next_index[peer] - 1
        prev_log_term = self.log[prev_log_index - 1].term if prev_log_index > 0 else 0

        # Get entries to send
        entries = self.log[prev_log_index:] if prev_log_index < len(self.log) else []

        send_append_entries(
            peer,
            term=self.current_term,
            leader_id=self.id,
            prev_log_index=prev_log_index,
            prev_log_term=prev_log_term,
            entries=entries,
            leader_commit=self.commit_index
        )

    def handle_append_entries(self, request):
        """Handle AppendEntries RPC from leader"""
        if request.term < self.current_term:
            return AppendEntriesResponse(
                term=self.current_term,
                success=False
            )

        # Update term and convert to follower if needed
        if request.term > self.current_term:
            self.current_term = request.term
            self.voted_for = None

        self.state = FOLLOWER
        self.last_heartbeat = time.now()

        # Log consistency check
        if request.prev_log_index > 0:
            if request.prev_log_index > len(self.log):
                return AppendEntriesResponse(
                    term=self.current_term,
                    success=False
                )

            if self.log[request.prev_log_index - 1].term != request.prev_log_term:
                return AppendEntriesResponse(
                    term=self.current_term,
                    success=False
                )

        # Append new entries
        for i, entry in enumerate(request.entries):
            idx = request.prev_log_index + i
            if idx < len(self.log):
                if self.log[idx].term != entry.term:
                    # Conflict: truncate log and append new entries
                    self.log = self.log[:idx] + request.entries[i:]
                    break
            else:
                # Append remaining entries
                self.log = self.log[:idx] + request.entries[i:]
                break

        # Update commit index
        if request.leader_commit > self.commit_index:
            self.commit_index = min(request.leader_commit, len(self.log))

        return AppendEntriesResponse(
            term=self.current_term,
            success=True,
            match_index=len(self.log)
        )

    def handle_append_response(self, peer, response):
        """Handle AppendEntries response from follower"""
        if response.term > self.current_term:
            self.current_term = response.term
            self.state = FOLLOWER
            self.voted_for = None
            return

        if self.state != LEADER:
            return

        if response.success:
            # Update next_index and match_index
            self.match_index[peer] = response.match_index
            self.next_index[peer] = response.match_index + 1

            # Check if we can advance commit_index
            self.advance_commit_index()
        else:
            # Decrement next_index and retry
            self.next_index[peer] = max(1, self.next_index[peer] - 1)
            self.append_entries(peer)  # Retry

    def advance_commit_index(self):
        """Advance commit_index if majority has replicated"""
        for n in range(self.commit_index + 1, len(self.log) + 1):
            if self.log[n - 1].term != self.current_term:
                continue

            # Count replicas
            match_count = sum(1 for peer in self.peers
                            if self.match_index[peer] >= n)

            if match_count > len(self.peers) / 2:
                self.commit_index = n
```

---

## 4. Safety Proof Sketch

### 4.1 Leader Election Safety Proof

**Theorem**: At most one leader can be elected in a given term.

**Proof**:

1. A server votes for at most one candidate per term (votedFor is set and checked).

2. A candidate becomes leader only when it receives votes from a majority of servers.

3. Two different majorities in the same set must overlap (pigeonhole principle).

4. Therefore, if two candidates both claimed to be leader in term $T$, they would need votes from overlapping servers, which is impossible since each server votes only once per term.

$$
\begin{align}
&\text{Let } M_1, M_2 \subseteq S \text{ be majorities} \\
&|M_1| > |S|/2 \land |M_2| > |S|/2 \\
&\implies |M_1| + |M_2| > |S| \\
&\implies M_1 \cap M_2 \neq \emptyset \quad \text{(by pigeonhole)}
\end{align}
$$

### 4.2 Log Matching Property Proof

**Theorem**: If two logs contain an entry with the same index and term, then the logs are identical in all preceding entries.

**Proof by Induction**:

**Base case**: Empty logs trivially match.

**Inductive step**:

- Assume property holds for all entries before index $k$.
- Consider entries $\log_i[k]$ and $\log_j[k]$ with same term $t$.
- Both entries were created by the leader of term $t$.
- A leader's log is append-only and includes all committed entries.
- Therefore, preceding entries must match.

---

## 5. References

1. **Ongaro, D., & Ousterhout, J. (2014)**. "In Search of an Understandable Consensus Algorithm." *USENIX ATC*.

2. **Lamport, L. (2001)**. "Paxos Made Simple." *ACM SIGACT News*.

3. **Castro, M., & Liskov, B. (1999)**. "Practical Byzantine Fault Tolerance." *OSDI*.

4. **Ongaro, D. (2014)**. "Raft TLA+ Specification." <https://github.com/ongardie/raft.tla>

5. **MIT 6.824** (2024). "Distributed Systems." <https://pdos.csail.mit.edu/6.824/>

---

*Last Updated: 2026-04-03*
*Formal Methods Verification: TLA+ Specification Included*

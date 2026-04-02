# FT-019: Operational Transformation - Formal Specification

## Overview

Operational Transformation (OT) is a technology for supporting collaborative editing of shared documents. It allows multiple users to concurrently edit a document and synchronizes their changes while maintaining consistency. Unlike CRDTs, OT requires a central server or coordination point to transform operations before applying them.

## Theoretical Foundations

### 1.1 OT Model

**Basic Concepts**:

```
Document: A sequence of characters D = c₁c₂...cₙ
Operation: A transformation on a document op: D → D'
Operation Types:
  - Insert(pos, text): Insert text at position
  - Delete(pos, length): Delete length characters from position
  - Retain(length): Skip length characters (identity)
```

**Operation Context**:

```
Each operation has a context - the document state against which it was defined.
op = (type, params, context)

Two operations are concurrent if they have the same context.
```

### 1.2 Transformation Functions

**Core Transformation Property (TP1)**:

```
For any two concurrent operations O₁ and O₂ defined against document D:

O₁' = T(O₁, O₂)
O₂' = T(O₂, O₁)

Such that:
apply(apply(D, O₁), O₂') = apply(apply(D, O₂), O₁')

This ensures that applying O₁ then O₂' yields the same document
as applying O₂ then O₁'.
```

**Convergence Property (TP2)**:

```
For any three concurrent operations O₁, O₂, O₃:

T(T(O₁, O₂), T(O₃, O₂)) = T(T(O₁, O₃), T(O₂, O₃))

This ensures that transformation is associative and transformations
converge regardless of the order they are applied.
```

### 1.3 Transformation Functions Formal Definition

```
Transformation Function T:
T: Operation × Operation → Operation

Insert-Insert Transformation:
T(Insert(p₁, s₁), Insert(p₂, s₂)):
  if p₁ < p₂:
    return Insert(p₁, s₁)
  else if p₁ > p₂:
    return Insert(p₁ + |s₂|, s₁)
  else:  // p₁ = p₂
    if tiebreaker(s₁) < tiebreaker(s₂):
      return Insert(p₁, s₁)
    else:
      return Insert(p₁ + |s₂|, s₁)

Insert-Delete Transformation:
T(Insert(pᵢ, s), Delete(pᵈ, l)):
  if pᵢ ≤ pᵈ:
    return Insert(pᵢ, s)
  else if pᵢ ≥ pᵈ + l:
    return Insert(pᵢ - l, s)
  else:  // pᵢ is within deleted range
    return Insert(pᵈ, s)

T(Delete(pᵈ, lᵈ), Insert(pⁱ, s)):
  if pⁱ ≥ pᵈ + lᵈ:
    return Delete(pᵈ, lᵈ)
  else if pⁱ ≤ pᵈ:
    return Delete(pᵈ + |s|, lᵈ)
  else:  // insertion within deletion
    return Delete(pᵈ, lᵈ + |s|)

Delete-Delete Transformation:
T(Delete(p₁, l₁), Delete(p₂, l₂)):
  if p₁ + l₁ ≤ p₂:
    return Delete(p₁, l₁)
  else if p₂ + l₂ ≤ p₁:
    return Delete(p₁ - l₂, l₁)
  else:  // overlapping deletions
    // Calculate intersection
    start = max(p₁, p₂)
    end = min(p₁ + l₁, p₂ + l₂)
    newLen = l₁ - (end - start)
    newPos = min(p₁, p₂)
    return Delete(newPos, newLen)
```

### 1.4 Correctness Proof

```
Theorem 1 (Transformation Correctness): The transformation functions
satisfy property TP1.

Proof (Insert-Insert case):
Let D be the initial document, |D| = n.
Let O₁ = Insert(p₁, s₁) and O₂ = Insert(p₂, s₂) with 0 ≤ p₁, p₂ ≤ n.

Case 1: p₁ < p₂
  O₁' = T(O₁, O₂) = Insert(p₁, s₁)
  O₂' = T(O₂, O₁) = Insert(p₂ + |s₁|, s₂)

  apply(apply(D, O₁), O₂') = apply(D[0:p₁] + s₁ + D[p₁:], O₂')
                           = D[0:p₁] + s₁ + D[p₁:p₂] + s₂ + D[p₂:]

  apply(apply(D, O₂), O₁') = apply(D[0:p₂] + s₂ + D[p₂:], O₁')
                           = D[0:p₁] + s₁ + D[p₁:p₂] + s₂ + D[p₂:]

  Therefore: apply(apply(D, O₁), O₂') = apply(apply(D, O₂), O₁') ∎

Case 2: p₁ > p₂
  Symmetric to Case 1.

Case 3: p₁ = p₂ (tie-breaker required)
  Let t₁ = tiebreaker(s₁), t₂ = tiebreaker(s₂)

  If t₁ < t₂:
    O₁' = Insert(p₁, s₁)
    O₂' = Insert(p₂ + |s₁|, s₂)

    Both executions: D[0:p] + s₁ + s₂ + D[p:]

  Result is the same regardless of execution order. ∎

Theorem 2 (Convergence): The transformation functions satisfy property TP2.

Proof Sketch:
TP2 requires that transformation paths converge:
  O₁ → O₁(O₂) → O₁(O₂)(O₃(O₂))
  O₁ → O₁(O₃) → O₁(O₃)(O₂(O₃))

For insert operations, the final position depends on the sum of
lengths of all preceding concurrent insertions, which is independent
of transformation order. Therefore TP2 holds. ∎
```

## TLA+ Specification

```tla
--------------------------- MODULE OperationalTransformation ---------------------------
EXTENDS Integers, Sequences, FiniteSets, TLC

CONSTANTS Replicas,        \* Set of client replicas
            Server,         \* Central server
            Characters,     \* Set of possible characters
            MaxOps          \* Maximum operations for model checking

VARIABLES clientStates,    \* State at each client
          serverState,     \* State at server
          opLog,           \* Log of all operations
          pendingOps       \* Operations waiting for server

\* Document is a sequence of characters
Document == Seq(Characters)

\* Operation types
OpType == {"insert", "delete", "retain"}

\* Operation structure
Operation == [type: OpType,
              pos: Nat,
              len: Nat,          \* for delete/retain
              text: Seq(Characters),  \* for insert
              client: Replicas,
              timestamp: Nat,
              context: Nat]       \* Document version

\* Initial state
Init ==
  ∧ clientStates = [r ∈ Replicas ↦ [doc ↦ ⟨⟩, version ↦ 0]]
  ∧ serverState = [doc ↦ ⟨⟩, version ↦ 0, history ↦ ⟨⟩]
  ∧ opLog = ⟨⟩
  ∧ pendingOps = [r ∈ Replicas ↦ ⟨⟩]

\* Client generates an operation
ClientOperation(client, type, pos, content) ==
  LET ctx == clientStates[client].version
      op == [type ↦ type,
             pos ↦ pos,
             len ↦ Len(content),
             text ↦ content,
             client ↦ client,
             timestamp ↦ Len(opLog) + 1,
             context ↦ ctx]
  IN ∧ pendingOps' = [pendingOps EXCEPT ![client] = Append(@, op)]
     ∧ UNCHANGED ⟨clientStates, serverState, opLog⟩

\* Transform operation against all operations in a log
TransformAgainstLog(op, log) ==
  IF Len(log) = 0
  THEN op
  ELSE LET transformed == Transform(op, log[1])
       IN TransformAgainstLog(transformed, SubSeq(log, 2, Len(log)))

\* Core transformation function (simplified)
Transform(op1, op2) ==
  IF op1.type = "insert" ∧ op2.type = "insert"
  THEN IF op1.pos < op2.pos ∨ (op1.pos = op2.pos ∧ op1.client < op2.client)
       THEN op1
       ELSE [op1 EXCEPT !.pos = @ + Len(op2.text)]
  ELSE IF op1.type = "insert" ∧ op2.type = "delete"
  THEN IF op1.pos ≤ op2.pos
       THEN op1
       ELSE IF op1.pos ≥ op2.pos + op2.len
       THEN [op1 EXCEPT !.pos = @ - op2.len]
       ELSE [op1 EXCEPT !.pos = op2.pos]
  ELSE IF op1.type = "delete" ∧ op2.type = "insert"
  THEN IF op2.pos ≥ op1.pos + op1.len
       THEN op1
       ELSE IF op2.pos ≤ op1.pos
       THEN [op1 EXCEPT !.pos = @ + Len(op2.text)]
       ELSE [op1 EXCEPT !.len = @ + Len(op2.text)]
  ELSE  \* delete-delete
       IF op1.pos + op1.len ≤ op2.pos
       THEN op1
       ELSE IF op2.pos + op2.len ≤ op1.pos
       THEN [op1 EXCEPT !.pos = @ - op2.len]
       ELSE  \* overlapping
            LET newStart == Min(op1.pos, op2.pos)
                newEnd == Max(op1.pos + op1.len, op2.pos + op2.len)
            IN [op1 EXCEPT !.pos = newStart, !.len = newEnd - newStart]

\* Server processes pending operation
ServerProcess(client) ==
  ∧ Len(pendingOps[client]) > 0
  ∧ LET op == Head(pendingOps[client])
        history == serverState.history
        transformed == TransformAgainstLog(op, history)
        newDoc == ApplyOperation(serverState.doc, transformed)
    IN ∧ serverState' = [doc ↦ newDoc,
                         version ↦ serverState.version + 1,
                         history ↦ Append(history, transformed)]
       ∧ pendingOps' = [pendingOps EXCEPT ![client] = Tail(@)]
       ∧ opLog' = Append(opLog, transformed)
       ∧ UNCHANGED clientStates

\* Apply operation to document
ApplyOperation(doc, op) ==
  CASE op.type = "insert" → SubSeq(doc, 1, op.pos) ∘ op.text ∘ SubSeq(doc, op.pos + 1, Len(doc))
    [] op.type = "delete" → SubSeq(doc, 1, op.pos) ∘ SubSeq(doc, op.pos + op.len + 1, Len(doc))
    [] op.type = "retain" → doc

\* Client receives acknowledgment and updates state
ClientAcknowledge(client) ==
  ∧ Len(opLog) > clientStates[client].version
  ∧ LET serverVersion == clientStates[client].version + 1
        newOp == opLog[serverVersion]
    IN ∧ clientStates' = [clientStates EXCEPT
           ![client].doc = ApplyOperation(@.doc, newOp),
           ![client].version = serverVersion]
       ∧ UNCHANGED ⟨serverState, opLog, pendingOps⟩

\* Next state
Next ==
  ∨ ∃r ∈ Replicas, pos ∈ 0..10, text ∈ Seq(Characters):
      ∧ Len(text) > 0
      ∧ ClientOperation(r, "insert", pos, text)
  ∨ ∃r ∈ Replicas: ServerProcess(r)
  ∨ ∃r ∈ Replicas: ClientAcknowledge(r)

\* Invariants

\* All clients eventually see the same document
EventualConsistency ==
  ◇(∀r1, r2 ∈ Replicas: clientStates[r1].doc = clientStates[r2].doc)

\* Document length is bounded
DocumentBound ==
  ∀r ∈ Replicas: Len(clientStates[r].doc) ≤ 20

\* Server document is the canonical version
ServerCanonical ==
  ∀r ∈ Replicas:
    clientStates[r].version = serverState.version
    ⇒ clientStates[r].doc = serverState.doc

=============================================================================
```

## Algorithm Pseudocode

### Basic OT Algorithm

```
Algorithm: Basic Operational Transformation (Client-Server)

Types:
  Document: sequence of characters
  Operation:
    - Insert(position: int, text: string)
    - Delete(position: int, length: int)
    - Retain(length: int)
  Operation has: clientId, timestamp, contextVersion

State at Client:
  doc: Document
  serverVersion: int  // Last acknowledged server version
  pending: Queue<Operation>  // Operations not yet sent
  buffer: Queue<Operation>   // Operations sent but not acked

State at Server:
  doc: Document
  version: int
  history: List<Operation>  // All operations in order

ClientGenerate(op):
  op.contextVersion = serverVersion
  pending.enqueue(op)
  localDoc = apply(localDoc, op)
  ClientSend()

ClientSend():
  while pending not empty and buffer empty:
    op = pending.dequeue()
    buffer.enqueue(op)
    sendToServer(op)

ClientReceive(ackOp):
  // Server acknowledged our operation
  sentOp = buffer.dequeue()
  serverVersion = ackOp.timestamp

  // Transform pending operations against acked operation
  for op in pending:
    op = Transform(op, sentOp)

  ClientSend()

ClientReceiveServerOp(serverOp):
  // Another client's operation from server
  if buffer not empty:
    // Transform server op against our unacked ops
    for ourOp in buffer:
      serverOp = Transform(serverOp, ourOp)

  localDoc = apply(localDoc, serverOp)
  serverVersion = serverOp.timestamp

ServerReceive(op):
  // Transform against all ops since op's context
  opsSinceContext = history[op.contextVersion+1 ... version]

  for histOp in opsSinceContext:
    op = Transform(op, histOp)

  // Apply and broadcast
  doc = apply(doc, op)
  version = version + 1
  op.timestamp = version
  history.append(op)

  broadcastToAllClients(op)

Transform(op1, op2):
  // See formal definition above
  return transformed operation

Apply(doc, op):
  switch op.type:
    case Insert:
      return doc[0:op.pos] + op.text + doc[op.pos:]
    case Delete:
      return doc[0:op.pos] + doc[op.pos + op.length:]
    case Retain:
      return doc
```

### Undo/Redo with OT

```
Algorithm: Undo/Redo Support in Operational Transformation

Additional Types:
  InverseOperation:
    - inverse(Insert(p, s)) = Delete(p, |s|)
    - inverse(Delete(p, l)) = Insert(p, deleted_text)

State:
  undoStack: Stack<Operation>
  redoStack: Stack<Operation>

Undo():
  if undoStack empty: return

  op = undoStack.pop()
  inverseOp = inverse(op)

  // Transform inverse against concurrent ops
  for pendingOp in pending:
    inverseOp = Transform(inverseOp, pendingOp)

  // Apply locally
  localDoc = apply(localDoc, inverseOp)

  // Send to server (transformed)
  pending.enqueue(inverseOp)
  redoStack.push(op)
  ClientSend()

Redo():
  if redoStack empty: return

  op = redoStack.pop()

  // Transform against concurrent ops
  for pendingOp in pending:
    op = Transform(op, pendingOp)

  // Apply locally
  localDoc = apply(localDoc, op)

  // Send to server
  pending.enqueue(op)
  undoStack.push(op)
  ClientSend()

Inverse(op):
  switch op.type:
    case Insert(p, s):
      return Delete(p, len(s))
    case Delete(p, l):
      // Need to capture deleted text at the time of deletion
      deletedText = captureTextBeforeDeletion(p, l)
      return Insert(p, deletedText)
```

### Jupiter Algorithm (Google Docs)

```
Algorithm: Jupiter (Google Docs-style Collaborative Editing)

State at Client i:
  doc_i: Document
  outQueue: Queue<Operation>  // Operations to send to server
  inQueue: Queue<Operation>   // Operations from server
  counter: int  // Number of local ops since last server op

State at Server:
  doc_s: Document
  clientCounters: Map<Client, int>
  history: List<Operation>

ClientGenerate(op):
  // Mark with current counter
  op.clientCounter = counter
  counter = counter + 1

  // Apply locally immediately
  doc_i = apply(doc_i, op)
  outQueue.enqueue(op)
  sendToServer(op)

ClientReceive(serverOp):
  // Transform against all local ops with lower or equal counter
  for localOp in outQueue:
    if localOp.clientCounter < serverOp.serverCounter:
      serverOp = Transform(serverOp, localOp)
    else:
      localOp = Transform(localOp, serverOp)

  // Apply server op
  doc_i = apply(doc_i, serverOp)
  counter = 0  // Reset after server acknowledgment

ServerReceive(clientOp):
  client = clientOp.clientId

  // Transform against all ops from other clients since last ack
  for otherOp in history:
    if otherOp.clientId ≠ client:
      clientOp = Transform(clientOp, otherOp)

  // Mark with server counter
  clientOp.serverCounter = clientCounters[client]
  clientCounters[client] = clientCounters[client] + 1

  // Apply and broadcast
  doc_s = apply(doc_s, clientOp)
  history.append(clientOp)

  broadcastToAllClientsExcept(clientOp, client)

Key Insight:
- Server maintains separate counter for each client
- Operations are transformed based on their "distance" from the server state
- Allows O(n) transformation complexity instead of O(n²)
```

## Go Implementation

```go
// Package ot implements Operational Transformation algorithms
package ot

import (
 "fmt"
 "strings"
 "sync"
)

// OpType represents the type of operation
type OpType int

const (
 OpInsert OpType = iota
 OpDelete
 OpRetain
)

// Operation represents a document operation
type Operation struct {
 Type      OpType
 Position  int
 Length    int
 Text      string
 ClientID  string
 Timestamp int64
 Context   int // Document version this op was created against
}

// String returns string representation
func (op Operation) String() string {
 switch op.Type {
 case OpInsert:
  return fmt.Sprintf("Insert(%d, %q)", op.Position, op.Text)
 case OpDelete:
  return fmt.Sprintf("Delete(%d, %d)", op.Position, op.Length)
 case OpRetain:
  return fmt.Sprintf("Retain(%d)", op.Length)
 default:
  return "Unknown"
 }
}

// Document represents a text document
type Document string

// Apply applies an operation to a document
func (d Document) Apply(op Operation) (Document, error) {
 switch op.Type {
 case OpInsert:
  if op.Position < 0 || op.Position > len(d) {
   return d, fmt.Errorf("insert position %d out of bounds", op.Position)
  }
  return d[:op.Position] + Document(op.Text) + d[op.Position:], nil

 case OpDelete:
  if op.Position < 0 || op.Position+op.Length > len(d) {
   return d, fmt.Errorf("delete range [%d, %d) out of bounds", op.Position, op.Position+op.Length)
  }
  return d[:op.Position] + d[op.Position+op.Length:], nil

 case OpRetain:
  return d, nil

 default:
  return d, fmt.Errorf("unknown operation type: %v", op.Type)
 }
}

// Transform transforms op1 against op2
func Transform(op1, op2 Operation) (Operation, error) {
 // Handle different operation type combinations
 switch {
 case op1.Type == OpInsert && op2.Type == OpInsert:
  return transformInsertInsert(op1, op2)

 case op1.Type == OpInsert && op2.Type == OpDelete:
  return transformInsertDelete(op1, op2)

 case op1.Type == OpDelete && op2.Type == OpInsert:
  return transformDeleteInsert(op1, op2)

 case op1.Type == OpDelete && op2.Type == OpDelete:
  return transformDeleteDelete(op1, op2)

 default:
  return op1, nil // Retain doesn't affect other ops
 }
}

func transformInsertInsert(op1, op2 Operation) (Operation, error) {
 if op1.Position < op2.Position {
  return op1, nil
 } else if op1.Position > op2.Position {
  return Operation{
   Type:      op1.Type,
   Position:  op1.Position + len(op2.Text),
   Text:      op1.Text,
   ClientID:  op1.ClientID,
   Timestamp: op1.Timestamp,
   Context:   op1.Context,
  }, nil
 } else {
  // Same position - use client ID as tiebreaker
  if op1.ClientID < op2.ClientID {
   return op1, nil
  }
  return Operation{
   Type:      op1.Type,
   Position:  op1.Position + len(op2.Text),
   Text:      op1.Text,
   ClientID:  op1.ClientID,
   Timestamp: op1.Timestamp,
   Context:   op1.Context,
  }, nil
 }
}

func transformInsertDelete(op1, op2 Operation) (Operation, error) {
 if op1.Position <= op2.Position {
  return op1, nil
 } else if op1.Position >= op2.Position+op2.Length {
  return Operation{
   Type:      op1.Type,
   Position:  op1.Position - op2.Length,
   Text:      op1.Text,
   ClientID:  op1.ClientID,
   Timestamp: op1.Timestamp,
   Context:   op1.Context,
  }, nil
 } else {
  // Insert position is within deleted range
  return Operation{
   Type:      op1.Type,
   Position:  op2.Position,
   Text:      op1.Text,
   ClientID:  op1.ClientID,
   Timestamp: op1.Timestamp,
   Context:   op1.Context,
  }, nil
 }
}

func transformDeleteInsert(op1, op2 Operation) (Operation, error) {
 if op2.Position >= op1.Position+op1.Length {
  return op1, nil
 } else if op2.Position <= op1.Position {
  return Operation{
   Type:      op1.Type,
   Position:  op1.Position + len(op2.Text),
   Length:    op1.Length,
   ClientID:  op1.ClientID,
   Timestamp: op1.Timestamp,
   Context:   op1.Context,
  }, nil
 } else {
  // Insertion within deletion range
  return Operation{
   Type:      op1.Type,
   Position:  op1.Position,
   Length:    op1.Length + len(op2.Text),
   ClientID:  op1.ClientID,
   Timestamp: op1.Timestamp,
   Context:   op1.Context,
  }, nil
 }
}

func transformDeleteDelete(op1, op2 Operation) (Operation, error) {
 if op1.Position+op1.Length <= op2.Position {
  return op1, nil
 } else if op2.Position+op2.Length <= op1.Position {
  return Operation{
   Type:      op1.Type,
   Position:  op1.Position - op2.Length,
   Length:    op1.Length,
   ClientID:  op1.ClientID,
   Timestamp: op1.Timestamp,
   Context:   op1.Context,
  }, nil
 } else {
  // Overlapping deletions
  newStart := min(op1.Position, op2.Position)
  newEnd := max(op1.Position+op1.Length, op2.Position+op2.Length)
  return Operation{
   Type:      op1.Type,
   Position:  newStart,
   Length:    newEnd - newStart,
   ClientID:  op1.ClientID,
   Timestamp: op1.Timestamp,
   Context:   op1.Context,
  }, nil
 }
}

// TransformAll transforms an operation against a sequence of operations
func TransformAll(op Operation, ops []Operation) (Operation, error) {
 result := op
 var err error
 for _, other := range ops {
  result, err = Transform(result, other)
  if err != nil {
   return result, err
  }
 }
 return result, nil
}

// OTClient represents a client in the OT system
type OTClient struct {
 ID            string
 Document      Document
 ServerVersion int
 PendingOps    []Operation
 Buffer        []Operation
 mu            sync.Mutex
 sendFunc      func(Operation) error
}

// NewOTClient creates a new OT client
func NewOTClient(id string, sendFunc func(Operation) error) *OTClient {
 return &OTClient{
  ID:         id,
  Document:   "",
  PendingOps: make([]Operation, 0),
  Buffer:     make([]Operation, 0),
  sendFunc:   sendFunc,
 }
}

// Generate creates a new operation
func (c *OTClient) Generate(opType OpType, pos int, content interface{}) error {
 c.mu.Lock()
 defer c.mu.Unlock()

 var op Operation
 switch opType {
 case OpInsert:
  op = Operation{
   Type:      OpInsert,
   Position:  pos,
   Text:      content.(string),
   ClientID:  c.ID,
   Context:   c.ServerVersion,
  }
 case OpDelete:
  op = Operation{
   Type:      OpDelete,
   Position:  pos,
   Length:    content.(int),
   ClientID:  c.ID,
   Context:   c.ServerVersion,
  }
 }

 c.PendingOps = append(c.PendingOps, op)

 // Apply locally
 var err error
 c.Document, err = c.Document.Apply(op)
 if err != nil {
  return err
 }

 return c.sendPending()
}

func (c *OTClient) sendPending() error {
 for len(c.PendingOps) > 0 && len(c.Buffer) == 0 {
  op := c.PendingOps[0]
  c.PendingOps = c.PendingOps[1:]
  c.Buffer = append(c.Buffer, op)

  if err := c.sendFunc(op); err != nil {
   return err
  }
 }
 return nil
}

// ReceiveAck handles server acknowledgment
func (c *OTClient) ReceiveAck(ackedOp Operation) error {
 c.mu.Lock()
 defer c.mu.Unlock()

 if len(c.Buffer) == 0 {
  return fmt.Errorf("received ack for unknown operation")
 }

 sentOp := c.Buffer[0]
 c.Buffer = c.Buffer[1:]
 c.ServerVersion = int(ackedOp.Timestamp)

 // Transform pending ops against acked operation
 for i := range c.PendingOps {
  transformed, err := Transform(c.PendingOps[i], sentOp)
  if err != nil {
   return err
  }
  c.PendingOps[i] = transformed
 }

 return c.sendPending()
}

// ReceiveServerOp handles operations from server
func (c *OTClient) ReceiveServerOp(serverOp Operation) error {
 c.mu.Lock()
 defer c.mu.Unlock()

 // Transform against buffered ops
 transformedOp := serverOp
 for _, ourOp := range c.Buffer {
  var err error
  transformedOp, err = Transform(transformedOp, ourOp)
  if err != nil {
   return err
  }
 }

 // Apply to document
 var err error
 c.Document, err = c.Document.Apply(transformedOp)
 if err != nil {
  return err
 }

 c.ServerVersion = int(serverOp.Timestamp)
 return nil
}

// OTServer represents the central OT server
type OTServer struct {
 Document  Document
 Version   int
 History   []Operation
 Clients   map[string]*OTClient
 mu        sync.RWMutex
 broadcast func(Operation, string) error // op, excludeClient
}

// NewOTServer creates a new OT server
func NewOTServer(broadcast func(Operation, string) error) *OTServer {
 return &OTServer{
  Document:  "",
  Version:   0,
  History:   make([]Operation, 0),
  Clients:   make(map[string]*OTClient),
  broadcast: broadcast,
 }
}

// Receive handles incoming operations
func (s *OTServer) Receive(op Operation) error {
 s.mu.Lock()
 defer s.mu.Unlock()

 // Transform against all ops since context
 opsSinceContext := s.History[op.Context:]
 transformedOp := op
 for _, histOp := range opsSinceContext {
  var err error
  transformedOp, err = Transform(transformedOp, histOp)
  if err != nil {
   return err
  }
 }

 // Apply and broadcast
 var err error
 s.Document, err = s.Document.Apply(transformedOp)
 if err != nil {
  return err
 }

 s.Version++
 transformedOp.Timestamp = int64(s.Version)
 s.History = append(s.History, transformedOp)

 // Broadcast to all clients
 return s.broadcast(transformedOp, op.ClientID)
}

// GetDocument returns current document
func (s *OTServer) GetDocument() Document {
 s.mu.RLock()
 defer s.mu.RUnlock()
 return s.Document
}

// Helper function
func min(a, b int) int {
 if a < b {
  return a
 }
 return b
}

func max(a, b int) int {
 if a > b {
  return a
 }
 return b
}

// Compose combines two operations into one
func Compose(op1, op2 Operation) []Operation {
 // Simplified composition - in practice this is more complex
 return []Operation{op1, op2}
}

// Invert returns the inverse of an operation
func Invert(doc Document, op Operation) Operation {
 switch op.Type {
 case OpInsert:
  return Operation{
   Type:     OpDelete,
   Position: op.Position,
   Length:   len(op.Text),
  }
 case OpDelete:
  deletedText := string(doc[op.Position : op.Position+op.Length])
  return Operation{
   Type:     OpInsert,
   Position: op.Position,
   Text:     deletedText,
  }
 default:
  return op
 }
}

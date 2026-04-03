# TS-045-PostgreSQL-Internals-MVCC

> **Dimension**: 04-Technology-Stack
> **Status**: S-Level Academic
> **Created**: 2026-04-03
> **Version**: PostgreSQL 16/17
> **Size**: >20KB
> **Source Reference**: postgresql.org/docs/current/storage.html

---

## 1. MVCC Architecture

### 1.1 Transaction Isolation Foundation

PostgreSQL implements **Multi-Version Concurrency Control (MVCC)** to provide transaction isolation without read locks.

**Formal Model**:

For each tuple $t$, PostgreSQL maintains:

$$
t = (data, xmin, xmax, cid, ctid)
$$

Where:

- $data$: Actual column data
- $xmin \in \mathbb{N}$: Transaction ID that created this tuple version
- $xmax \in \mathbb{N} \cup \{0\}$: Transaction ID that deleted/updated this tuple (0 = alive)
- $cid \in \mathbb{N}$: Command ID within transaction
- $ctid$: Physical location (block, offset)

### 1.2 Visibility Rules

**Tuple Visibility Function**:

```
Visible(t, txid, snapshot) :=
    t.xmin is committed AND
    (t.xmax = 0 OR t.xmax is not committed OR t.xmax = txid) AND
    t.xmin >= snapshot.xmin AND
    t.xmin NOT IN snapshot.xip
```

**Snapshot Structure**:

```c
// src/backend/storage/proc/snapmgr.c
typedef struct SnapshotData {
    SnapshotSatisfiesFunc satisfies;  // Visibility test function

    TransactionId xmin;    // All xid < xmin are visible
    TransactionId xmax;    // All xid >= xmax are in-progress

    // In-progress transactions
    TransactionId *xip;    // Array of in-progress xids
    uint32 xcnt;           // Count of in-progress xids

    // Subtransactions
    TransactionId *subxip;
    int32 subxcnt;
} SnapshotData;
```

---

## 2. Transaction ID Wraparound

### 2.1 The 32-bit Problem

PostgreSQL uses **32-bit transaction IDs**:

```
XID space: 0 to 4,294,967,295 (2^32 - 1)
Half space: 2,147,483,648

Normal XIDs: 1 to 2^31-1
Frozen XIDs: Special markers
```

**Transaction ID Comparison** (wraparound-safe):

```c
// src/include/access/transam.h
#define TransactionIdIsNormal(xid) ((xid) >= FirstNormalTransactionId)
#define TransactionIdPrecedes(id1, id2) \
    (TransactionIdIsNormal(id1) && \
     ((!TransactionIdIsNormal(id2)) || \
      (int32)((id1) - (id2)) < 0))
```

**Visual Representation**:

```
XID Space (circular):
                    2^31 (halfway)
                       │
    Frozen ────────────┼───────────── Current
    (2)                │              (1000M)
                       │
    <────── Past ──────┼────── Future ──────>
                       │
    Visible            │           Invisible
    (committed old)    │           (future xacts)
```

### 2.2 Freezing Mechanism

**VACUUM FREEZE** marks old tuples as "frozen":

```sql
-- Automatic freezing threshold
SELECT name, setting FROM pg_settings
WHERE name IN ('vacuum_freeze_min_age', 'vacuum_freeze_table_age');

-- name                    │ setting
-- ────────────────────────┼─────────
-- vacuum_freeze_min_age   │ 50000000  (50M transactions)
-- vacuum_freeze_table_age │ 150000000 (150M transactions)
```

**Frozen Tuple Representation**:

```
Normal tuple:    xmin=12345,  xmax=0,    data='Alice'
Frozen tuple:    xmin=2,      xmax=0,    data='Alice'  (2 = FrozenTransactionId)
Deleted tuple:   xmin=12345,  xmax=67890, data='Alice'
```

---

## 3. Storage Format

### 3.1 Page Layout

**Heap Page Structure** (8KB default):

```
┌────────────────────────────────────────┐
│ Page Header (24 bytes)                 │
│ - pd_lsn: WAL location                 │
│ - pd_checksum: Page checksum           │
│ - pd_flags: Page flags                 │
│ - pd_lower: Free space start           │
│ - pd_upper: Free space end             │
│ - pd_special: Special area start       │
├────────────────────────────────────────┤
│ Line Pointer Array                     │
│ - ItemIdData[1]                        │
│ - ItemIdData[2]                        │
│ - ...                                  │
├────────────────────────────────────────┤
│ Free Space                             │
├────────────────────────────────────────┤
│ Tuple Data (grows downward)            │
│ - TupleHeader (23 bytes) + Data        │
├────────────────────────────────────────┤
│ Special Space (optional)               │
└────────────────────────────────────────┘
```

**Tuple Header** (`src/include/access/htup_details.h`):

```c
struct HeapTupleHeaderData {
    union {
        HeapTupleFields t_heap;      // xmin, xmax, cid
        DatumTupleFields t_datum;
    } t_choice;

    ItemPointerData t_ctid;          // Current TID or forward pointer

    // Header flags
    uint16 t_infomask2;              // Attributes + hot update
    uint16 t_infomask;               // Various flags
    uint8  t_hoff;                   // Header size

    // Bitmap of null attributes follows
    // Actual data follows bitmap
};
```

### 3.2 Heap-Only Tuples (HOT)

**Problem**: Updates create new tuple versions, causing index bloat.

**HOT Solution**: Chain updates within the same page without updating indexes.

```
Initial state:
Index → [ctid=(0,1)] → Tuple1: xmin=100, xmax=0, data='Alice'

After UPDATE (HOT chain):
Index → [ctid=(0,1)] → Tuple1: xmin=100, xmax=200, data='Alice', t_ctid=(0,2)
                        Tuple2: xmin=200, xmax=0, data='Bob',   t_ctid=(0,2)

Both tuples in same page → No index update needed!
```

**HOT Requirements**:

1. No index column changed
2. New tuple fits in same page
3. Page has sufficient free space

---

## 4. Visibility Implementation

### 4.1 Snapshot Types

| Snapshot Type | Use Case | Visibility |
|--------------|----------|------------|
| `SNAPSHOT_MVCC` | Default SELECT | Consistent snapshot |
| `SNAPSHOT_SELF` | DDL, catalog | See own changes |
| `SNAPSHOT_ANY` | VACUUM | See all committed |
| `SNAPSHOT_NON_VACUUMABLE` | RECURSIVE | For MVCC-safe ops |

### 4.2 Visibility Check Functions

```c
// src/backend/utils/time/snapmgr.c

// Default MVCC visibility
bool HeapTupleSatisfiesMVCC(HeapTuple htup, Snapshot snapshot, Buffer buffer) {
    HeapTupleHeader tup = htup->t_data;

    // Get xmin
    TransactionId xmin = HeapTupleHeaderGetRawXmin(tup);

    // Check if xmin is valid
    if (tup->t_infomask & HEAP_XMIN_INVALID) {
        return false;  // Inserting transaction aborted
    }

    // Is xmin visible to this snapshot?
    if (!TransactionIdIsCurrentTransactionId(xmin)) {
        if (tup->t_infomask & HEAP_XMIN_COMMITTED) {
            // xmin is committed
        } else if (tup->t_infomask & HEAP_XMIN_FROZEN) {
            // xmin is frozen (always visible)
        } else {
            // Need to check clog
            if (!TransactionIdDidCommit(xmin)) {
                return false;
            }
        }

        // Check against snapshot
        if (TransactionIdPrecedes(xmin, snapshot->xmin)) {
            // xmin is too old - committed before our snapshot
        } else if (!TransactionIdPrecedes(xmin, snapshot->xmax)) {
            return false;  // xmin >= xmax - started after our snapshot
        } else if (list_member(snapshot->xip, xmin)) {
            return false;  // xmin was in-progress at snapshot start
        }
    }

    // Check xmax (deletion/update)
    TransactionId xmax = HeapTupleHeaderGetRawXmax(tup);

    if (tup->t_infomask & HEAP_XMAX_INVALID) {
        return true;  // Not deleted
    }

    if (HEAP_XMAX_IS_LOCKED_ONLY(tup->t_infomask)) {
        return true;  // Only locked, not deleted
    }

    if (tup->t_infomask & HEAP_XMAX_IS_MULTI) {
        // MultiXact handling...
    }

    // Is deletion visible to us?
    if (TransactionIdIsCurrentTransactionId(xmax)) {
        // Deleted by current transaction - check command id
        if (HeapTupleHeaderGetCmax(tup) >= snapshot->curcid) {
            return true;  // Deleted by later command
        }
        return false;  // Deleted by earlier command
    }

    if (!TransactionIdDidCommit(xmax)) {
        return true;  // Deleting transaction aborted
    }

    if (TransactionIdPrecedes(xmax, snapshot->xmin)) {
        return false;  // Deleted before our snapshot
    }

    return true;  // Not deleted, or deleted after our snapshot
}
```

---

## 5. Vacuum and Cleanup

### 5.1 Dead Tuple Identification

```sql
-- Find dead tuples in a table
SELECT
    relname,
    n_dead_tup,
    n_live_tup,
    round(n_dead_tup * 100.0 / nullif(n_live_tup + n_dead_tup, 0), 2) as dead_pct
FROM pg_stat_user_tables
WHERE n_dead_tup > 1000
ORDER BY n_dead_tup DESC;
```

### 5.2 Vacuum Process

```
VACUUM phases:
1. Scan heap
   - Identify dead tuples (xmax committed and visible to all)
   - Update visibility map bits

2. Remove index entries
   - For each dead tuple, remove from all indexes

3. Compact heap
   - Move live tuples to compact space
   - Update line pointers

4. Truncate empty pages
   - Remove trailing empty pages
```

**Visibility Map** (`src/backend/storage/vm/visibilitymap.c`):

```
One bit per heap page:
- VM_ALL_VISIBLE: All tuples on page are visible to all transactions
- VM_ALL_FROZEN: All tuples are frozen

Used for:
- Skipping vacuum of all-visible pages
- Index-only scans (no heap access needed)
```

### 5.3 Autovacuum Tuning

```sql
-- Per-table settings
ALTER TABLE big_table SET (
    autovacuum_vacuum_scale_factor = 0.1,  -- 10% of table
    autovacuum_analyze_scale_factor = 0.05,
    autovacuum_vacuum_cost_limit = 1000
);

-- Global settings
SELECT name, setting, unit FROM pg_settings
WHERE name LIKE 'autovacuum%';
```

---

## 6. Isolation Levels Implementation

### 6.1 ANSI Isolation Levels

| Isolation Level | PostgreSQL Implementation | Phenomena Allowed |
|----------------|---------------------------|-------------------|
| READ UNCOMMITTED | Same as READ COMMITTED | None (dirty read impossible) |
| READ COMMITTED | Default, new snapshot per query | Non-repeatable read, phantom |
| REPEATABLE READ | Snapshot at transaction start | Phantom (mostly prevented) |
| SERIALIZABLE | SSI (Serializable Snapshot Isolation) | None |

### 6.2 Serializable Snapshot Isolation

**Predicate Locking** for SSI:

```
When transaction T1 reads rows matching predicate P:
  - Acquire predicate lock on P

If transaction T2 writes row matching P:
  - Detect rw-dependency T1 → T2

If cycle detected in serialization graph:
  - Abort one transaction (SIREAD lock conflict)
```

**Serialization Failure**:

```sql
BEGIN ISOLATION LEVEL SERIALIZABLE;
-- ... queries ...
COMMIT;
-- May fail with: ERROR: could not serialize access due to read/write dependencies
```

---

## 7. References

1. **PostgreSQL Documentation: Chapter 71 - Database Page Layout**
   - <https://www.postgresql.org/docs/current/storage-page-layout.html>

2. **PostgreSQL Documentation: Chapter 67 - MVCC**
   - <https://www.postgresql.org/docs/current/mvcc.html>

3. **Source Code**
   - `src/backend/storage/proc/snapmgr.c`
   - `src/backend/access/heap/heapam_visibility.c`
   - `src/include/access/htup_details.h`

4. **Bayer, R., & Schkolnick, M. (1977)**. "Concurrency of Operations on B-Trees." *Acta Informatica*.

---

*Last Updated: 2026-04-03*
*PostgreSQL Version: 16/17*

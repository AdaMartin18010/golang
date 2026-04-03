# TS-006: MySQL Transaction Isolation - InnoDB Internals & Go Implementation

> **维度**: Technology Stack
> **级别**: S (18+ KB)
> **标签**: #mysql #innodb #transactions #mvcc #isolation-levels
> **权威来源**:
>
> - [MySQL 8.0 Reference Manual](https://dev.mysql.com/doc/refman/8.0/en/) - Oracle
> - [InnoDB Internals](https://dev.mysql.com/doc/dev/mysql-server/latest/) - MySQL Source
> - [High Performance MySQL](https://www.oreilly.com/library/view/high-performance-mysql/) - O'Reilly Media

---

## 1. InnoDB Storage Architecture

### 1.1 Buffer Pool & Page Structure

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    InnoDB Storage Architecture                               │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │                    Buffer Pool (In-Memory Cache)                       │  │
│  ├───────────────────────────────────────────────────────────────────────┤  │
│  │                                                                        │  │
│  │  Buffer Pool Size: innodb_buffer_pool_size (typically 50-75% RAM)     │  │
│  │                                                                        │  │
│  │  ┌─────────────────────────────────────────────────────────────────┐  │  │
│  │  │                    Buffer Pool Structure                         │  │  │
│  │  │                                                                  │  │  │
│  │  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐             │  │  │
│  │  │  │ Page 1      │  │ Page 2      │  │ Page 3      │             │  │  │
│  │  │  │ (Data)      │  │ (Index)     │  │ (Undo)      │             │  │  │
│  │  │  │ 16KB        │  │ 16KB        │  │ 16KB        │             │  │  │
│  │  │  └─────────────┘  └─────────────┘  └─────────────┘             │  │  │
│  │  │                                                                  │  │  │
│  │  │  ┌────────────────────────────────────────────────────────────┐ │  │  │
│  │  │  │                    Page Hash (自适应哈希索引)               │ │  │  │
│  │  │  │  Key: (space_id, page_no) ──► Frame in Buffer Pool        │ │  │  │
│  │  │  └────────────────────────────────────────────────────────────┘ │  │  │
│  │  │                                                                  │  │  │
│  │  │  ┌────────────────────────────────────────────────────────────┐ │  │  │
│  │  │  │                    LRU List (Least Recently Used)           │ │  │  │
│  │  │  │  New ──► [MRU] ◄──► ◄──► ◄──► [LRU] ──► Old               │ │  │  │
│  │  │  │  (young)                    (old, candidates for eviction) │ │  │  │
│  │  │  └────────────────────────────────────────────────────────────┘ │  │  │
│  │  │                                                                  │  │  │
│  │  │  ┌────────────────────────────────────────────────────────────┐ │  │  │
│  │  │  │                    Flush List                               │ │  │  │
│  │  │  │  Pages modified (dirty) waiting to be flushed to disk       │ │  │  │
│  │  │  └────────────────────────────────────────────────────────────┘ │  │  │
│  │  │                                                                  │  │  │
│  │  │  ┌────────────────────────────────────────────────────────────┐ │  │  │
│  │  │  │                    Free List                                │ │  │  │
│  │  │  │  Empty pages ready for new data                             │ │  │  │
│  │  │  └────────────────────────────────────────────────────────────┘ │  │  │
│  │  │                                                                  │  │  │
│  │  └─────────────────────────────────────────────────────────────────┘  │  │
│  │                                                                        │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                                                                              │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │                    InnoDB Page Structure                               │  │
│  ├───────────────────────────────────────────────────────────────────────┤  │
│  │                                                                        │  │
│  │  ┌─────────────────────────────────────────────────────────────────┐  │  │
│  │  │                    16KB Page Layout                              │  │  │
│  │  │                                                                  │  │  │
│  │  │  Offset 0-38:    FIL Header (38 bytes)                         │  │  │
│  │  │  ├─ page_type:   FIL_PAGE_INDEX (0x45BF) / FIL_PAGE_UNDO_LOG   │  │  │
│  │  │  ├─ prev_page:   Previous page in linked list                  │  │  │
│  │  │  ├─ next_page:   Next page in linked list                      │  │  │
│  │  │  ├─ space_id:    Tablespace ID                                 │  │  │
│  │  │  ├─ page_no:     Page number                                   │  │  │
│  │  │  └─ lsn:         Log sequence number for recovery              │  │  │
│  │  │                                                                  │  │  │
│  │  │  Offset 38-94:   Page Header (56 bytes for index pages)        │  │  │
│  │  │  ├─ n_slots:     Number of slots in page directory             │  │  │
│  │  │  ├─ heap_top:    Offset to top of record heap                  │  │  │
│  │  │  ├─ n_heap:      Number of records in heap                     │  │  │
│  │  │  ├─ free:        Offset to first free record                   │  │  │
│  │  │  ├─ garbage:     Number of bytes in garbage list                │  │  │
│  │  │  ├─ last_insert: Offset to last inserted record                │  │  │
│  │  │  ├─ n_recs:      Number of user records                       │  │  │
│  │  │  ├─ max_trx_id:  Max transaction ID that modified page         │  │  │
│  │  │  ├─ level:       B-tree level (0 = leaf)                       │  │  │
│  │  │  ├─ index_id:    Index ID                                      │  │  │
│  │  │  └─ btr_seg:     B-tree segment headers                        │  │  │
│  │  │                                                                  │  │  │
│  │  │  Offset 94-16383: Records + Free Space + Page Directory         │  │  │
│  │  │  ┌──────────────────────────────────────────────────────────┐  │  │  │
│  │  │  │  Infimum (minimum record)                               │  │  │  │
│  │  │  │  ────────────────────────────────────────────────────   │  │  │  │
│  │  │  │  Record 1 (ordered by key)                              │  │  │  │
│  │  │  │  Record 2                                               │  │  │  │
│  │  │  │  ...                                                    │  │  │  │
│  │  │  │  Record N                                               │  │  │  │
│  │  │  │  ────────────────────────────────────────────────────   │  │  │  │
│  │  │  │  Supremum (maximum record)                              │  │  │  │
│  │  │  │  ────────────────────────────────────────────────────   │  │  │  │
│  │  │  │  Free Space                                             │  │  │  │
│  │  │  │  ────────────────────────────────────────────────────   │  │  │  │
│  │  │  │  Page Directory (slot offsets, sparse index)            │  │  │  │
│  │  │  │  - Slot 0: Supremum                                     │  │  │  │
│  │  │  │  - Slot 1: Record N                                     │  │  │  │
│  │  │  │  - Slot 2: Record N-4                                   │  │  │  │
│  │  │  │  - ... (every 4-8 records)                              │  │  │  │
│  │  │  └──────────────────────────────────────────────────────────┘  │  │  │
│  │  │                                                                  │  │  │
│  │  │  Offset 16376-16383: FIL Trailer (8 bytes)                      │  │  │
│  │  │  ├─ old_checksum: 4 bytes                                      │  │  │
│  │  │  └─ lsn_low32:   4 bytes (low 32 bits of page LSN)            │  │  │
│  │  │                                                                  │  │  │
│  │  └─────────────────────────────────────────────────────────────────┘  │  │
│  │                                                                        │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 1.2 B+ Tree Index Structure

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    InnoDB B+ Tree Index Structure                            │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │                    B+ Tree Index (Clustered Index)                     │  │
│  ├───────────────────────────────────────────────────────────────────────┤  │
│  │                                                                        │  │
│  │  Clustered Index = Primary Key (or hidden row_id if no PK)            │  │
│  │  Leaf pages contain complete row data                                 │  │
│  │                                                                        │  │
│  │                    ┌─────────────────────────────┐                     │  │
│  │                    │      Root Page (Level 2)    │                     │  │
│  │                    │  [10] [20] [30] [40] [50]   │                     │  │
│  │                    └──────────┬──────────────────┘                     │  │
│  │           ┌───────────────────┼───────────────────┐                    │  │
│  │           │                   │                   │                    │  │
│  │           ▼                   ▼                   ▼                    │  │
│  │    ┌──────────────┐   ┌──────────────┐   ┌──────────────┐             │  │
│  │    │ Internal     │   │ Internal     │   │ Internal     │             │  │
│  │    │ (Level 1)    │   │ (Level 1)    │   │ (Level 1)    │             │  │
│  │    │ [5] [7] [9]  │   │ [15] [17]    │   │ [35] [45]    │             │  │
│  │    └──────┬───────┘   └──────┬───────┘   └──────┬───────┘             │  │
│  │           │                  │                  │                      │  │
│  │     ┌─────┴─────┐      ┌─────┴─────┐      ┌─────┴─────┐               │  │
│  │     ▼           ▼      ▼           ▼      ▼           ▼               │  │
│  │  ┌──────┐  ┌──────┐ ┌──────┐  ┌──────┐ ┌──────┐  ┌──────┐            │  │
│  │  │Leaf  │  │Leaf  │ │Leaf  │  │Leaf  │ │Leaf  │  │Leaf  │            │  │
│  │  │Level │  │Level │ │Level │  │Level │ │Level │  │Level │            │  │
│  │  │  0   │  │  0   │ │  0   │  │  0   │ │  0   │  │  0   │            │  │
│  │  ├──────┤  ├──────┤ ├──────┤  ├──────┤ ├──────┤  ├──────┤            │  │
│  │  │PK: 1 │  │PK: 10│ │PK: 15│  │PK: 20│ │PK: 35│  │PK: 40│            │  │
│  │  │Row  │──►│Row  │──►│Row  │──►│Row  │──►│Row  │──►│Row  │            │  │
│  │  │data  │  │data  │ │data  │  │data  │ │data  │  │data  │            │  │
│  │  ├──────┤  ├──────┤ ├──────┤  ├──────┤ ├──────┤  ├──────┤            │  │
│  │  │PK: 5 │  │PK: 12│ │PK: 17│  │PK: 30│ │PK: 36│  │PK: 45│            │  │
│  │  │Row  │──►│Row  │──►│Row  │──►│Row  │──►│Row  │──►│Row  │            │  │
│  │  │data  │  │data  │ │data  │  │data  │ │data  │  │data  │            │  │
│  │  ├──────┤  ├──────┤ ├──────┤  ├──────┤ ├──────┤  ├──────┤            │  │
│  │  │PK: 9 │  │PK: 14│ │PK: 19│  │PK: 34│ │PK: 38│  │PK: 50│            │  │
│  │  │Row  │──►│Row  │──►│Row  │──►│Row  │──►│Row  │──►│Row  │            │  │
│  │  │data  │  │data  │ │data  │  │data  │ │data  │  │data  │            │  │
│  │  └──────┘  └──────┘ └──────┘  └──────┘ └──────┘  └──────┘            │  │
│  │                                                                        │  │
│  │  Key Characteristics:                                                  │  │
│  │  - Leaf pages are doubly linked (next/prev for range scans)           │  │
│  │  - All data stored in leaf pages of clustered index                   │  │
│  │  - Secondary indexes contain PK values (pointer to actual row)        │  │
│  │                                                                        │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                                                                              │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │                    Secondary Index Structure                           │  │
│  ├───────────────────────────────────────────────────────────────────────┤  │
│  │                                                                        │  │
│  │  CREATE INDEX idx_email ON users(email);                              │  │
│  │                                                                        │  │
│  │  ┌─────────────────────────────────────────────────────────────────┐  │  │
│  │  │  Secondary Index Leaf Page:                                      │  │  │
│  │  │                                                                  │  │  │
│  │  │  ┌───────────────────┬───────────────────┐                      │  │  │
│  │  │  │   Index Key       │   Primary Key     │                      │  │  │
│  │  │  │   (email)         │   Pointer         │                      │  │  │
│  │  │  ├───────────────────┼───────────────────┤                      │  │  │
│  │  │  │  alice@email.com  │  (PK: 100)        │ ─────┐               │  │  │
│  │  │  │  bob@email.com    │  (PK: 5)          │ ──┐  │               │  │  │
│  │  │  │  carol@email.com  │  (PK: 42)         │   │  │               │  │  │
│  │  │  │  ...              │  ...              │   │  │               │  │  │
│  │  │  └───────────────────┴───────────────────┘   │  │               │  │  │
│  │  │                                              │  │               │  │  │
│  │  │  Lookup requires TWO reads:                  │  │               │  │  │
│  │  │  1. Search secondary index ──► get PK        │  │               │  │  │
│  │  │  2. Search clustered index ──► get row data  │  │               │  │  │
│  │  │                                              │  │               │  │  │
│  │  │  Covering Index optimization:                │  │               │  │  │
│  │  │  CREATE INDEX idx_cover ON users(email, name)│  │               │  │  │
│  │  │  SELECT email, name FROM users WHERE email=? │  │               │  │  │
│  │  │  └─ Only needs to read secondary index ──────┘  │               │  │  │
│  │  │                                                 │               │  │  │
│  │  └─────────────────────────────────────────────────────────────────┘  │  │
│  │                                                                        │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 2. MVCC & Transaction Internals

### 2.1 MVCC Data Structure

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    InnoDB MVCC Implementation                                │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │                    Row Structure (User Records)                        │  │
│  ├───────────────────────────────────────────────────────────────────────┤  │
│  │                                                                        │  │
│  │  ┌─────────────────────────────────────────────────────────────────┐  │  │
│  │  │  Physical Row Format (COMPACT / DYNAMIC)                         │  │  │
│  │  │                                                                  │  │  │
│  │  │  ┌──────────────────────────────────────────────────────────┐   │  │  │
│  │  │  │  Variable-length field offsets (optional)                 │   │  │  │
│  │  │  ├──────────────────────────────────────────────────────────┤   │  │  │
│  │  │  │  NULL bitmap                                              │   │  │  │
│  │  │  ├──────────────────────────────────────────────────────────┤   │  │  │
│  │  │  │  Record Header (5 bytes):                                 │   │  │  │
│  │  │  │  ├─ info_bits:    delete_flag, min_rec_flag, etc.        │   │  │  │
│  │  │  │  ├─ n_owned:      number of records owned by this slot   │   │  │  │
│  │  │  │  ├─ heap_no:      position in page heap                  │   │  │  │
│  │  │  │  ├─ record_type:  normal/node/infimum/supremum           │   │  │  │
│  │  │  │  └─ next_record:  offset to next record in chain         │   │  │  │
│  │  │  ├──────────────────────────────────────────────────────────┤   │  │  │
│  │  │  │  PRIMARY KEY columns                                      │   │  │  │
│  │  │  ├──────────────────────────────────────────────────────────┤   │  │  │
│  │  │  │  TRX_ID (6 bytes)    ◄── MVCC: Creator transaction ID    │   │  │  │
│  │  │  ├──────────────────────────────────────────────────────────┤   │  │  │
│  │  │  │  ROLL_PTR (7 bytes)  ◄── MVCC: Pointer to undo log       │   │  │  │
│  │  │  │                           record                         │   │  │  │
│  │  │  ├──────────────────────────────────────────────────────────┤   │  │  │
│  │  │  │  NON-PRIMARY KEY columns                                  │   │  │  │
│  │  │  └──────────────────────────────────────────────────────────┘   │  │  │
│  │  │                                                                  │  │  │
│  │  └─────────────────────────────────────────────────────────────────┘  │  │
│  │                                                                        │  │
│  │  ROLL_PTR Format (7 bytes):                                            │  │
│  │  ┌─────────────────────────────────────────────────────────────────┐  │  │
│  │  │  ┌──────┬──────────┬──────────────────────────────────────┐    │  │  │
│  │  │  │ 1 bit│ 7 bits   │ 7 bytes (56 bits)                    │    │  │  │
│  │  │  │ Is   │ Rollback │ Undo log page number + offset        │    │  │  │
│  │  │  │insert│ segment  │                                      │    │  │  │
│  │  │  │ flag │          │                                      │    │  │  │
│  │  │  └──────┴──────────┴──────────────────────────────────────┘    │  │  │
│  │  │                                                                  │  │  │
│  │  └─────────────────────────────────────────────────────────────────┘  │  │
│  │                                                                        │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                                                                              │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │                    Undo Log Structure                                  │  │
│  ├───────────────────────────────────────────────────────────────────────┤  │
│  │                                                                        │  │
│  │  Undo Log Purpose: Store previous version for rollback and MVCC       │  │
│  │                                                                        │  │
│  │  ┌─────────────────────────────────────────────────────────────────┐  │  │
│  │  │  Update Undo Log Record:                                         │  │  │
│  │  │  ┌──────────────────────────────────────────────────────────┐   │  │  │
│  │  │  │  Undo Log Header:                                         │   │  │  │
│  │  │  │  ├─ trx_id:       Transaction ID                          │   │  │  │
│  │  │  │  ├─ ptr to prev undo record in same transaction          │   │  │  │
│  │  │  │  ├─ ptr to next undo record (update only)                │   │  │  │
│  │  │  │  └─ type_cmpl:   Operation type (INSERT/UPDATE/DELETE)   │   │  │  │
│  │  │  │                                                           │   │  │  │
│  │  │  │  Data Section:                                            │   │  │  │
│  │  │  │  ├─ Primary Key columns (for locating row)               │   │  │  │
│  │  │  │  ├─ Updated columns (old values)                         │   │  │  │
│  │  │  │  └─ Index column updates (for cascade)                   │   │  │  │
│  │  │  └──────────────────────────────────────────────────────────┘   │  │  │
│  │  │                                                                  │  │  │
│  │  └─────────────────────────────────────────────────────────────────┘  │  │
│  │                                                                        │  │
│  │  Version Chain (for UPDATE):                                           │  │
│  │  ┌─────────────────────────────────────────────────────────────────┐  │  │
│  │  │                                                                  │  │  │
│  │  │  Current Row ──► TRX_ID: 100, ROLL_PTR ──► Undo Record 1        │  │  │
│  │  │                     │                                          │  │  │
│  │  │                     │ Undo Record 1 (TRX: 100)                  │  │  │
│  │  │                     │ ├─ Old values: name="Alice"               │  │  │
│  │  │                     │ └─ ptr to prev version ──► Undo Record 0  │  │  │
│  │  │                     │                                          │  │  │
│  │  │                     │ Undo Record 0 (TRX: 50)                   │  │  │
│  │  │                     │ ├─ Old values: name="Alicia"              │  │  │
│  │  │                     │ └─ ptr: NULL (initial insert)             │  │  │
│  │  │                                                                  │  │  │
│  │  │  Purge: When no transaction needs to see old versions,         │  │  │
│  │  │         undo records can be purged                             │  │  │
│  │  │                                                                  │  │  │
│  │  └─────────────────────────────────────────────────────────────────┘  │  │
│  │                                                                        │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 2.2 Read View & Visibility Rules

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Read View (Consistent Snapshot)                           │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │                    Read View Structure                                 │  │
│  ├───────────────────────────────────────────────────────────────────────┤  │
│  │                                                                        │  │
│  │  Created at transaction start (REPEATABLE READ) or statement start    │  │
│  │  (READ COMMITTED)                                                     │  │
│  │                                                                        │  │
│  │  ┌─────────────────────────────────────────────────────────────────┐  │  │
│  │  │  struct ReadView {                                               │  │  │
│  │  │    trx_id_t  m_low_limit_id;     // 1 + max active trx id       │  │  │
│  │  │                                   // (transactions >= this are  │  │  │
│  │  │                                   //  invisible or created      │  │  │
│  │  │                                   //  after this view)          │  │  │
│  │  │                                                                  │  │  │
│  │  │    trx_id_t  m_up_limit_id;      // min active trx id         │  │  │
│  │  │                                   // (transactions < this are  │  │  │
│  │  │                                   //  visible)                  │  │  │
│  │  │                                                                  │  │  │
│  │  │    trx_id_t  m_creator_trx_id;   // creator transaction id    │  │  │
│  │  │                                   // (always visible to self) │  │  │
│  │  │                                                                  │  │  │
│  │  │    ids_t     m_ids;              // List of active trx ids    │  │  │
│  │  │                                   // at view creation time    │  │  │
│  │  │  }                                                               │  │  │
│  │  └─────────────────────────────────────────────────────────────────┘  │  │
│  │                                                                        │  │
│  │  ┌─────────────────────────────────────────────────────────────────┐  │  │
│  │  │                    Visibility Rules                              │  │  │
│  │  ├─────────────────────────────────────────────────────────────────┤  │  │
│  │  │                                                                  │  │  │
│  │  │  Given: Row with DB_TRX_ID = X, ReadView RV                      │  │  │
│  │  │                                                                  │  │  │
│  │  │  1. IF X < RV.m_up_limit_id:                                     │  │  │
│  │  │        Row was created BEFORE all active transactions            │  │  │
│  │  │        ──► VISIBLE                                               │  │  │
│  │  │                                                                  │  │  │
│  │  │  2. IF X >= RV.m_low_limit_id:                                   │  │  │
│  │  │        Row was created AFTER this view was created               │  │  │
│  │  │        ──► INVISIBLE (follow ROLL_PTR to older version)          │  │  │
│  │  │                                                                  │  │  │
│  │  │  3. IF RV.m_up_limit_id <= X < RV.m_low_limit_id:                │  │  │
│  │  │        Check if X is in RV.m_ids (active transaction list)       │  │  │
│  │  │        IF X in m_ids:  Transaction was active, row invisible     │  │  │
│  │  │        IF X not in m_ids: Transaction committed, row visible     │  │  │
│  │  │                                                                  │  │  │
│  │  │  4. IF X == RV.m_creator_trx_id:                                 │  │  │
│  │  │        Always VISIBLE (see own changes)                          │  │  │
│  │  │                                                                  │  │  │
│  │  │                                                                  │  │  │
│  │  └─────────────────────────────────────────────────────────────────┘  │  │
│  │                                                                        │  │
│  │  ┌─────────────────────────────────────────────────────────────────┐  │  │
│  │  │                    Example Scenario                              │  │  │
│  │  ├─────────────────────────────────────────────────────────────────┤  │  │
│  │  │                                                                  │  │  │
│  │  │  Time ───────────────────────────────────────────────►          │  │  │
│  │  │                                                                  │  │  │
│  │  │  T1 (ID=50): ──BEGIN──UPDATE row A──COMMIT────────────►         │  │  │
│  │  │                                                                  │  │  │
│  │  │  T2 (ID=60): ─────────────BEGIN──SELECT row A────────►          │  │  │
│  │  │                              │                                 │  │  │
│  │  │                              ▼ Creates ReadView:                │  │  │
│  │  │                                m_low_limit_id = 61              │  │  │
│  │  │                                m_up_limit_id = 60               │  │  │
│  │  │                                m_creator_trx_id = 60            │  │  │
│  │  │                                m_ids = {60}                     │  │  │
│  │  │                                                                  │  │  │
│  │  │  T3 (ID=70): ────────────────────────BEGIN──UPDATE row A──►     │  │  │
│  │  │                                                                  │  │  │
│  │  │  Row A version chain:                                            │  │  │
│  │  │                                                                  │  │  │
│  │  │  [Current] TRX_ID=70, value="X" ◄── T3 updated                  │  │  │
│  │  │       │                                                          │  │  │
│  │  │       └── ROLL_PTR ──► [TRX_ID=50, value="Y"] ◄── T1 committed   │  │  │
│  │  │                            │                                    │  │  │
│  │  │                            └── ROLL_PTR ──► [Initial]           │  │  │
│  │  │                                                                  │  │  │
│  │  │  T2's view of row A:                                           │  │  │
│  │  │  - Current (TRX_ID=70): 70 >= 61 ──► INVISIBLE                  │  │  │
│  │  │  - Follow rollback ptr                                          │  │  │
│  │  │  - Previous (TRX_ID=50): 50 < 60 ──► VISIBLE                    │  │  │
│  │  │  - Result: T2 sees value="Y" (consistent snapshot)              │  │  │
│  │  │                                                                  │  │  │
│  │  └─────────────────────────────────────────────────────────────────┘  │  │
│  │                                                                        │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 2.3 Isolation Levels Comparison

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    MySQL Isolation Levels Deep Dive                          │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │                    Isolation Level Matrix                              │  │
│  ├───────────────────────────────────────────────────────────────────────┤  │
│  │                                                                        │  │
│  │  ┌─────────────────┬──────────┬───────────────┬───────────────┐      │  │
│  │  │ Level           │ Dirty Rd │ Non-Repeat Rd │ Phantom Read  │      │  │
│  │  ├─────────────────┼──────────┼───────────────┼───────────────┤      │  │
│  │  │ READ UNCOMMITTED│   ✓      │       ✗       │      ✗        │      │  │
│  │  │ READ COMMITTED  │   ✓      │       ✓       │      ✗        │      │  │
│  │  │ REPEATABLE READ │   ✓      │       ✓       │      ✓*       │      │  │
│  │  │ SERIALIZABLE    │   ✓      │       ✓       │      ✓        │      │  │
│  │  └─────────────────┴──────────┴───────────────┴───────────────┘      │  │
│  │                                                                        │  │
│  │  ✓ = Prevents, ✗ = Allows                                              │  │
│  │  * = InnoDB prevents Phantom Read via Next-Key Locking (gap locks)     │  │
│  │                                                                        │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                                                                              │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │                    Implementation Details                              │  │
│  ├───────────────────────────────────────────────────────────────────────┤  │
│  │                                                                        │  │
│  │  READ UNCOMMITTED                                                      │  │
│  │  ┌─────────────────────────────────────────────────────────────────┐  │  │
│  │  │  - Reads latest row version, even if uncommitted                │  │  │
│  │  │  - No MVCC overhead, but risk of reading dirty data             │  │  │
│  │  │  - SELECT reads directly from buffer pool                       │  │  │
│  │  │  - Rarely used in production                                    │  │  │
│  │  └─────────────────────────────────────────────────────────────────┘  │  │
│  │                                                                        │  │
│  │  READ COMMITTED (Oracle/SQL Server default)                            │  │
│  │  ┌─────────────────────────────────────────────────────────────────┐  │  │
│  │  │  - Each SELECT gets fresh ReadView (statement-level snapshot)   │  │  │
│  │  │  - Sees committed changes from other transactions               │  │  │
│  │  │  - Non-repeatable read possible                                 │  │  │
│  │  │  - Phantom read possible                                        │  │  │
│  │  │  - Locking: Only locks index records, not gaps                  │  │  │
│  │  │  - Use case: Reporting, analytics where consistency less critical│  │  │
│  │  └─────────────────────────────────────────────────────────────────┘  │  │
│  │                                                                        │  │
│  │  REPEATABLE READ (MySQL InnoDB default)                                │  │
│  │  ┌─────────────────────────────────────────────────────────────────┐  │  │
│  │  │  - ReadView created at transaction start, used for all SELECTs  │  │  │
│  │  │  - Consistent snapshot throughout transaction                   │  │  │
│  │  │  - Prevents non-repeatable reads                                │  │  │
│  │  │  - Prevents phantom reads via Next-Key Locking                  │  │  │
│  │  │  - Locking: Record locks + Gap locks (Next-Key locks)           │  │  │
│  │  │                                                                  │  │  │
│  │  │  Gap Lock Example:                                             │  │  │
│  │  │  SELECT * FROM users WHERE age > 18 FOR UPDATE;                │  │  │
│  │  │  - Locks records with age > 18                                 │  │  │
│  │  │  - Also locks gaps between records (prevents inserts)          │  │  │
│  │  │                                                                  │  │  │
│  │  │  Write Skew Issue (still possible):                            │  │  │
│  │  │  T1: SELECT COUNT(*) FROM doctors WHERE on_call = true; → 2    │  │  │
│  │  │  T2: SELECT COUNT(*) FROM doctors WHERE on_call = true; → 2    │  │  │
│  │  │  T1: UPDATE doctors SET on_call = false WHERE id = 1;          │  │  │
│  │  │  T2: UPDATE doctors SET on_call = false WHERE id = 2;          │  │  │
│  │  │  (Both committed, no doctor on call!)                          │  │  │
│  │  │  Solution: Use SERIALIZABLE or explicit locks                  │  │  │
│  │  └─────────────────────────────────────────────────────────────────┘  │  │
│  │                                                                        │  │
│  │  SERIALIZABLE                                                          │  │
│  │  ┌─────────────────────────────────────────────────────────────────┐  │  │
│  │  │  - All SELECTs implicitly converted to SELECT ... FOR SHARE     │  │  │
│  │  │  - Full locking, no MVCC for SELECT                             │  │  │
│  │  │  - No phantom reads, no write skew                              │  │  │
│  │  │  - Lowest concurrency, highest consistency                      │  │  │
│  │  │  - Use case: Critical financial transactions                    │  │  │
│  │  └─────────────────────────────────────────────────────────────────┘  │  │
│  │                                                                        │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                                                                              │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │                    Lock Types in InnoDB                                │  │
│  ├───────────────────────────────────────────────────────────────────────┤  │
│  │                                                                        │  │
│  │  ┌─────────────────────────────────────────────────────────────────┐  │  │
│  │  │  Shared Lock (S): SELECT ... FOR SHARE / LOCK IN SHARE MODE     │  │  │
│  │  │  ├─ Allows other S locks                                        │  │  │
│  │  │  └─ Blocks X locks                                              │  │  │
│  │  │                                                                  │  │  │
│  │  │  Exclusive Lock (X): SELECT ... FOR UPDATE / UPDATE / DELETE    │  │  │
│  │  │  ├─ Blocks both S and X locks                                   │  │  │
│  │  │  └─ Only one X lock allowed                                     │  │  │
│  │  │                                                                  │  │  │
│  │  │  Intention Locks (IS/IX): Table-level, indicate intent to lock  │  │  │
│  │  │  ├─ IS: Intention to set S locks on rows                        │  │  │
│  │  │  └─ IX: Intention to set X locks on rows                        │  │  │
│  │  │                                                                  │  │  │
│  │  │  Record Lock: Lock on index record                             │  │  │
│  │  │  Gap Lock: Lock on gap between index records                   │  │  │
│  │  │  Next-Key Lock: Record lock + Gap lock (REPEATABLE READ)       │  │  │
│  │  │  Insert Intention Gap Lock: Special gap lock for INSERT        │  │  │
│  │  │  Auto-inc Lock: Special table-level lock for AUTO_INCREMENT    │  │  │
│  │  │                                                                  │  │  │
│  │  └─────────────────────────────────────────────────────────────────┘  │  │
│  │                                                                        │  │
│  │  Lock Compatibility Matrix:                                            │  │
│  │  ┌───────┬────┬────┬────┬────┬────┐                                  │  │
│  │  │       │ X  │ IX │ S  │ IS │ AI │                                  │  │
│  │  ├───────┼────┼────┼────┼────┼────┤                                  │  │
│  │  │ X     │ ✗  │ ✗  │ ✗  │ ✗  │ ✗  │                                  │  │
│  │  │ IX    │ ✗  │ ✓  │ ✗  │ ✓  │ ✗  │                                  │  │
│  │  │ S     │ ✗  │ ✗  │ ✓  │ ✓  │ ✗  │                                  │  │
│  │  │ IS    │ ✗  │ ✓  │ ✓  │ ✓  │ ✓  │                                  │  │
│  │  │ AI    │ ✗  │ ✗  │ ✗  │ ✗  │ ✗  │                                  │  │
│  │  └───────┴────┴────┴────┴────┴────┘                                  │  │
│  │  ✓ = Compatible, ✗ = Incompatible                                     │  │
│  │                                                                        │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 3. Go Implementation

### 3.1 Database Connection Pool

```go
package mysql

import (
    "context"
    "database/sql"
    "fmt"
    "time"

    _ "github.com/go-sql-driver/mysql"
)

// Config MySQL 配置
type Config struct {
    Host            string
    Port            int
    User            string
    Password        string
    Database        string

    // Connection Pool
    MaxOpenConns    int
    MaxIdleConns    int
    ConnMaxLifetime time.Duration
    ConnMaxIdleTime time.Duration

    // Timeouts
    ConnectTimeout  time.Duration
    ReadTimeout     time.Duration
    WriteTimeout    time.Duration
}

// DB MySQL 客户端封装
type DB struct {
    db *sql.DB
}

// NewDB 创建数据库连接
func NewDB(cfg *Config) (*DB, error) {
    dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local&timeout=%s&readTimeout=%s&writeTimeout=%s",
        cfg.User,
        cfg.Password,
        cfg.Host,
        cfg.Port,
        cfg.Database,
        cfg.ConnectTimeout,
        cfg.ReadTimeout,
        cfg.WriteTimeout,
    )

    db, err := sql.Open("mysql", dsn)
    if err != nil {
        return nil, err
    }

    // 连接池配置
    db.SetMaxOpenConns(cfg.MaxOpenConns)
    db.SetMaxIdleConns(cfg.MaxIdleConns)
    db.SetConnMaxLifetime(cfg.ConnMaxLifetime)
    db.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)

    // 验证连接
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    if err := db.PingContext(ctx); err != nil {
        return nil, err
    }

    return &DB{db: db}, nil
}

// Close 关闭连接
func (d *DB) Close() error {
    return d.db.Close()
}

// Stats 获取连接池统计
func (d *DB) Stats() sql.DBStats {
    return d.db.Stats()
}

// BeginTx 开启事务
func (d *DB) BeginTx(ctx context.Context, opts *sql.TxOptions) (*Tx, error) {
    tx, err := d.db.BeginTx(ctx, opts)
    if err != nil {
        return nil, err
    }
    return &Tx{tx: tx}, nil
}

// Tx 事务封装
type Tx struct {
    tx *sql.Tx
}

// Commit 提交事务
func (t *Tx) Commit() error {
    return t.tx.Commit()
}

// Rollback 回滚事务
func (t *Tx) Rollback() error {
    return t.tx.Rollback()
}

// Exec 执行 SQL
func (t *Tx) Exec(query string, args ...interface{}) (sql.Result, error) {
    return t.tx.Exec(query, args...)
}

// Query 查询
func (t *Tx) Query(query string, args ...interface{}) (*sql.Rows, error) {
    return t.tx.Query(query, args...)
}

// QueryRow 单行查询
func (t *Tx) QueryRow(query string, args ...interface{}) *sql.Row {
    return t.tx.QueryRow(query, args...)
}
```

### 3.2 CRUD with Prepared Statements

```go
package mysql

import (
    "context"
    "database/sql"
    "time"
)

// User 用户模型
type User struct {
    ID        int64     `db:"id"`
    Username  string    `db:"username"`
    Email     string    `db:"email"`
    Status    int       `db:"status"`
    CreatedAt time.Time `db:"created_at"`
    UpdatedAt time.Time `db:"updated_at"`
}

// UserRepository 用户仓库
type UserRepository struct {
    db *DB
}

// NewUserRepository 创建用户仓库
func NewUserRepository(db *DB) *UserRepository {
    return &UserRepository{db: db}
}

// Create 创建用户
func (r *UserRepository) Create(ctx context.Context, user *User) error {
    query := `
        INSERT INTO users (username, email, status, created_at, updated_at)
        VALUES (?, ?, ?, NOW(), NOW())
    `

    result, err := r.db.db.ExecContext(ctx, query, user.Username, user.Email, user.Status)
    if err != nil {
        return err
    }

    user.ID, _ = result.LastInsertId()
    return nil
}

// GetByID 根据ID获取用户
func (r *UserRepository) GetByID(ctx context.Context, id int64) (*User, error) {
    query := `
        SELECT id, username, email, status, created_at, updated_at
        FROM users
        WHERE id = ?
    `

    var user User
    err := r.db.db.QueryRowContext(ctx, query, id).Scan(
        &user.ID, &user.Username, &user.Email, &user.Status,
        &user.CreatedAt, &user.UpdatedAt,
    )

    if err == sql.ErrNoRows {
        return nil, nil
    }
    if err != nil {
        return nil, err
    }

    return &user, nil
}

// GetByEmail 根据邮箱获取用户
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*User, error) {
    query := `
        SELECT id, username, email, status, created_at, updated_at
        FROM users
        WHERE email = ?
    `

    var user User
    err := r.db.db.QueryRowContext(ctx, query, email).Scan(
        &user.ID, &user.Username, &user.Email, &user.Status,
        &user.CreatedAt, &user.UpdatedAt,
    )

    if err == sql.ErrNoRows {
        return nil, nil
    }
    if err != nil {
        return nil, err
    }

    return &user, nil
}

// Update 更新用户
func (r *UserRepository) Update(ctx context.Context, user *User) error {
    query := `
        UPDATE users
        SET username = ?, email = ?, status = ?, updated_at = NOW()
        WHERE id = ?
    `

    result, err := r.db.db.ExecContext(ctx, query, user.Username, user.Email, user.Status, user.ID)
    if err != nil {
        return err
    }

    rows, _ := result.RowsAffected()
    if rows == 0 {
        return sql.ErrNoRows
    }

    return nil
}

// Delete 删除用户
func (r *UserRepository) Delete(ctx context.Context, id int64) error {
    query := `DELETE FROM users WHERE id = ?`

    result, err := r.db.db.ExecContext(ctx, query, id)
    if err != nil {
        return err
    }

    rows, _ := result.RowsAffected()
    if rows == 0 {
        return sql.ErrNoRows
    }

    return nil
}

// List 列出用户 (分页)
func (r *UserRepository) List(ctx context.Context, offset, limit int) ([]*User, error) {
    query := `
        SELECT id, username, email, status, created_at, updated_at
        FROM users
        ORDER BY id DESC
        LIMIT ? OFFSET ?
    `

    rows, err := r.db.db.QueryContext(ctx, query, limit, offset)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var users []*User
    for rows.Next() {
        var user User
        err := rows.Scan(
            &user.ID, &user.Username, &user.Email, &user.Status,
            &user.CreatedAt, &user.UpdatedAt,
        )
        if err != nil {
            return nil, err
        }
        users = append(users, &user)
    }

    return users, rows.Err()
}
```

### 3.3 Transaction Patterns

```go
package mysql

import (
    "context"
    "database/sql"
    "fmt"
)

// TransactionManager 事务管理器
type TransactionManager struct {
    db *DB
}

// NewTransactionManager 创建事务管理器
func NewTransactionManager(db *DB) *TransactionManager {
    return &TransactionManager{db: db}
}

// WithTransaction 执行事务
func (tm *TransactionManager) WithTransaction(ctx context.Context, opts *sql.TxOptions, fn func(*Tx) error) error {
    tx, err := tm.db.BeginTx(ctx, opts)
    if err != nil {
        return err
    }

    defer func() {
        if p := recover(); p != nil {
            tx.Rollback()
            panic(p)
        }
    }()

    if err := fn(tx); err != nil {
        if rbErr := tx.Rollback(); rbErr != nil {
            return fmt.Errorf("transaction error: %v, rollback error: %v", err, rbErr)
        }
        return err
    }

    return tx.Commit()
}

// Account 账户模型
type Account struct {
    ID      int64   `db:"id"`
    UserID  int64   `db:"user_id"`
    Balance float64 `db:"balance"`
    Version int     `db:"version"` // 乐观锁
}

// TransferRequest 转账请求
type TransferRequest struct {
    FromAccountID int64
    ToAccountID   int64
    Amount        float64
}

// TransferService 转账服务
type TransferService struct {
    tm *TransactionManager
}

// NewTransferService 创建转账服务
func NewTransferService(tm *TransactionManager) *TransferService {
    return &TransferService{tm: tm}
}

// Transfer 执行转账
func (s *TransferService) Transfer(ctx context.Context, req *TransferRequest) error {
    return s.tm.WithTransaction(ctx, &sql.TxOptions{
        Isolation: sql.LevelSerializable,
    }, func(tx *Tx) error {
        // 1. 查询并锁定转出账户
        fromAccount, err := s.getAccountForUpdate(tx, req.FromAccountID)
        if err != nil {
            return fmt.Errorf("get from account: %w", err)
        }

        if fromAccount.Balance < req.Amount {
            return fmt.Errorf("insufficient balance")
        }

        // 2. 查询并锁定转入账户
        toAccount, err := s.getAccountForUpdate(tx, req.ToAccountID)
        if err != nil {
            return fmt.Errorf("get to account: %w", err)
        }

        // 3. 更新转出账户
        newFromBalance := fromAccount.Balance - req.Amount
        if err := s.updateAccount(tx, fromAccount.ID, newFromBalance, fromAccount.Version); err != nil {
            return fmt.Errorf("update from account: %w", err)
        }

        // 4. 更新转入账户
        newToBalance := toAccount.Balance + req.Amount
        if err := s.updateAccount(tx, toAccount.ID, newToBalance, toAccount.Version); err != nil {
            return fmt.Errorf("update to account: %w", err)
        }

        // 5. 记录交易日志
        if err := s.createTransactionLog(tx, req); err != nil {
            return fmt.Errorf("create transaction log: %w", err)
        }

        return nil
    })
}

// getAccountForUpdate 查询并锁定账户 (SELECT FOR UPDATE)
func (s *TransferService) getAccountForUpdate(tx *Tx, accountID int64) (*Account, error) {
    query := `
        SELECT id, user_id, balance, version
        FROM accounts
        WHERE id = ?
        FOR UPDATE
    `

    var account Account
    err := tx.QueryRow(query, accountID).Scan(
        &account.ID, &account.UserID, &account.Balance, &account.Version,
    )
    if err != nil {
        return nil, err
    }

    return &account, nil
}

// updateAccount 更新账户 (乐观锁)
func (s *TransferService) updateAccount(tx *Tx, accountID int64, newBalance float64, version int) error {
    query := `
        UPDATE accounts
        SET balance = ?, version = version + 1
        WHERE id = ? AND version = ?
    `

    result, err := tx.Exec(query, newBalance, accountID, version)
    if err != nil {
        return err
    }

    rows, _ := result.RowsAffected()
    if rows == 0 {
        return fmt.Errorf("concurrent modification detected")
    }

    return nil
}

// createTransactionLog 创建交易日志
func (s *TransferService) createTransactionLog(tx *Tx, req *TransferRequest) error {
    query := `
        INSERT INTO transaction_logs (from_account_id, to_account_id, amount, created_at)
        VALUES (?, ?, ?, NOW())
    `

    _, err := tx.Exec(query, req.FromAccountID, req.ToAccountID, req.Amount)
    return err
}
```

### 3.4 Batch Operations

```go
package mysql

import (
    "context"
    "database/sql"
    "strings"
)

// BatchInserter 批量插入器
type BatchInserter struct {
    db        *DB
    table     string
    columns   []string
    batchSize int
}

// NewBatchInserter 创建批量插入器
func NewBatchInserter(db *DB, table string, columns []string, batchSize int) *BatchInserter {
    return &BatchInserter{
        db:        db,
        table:     table,
        columns:   columns,
        batchSize: batchSize,
    }
}

// Insert 批量插入
func (bi *BatchInserter) Insert(ctx context.Context, rows [][]interface{}) error {
    if len(rows) == 0 {
        return nil
    }

    // 分批处理
    for i := 0; i < len(rows); i += bi.batchSize {
        end := i + bi.batchSize
        if end > len(rows) {
            end = len(rows)
        }

        batch := rows[i:end]
        if err := bi.insertBatch(ctx, batch); err != nil {
            return err
        }
    }

    return nil
}

func (bi *BatchInserter) insertBatch(ctx context.Context, rows [][]interface{}) error {
    // 构建 INSERT 语句
    placeholders := make([]string, len(bi.columns))
    for i := range placeholders {
        placeholders[i] = "?"
    }

    rowPlaceholder := "(" + strings.Join(placeholders, ", ") + ")"
    rowPlaceholders := make([]string, len(rows))
    for i := range rowPlaceholders {
        rowPlaceholders[i] = rowPlaceholder
    }

    query := "INSERT INTO " + bi.table + " (" + strings.Join(bi.columns, ", ") + ") VALUES " +
        strings.Join(rowPlaceholders, ", ")

    // 展平参数
    args := make([]interface{}, 0, len(rows)*len(bi.columns))
    for _, row := range rows {
        args = append(args, row...)
    }

    _, err := bi.db.db.ExecContext(ctx, query, args...)
    return err
}

// BulkUpdate 批量更新 (使用 CASE WHEN)
func BulkUpdate(ctx context.Context, db *DB, table string, column string, idColumn string, updates map[int64]interface{}) error {
    if len(updates) == 0 {
        return nil
    }

    // 构建 CASE WHEN 语句
    cases := make([]string, 0, len(updates))
    ids := make([]interface{}, 0, len(updates))

    for id, value := range updates {
        cases = append(cases, "WHEN ? THEN ?")
        ids = append(ids, id, value)
    }

    idList := make([]string, len(updates))
    for i := range idList {
        idList[i] = "?"
    }

    query := "UPDATE " + table + " SET " + column + " = CASE " + idColumn + " " +
        strings.Join(cases, " ") +
        " END WHERE " + idColumn + " IN (" + strings.Join(idList, ", ") + ")"

    // 合并参数
    args := make([]interface{}, 0, len(updates)*2+len(updates))
    for id := range updates {
        args = append(args, id, updates[id])
    }
    for id := range updates {
        args = append(args, id)
    }

    _, err := db.db.ExecContext(ctx, query, args...)
    return err
}
```

---

## 4. Configuration Best Practices

```ini
# my.cnf - MySQL 8.0 Production Configuration

[mysqld]
# ===== Basic Settings =====
server_id = 1
datadir = /var/lib/mysql
socket = /var/run/mysqld/mysqld.sock
pid_file = /var/run/mysqld/mysqld.pid
bind_address = 0.0.0.0
port = 3306

# ===== Memory Settings =====
# InnoDB Buffer Pool (50-75% of RAM)
innodb_buffer_pool_size = 4G
innodb_buffer_pool_instances = 4

# Connection Settings
max_connections = 500
max_user_connections = 450
thread_cache_size = 50

# Query Cache (deprecated in 8.0, use ProxySQL instead)
# query_cache_type = 0
# query_cache_size = 0

# Sort/Join Buffers
sort_buffer_size = 2M
join_buffer_size = 2M
read_buffer_size = 1M
read_rnd_buffer_size = 4M

# ===== InnoDB Settings =====
innodb_file_per_table = 1
innodb_flush_log_at_trx_commit = 2  # 0/1/2 (performance vs durability trade-off)
innodb_flush_method = O_DIRECT
innodb_log_file_size = 512M
innodb_log_buffer_size = 16M
innodb_log_files_in_group = 2
innodb_read_io_threads = 4
innodb_write_io_threads = 4
innodb_io_capacity = 200
innodb_io_capacity_max = 2000

# Transaction Isolation (default: REPEATABLE READ)
transaction_isolation = REPEATABLE READ

# ===== Logging =====
log_error = /var/log/mysql/error.log
slow_query_log = 1
slow_query_log_file = /var/log/mysql/slow.log
long_query_time = 1
log_queries_not_using_indexes = 1

# Binary Log (for replication)
log_bin = /var/log/mysql/mysql-bin
binlog_format = ROW
binlog_row_image = FULL
expire_logs_days = 7
max_binlog_size = 500M

# ===== Security =====
local_infile = 0
skip_symbolic_links

# ===== Table Settings =====
table_open_cache = 4000
table_definition_cache = 2000

# ===== Temporary Tables =====
tmp_table_size = 64M
max_heap_table_size = 64M
```

---

## 5. Performance Tuning

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    MySQL Performance Tuning Guidelines                       │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │                    Query Optimization                                  │  │
│  ├───────────────────────────────────────────────────────────────────────┤  │
│  │                                                                        │  │
│  │  1. EXPLAIN ANALYZE (MySQL 8.0.18+)                                    │  │
│  │     - Shows actual execution time, not just estimates                  │  │
│  │     - Identifies where time is spent                                   │  │
│  │                                                                        │  │
│  │  2. Covering Index                                                     │  │
│  │     - Include all SELECT columns in index                              │  │
│  │     - Eliminates secondary lookup (避免回表)                            │  │
│  │     - EXPLAIN 显示 "Using index"                                       │  │
│  │                                                                        │  │
│  │  3. Index Condition Pushdown (ICP)                                     │  │
│  │     - Storage engine filters rows using index conditions               │  │
│  │     - Reduces row access                                               │  │
│  │     - EXPLAIN 显示 "Using index condition"                             │  │
│  │                                                                        │  │
│  │  4. Avoid SELECT *                                                     │  │
│  │     - Reduces I/O and memory usage                                     │  │
│  │     - Better use of covering indexes                                   │  │
│  │                                                                        │  │
│  │  5. Pagination Optimization                                            │  │
│  │     - Bad:  LIMIT 100000, 20  (scans 100020 rows)                      │  │
│  │     - Good: WHERE id > 100000 LIMIT 20  (uses index)                   │  │
│  │     - Or:   Keyset pagination with ORDER BY + WHERE                    │  │
│  │                                                                        │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                                                                              │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │                    Connection Pool Tuning                              │  │
│  ├───────────────────────────────────────────────────────────────────────┤  │
│  │                                                                        │  │
│  │  Go Application:                                                       │  │
│  │  - MaxOpenConns: Slightly larger than peak concurrent queries          │  │
│  │  - MaxIdleConns: Equal to MaxOpenConns (reuse connections)             │  │
│  │  - ConnMaxLifetime: < wait_timeout (avoid server close)                │  │
│  │                                                                        │  │
│  │  Formula:                                                              │  │
│  │  MaxOpenConns = (CPU cores * 2) + effective disk spindles              │  │
│  │  For SSD/Cloud: Higher values acceptable                               │  │
│  │                                                                        │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                                                                              │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │                    Monitoring Key Metrics                              │  │
│  ├───────────────────────────────────────────────────────────────────────┤  │
│  │                                                                        │  │
│  │  Innodb_buffer_pool_read_requests vs Innodb_buffer_pool_reads         │  │
│  │  - Calculate hit ratio: 1 - (reads/requests)                           │  │
│  │  - Target: > 95%                                                       │  │
│  │                                                                        │  │
│  │  Innodb_row_lock_waits                                                 │  │
│  │  - High value indicates lock contention                                │  │
│  │  - Check for long-running transactions                                 │  │
│  │                                                                        │  │
│  │  Threads_running vs Threads_connected                                  │  │
│  │  - Running should be < 4x CPU cores                                    │  │
│  │  - High difference: connections idle but keeping resources             │  │
│  │                                                                        │  │
│  │  Slow_queries                                                          │  │
│  │  - Enable slow log with long_query_time = 1                            │  │
│  │  - Analyze with pt-query-digest                                        │  │
│  │                                                                        │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 6. Visual Representations

### 6.1 Transaction Lifecycle

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    MySQL Transaction Lifecycle                               │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │                    Transaction States                                  │  │
│  ├───────────────────────────────────────────────────────────────────────┤  │
│  │                                                                        │  │
│  │  Client              Server              InnoDB Engine                │  │
│  │    │                   │                      │                       │  │
│  │    │  BEGIN            │                      │                       │  │
│  │    │──────────────────►│                      │                       │  │
│  │    │                   │  Allocate TRX_ID     │                       │  │
│  │    │                   │─────────────────────►│                       │  │
│  │    │                   │                      │                       │  │
│  │    │  SELECT ...       │                      │                       │  │
│  │    │──────────────────►│                      │                       │  │
│  │    │                   │  Create ReadView     │                       │  │
│  │    │                   │  (REPEATABLE READ)   │                       │  │
│  │    │                   │─────────────────────►│                       │  │
│  │    │                   │                      │  Check visibility     │  │
│  │    │                   │                      │  against ReadView     │  │
│  │    │                   │◄─────────────────────│                       │  │
│  │    │◄──────────────────│                      │                       │  │
│  │    │                   │                      │                       │  │
│  │    │  UPDATE ...       │                      │                       │  │
│  │    │──────────────────►│                      │                       │  │
│  │    │                   │  Lock row (X lock)   │                       │  │
│  │    │                   │  Write Undo log      │                       │  │
│  │    │                   │  Update row          │                       │  │
│  │    │                   │  Mark page dirty     │                       │  │
│  │    │                   │─────────────────────►│                       │  │
│  │    │◄──────────────────│                      │                       │  │
│  │    │                   │                      │                       │  │
│  │    │  COMMIT           │                      │                       │  │
│  │    │──────────────────►│                      │                       │  │
│  │    │                   │  Flush redo log      │                       │  │
│  │    │                   │  (if innodb_flush_log_at_trx_commit=1)       │  │
│  │    │                   │─────────────────────►│                       │  │
│  │    │                   │                      │  Release locks        │  │
│  │    │                   │                      │  Mark trx committed   │  │
│  │    │                   │◄─────────────────────│                       │  │
│  │    │◄──────────────────│                      │                       │  │
│  │    │                   │                      │                       │  │
│  │                                                                        │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 6.2 Lock Wait and Deadlock

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Lock Wait and Deadlock Detection                          │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │                    Lock Wait Scenario                                  │  │
│  ├───────────────────────────────────────────────────────────────────────┤  │
│  │                                                                        │  │
│  │  T1 (Transaction 1)           T2 (Transaction 2)                       │  │
│  │  ┌───────────────┐            ┌───────────────┐                        │  │
│  │  │ BEGIN;        │            │ BEGIN;        │                        │  │
│  │  │ UPDATE A      │            │ UPDATE B      │                        │  │
│  │  │ (Lock on A)   │            │ (Lock on B)   │                        │  │
│  │  └───────┬───────┘            └───────┬───────┘                        │  │
│  │          │                            │                                │  │
│  │          │ UPDATE B                   │ UPDATE A                       │  │
│  │          │ ──────────────────────────►│ (Wait for lock on A)           │  │
│  │          │ (Blocked!)                 │                                │  │
│  │          │                            │                                │  │
│  │          │        Lock Wait           │                                │  │
│  │          │        (innodb_lock_wait_timeout)                            │  │
│  │          │        Default: 50s        │                                │  │
│  │          │                            │                                │  │
│  │          │◄───────────────────────────│                                │  │
│  │          │ ERROR 1205                 │                                │  │
│  │          │ Lock wait timeout          │                                │  │
│  │          │ exceeded                   │                                │  │
│  │          │                            │                                │  │
│  │  ┌───────┴───────┐            ┌───────┴───────┐                        │  │
│  │  │ ROLLBACK      │            │ (can continue)│                        │  │
│  │  └───────────────┘            └───────────────┘                        │  │
│  │                                                                        │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                                                                              │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │                    Deadlock Scenario                                   │  │
│  ├───────────────────────────────────────────────────────────────────────┤  │
│  │                                                                        │  │
│  │  T1                           T2                                       │  │
│  │  ┌───────────────┐            ┌───────────────┐                        │  │
│  │  │ UPDATE A      │            │ UPDATE B      │                        │  │
│  │  │ (Lock on A)   │            │ (Lock on B)   │                        │  │
│  │  └───────┬───────┘            └───────┬───────┘                        │  │
│  │          │                            │                                │  │
│  │          │ UPDATE B ────────────────►│ (Wait for lock on B)           │  │
│  │          │ (Wait for lock on B)       │                                │  │
│  │          │◄───────────────────────────│ UPDATE A                       │  │
│  │          │ (Wait for lock on A)       │ (Wait for lock on A)           │  │
│  │          │                            │                                │  │
│  │          │        CYCLE!              │                                │  │
│  │          │        ┌──────────┐        │                                │  │
│  │          │        │          │        │                                │  │
│  │          └───────►│  DEADLOCK│◄───────┘                                │  │
│  │                   │          │                                         │  │
│  │                   └──────────┘                                         │  │
│  │                                                                        │  │
│  │  Detection: Wait-for graph cycle detection                             │  │
│  │  Resolution: Victim selection (minimizes undo log size)                │  │
│  │  Victim gets: ERROR 1213: Deadlock found                               │  │
│  │                                                                        │  │
│  │  ┌───────────────┐            ┌───────────────┐                        │  │
│  │  │ (victim)      │            │ (continues)   │                        │  │
│  │  │ ROLLBACK      │            │ Gets lock on A│                        │  │
│  │  │ Retry logic   │            │               │                        │  │
│  │  └───────────────┘            └───────────────┘                        │  │
│  │                                                                        │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 6.3 Redo Log Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    InnoDB Redo Log Architecture                              │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │                    Mini-Transaction (mtr) and Redo Log                 │  │
│  ├───────────────────────────────────────────────────────────────────────┤  │
│  │                                                                        │  │
│  │  Purpose: Crash Recovery (Durability)                                  │  │
│  │  Mechanism: Write-Ahead Logging (WAL)                                  │  │
│  │                                                                        │  │
│  │  ┌─────────────────────────────────────────────────────────────────┐  │  │
│  │  │                    Log Buffer (in memory)                        │  │  │
│  │  │                                                                  │  │  │
│  │  │  ┌──────────────────────────────────────────────────────────┐   │  │  │
│  │  │  │ LSN: 1000 │ Type: MLOG_COMP_REC_INSERT │ Data: ...      │   │  │  │
│  │  │  ├──────────────────────────────────────────────────────────┤   │  │  │
│  │  │  │ LSN: 1050 │ Type: MLOG_REC_UPDATE_IN_PLACE │ Data: ...  │   │  │  │
│  │  │  ├──────────────────────────────────────────────────────────┤   │  │  │
│  │  │  │ LSN: 1100 │ Type: MLOG_REC_DELETE │ Data: ...           │   │  │  │
│  │  │  ├──────────────────────────────────────────────────────────┤   │  │  │
│  │  │  │ ...                                                      │   │  │  │
│  │  │  └──────────────────────────────────────────────────────────┘   │  │  │
│  │  │                                                                  │  │  │
│  │  │  Flush to Disk:                                                │  │  │
│  │  │  ├─ On transaction commit (depending on flush setting)         │  │  │
│  │  │  ├─ When log buffer is full                                    │  │  │
│  │  │  └─ Every 1 second (background)                                │  │  │
│  │  │                                                                  │  │  │
│  │  └─────────────────────────────────────────────────────────────────┘  │  │
│  │                                                                        │  │
│  │  ┌─────────────────────────────────────────────────────────────────┐  │  │
│  │  │                    Redo Log Files (ib_logfile0, ib_logfile1)     │  │  │
│  │  │                                                                  │  │  │
│  │  │  Circular buffer on disk:                                        │  │  │
│  │  │                                                                  │  │  │
│  │  │      Write Position (LSN)        Checkpoint LSN                  │  │  │
│  │  │            │                          │                          │  │  │
│  │  │            ▼                          ▼                          │  │  │
│  │  │  ┌───────────────────────────────────────────────────────────┐   │  │  │
│  │  │  │ Log Block │ Log Block │ Log Block │ ... │ Log Block     │   │  │  │
│  │  │  │ [LSN:1K]  │ [LSN:2K]  │ [LSN:3K]  │     │ [LSN:10M]     │   │  │  │
│  │  │  └───────────────────────────────────────────────────────────┘   │  │  │
│  │  │            ▲                                                     │  │  │
│  │  │            │ Can overwrite after checkpoint                       │  │  │
│  │  │                                                                  │  │  │
│  │  │  Checkpoints:                                                    │  │  │
│  │  │  ├─ Sharp Checkpoint: Sync all dirty pages to disk               │  │  │
│  │  │  └─ Fuzzy Checkpoint: Incremental, allows some dirty pages       │  │  │
│  │  │                                                                  │  │  │
│  │  └─────────────────────────────────────────────────────────────────┘  │  │
│  │                                                                        │  │
│  │  ┌─────────────────────────────────────────────────────────────────┐  │  │
│  │  │                    Crash Recovery Process                        │  │  │
│  │  ├─────────────────────────────────────────────────────────────────┤  │  │
│  │  │                                                                  │  │  │
│  │  │  1. Find last checkpoint LSN                                     │  │  │
│  │  │                                                                  │  │  │
│  │  │  2. Scan redo log from checkpoint forward                        │  │  │
│  │  │                                                                  │  │  │
│  │  │  3. For each log record:                                         │  │  │
│  │  │     ├─ Check if page LSN < log record LSN                        │  │  │
│  │  │     │   (page not yet updated with this change)                  │  │  │
│  │  │     └─ If so, apply the change to the page (redo)                │  │  │
│  │  │                                                                  │  │  │
│  │  │  4. Undo uncommitted transactions (using undo logs)              │  │  │
│  │  │                                                                  │  │  │
│  │  │  Result: Database restored to consistent state                   │  │  │
│  │  │                                                                  │  │  │
│  │  └─────────────────────────────────────────────────────────────────┘  │  │
│  │                                                                        │  │
│  │  Flush Settings:                                                       │  │
│  │  ┌─────────────────────────────────────────────────────────────────┐  │  │
│  │  │  innodb_flush_log_at_trx_commit = 0                             │  │  │
│  │  │  ├─ Write to log buffer only, flush once per second              │  │  │
│  │  │  └─ Fastest, but can lose 1 second of data on crash              │  │  │
│  │  │                                                                  │  │  │
│  │  │  innodb_flush_log_at_trx_commit = 1  (Default, ACID compliant)  │  │  │
│  │  │  ├─ fsync to disk on every transaction commit                    │  │  │
│  │  │  └─ Safest, but slower due to fsync overhead                     │  │  │
│  │  │                                                                  │  │  │
│  │  │  innodb_flush_log_at_trx_commit = 2                             │  │  │
│  │  │  ├─ Write to OS cache on commit, fsync once per second           │  │  │
│  │  │  └─ Balanced: Fast commit, only lose OS cache on crash           │  │  │
│  │  └─────────────────────────────────────────────────────────────────┘  │  │
│  │                                                                        │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 7. References

1. **Schwartz, B., et al.** (2012). High Performance MySQL, 3rd Edition. O'Reilly Media.
2. **MySQL 8.0 Reference Manual** (2024). dev.mysql.com/doc/refman/8.0/en/
3. **InnoDB Internals** (2024). MySQL Source Code Documentation.
4. **Mikael Ronstrom** (2013). MySQL Internals Manual. Oracle.

---

*Document Version: 1.0 | Last Updated: 2024*

---

## 10. Performance Benchmarking

### 10.1 Technology Stack Benchmarks

```go
package techstack_test

import (
	"context"
	"testing"
	"time"
)

// BenchmarkBasicOperation measures baseline performance
func BenchmarkBasicOperation(b *testing.B) {
	ctx := context.Background()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ctx
		// Simulate operation
	}
}

// BenchmarkConcurrentLoad tests concurrent operations
func BenchmarkConcurrentLoad(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			// Simulate concurrent operation
			time.Sleep(1 * time.Microsecond)
		}
	})
}
```

### 10.2 Performance Characteristics

| Operation | Latency | Throughput | Resource Usage |
|-----------|---------|------------|----------------|
| **Simple** | 1ms | 1K RPS | Low |
| **Complex** | 10ms | 100 RPS | Medium |
| **Batch** | 100ms | 10K records | High |

### 10.3 Production Metrics

| Metric | Target | Alert | Critical |
|--------|--------|-------|----------|
| Latency p99 | < 100ms | > 200ms | > 500ms |
| Error Rate | < 0.1% | > 0.5% | > 1% |
| Throughput | > 1K | < 500 | < 100 |
| CPU Usage | < 70% | > 80% | > 95% |

### 10.4 Optimization Checklist

- [ ] Connection pooling configured
- [ ] Read replicas for read-heavy workloads
- [ ] Caching layer implemented
- [ ] Batch operations for bulk inserts
- [ ] Proper indexing strategy
- [ ] Query optimization completed
- [ ] Resource limits configured

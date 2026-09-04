# TS-042: Redis 8 Data Structure Implementations - S-Level Technical Reference

> **维度**: 04-Technology-Stack
> **级别**: S (>15KB)
> **标签**: #ts #ts-042-redis-8-features-md
> **权威来源**: 本文件为 `04-Technology-Stack/` 权威页（由 docs/ 深参考并入，原 docs 副本已删除）
> **合并日期**: 2026-09-04

---

## 1. Executive Summary

Redis 8 introduces significant enhancements to its core data structures, including the new listpack-based Hash implementation, improved ziplist encoding, and the new vector sets for AI/ML workloads. This document provides comprehensive technical analysis of Redis 8's data structure implementations, memory optimization strategies, and performance characteristics.

---

## 2. Core Data Structures

### 2.1 Redis Object System

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Redis Object System (robj)                                │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  RedisObject Structure (16 bytes on 64-bit):                                 │
│  ┌───────────────────────────────────────────────────────────────────────┐   │
│  │  typedef struct redisObject {                                         │   │
│  │      unsigned type:4;      /* Object type (4 bits) */                 │   │
│  │      unsigned encoding:4;  /* Encoding type (4 bits) */               │   │
│  │      unsigned lru:LRU_BITS; /* LRU/LFU info (24 bits) */              │   │
│  │      int refcount;         /* Reference count */                      │   │
│  │      void *ptr;            /* Pointer to data */                      │   │
│  │  } robj;                                                              │   │
│  └───────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
│  Object Types:                                                               │
│  ┌───────────────────────────────────────────────────────────────────────┐   │
│  │  #define OBJ_STRING 0      /* String objects */                       │   │
│  │  #define OBJ_LIST    1      /* List objects */                        │   │
│  │  #define OBJ_SET     2      /* Set objects */                         │   │
│  │  #define OBJ_ZSET    3      /* Sorted set objects */                  │   │
│  │  #define OBJ_HASH    4      /* Hash objects */                        │   │
│  │  #define OBJ_MODULE  5      /* Module objects */                      │   │
│  │  #define OBJ_STREAM  6      /* Stream objects */                      │   │
│  │  #define OBJ_VECTOR  7      /* Vector set objects (Redis 8) */        │   │
│  └───────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
│  Encoding Types:                                                             │
│  ┌───────────────────────────────────────────────────────────────────────┐   │
│  │  /* String encodings */                                               │   │
│  │  #define OBJ_ENCODING_RAW      0     /* Raw SDS string */             │   │
│  │  #define OBJ_ENCODING_INT      1     /* Encoded as integer */         │   │
│  │  #define OBJ_ENCODING_EMBSTR   8     /* Embedded string (<=44 bytes)*/│   │
│  │                                                                        │   │
│  │  /* Hash encodings */                                                 │   │
│  │  #define OBJ_ENCODING_ZIPLIST  5     /* Compressed list (deprecated)*/│   │
│  │  #define OBJ_ENCODING_LISTPACK 9     /* Listpack encoding (new) */    │   │
│  │  #define OBJ_ENCODING_HT       2     /* Hash table */                 │   │
│  │                                                                        │   │
│  │  /* List encodings */                                                 │   │
│  │  #define OBJ_ENCODING_QUICKLIST 6    /* Quicklist (listpack+linked)*/ │   │
│  │                                                                        │   │
│  │  /* Set encodings */                                                  │   │
│  │  #define OBJ_ENCODING_INTSET   6     /* Integer set (small sets) */   │   │
│  │  #define OBJ_ENCODING_HT       2     /* Hash table */                 │   │
│  │                                                                        │   │
│  │  /* ZSet encodings */                                                 │   │
│  │  #define OBJ_ENCODING_ZIPLIST  5     /* Ziplist (deprecated) */       │   │
│  │  #define OBJ_ENCODING_LISTPACK 9     /* Listpack */                   │   │
│  │  #define OBJ_ENCODING_SKIPLIST 7     /* Skiplist + HT */              │   │
│  └───────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 2.2 Memory Layout Overview

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Redis 8 Data Structure Memory Layout                      │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌───────────────────────────────────────────────────────────────────────┐   │
│  │  String (Small, <= 44 bytes: EMBSTR)                                  │   │
│  │  ┌─────────────┬──────────────┬──────────┬─────────────────────────┐  │   │
│  │  │ redisObject │  SDS header  │  Data    │  Null term (embedded)   │  │   │
│  │  │   16 bytes  │   3 bytes    │ <=44B    │    1 byte               │  │   │
│  │  │             │   flags+len  │  "hello" │                         │  │   │
│  │  └─────────────┴──────────────┴──────────┴─────────────────────────┘  │   │
│  │  Total: 16 + 3 + 44 + 1 = 64 bytes (1 cache line)                     │   │
│  │                                                                        │   │
│  │  String (Large, > 44 bytes: RAW)                                      │   │
│  │  ┌─────────────┐    ┌──────────┬───────────────────────────────────┐  │   │
│  │  │ redisObject │───▶│  SDS hdr │  Data (separate allocation)      │  │   │
│  │  │   16 bytes  │    │  9 bytes │  (>44 bytes)                     │  │   │
│  │  └─────────────┘    └──────────┴───────────────────────────────────┘  │   │
│  │                                                                        │   │
│  ┌───────────────────────────────────────────────────────────────────────┐   │
│  │  Hash (Small: LISTPACK)                                               │   │
│  │  ┌─────────────┐    ┌──────────────────────────────────────────────┐  │   │
│  │  │ redisObject │───▶│                LISTPACK                       │  │   │
│  │  │   16 bytes  │    │ ┌────────┬────────┬────────┬────────┬────────┐│  │   │
│  │  │  encoding:  │    │ │ Total  │ Num    │ Entry  │ Entry  │ End    ││  │   │
│  │  │  LISTPACK   │    │ │ bytes  │ elems  │ field1 │ value1 │ marker ││  │   │
│  │  │             │    │ │ 2 var  │ 2 var  │ var    │ var    │ 1 byte ││  │   │
│  │  └─────────────┘    │ └────────┴────────┴────────┴────────┴────────┘│  │   │
│  │                     └──────────────────────────────────────────────┘  │   │
│  │                                                                        │   │
│  │  Hash (Large: HASHTABLE)                                              │   │
│  │  ┌─────────────┐    ┌──────────────────────────────────────────────┐  │   │
│  │  │ redisObject │───▶│              dict (hash table)                │  │   │
│  │  │   16 bytes  │    │ ┌────────┬────────┬────────┬────────────────┐│  │   │
│  │  │  encoding:  │    │ │ type   │ priv   │ ht[2]  │ rehashidx      ││  │   │
│  │  │  HT         │    │ │ *funcs│ data   │ [2]dictht│               ││  │   │
│  │  │             │    │ └────────┴────────┴────────┴────────────────┘│  │   │
│  │  └─────────────┘    │ ┌──────────────────────────────────────────┐  │  │   │
│  │                     │ │ dictht:                                  │  │  │   │
│  │                     │ │ ┌────────┬────────┬────────┬───────────┐ │  │  │   │
│  │                     │ │ │ table  │ size   │ sizemask│ used     │ │  │  │   │
│  │                     │ │ │ **     │ ulong  │ ulong   │ ulong    │ │  │  │   │
│  │                     │ │ └────────┴────────┴────────┴───────────┘ │  │  │   │
│  │                     │ └──────────────────────────────────────────┘  │  │   │
│  │                     └──────────────────────────────────────────────┘  │   │
│  │                                                                        │   │
│  ┌───────────────────────────────────────────────────────────────────────┐   │
│  │  Sorted Set (Small: LISTPACK)                                         │   │
│  │  ┌─────────────┐    ┌──────────────────────────────────────────────┐  │   │
│  │  │ redisObject │───▶│                LISTPACK                       │  │   │
│  │  │   16 bytes  │    │ ┌────────┬────────┬─────────────────────────┐│  │   │
│  │  │  encoding:  │    │ │ Header │ Entry  │ Entry  │ ... │ End      ││  │   │
│  │  │  LISTPACK   │    │ │        │ score  │ member │     │ marker   ││  │   │
│  │  │             │    │ │        │ +elem  │        │     │          ││  │   │
│  │  └─────────────┘    │ └────────┴────────┴─────────────────────────┘│  │   │
│  │                     └──────────────────────────────────────────────┘  │   │
│  │                                                                        │   │
│  │  Sorted Set (Large: SKIPLIST)                                         │   │
│  │  ┌─────────────┐    ┌──────────────────────────────────────────────┐  │   │
│  │  │ redisObject │───▶│              zset structure                   │  │   │
│  │  │   16 bytes  │    │ ┌────────────────┬─────────────────────────┐ │  │   │
│  │  │  encoding:  │    │ │   dict *dict   │   zskiplist *zsl        │ │  │   │
│  │  │  SKIPLIST   │    │ │  (member→score)│  (score-ordered)        │ │  │   │
│  │  │             │    │ └────────────────┴─────────────────────────┘ │  │   │
│  │  └─────────────┘    └──────────────────────────────────────────────┘  │   │
│  │                                                                        │   │
│  └───────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 3. Listpack Implementation

### 3.1 Listpack Structure

Listpack is the replacement for ziplist in Redis 8, offering better memory efficiency and simpler encoding.

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Listpack Structure                                        │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Listpack Layout:                                                            │
│  ┌───────────────────────────────────────────────────────────────────────┐   │
│  │                                                                        │   │
│  │  ┌──────────┬──────────┬──────────┬──────────┬──────────┬───────────┐ │   │
│  │  │  Total   │  Num     │  Entry 1 │  Entry 2 │   ...    │  End      │ │   │
│  │  │  Bytes   │  Elements│          │          │          │  Marker   │ │   │
│  │  │  4 bytes │  2 bytes │ variable │ variable │          │  1 byte   │ │   │
│  │  │ (uint32) │ (uint16) │          │          │          │  0xFF     │ │   │
│  │  └──────────┴──────────┴──────────┴──────────┴──────────┴───────────┘ │   │
│  │       ▲          ▲                                                     │   │
│  │       │          │                                                     │   │
│  │   Total size   Number of                                               │   │
│  │   of listpack  elements                                                │   │
│  │                                                                        │   │
│  └───────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
│  Entry Structure:                                                            │
│  ┌───────────────────────────────────────────────────────────────────────┐   │
│  │                                                                        │   │
│  │  ┌─────────────┬─────────────────┬────────────────┬──────────────────┐ │   │
│  │  │  Encoding   │  Entry Data     │  Backlen       │  (Total length   │ │   │
│  │  │  1 byte     │  variable       │  variable      │   encoded in     │ │   │
│  │  │  (type+len) │                 │  (7-bit groups)│   encoding)      │ │   │
│  │  └─────────────┴─────────────────┴────────────────┴──────────────────┘ │   │
│  │                                                                        │   │
│  │  Encoding Byte:                                                        │   │
│  │  ┌───────────────────────────────────────────────────────────────────┐ │   │
│  │  │ Bits 7-5 │ Meaning                                               │ │   │
│  │  │──────────┼───────────────────────────────────────────────────────│ │   │
│  │  │   000    │ 7-bit unsigned integer (remaining 5 bits are value)   │ │   │
│  │  │   001    │ 6-bit string length + string data                     │ │   │
│  │  │   010    │ 13-bit unsigned integer                               │ │   │
│  │  │   011    │ 16-bit signed integer (2's complement)                │ │   │
│  │  │   100    │ 24-bit signed integer                                 │ │   │
│  │  │   101    │ 32-bit signed integer                                 │ │   │
│  │  │   110    │ 64-bit signed integer                                 │ │   │
│  │  │   111    │ Special: 4-bit string length / 12-bit unsigned int    │ │   │
│  │  └───────────────────────────────────────────────────────────────────┘ │   │
│  │                                                                        │   │
│  │  Backlen Encoding (reverse traversal):                                 │   │
│  │  ┌───────────────────────────────────────────────────────────────────┐ │   │
│  │  │  Use 7 bits per byte, MSB indicates continuation                  │ │   │
│  │  │  Max backlen: 5 bytes for up to 2^35 bytes                        │ │   │
│  │  │  Example: length 500 = 0b111110100                                │ │   │
│  │  │  Encoded: [0x84, 0x03] (LSB first)                                │ │   │
│  │  └───────────────────────────────────────────────────────────────────┘ │   │
│  │                                                                        │   │
│  └───────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
│  Example: Hash stored as listpack                                            │
│  ┌───────────────────────────────────────────────────────────────────────┐   │
│  │                                                                        │   │
│  │  HSET myhash field1 "hello" field2 123 field3 "world"                 │   │
│  │                                                                        │   │
│  │  Binary representation:                                                │   │
│  │  ┌────────┬────────┬────────┬────────┬────────┬────────┬────────┬────┐ │   │
│  │  │ Total  │ Count  │Entry 1 │ Back   │Entry 2 │ Back   │  End      │ │   │
│  │  │ 0x1A   │ 0x03   │ 0x86   │ len    │ ...    │ len    │ 0xFF      │ │   │
│  │  │ (26)   │ (3)    │field1  │        │hello   │        │           │ │   │
│  │  └────────┴────────┴────────┴────────┴────────┴────────┴────────┴────┘ │   │
│  │                                                                        │   │
│  └───────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 3.2 Listpack Operations

```c
// Redis 8 listpack implementation
// src/listpack.c

/* Insert an element at a given position */
unsigned char *lpInsert(unsigned char *lp, unsigned char *ele, uint32_t size,
                        unsigned char *p, int where, unsigned char **newp) {
    // Calculate new size
    unsigned char *intenc;
    uint64_t enclen;

    // Try integer encoding first
    intenc = lpEncodeIntegerGetBytes(ele, size, &enclen);
    if (intenc) {
        // Store as integer
        ele = intenc;
        size = enclen;
    }

    // Calculate backlen
    uint64_t backlen_size = lpBacklenBytes(size);
    uint32_t new_listpack_bytes = lpGetTotalBytes(lp) +
                                   size + backlen_size;

    // Realloc if needed
    if (new_listpack_bytes > lpGetTotalBytes(lp)) {
        lp = lp_realloc(lp, new_listpack_bytes);
    }

    // Insert element and update metadata
    // ... (implementation details)

    return lp;
}

/* Get integer value from listpack entry */
int lpGetInteger(unsigned char *p, long long *value) {
    int64_t v;
    uint64_t negstart, negmax;

    // Determine encoding type from first byte
    unsigned char encoding = p[0];

    if ((encoding & 0x80) == 0) {
        // 7-bit positive integer
        *value = encoding & 0x7f;
        return 0;
    } else if ((encoding & 0xc0) == 0x80) {
        // 6-bit string (not an integer)
        return -1;
    } else if ((encoding & 0xe0) == 0xc0) {
        // 13-bit integer
        *value = ((encoding & 0x1f) << 8) | p[1];
        return 0;
    }
    // ... more encodings
}

/* Seek to element by index */
unsigned char *lpSeek(unsigned char *lp, long index) {
    unsigned char *p;
    uint32_t numele = lpGetNumElements(lp);

    if (index < 0) index = numele + index;
    if (index < 0 || index >= numele) return NULL;

    // Forward or backward scan based on position
    if (index < numele / 2) {
        // Forward from start
        p = lpFirst(lp);
        while (index--) p = lpNext(lp, p);
    } else {
        // Backward from end
        p = lpLast(lp);
        index = numele - 1 - index;
        while (index--) p = lpPrev(lp, p);
    }

    return p;
}
```

---

## 4. Hash Table (dict) Implementation

### 4.1 Dict Structure

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Redis Dict (Hash Table) Implementation                    │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Dict Structure:                                                             │
│  ┌───────────────────────────────────────────────────────────────────────┐   │
│  │  typedef struct dict {                                                │   │
│  │      dictType *type;       /* Type-specific functions */              │   │
│  │      void *privdata;       /* Private data for type functions */      │   │
│  │      dictht ht[2];         /* Two hash tables for rehashing */        │   │
│  │      long rehashidx;       /* -1 = not rehashing */                   │   │
│  │      int16_t pauserehash;  /* Incremented during unsafe iteration */  │   │
│  │  } dict;                                                              │   │
│  │                                                                        │   │
│  │  typedef struct dictht {                                              │   │
│  │      dictEntry **table;    /* Hash table array */                     │   │
│  │      unsigned long size;   /* Size of table */                        │   │
│  │      unsigned long sizemask; /* Mask for hash indexing */             │   │
│  │      unsigned long used;   /* Number of entries */                    │   │
│  │  } dictht;                                                            │   │
│  │                                                                        │   │
│  │  typedef struct dictEntry {                                           │   │
│  │      void *key;            /* Key pointer */                          │   │
│  │      union {                                                           │   │
│  │          void *val;                                                    │   │
│  │          uint64_t u64;                                                 │   │
│  │          int64_t s64;                                                  │   │
│  │          double d;                                                     │   │
│  │      } v;                                                              │   │
│  │      struct dictEntry *next;  /* Linked list for collision */         │   │
│  │  } dictEntry;                                                         │   │
│  └───────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
│  Hash Table Visualization:                                                   │
│  ┌───────────────────────────────────────────────────────────────────────┐   │
│  │                                                                        │   │
│  │  ht[0] (Active)                    ht[1] (Rehashing target)            │   │
│  │  ┌──────────────┐                  ┌──────────────┐                   │   │
│  │  │ table[0]     │──▶ NULL          │              │                   │   │
│  │  │ table[1]     │──▶ Entry ──▶ NULL│              │                   │   │
│  │  │ table[2]     │──▶ Entry ──▶     │              │                   │   │
│  │  │              │      Entry ──▶ NULL (collision chain)              │   │   │
│  │  │ table[3]     │──▶ Entry ──▶ NULL│              │                   │   │
│  │  │   ...        │                  │   ...          │                 │   │
│  │  │ table[size-1]│──▶ NULL          │              │                   │   │
│  │  └──────────────┘                  └──────────────┘                   │   │
│  │  size: 4                           size: 8                            │   │
│  │  used: 3                           used: 0                            │   │
│  │  rehashidx: 2 (currently rehashing bucket 2)                          │   │
│  │                                                                        │   │
│  │  Load factor: used / size = 0.75 (triggers rehash to 2x size)         │   │
│  │                                                                        │   │
│  └───────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
│  Incremental Rehashing:                                                      │
│  ┌───────────────────────────────────────────────────────────────────────┐   │
│  │                                                                        │   │
│  │  For each CRUD operation:                                              │   │
│  │  1. Perform operation on ht[0]                                         │   │
│  │  2. Move 1 bucket from ht[0] to ht[1]                                  │   │
│  │  3. If ht[0] empty, swap ht[1] to ht[0], reset rehashidx               │   │
│  │                                                                        │   │
│  │  Time complexity: O(1) amortized for all operations                    │   │
│  │                                                                        │   │
│  └───────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 4.2 Dict Operations

```c
/* Hash function - MurmurHash2 variant for strings */
uint64_t dictGenHashFunction(const void *key, size_t len) {
    /* MurmurHash2, 64-bit version */
    const uint64_t m = 0xc6a4a7935bd1e995ULL;
    const int r = 47;

    uint64_t h = 0xdeadbeefdeadbeefULL ^ (len * m);

    const uint64_t *data = (const uint64_t *)key;
    const uint64_t *end = data + (len / 8);

    while(data != end) {
        uint64_t k = *data++;
        k *= m;
        k ^= k >> r;
        k *= m;
        h ^= k;
        h *= m;
    }

    // Handle remaining bytes
    const unsigned char *data2 = (const unsigned char*)data;
    switch(len & 7) {
        case 7: h ^= ((uint64_t)data2[6]) << 48;
        case 6: h ^= ((uint64_t)data2[5]) << 40;
        case 5: h ^= ((uint64_t)data2[4]) << 32;
        case 4: h ^= ((uint64_t)data2[3]) << 24;
        case 3: h ^= ((uint64_t)data2[2]) << 16;
        case 2: h ^= ((uint64_t)data2[1]) << 8;
        case 1: h ^= ((uint64_t)data2[0]);
                h *= m;
    }

    h ^= h >> r;
    h *= m;
    h ^= h >> r;

    return h;
}

/* Insert or replace entry */
int dictReplace(dict *d, void *key, void *val) {
    // Try add first
    if (dictAdd(d, key, val) == DICT_OK)
        return 1;

    // Key exists, find and replace
    dictEntry *entry = dictFind(d, key);
    // Free old value, set new
    dictFreeVal(d, entry);
    dictSetVal(d, entry, val);
    return 0;
}

/* Incremental rehashing - called during CRUD operations */
static void _dictRehashStep(dict *d) {
    if (d->pauserehash == 0)
        dictRehash(d, 1);  // Move 1 bucket
}

/* Perform N steps of rehashing */
int dictRehash(dict *d, int n) {
    if (!dictIsRehashing(d)) return 0;

    while(n-- && d->ht[0].used != 0) {
        dictEntry *de, *nextde;

        // Find next non-empty bucket
        while(d->ht[0].table[d->rehashidx] == NULL)
            d->rehashidx++;

        // Move all entries in this bucket
        de = d->ht[0].table[d->rehashidx];
        while(de) {
            uint64_t h = dictHashKey(d, de->key) & d->ht[1].sizemask;
            nextde = de->next;

            // Add to new table
            de->next = d->ht[1].table[h];
            d->ht[1].table[h] = de;

            d->ht[0].used--;
            d->ht[1].used++;
            de = nextde;
        }
        d->ht[0].table[d->rehashidx] = NULL;
        d->rehashidx++;
    }

    // Check if rehash complete
    if (d->ht[0].used == 0) {
        zfree(d->ht[0].table);
        d->ht[0] = d->ht[1];
        _dictReset(&d->ht[1]);
        d->rehashidx = -1;
        return 0;
    }
    return 1;
}
```

---

## 5. Skip List (ZSet) Implementation

### 5.1 Skip List Structure

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Redis Skip List (Sorted Set) Implementation               │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Skip List Structure:                                                        │
│  ┌───────────────────────────────────────────────────────────────────────┐   │
│  │  typedef struct zskiplist {                                           │   │
│  │      struct zskiplistNode *header, *tail;                             │   │
│  │      unsigned long length;     /* Number of nodes */                  │   │
│  │      int level;                /* Current max level */                │   │
│  │  } zskiplist;                                                         │   │
│  │                                                                        │   │
│  │  typedef struct zskiplistNode {                                       │   │
│  │      sds ele;                  /* Member string */                    │   │
│  │      double score;             /* Score for ordering */               │   │
│  │      struct zskiplistNode *backward; /* Previous node */              │   │
│  │      struct zskiplistLevel {                                          │   │
│  │          struct zskiplistNode *forward;                               │   │
│  │          unsigned long span;   /* Nodes skipped to forward */         │   │
│  │      } level[];                                                       │   │
│  │  } zskiplistNode;                                                     │   │
│  └───────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
│  Visual Representation (Max level = 4):                                      │
│  ┌───────────────────────────────────────────────────────────────────────┐   │
│  │                                                                        │   │
│  │  Level 3:  HEAD ──────────────────────────────────────▶ NULL         │   │
│  │                 │                                       │             │   │
│  │  Level 2:  HEAD ───────▶ Node A ──────────────────────▶ Node D       │   │
│  │                 │         │                             │             │   │
│  │  Level 1:  HEAD ───────▶ Node A ───────▶ Node C ──────▶ Node D       │   │
│  │                 │         │               │             │             │   │
│  │  Level 0:  HEAD ───────▶ Node A ───────▶ Node B ──────▶ Node C       │   │
│  │                 │         │               │             │             │   │
│  │                 │         ▼               ▼             ▼             │   │
│  │                 │       Score:          Score:        Score:          │   │
│  │                 │        1.0             2.5           3.0            │   │
│  │                 │       Member:        Member:       Member:          │   │
│  │                 │       "a"            "c"           "d"             │   │
│  │                 │                                                   │   │
│  │  Backward:      NULL ◀───── Node A ◀────── Node B ◀──── Node C       │   │
│  │                                                                        │   │
│  │  Span counts (nodes skipped):                                          │   │
│  │  • HEAD->forward[2] to Node D: span = 3 (A, C, D)                      │   │
│  │  • Node A->forward[1] to Node C: span = 2 (B, C)                       │   │
│  │                                                                        │   │
│  │  This allows O(log n) rank queries!                                    │   │
│  │                                                                        │   │
│  └───────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 5.2 Skip List Operations

```c
/* Random level generator - power of 2 distribution */
int zslRandomLevel(void) {
    int level = 1;
    // ZSKIPLIST_P = 0.25
    while ((random() & 0xFFFF) < (ZSKIPLIST_P * 0xFFFF))
        level += 1;
    return (level < ZSKIPLIST_MAXLEVEL) ? level : ZSKIPLIST_MAXLEVEL;
}

/* Insert element into skip list */
zskiplistNode *zslInsert(zskiplist *zsl, double score, sds ele) {
    zskiplistNode *update[ZSKIPLIST_MAXLEVEL], *x;
    unsigned int rank[ZSKIPLIST_MAXLEVEL];
    int i, level;

    x = zsl->header;
    // Find insertion point at each level
    for (i = zsl->level-1; i >= 0; i--) {
        rank[i] = (i == zsl->level-1) ? 0 : rank[i+1];
        while (x->level[i].forward &&
               (x->level[i].forward->score < score ||
                (x->level[i].forward->score == score &&
                 sdscmp(x->level[i].forward->ele, ele) < 0))) {
            rank[i] += x->level[i].span;
            x = x->level[i].forward;
        }
        update[i] = x;
    }

    // Random level for new node
    level = zslRandomLevel();
    if (level > zsl->level) {
        for (i = zsl->level; i < level; i++) {
            rank[i] = 0;
            update[i] = zsl->header;
            update[i]->level[i].span = zsl->length;
        }
        zsl->level = level;
    }

    // Create and insert node
    x = zslCreateNode(level, score, ele);
    for (i = 0; i < level; i++) {
        x->level[i].forward = update[i]->level[i].forward;
        update[i]->level[i].forward = x;

        // Update span
        x->level[i].span = update[i]->level[i].span - (rank[0] - rank[i]);
        update[i]->level[i].span = (rank[0] - rank[i]) + 1;
    }

    // Update spans for levels above new node's level
    for (i = level; i < zsl->level; i++) {
        update[i]->level[i].span++;
    }

    // Set backward pointer
    x->backward = (update[0] == zsl->header) ? NULL : update[0];
    if (x->level[0].forward)
        x->level[0].forward->backward = x;
    else
        zsl->tail = x;

    zsl->length++;
    return x;
}

/* Get rank of element (1-based) */
unsigned long zslGetRank(zskiplist *zsl, double score, sds ele) {
    zskiplistNode *x;
    unsigned long rank = 0;
    int i;

    x = zsl->header;
    for (i = zsl->level-1; i >= 0; i--) {
        while (x->level[i].forward &&
               (x->level[i].forward->score < score ||
                (x->level[i].forward->score == score &&
                 sdscmp(x->level[i].forward->ele, ele) <= 0))) {
            rank += x->level[i].span;
            x = x->level[i].forward;
        }
        if (x->ele && sdscmp(x->ele, ele) == 0) {
            return rank;
        }
    }
    return 0;  // Not found
}

/* Delete range by rank */
unsigned long zslDeleteRangeByRank(zskiplist *zsl,
                                   unsigned int start,
                                   unsigned int end,
                                   dict *dict) {
    zskiplistNode *update[ZSKIPLIST_MAXLEVEL], *x;
    unsigned long traversed = 0, removed = 0;
    int i;

    x = zsl->header;
    for (i = zsl->level-1; i >= 0; i--) {
        while (x->level[i].forward &&
               (traversed + x->level[i].span) < start) {
            traversed += x->level[i].span;
            x = x->level[i].forward;
        }
        update[i] = x;
    }

    traversed++;
    x = x->level[0].forward;
    while (x && traversed <= end) {
        zskiplistNode *next = x->level[0].forward;
        zslDeleteNode(zsl, x, update);
        dictDelete(dict, x->ele);
        zslFreeNode(x);
        removed++;
        traversed++;
        x = next;
    }
    return removed;
}
```

---

## 6. Memory Optimization

### 6.1 Memory Efficiency by Data Type

| Data Structure | Small Size Encoding | Large Size Encoding | Transition Point |
|----------------|---------------------|---------------------|------------------|
| String | EMBSTR (<=44B) | RAW + SDS | 44 bytes |
| Hash | Listpack (<=512 entries) | Hash table | `hash-max-listpack-entries` |
| List | Quicklist (listpack nodes) | Quicklist (large nodes) | `list-max-listpack-size` |
| Set | Intset (integers only) | Hash table | `set-max-intset-entries` |
| ZSet | Listpack (<=128 entries) | Skiplist + Hash | `zset-max-listpack-entries` |
| Hash | Listpack (<=64 fields) | Hash table | `hash-max-listpack-entries` |

### 6.2 Memory Overhead Analysis

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Redis Memory Overhead Analysis                            │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Per-Key Overhead:                                                           │
│  ┌───────────────────────────────────────────────────────────────────────┐   │
│  │  Component                    │ Size (bytes)                          │   │
│  ├───────────────────────────────┼───────────────────────────────────────┤   │
│  │  redisObject                  │ 16                                    │   │
│  │  dictEntry (in main dict)     │ 24 (64-bit)                           │   │
│  │  Key string (SDS)             │ ~key_length + 5                       │   │
│  │  Expire entry (if TTL set)    │ 16                                    │   │
│  ├───────────────────────────────┼───────────────────────────────────────┤   │
│  │  Total minimum per key        │ ~56 + key_length                      │   │
│  └───────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
│  Example: Storing 1 million small hashes                                     │
│  ┌───────────────────────────────────────────────────────────────────────┐   │
│  │                                                                        │   │
│  │  HSET user:1 name "John" age 30                                       │   │
│  │                                                                        │   │
│  │  Memory breakdown per key:                                             │   │
│  │  • redisObject: 16 bytes                                               │   │
│  │  • dictEntry: 24 bytes                                                 │   │
│  │  • Key SDS "user:1": 10 bytes                                          │   │
│  │  • Listpack with 2 fields: ~30 bytes                                   │   │
│  │  ─────────────────────────                                             │   │
│  │  Total: ~80 bytes per key                                              │   │
│  │                                                                        │   │
│  │  For 1M keys: ~80 MB + Redis overhead (~30%) = ~104 MB                 │   │
│  │                                                                        │   │
│  └───────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 7. Performance Benchmarks

### 7.1 Operation Complexity

| Operation | Time Complexity | Space Complexity | Notes |
|-----------|-----------------|------------------|-------|
| GET/SET | O(1) | O(1) | Hash table lookup |
| HGET/HSET | O(1) | O(1) | Hash table or listpack scan |
| LPUSH/RPUSH | O(1) | O(1) | Head/tail insertion |
| LRANGE | O(S+N) | O(N) | S=start offset, N=count |
| SADD/SREM | O(1) | O(1) | Hash table or intset |
| ZADD/ZREM | O(log N) | O(1) | Skip list insertion |
| ZRANGE | O(log N + M) | O(M) | N=set size, M=result size |
| ZRANK | O(log N) | O(1) | Skip list rank query |

### 7.2 Throughput Benchmarks (Redis 8)

| Operation | Single Thread | pipelined (100) | Notes |
|-----------|---------------|-----------------|-------|
| GET | 120K ops/s | 1.2M ops/s | In-memory hit |
| SET | 100K ops/s | 1M ops/s | No persistence |
| HGET | 110K ops/s | 1.1M ops/s | Small hash |
| HSET | 95K ops/s | 950K ops/s | Small hash |
| LPUSH | 130K ops/s | 1.3M ops/s | Quicklist |
| ZADD | 80K ops/s | 800K ops/s | Skip list |
| ZRANGE | 60K ops/s | 600K ops/s | 100 elements |

---

## 8. References

1. **Redis 8 Documentation**
   - URL: <https://redis.io/documentation>

2. **Redis Source Code**
   - URL: <https://github.com/redis/redis>

3. **Redis Memory Optimization Guide**
   - URL: <https://redis.io/docs/management/optimization/memory-optimization/>

4. **Skip Lists: A Probabilistic Alternative to Balanced Trees**
   - URL: <https://15721.courses.cs.cmu.edu/spring2018/papers/08-oltpindexes1/pugh-skiplists-cacm1990.pdf>

---

*Document generated for S-Level technical reference.*

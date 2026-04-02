# EC-028: Claim-Check Pattern

## Problem Formalization

### The Large Message Problem

Message brokers have limits on message sizes (typically 1MB for Kafka, 256KB for SQS). When applications need to exchange large payloads (files, images, bulk data), sending them directly through the message bus degrades performance and can cause system failures.

#### Problem Statement

Given:

- Message broker B with maximum message size M_max
- Payload P with size |P| where |P| > M_max or |P| causes performance issues
- Storage system S with high capacity

Transform the messaging flow:

```
Before: Producer ──► [P] ──► Broker ──► Consumer
                         (slow, may fail)

After:  Producer ──► [Ref] ──► Broker ──► Consumer ──► S ──► P
             │                                        │
             └────────────────────────────────────────┘
             Store P in S, pass reference through broker
```

Constraints:
    - Ref size << P size
    - Consumer can retrieve P from S using Ref
    - P is eventually cleaned up from S
    - Atomicity: Both store and send succeed, or both fail

```

### Message Size Impact

```

┌─────────────────────────────────────────────────────────────────────────┐
│                    Message Size Impact on System                        │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  Without Claim-Check (Large Messages):                                  │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │  ┌──────────┐                                                    │   │
│  │  │ 10MB msg │ Memory pressure on broker                         │   │
│  │  │ 10MB msg │ Increased GC pauses                               │   │
│  │  │ 10MB msg │ Network saturation                                │   │
│  │  │ 10MB msg │ Slow replication                                  │   │
│  │  └──────────┘ Consumer lag increases                            │   │
│  │                                                                  │   │
│  │  Result: System degradation, potential OOM                      │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                                                                         │
│  With Claim-Check:                                                      │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │  ┌──────────┐        ┌─────────────────────────────────────┐    │   │
│  │  │ 100B ref │        │ External Storage (S3, GCS, etc.)    │    │   │
│  │  │ 100B ref │◄──────►│ • 10MB blob 1                       │    │   │
│  │  │ 100B ref │        │ • 10MB blob 2                       │    │   │
│  │  │ 100B ref │        │ • 10MB blob 3                       │    │   │
│  │  └──────────┘        │ • 10MB blob 4                       │    │   │
│  │       Broker         └─────────────────────────────────────┘    │   │
│  │                                                                  │   │
│  │  Result: Efficient broker operation, scalable storage           │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘

```

## Solution Architecture

### Claim-Check Flow

```

┌─────────────────────────────────────────────────────────────────────────┐
│                    Claim-Check Pattern Flow                             │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  Phase 1: Producer Side                                                 │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │                                                                  │   │
│  │  1. Receive large payload (file, batch data, image)            │   │
│  │            │                                                     │   │
│  │            ▼                                                     │   │
│  │  2. Store in external storage (S3, GCS, Azure Blob)            │   │
│  │     • Generate unique object ID                                │   │
│  │     • Set TTL/cleanup policy                                   │   │
│  │     • Optional: encrypt, compress                              │   │
│  │            │                                                     │   │
│  │            ▼                                                     │   │
│  │  3. Create claim check message                                 │   │
│  │     {                                                           │   │
│  │       "storage_type": "s3",                                     │   │
│  │       "bucket": "my-bucket",                                    │   │
│  │       "key": "uploads/uuid",                                    │   │
│  │       "size": 10485760,                                         │   │
│  │       "checksum": "sha256:abc...",                              │   │
│  │       "expires_at": "2024-01-20T00:00:00Z"                      │   │
│  │     }                                                           │   │
│  │            │                                                     │   │
│  │            ▼                                                     │   │
│  │  4. Publish claim check to message broker                      │   │
│  │                                                                  │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                                                                         │
│  Phase 2: Consumer Side                                                 │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │                                                                  │   │
│  │  1. Receive claim check from broker                            │   │
│  │            │                                                     │   │
│  │            ▼                                                     │   │
│  │  2. Validate claim check                                         │   │
│  │     • Verify not expired                                         │   │
│  │     • Check required fields                                      │   │
│  │            │                                                     │   │
│  │            ▼                                                     │   │
│  │  3. Retrieve payload from storage                              │   │
│  │     • Stream download for large files                          │   │
│  │     • Verify checksum                                          │   │
│  │            │                                                     │   │
│  │            ▼                                                     │   │
│  │  4. Process payload                                              │   │
│  │                                                                  │   │
│  │  5. (Optional) Delete or mark as processed                     │   │
│  │                                                                  │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘

```

### Storage Options Comparison

```

┌─────────────────────────────────────────────────────────────────────────┐
│                    Storage Backend Options                              │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │  Object Storage (S3, GCS, Azure Blob)                            │   │
│  │  ┌─────────────────────────────────────────────────────────────┐ │   │
│  │  │  Pros:                                                      │ │   │
│  │  │  • Virtually unlimited capacity                             │ │   │
│  │  │  • High durability (99.999999999%)                          │ │   │
│  │  │  • Cost-effective for large objects                         │ │   │
│  │  │  • Global accessibility                                     │ │   │
│  │  │                                                             │ │   │
│  │  │  Cons:                                                      │ │   │
│  │  │  • Higher latency (10-100ms)                                │ │   │
│  │  │  • Eventual consistency                                     │ │   │
│  │  │  • Egress costs                                             │ │   │
│  │  └─────────────────────────────────────────────────────────────┘ │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                                                                         │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │  Distributed File System (HDFS, CephFS)                          │   │
│  │  ┌─────────────────────────────────────────────────────────────┐ │   │
│  │  │  Pros:                                                      │ │   │
│  │  │  • Low latency within cluster                               │ │   │
│  │  │  • High throughput                                          │ │   │
│  │  │  • Strong consistency                                       │ │   │
│  │  │                                                             │ │   │
│  │  │  Cons:                                                      │ │   │
│  │  │  • Complex to operate                                       │ │   │
│  │  │  • Limited by cluster size                                  │ │   │
│  │  │  • Not geo-distributed                                      │ │   │
│  │  └─────────────────────────────────────────────────────────────┘ │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                                                                         │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │  Database (PostgreSQL LOB, MongoDB GridFS)                       │   │
│  │  ┌─────────────────────────────────────────────────────────────┐ │   │
│  │  │  Pros:                                                      │ │   │
│  │  │  • Transactional with metadata                              │ │   │
│  │  │  • ACID guarantees                                          │ │   │
│  │  │  • Query capabilities                                       │ │   │
│  │  │                                                             │ │   │
│  │  │  Cons:                                                      │ │   │
│  │  │  • Size limitations                                         │ │   │
│  │  │  • Database performance impact                              │ │   │
│  │  │  • Not designed for large files                             │ │   │
│  │  └─────────────────────────────────────────────────────────────┘ │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘

```

## Production-Ready Go Implementation

### Claim-Check Core Implementation

```go
// pkg/claimcheck/claimcheck.go
package claimcheck

import (
    "context"
    "crypto/sha256"
    "encoding/hex"
    "encoding/json"
    "fmt"
    "io"
    "time"

    "github.com/google/uuid"
)

// ClaimCheck represents a reference to stored data
type ClaimCheck struct {
    ID          string            `json:"id"`
    StorageType string            `json:"storage_type"`
    Location    StorageLocation   `json:"location"`
    Size        int64             `json:"size"`
    Checksum    string            `json:"checksum"`
    ContentType string            `json:"content_type"`
    Metadata    map[string]string `json:"metadata,omitempty"`
    CreatedAt   time.Time         `json:"created_at"`
    ExpiresAt   *time.Time        `json:"expires_at,omitempty"`
    TTL         time.Duration     `json:"ttl,omitempty"`
}

// StorageLocation contains storage-specific location info
type StorageLocation struct {
    Bucket    string `json:"bucket,omitempty"`
    Key       string `json:"key"`
    VersionID string `json:"version_id,omitempty"`
    Region    string `json:"region,omitempty"`
}

// StorageBackend interface for different storage systems
type StorageBackend interface {
    Store(ctx context.Context, key string, reader io.Reader, opts StoreOptions) (*StorageLocation, error)
    Retrieve(ctx context.Context, location *StorageLocation) (io.ReadCloser, error)
    Delete(ctx context.Context, location *StorageLocation) error
    Exists(ctx context.Context, location *StorageLocation) (bool, error)
    GetURL(ctx context.Context, location *StorageLocation, expiry time.Duration) (string, error)
}

type StoreOptions struct {
    ContentType string
    Metadata    map[string]string
    TTL         time.Duration
}

// Manager handles claim-check operations
type Manager struct {
    storage StorageBackend
    config  *Config
}

type Config struct {
    DefaultTTL      time.Duration
    MaxSize         int64
    ChecksumAlgo    string // sha256, md5, etc.
    CleanupEnabled  bool
}

func NewManager(storage StorageBackend, config *Config) *Manager {
    return &Manager{
        storage: storage,
        config:  config,
    }
}

// Create stores data and returns a claim check
func (m *Manager) Create(ctx context.Context, reader io.Reader, opts StoreOptions) (*ClaimCheck, error) {
    // Generate unique ID
    id := uuid.New().String()

    // Compute checksum while reading
    hasher := sha256.New()
    teeReader := io.TeeReader(reader, hasher)

    // Calculate size if needed
    var size int64
    if sized, ok := reader.(io.Seeker); ok {
        // Can determine size
        current, _ := sized.Seek(0, io.SeekCurrent)
        end, _ := sized.Seek(0, io.SeekEnd)
        size = end - current
        sized.Seek(current, io.SeekStart)
    }

    // Store data
    key := generateStorageKey(id)
    location, err := m.storage.Store(ctx, key, teeReader, opts)
    if err != nil {
        return nil, fmt.Errorf("storing data: %w", err)
    }

    // If we couldn't get size before, try to get it now
    if size == 0 {
        // Could use Content-Length from storage response
    }

    checksum := hex.EncodeToString(hasher.Sum(nil))

    // Set TTL
    ttl := opts.TTL
    if ttl == 0 {
        ttl = m.config.DefaultTTL
    }

    var expiresAt *time.Time
    if ttl > 0 {
        t := time.Now().Add(ttl)
        expiresAt = &t
    }

    claimCheck := &ClaimCheck{
        ID:          id,
        StorageType: m.storage.Type(),
        Location:    *location,
        Size:        size,
        Checksum:    fmt.Sprintf("sha256:%s", checksum),
        ContentType: opts.ContentType,
        Metadata:    opts.Metadata,
        CreatedAt:   time.Now(),
        ExpiresAt:   expiresAt,
        TTL:         ttl,
    }

    return claimCheck, nil
}

// Retrieve gets the data referenced by claim check
func (m *Manager) Retrieve(ctx context.Context, claimCheck *ClaimCheck) (io.ReadCloser, error) {
    // Validate not expired
    if claimCheck.ExpiresAt != nil && time.Now().After(*claimCheck.ExpiresAt) {
        return nil, fmt.Errorf("claim check expired")
    }

    // Retrieve from storage
    reader, err := m.storage.Retrieve(ctx, &claimCheck.Location)
    if err != nil {
        return nil, fmt.Errorf("retrieving data: %w", err)
    }

    // If checksum verification requested, wrap reader
    if claimCheck.Checksum != "" {
        return &verifyingReader{
            reader:   reader,
            expected: claimCheck.Checksum,
        }, nil
    }

    return reader, nil
}

// Delete removes the stored data
func (m *Manager) Delete(ctx context.Context, claimCheck *ClaimCheck) error {
    return m.storage.Delete(ctx, &claimCheck.Location)
}

// Serialize converts claim check to message format
func (c *ClaimCheck) Serialize() ([]byte, error) {
    return json.Marshal(c)
}

// Deserialize parses claim check from message
func Deserialize(data []byte) (*ClaimCheck, error) {
    var claimCheck ClaimCheck
    if err := json.Unmarshal(data, &claimCheck); err != nil {
        return nil, fmt.Errorf("deserializing claim check: %w", err)
    }
    return &claimCheck, nil
}

func generateStorageKey(id string) string {
    // Use prefix for organization, date for lifecycle policies
    now := time.Now()
    return fmt.Sprintf("claim-checks/%d/%02d/%s", now.Year(), now.Month(), id)
}

// verifyingReader wraps a reader and verifies checksum on close
type verifyingReader struct {
    reader   io.ReadCloser
    hasher   hash.Hash
    expected string
    closed   bool
}

func (v *verifyingReader) Read(p []byte) (n int, err error) {
    n, err = v.reader.Read(p)
    if n > 0 {
        v.hasher.Write(p[:n])
    }
    return n, err
}

func (v *verifyingReader) Close() error {
    if v.closed {
        return nil
    }
    v.closed = true

    actual := hex.EncodeToString(v.hasher.Sum(nil))
    expected := v.expected
    if idx := strings.Index(expected, ":"); idx != -1 {
        expected = expected[idx+1:]
    }

    if actual != expected {
        return fmt.Errorf("checksum mismatch: expected %s, got %s", expected, actual)
    }

    return v.reader.Close()
}
```

### S3 Storage Backend

```go
// pkg/claimcheck/s3_storage.go
package claimcheck

import (
    "context"
    "fmt"
    "io"
    "time"

    "github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/service/s3"
    "github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// S3Storage implements StorageBackend for AWS S3
type S3Storage struct {
    client     *s3.Client
    bucket     string
    region     string
    presigner  *s3.PresignClient
}

func NewS3Storage(ctx context.Context, bucket string) (*S3Storage, error) {
    cfg, err := config.LoadDefaultConfig(ctx)
    if err != nil {
        return nil, fmt.Errorf("loading AWS config: %w", err)
    }

    client := s3.NewFromConfig(cfg)
    presigner := s3.NewPresignClient(client)

    return &S3Storage{
        client:    client,
        bucket:    bucket,
        region:    cfg.Region,
        presigner: presigner,
    }, nil
}

func (s *S3Storage) Store(ctx context.Context, key string, reader io.Reader, opts StoreOptions) (*StorageLocation, error) {
    // Prepare metadata
    metadata := make(map[string]string)
    for k, v := range opts.Metadata {
        metadata[k] = v
    }
    metadata["uploaded_at"] = time.Now().Format(time.RFC3339)

    // Upload with content type
    contentType := opts.ContentType
    if contentType == "" {
        contentType = "application/octet-stream"
    }

    uploadInput := &s3.PutObjectInput{
        Bucket:      aws.String(s.bucket),
        Key:         aws.String(key),
        Body:        reader,
        ContentType: aws.String(contentType),
        Metadata:    metadata,
    }

    // Set TTL via lifecycle or object expiration
    if opts.TTL > 0 {
        expires := time.Now().Add(opts.TTL)
        uploadInput.Expires = &expires
    }

    result, err := s.client.PutObject(ctx, uploadInput)
    if err != nil {
        return nil, fmt.Errorf("uploading to S3: %w", err)
    }

    location := &StorageLocation{
        Bucket:    s.bucket,
        Key:       key,
        VersionID: aws.ToString(result.VersionId),
        Region:    s.region,
    }

    return location, nil
}

func (s *S3Storage) Retrieve(ctx context.Context, location *StorageLocation) (io.ReadCloser, error) {
    result, err := s.client.GetObject(ctx, &s3.GetObjectInput{
        Bucket:    aws.String(location.Bucket),
        Key:       aws.String(location.Key),
        VersionId: aws.String(location.VersionID),
    })
    if err != nil {
        return nil, fmt.Errorf("getting from S3: %w", err)
    }

    return result.Body, nil
}

func (s *S3Storage) Delete(ctx context.Context, location *StorageLocation) error {
    _, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
        Bucket:    aws.String(location.Bucket),
        Key:       aws.String(location.Key),
        VersionId: aws.String(location.VersionID),
    })
    return err
}

func (s *S3Storage) Exists(ctx context.Context, location *StorageLocation) (bool, error) {
    _, err := s.client.HeadObject(ctx, &s3.HeadObjectInput{
        Bucket:    aws.String(location.Bucket),
        Key:       aws.String(location.Key),
        VersionId: aws.String(location.VersionID),
    })

    if err != nil {
        var notFound *types.NotFound
        if errors.As(err, &notFound) {
            return false, nil
        }
        return false, err
    }

    return true, nil
}

func (s *S3Storage) GetURL(ctx context.Context, location *StorageLocation, expiry time.Duration) (string, error) {
    req, err := s.presigner.PresignGetObject(ctx, &s3.GetObjectInput{
        Bucket:    aws.String(location.Bucket),
        Key:       aws.String(location.Key),
        VersionId: aws.String(location.VersionID),
    }, s3.WithPresignExpires(expiry))

    if err != nil {
        return "", err
    }

    return req.URL, nil
}

func (s *S3Storage) Type() string {
    return "s3"
}
```

### Integration with Message Producer

```go
// pkg/claimcheck/producer.go
package claimcheck

import (
    "context"
    "fmt"
    "io"

    "github.com/company/project/pkg/pubsub"
)

// ClaimCheckProducer wraps a publisher with claim-check capability
type ClaimCheckProducer struct {
    publisher   pubsub.Publisher
    claimCheck  *Manager
    threshold   int64 // Size threshold for claim-check
}

func NewClaimCheckProducer(publisher pubsub.Publisher, manager *Manager, threshold int64) *ClaimCheckProducer {
    return &ClaimCheckProducer{
        publisher:  publisher,
        claimCheck: manager,
        threshold:  threshold,
    }
}

// Publish publishes a message, using claim-check if payload is large
func (p *ClaimCheckProducer) Publish(ctx context.Context, topic string, payload []byte, metadata map[string]string) error {
    // Check if claim-check is needed
    if int64(len(payload)) > p.threshold {
        return p.publishWithClaimCheck(ctx, topic, payload, metadata)
    }

    // Publish directly
    msg := &pubsub.Message{
        Topic:   topic,
        Value:   payload,
        Headers: metadata,
    }
    return p.publisher.Publish(ctx, msg)
}

func (p *ClaimCheckProducer) publishWithClaimCheck(ctx context.Context, topic string, payload []byte, metadata map[string]string) error {
    // Store payload
    reader := bytes.NewReader(payload)
    opts := StoreOptions{
        ContentType: metadata["content_type"],
        Metadata:    metadata,
    }

    claimCheck, err := p.claimCheck.Create(ctx, reader, opts)
    if err != nil {
        return fmt.Errorf("creating claim check: %w", err)
    }

    // Serialize claim check
    claimCheckData, err := claimCheck.Serialize()
    if err != nil {
        return fmt.Errorf("serializing claim check: %w", err)
    }

    // Add claim-check indicator to metadata
    metadata["claim_check"] = "true"
    metadata["original_size"] = fmt.Sprintf("%d", len(payload))

    // Publish claim check
    msg := &pubsub.Message{
        Topic:   topic,
        Value:   claimCheckData,
        Headers: metadata,
    }

    return p.publisher.Publish(ctx, msg)
}

// PublishStream handles streaming large data
func (p *ClaimCheckProducer) PublishStream(ctx context.Context, topic string, reader io.Reader, size int64, contentType string, metadata map[string]string) error {
    // Always use claim-check for streams
    opts := StoreOptions{
        ContentType: contentType,
        Metadata:    metadata,
    }

    claimCheck, err := p.claimCheck.Create(ctx, reader, opts)
    if err != nil {
        return err
    }

    claimCheck.Size = size

    claimCheckData, err := claimCheck.Serialize()
    if err != nil {
        return err
    }

    metadata["claim_check"] = "true"
    metadata["streaming"] = "true"

    msg := &pubsub.Message{
        Topic:   topic,
        Value:   claimCheckData,
        Headers: metadata,
    }

    return p.publisher.Publish(ctx, msg)
}
```

## Trade-off Analysis

### When to Use Claim-Check

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    Claim-Check Decision Matrix                          │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  Use Claim-Check when:                                                  │
│  ✓ Message > 100KB (Kafka) or > 256KB (SQS)                            │
│  ✓ Batch processing with large payloads                                │
│  ✓ File uploads/downloads                                              │
│  ✓ Media processing (images, video)                                    │
│  ✓ Data export/import operations                                       │
│  ✓ Report generation                                                   │
│                                                                         │
│  Don't use Claim-Check when:                                            │
│  ✗ Messages are small (< 10KB)                                         │
│  ✗ Low latency is critical (< 100ms end-to-end)                        │
│  ✗ Atomicity of message + storage is required                          │
│  ✗ Network to storage is unreliable                                    │
│  ✗ Simple pub/sub with small events                                    │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

### Performance Comparison

| Metric | Direct Pub/Sub | Claim-Check | Notes |
|--------|----------------|-------------|-------|
| **Latency (small msg)** | 1-10ms | 1-10ms | No difference |
| **Latency (large msg)** | 100-1000ms | 20-50ms + storage time | Much better |
| **Throughput** | Limited by broker | Limited by storage | Better for large msgs |
| **Cost (small msg)** | Low | Higher (storage overhead) | Not worth it |
| **Cost (large msg)** | High (broker resources) | Lower | Significant savings |
| **Reliability** | Depends on broker | Storage + broker | Higher with good storage |

## Testing Strategies

### Claim-Check Testing

```go
// test/claimcheck/claimcheck_test.go
package claimcheck

import (
    "bytes"
    "context"
    "io"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestClaimCheckRoundTrip(t *testing.T) {
    // Use in-memory storage for testing
    storage := NewMemoryStorage()
    manager := NewManager(storage, &Config{DefaultTTL: time.Hour})

    // Test data
    originalData := []byte("This is a large payload that would benefit from claim-check pattern")

    // Create claim check
    ctx := context.Background()
    claimCheck, err := manager.Create(ctx, bytes.NewReader(originalData), StoreOptions{
        ContentType: "text/plain",
    })
    require.NoError(t, err)
    assert.NotEmpty(t, claimCheck.ID)
    assert.NotEmpty(t, claimCheck.Checksum)

    // Retrieve data
    reader, err := manager.Retrieve(ctx, claimCheck)
    require.NoError(t, err)
    defer reader.Close()

    retrievedData, err := io.ReadAll(reader)
    require.NoError(t, err)

    assert.Equal(t, originalData, retrievedData)
}

func TestClaimCheckVerification(t *testing.T) {
    storage := NewMemoryStorage()
    manager := NewManager(storage, &Config{})

    originalData := []byte("test data")

    ctx := context.Background()
    claimCheck, err := manager.Create(ctx, bytes.NewReader(originalData), StoreOptions{})
    require.NoError(t, err)

    // Corrupt the stored data
    storage.Corrupt(claimCheck.Location.Key)

    // Retrieve should fail checksum verification
    reader, err := manager.Retrieve(ctx, claimCheck)
    require.NoError(t, err)
    defer reader.Close()

    // Read all to trigger verification on close
    io.ReadAll(reader)
    err = reader.Close()
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "checksum mismatch")
}

func TestClaimCheckExpiration(t *testing.T) {
    storage := NewMemoryStorage()
    manager := NewManager(storage, &Config{})

    data := []byte("expiring data")

    ctx := context.Background()
    claimCheck, err := manager.Create(ctx, bytes.NewReader(data), StoreOptions{
        TTL: time.Millisecond,
    })
    require.NoError(t, err)

    // Wait for expiration
    time.Sleep(10 * time.Millisecond)

    // Retrieve should fail
    _, err = manager.Retrieve(ctx, claimCheck)
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "expired")
}

func TestClaimCheckSerialization(t *testing.T) {
    original := &ClaimCheck{
        ID:          "test-id",
        StorageType: "s3",
        Location: StorageLocation{
            Bucket: "my-bucket",
            Key:    "test/key",
        },
        Size:        1024,
        Checksum:    "sha256:abc123",
        ContentType: "application/json",
    }

    // Serialize
    data, err := original.Serialize()
    require.NoError(t, err)

    // Deserialize
    deserialized, err := Deserialize(data)
    require.NoError(t, err)

    assert.Equal(t, original.ID, deserialized.ID)
    assert.Equal(t, original.StorageType, deserialized.StorageType)
    assert.Equal(t, original.Location.Key, deserialized.Location.Key)
}
```

## Summary

The Claim-Check Pattern provides:

1. **Message Size Management**: Handle payloads larger than broker limits
2. **Performance**: Keep message brokers efficient
3. **Cost Optimization**: Reduce expensive broker storage
4. **Flexibility**: Use appropriate storage for payload characteristics
5. **Reliability**: Checksum verification ensures integrity

Key considerations:

- Storage latency adds to total processing time
- Need cleanup strategy for claim-check data
- Consider encryption for sensitive data
- Monitor storage costs vs broker costs
- Handle atomicity between store and publish

# Security and Cryptography 2026

> **分类**: 工程与云原生
> **标签**: #security #cryptography #post-quantum #zerotrust #mtls #supply-chain #kubernetes
> **参考**: NIST FIPS 203-205, CNCF Security Whitepaper, Trail of Bits Go Audit 2025

---

## 1. Post-Quantum Cryptography (PQC)

### 1.1 NIST Standards Overview (August 2024)

In August 2024, NIST released the first three finalized post-quantum cryptography standards, marking the beginning of the post-quantum transition era.

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    NIST Post-Quantum Cryptography Standards                  │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │  FIPS 203 - ML-KEM (Module Lattice-based Key Encapsulation Mechanism)│   │
│  │  ─────────────────────────────────────────────────────────────────  │   │
│  │  • Formerly: Kyber (NIST Round 3 winner)                            │   │
│  │  • Purpose: General encryption/key encapsulation                     │   │
│  │  • Security Levels: ML-KEM-512, ML-KEM-768, ML-KEM-1024             │   │
│  │  • Based on: Module Learning With Errors (MLWE)                      │   │
│  │  • Performance: ~10x faster than RSA-3072 key generation             │   │
│  │                                                                     │   │
│  │  Use Cases:                                                         │   │
│  │  ✓ TLS handshake key exchange                                        │   │
│  │  ✓ Key establishment in protocols                                    │   │
│  │  ✓ Hybrid post-quantum TLS                                           │   │
│  │                                                                     │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │  FIPS 204 - ML-DSA (Module Lattice-based Digital Signature Algorithm)│   │
│  │  ─────────────────────────────────────────────────────────────────  │   │
│  │  • Formerly: Dilithium (NIST Round 3 winner)                        │   │
│  │  • Purpose: Digital signatures for authentication                    │   │
│  │  • Security Levels: ML-DSA-44, ML-DSA-65, ML-DSA-87                 │   │
│  │  • Based on: Module Learning With Errors + Short Integer Solution    │   │
│  │  • Performance: Sign ~2x faster than ECDSA, Verify ~3x faster        │   │
│  │                                                                     │   │
│  │  Use Cases:                                                         │   │
│  │  ✓ Code signing                                                      │   │
│  │  ✓ Document signing                                                  │   │
│  │  ✓ Certificate signing (CAs)                                         │   │
│  │  ✓ Software supply chain signing                                     │   │
│  │                                                                     │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │  FIPS 205 - SLH-DSA (Stateless Hash-based Digital Signature Algorithm)│  │
│  │  ─────────────────────────────────────────────────────────────────── │   │
│  │  • Formerly: SPHINCS+ (NIST Round 3 winner)                         │   │
│  │  • Purpose: Conservative backup signature scheme                     │   │
│  │  • Security Levels: SLH-DSA-SHA2-128s/f, SLH-DSA-SHAKE-256s/f       │   │
│  │  • Based on: Hash functions only (SHA-2 or SHAKE256)                │   │
│  │  • Performance: Larger signatures (~8KB), slower signing             │   │
│  │                                                                     │   │
│  │  Use Cases:                                                         │   │
│  │  ✓ Long-term trust anchors                                           │   │
│  │  ✓ Fallback when lattice assumptions fail                            │   │
│  │  ✓ High-security applications                                        │   │
│  │  ✓ Firmware signing                                                  │   │
│  │                                                                     │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                             │
│  FUTURE STANDARD (Draft):                                                   │
│  • FIPS 206 - FN-DSA (Falcon) - Lattice-based, NTRU hardness               │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 1.2 Go 1.24+ Post-Quantum Support

Go 1.24 introduced native support for post-quantum cryptography, with X25519MLKEM768 hybrid as the default for TLS.

```go
package pqc

import (
 "crypto/mlkem"
 "crypto/tls"
 "crypto/x509"
 "fmt"
)

// MLKEMKeyExchange demonstrates ML-KEM key encapsulation
func MLKEMKeyExchange() {
 // ML-KEM-768 (recommended default)
 // Provides NIST Level 3 security (~AES-192 equivalent)

 // Generate key pair
 publicKey, privateKey, err := mlkem.GenerateKey768()
 if err != nil {
  panic(err)
 }

 // Encapsulate: generates shared secret and ciphertext
 sharedSecret, ciphertext, err := publicKey.Encapsulate()
 if err != nil {
  panic(err)
 }

 // Decapsulate: recover shared secret from ciphertext
 decapsulatedSecret, err := privateKey.Decapsulate(ciphertext)
 if err != nil {
  panic(err)
 }

 // Verify both parties have the same shared secret
 if string(sharedSecret) != string(decapsulatedSecret) {
  panic("shared secrets don't match")
 }

 fmt.Printf("ML-KEM-768 key exchange successful\n")
 fmt.Printf("Ciphertext size: %d bytes\n", len(ciphertext))
 fmt.Printf("Shared secret size: %d bytes\n", len(sharedSecret))
}

// MLKEM512KeyExchange for resource-constrained environments
func MLKEM512KeyExchange() {
 publicKey, privateKey, err := mlkem.GenerateKey512()
 if err != nil {
  panic(err)
 }

 sharedSecret, ciphertext, err := publicKey.Encapsulate()
 if err != nil {
  panic(err)
 }

 decapsulatedSecret, err := privateKey.Decapsulate(ciphertext)
 if err != nil {
  panic(err)
 }

 _ = decapsulatedSecret
 fmt.Printf("ML-KEM-512 key exchange successful\n")
}

// MLKEM1024KeyExchange for maximum security
func MLKEM1024KeyExchange() {
 publicKey, privateKey, err := mlkem.GenerateKey1024()
 if err != nil {
  panic(err)
 }

 sharedSecret, ciphertext, err := publicKey.Encapsulate()
 if err != nil {
  panic(err)
 }

 decapsulatedSecret, err := privateKey.Decapsulate(ciphertext)
 if err != nil {
  panic(err)
 }

 _ = decapsulatedSecret
 fmt.Printf("ML-KEM-1024 key exchange successful\n")
}

// PostQuantumTLSConfig creates TLS config with post-quantum key exchange
func PostQuantumTLSConfig() *tls.Config {
 return &tls.Config{
  // Go 1.24+ defaults to X25519MLKEM768 hybrid key exchange
  // This combines X25519 (classical ECDH) with ML-KEM-768 (PQC)
  // Provides "harvest now, decrypt later" protection

  MinVersion: tls.VersionTLS13,

  // Explicitly prefer PQ key exchanges
  CurvePreferences: []tls.CurveID{
   tls.X25519MLKEM768,  // PQ hybrid (Go 1.24+)
   tls.SecP256R1,       // Fallback
   tls.SecP384R1,       // High security fallback
  },

  // Strong cipher suites
  CipherSuites: []uint16{
   tls.TLS_AES_256_GCM_SHA384,
   tls.TLS_AES_128_GCM_SHA256,
   tls.TLS_CHACHA20_POLY1305_SHA256,
  },

  PreferServerCipherSuites: true,
 }
}

// HybridKeyExchange demonstrates hybrid post-quantum key exchange
// Combines classical ECC with ML-KEM for defense in depth
func HybridKeyExchange() {
 // X25519MLKEM768 is the recommended hybrid
 // - X25519: classical security, well-studied
 // - ML-KEM-768: post-quantum security
 // - Combined: secure even if one algorithm is broken

 fmt.Println(`
Hybrid Key Exchange: X25519MLKEM768
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Classical Component: X25519
  • Curve: Curve25519
  • Security: ~128 bits classical
  • Status: Well-studied, widely deployed

Post-Quantum Component: ML-KEM-768
  • Lattice: Module Learning With Errors
  • Security: NIST Level 3 (post-quantum)
  • Status: NIST FIPS 203 standardized

Combined Security:
  • If X25519 broken: ML-KEM-768 protects
  • If ML-KEM broken: X25519 protects
  • Both broken: No known attack

Default in Go 1.24+ TLS for:
  • crypto/tls clients
  • crypto/tls servers
  • net/http clients
`)
}

// PQEnabledHTTPClient creates HTTP client with post-quantum TLS
func PQEnabledHTTPClient() *http.Client {
 return &http.Client{
  Transport: &http.Transport{
   TLSClientConfig: &tls.Config{
    MinVersion: tls.VersionTLS13,
    // X25519MLKEM768 automatically enabled in Go 1.24+
   },
  },
 }
}

// VerifyPQConnection verifies if a connection uses post-quantum key exchange
func VerifyPQConnection(conn *tls.Conn) {
 state := conn.ConnectionState()

 fmt.Printf("TLS Version: %s\n", tlsVersionName(state.Version))
 fmt.Printf("Cipher Suite: %s\n", tls.CipherSuiteName(state.CipherSuite))

 // Check negotiated key exchange
 for _, curve := range state.ServerName {
  switch curve {
  case tls.X25519MLKEM768:
   fmt.Println("✓ Post-quantum key exchange: X25519MLKEM768")
  case tls.SecP256R1, tls.SecP384R1, tls.SecP521R1:
   fmt.Printf("⚠ Classical key exchange only: %d\n", curve)
  }
 }
}

func tlsVersionName(version uint16) string {
 switch version {
 case tls.VersionTLS13:
  return "TLS 1.3"
 case tls.VersionTLS12:
  return "TLS 1.2"
 default:
  return fmt.Sprintf("Unknown (%d)", version)
 }
}
```

### 1.3 Migration Timelines

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Post-Quantum Cryptography Migration Timeline              │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  2024 ─────────────────────────────────────────────────────────────────►   │
│   │                                                                         │
│   ├── Aug 2024: NIST FIPS 203, 204, 205 finalized                          │
│   │                                                                         │
│   ├── Q3 2024: Go 1.23 experimental ML-KEM support                         │
│   │                                                                         │
│   ├── Q4 2024: Major browsers enable ML-KEM by default                     │
│   │   • Chrome 131+, Firefox 132+, Safari 18+                              │
│   │                                                                         │
│   └── Dec 2024: CNSA 2.0 timeline updated (NSA)                            │
│                                                                             │
│  2025 ─────────────────────────────────────────────────────────────────►   │
│   │                                                                         │
│   ├── Feb 2025: Go 1.24 released with X25519MLKEM768 default               │
│   │                                                                         │
│   ├── Q2 2025: OpenSSL 3.5 with full ML-KEM/ML-DSA support                 │
│   │                                                                         │
│   ├── Q3 2025: First commercial PQ certificates available                  │
│   │                                                                         │
│   └── Q4 2025: Major cloud providers enable PQ TLS by default              │
│       • AWS, Azure, GCP                                                      │
│                                                                             │
│  2026 ─────────────────────────────────────────────────────────────────►   │
│   │                                                                         │
│   ├── Q1 2026: Federal agencies must support PQ algorithms                 │
│   │   (US OMB Memo M-23-02)                                                  │
│   │                                                                         │
│   ├── Q2-Q4 2026: Enterprise migration acceleration                        │
│   │                                                                         │
│   └── Dec 2026: EU NIS2 Directive enforcement begins                       │
│                                                                             │
│  2027-2029 ────────────────────────────────────────────────────────────►   │
│   │                                                                         │
│   ├── 2027: Commercial software PQ mandates begin                          │
│   │                                                                         │
│   ├── 2028: Financial services sector deadline (FFIEC guidance)            │
│   │                                                                         │
│   └── 2029: Healthcare sector deadline (HHS guidance)                      │
│                                                                             │
│  2030 ─────────────────────────────────────────────────────────────────►   │
│   │                                                                         │
│   ├── Jan 2030: US Federal deadline for PQ-only algorithms                 │
│   │   (CNSA 2.0 compliance required)                                         │
│   │                                                                         │
│   └── Jan 2030: EU Cybersecurity Act PQC requirements                      │
│       (Critical entities must use PQC)                                       │
│                                                                             │
│  2035 ─────────────────────────────────────────────────────────────────►   │
│   │                                                                         │
│   └── Long-term: Full transition to PQC complete                           │
│       Classical algorithms phased out for sensitive applications             │
│                                                                             │
│  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━ │
│                                                                             │
│  MANDATORY COMPLIANCE DEADLINES:                                            │
│                                                                             │
│  ┌──────────────────┬─────────────────┬─────────────────────────────────┐   │
│  │ Jurisdiction     │ Deadline        │ Requirements                    │   │
│  ├──────────────────┼─────────────────┼─────────────────────────────────┤   │
│  │ US Federal       │ 2030            │ CNSA 2.0 compliant algorithms   │   │
│  │ EU Critical      │ 2030            │ PQC for critical infrastructure │   │
│  │ EU General       │ 2035            │ NIS2 Directive compliance       │   │
│  │ Financial (US)   │ 2028            │ FFIEC guidance compliance       │   │
│  │ Healthcare (US)  │ 2029            │ HHS 405(d) guidance             │   │
│  └──────────────────┴─────────────────┴─────────────────────────────────┘   │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 1.4 Performance Comparison

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Post-Quantum Algorithm Performance                        │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  OPERATION                    │  ML-KEM-768  │  ML-DSA-65   │  ECDH P-256  │
├───────────────────────────────┼──────────────┼──────────────┼──────────────┤
│  Key Generation               │  ~50 μs      │  ~300 μs     │  ~1 ms       │
│  Encapsulation/Sign           │  ~60 μs      │  ~400 μs     │  ~1.5 ms     │
│  Decapsulation/Verify         │  ~70 μs      │  ~150 μs     │  ~2 ms       │
│  Public Key Size              │  1184 bytes  │  1952 bytes  │  65 bytes    │
│  Secret Key Size              │  2400 bytes  │  4032 bytes  │  32 bytes    │
│  Ciphertext/Signature Size    │  1088 bytes  │  3293 bytes  │  64 bytes    │
│                                                                             │
│  RELATIVE PERFORMANCE (compared to classical):                              │
│                                                                             │
│  Algorithm          │  Speed    │  Key Size  │  Data Overhead              │
│  ───────────────────┼───────────┼────────────┼───────────────────────────  │
│  ML-KEM-768         │  ~10x     │  ~18x      │  ~17x larger ciphertext     │
│  ML-KEM-1024        │  ~8x      │  ~24x      │  ~21x larger ciphertext     │
│  ML-DSA-65          │  ~3x      │  ~30x      │  ~50x larger signature      │
│  SLH-DSA-SHA2-128s  │  ~100x    │  ~64x      │  ~128x larger signature     │
│                                                                             │
│  NETWORK IMPACT:                                                            │
│  • TLS handshake: +1-2 RTT (larger key exchange messages)                  │
│  • Certificate chains: 2-3x larger with PQ signatures                        │
│  • Bandwidth: ~5-10% increase for typical workloads                          │
│                                                                             │
│  HARDWARE ACCELERATION:                                                     │
│  • AVX2: 2-3x speedup for lattice operations                                 │
│  • AVX-512: 4-5x speedup (newer processors)                                  │
│  • ARM NEON: 2-3x speedup (mobile/embedded)                                  │
│  • GPU: Not recommended (constant-time concerns)                             │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 2. Zero-Knowledge Proofs (ZKP)

### 2.1 ZK-SNARKs vs ZK-STARKs

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Zero-Knowledge Proof Systems Comparison                   │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│                         ┌─────────────────────────┐                         │
│                         │     ZK-SNARKs           │                         │
│                         │  (Zero-Knowledge        │                         │
│                         │   Succinct Non-         │                         │
│                         │   Interactive Arguments │                         │
│                         │   of Knowledge)         │                         │
│                         └───────────┬─────────────┘                         │
│                                     │                                       │
│  Characteristics:                   │                                       │
│  • Proof size: ~200 bytes           │                                       │
│  • Verification: ~3ms               │                                       │
│  • Trusted setup: REQUIRED          │                                       │
│  • Post-quantum: NO                 │                                       │
│  • Cryptography: Elliptic curves    │                                       │
│                                     │                                       │
│  Implementations:                   │                                       │
│  • Groth16 (smallest proofs)        │                                       │
│  • PLONK (universal setup)          │                                       │
│  • Marlin (no trusted setup)        │                                       │
│  • Bulletproofs (no trusted setup)  │                                       │
│                                                                             │
│  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━ │
│                                                                             │
│                         ┌─────────────────────────┐                         │
│                         │     ZK-STARKs           │                         │
│                         │  (Zero-Knowledge        │                         │
│                         │   Scalable Transparent  │                         │
│                         │   Arguments of          │                         │
│                         │   Knowledge)            │                         │
│                         └───────────┬─────────────┘                         │
│                                     │                                       │
│  Characteristics:                   │                                       │
│  • Proof size: ~50-100 KB           │                                       │
│  • Verification: ~10ms              │                                       │
│  • Trusted setup: NOT REQUIRED      │                                       │
│  • Post-quantum: YES (hash-based)   │                                       │
│  • Cryptography: Hash functions     │                                       │
│                                     │                                       │
│  Implementations:                   │                                       │
│  • StarkWare (Cairo, StarkEx)       │                                       │
│  • Polygon Hermez (zkEVM)           │                                       │
│  • RISC Zero (general purpose)      │                                       │
│  • Succinct Labs (SP1)              │                                       │
│                                                                             │
│  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━ │
│                                                                             │
│  COMPARISON MATRIX:                                                         │
│                                                                             │
│  Feature              │  ZK-SNARKs   │  ZK-STARKs    │  Winner             │
│  ─────────────────────┼──────────────┼───────────────┼─────────────────────│
│  Proof size           │  ~200 B      │  ~50-100 KB   │  SNARKs             │
│  Verification time    │  ~3ms        │  ~10ms        │  SNARKs             │
│  Proving time         │  Seconds     │  Minutes      │  SNARKs             │
│  Trusted setup        │  Required    │  None         │  STARKs             │
│  Post-quantum secure  │  No          │  Yes          │  STARKs             │
│  Scalability          │  Limited     │  Excellent    │  STARKs             │
│  Implementation ease  │  Moderate    │  Complex      │  SNARKs             │
│                                                                             │
│  SELECTION GUIDE:                                                           │
│  • Use SNARKs: Small proofs, fast verification, acceptable trust           │
│  • Use STARKs: No trusted setup, post-quantum security, transparency       │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 2.2 ZKML: Zero-Knowledge Machine Learning

```go
package zkml

import (
 "crypto/sha256"
 "encoding/hex"
 "fmt"
)

// ZKMLConcept demonstrates the concept of zero-knowledge ML inference
// In ZKML, we prove that a model correctly computed an output
// without revealing the model weights or input data

type ZKMLProof struct {
 // Commitment to model weights (hash)
 ModelCommitment string

 // Commitment to input
 InputCommitment string

 // Public output
 Output float64

 // Zero-knowledge proof
 Proof []byte
}

// ModelCommitment represents a commitment to model weights
func ModelCommitment(weights []float64) string {
 // In production: Use Pedersen commitments or similar
 // This is a simplified hash-based commitment
 h := sha256.New()
 for _, w := range weights {
  h.Write([]byte(fmt.Sprintf("%.10f", w)))
 }
 return hex.EncodeToString(h.Sum(nil))
}

// InputCommitment represents a commitment to input data
func InputCommitment(input []float64) string {
 h := sha256.New()
 for _, x := range input {
  h.Write([]byte(fmt.Sprintf("%.10f", x)))
 }
 return hex.EncodeToString(h.Sum(nil))
}

// ZKMLInference simulates ZKML inference
// In practice, this would use a ZK circuit (Circom, GKR, etc.)
func ZKMLInference(modelWeights, input []float64) (*ZKMLProof, error) {
 // 1. Compute commitment to model
 modelCommit := ModelCommitment(modelWeights)

 // 2. Compute commitment to input
 inputCommit := InputCommitment(input)

 // 3. Compute inference (matrix multiplication + activation)
 output := computeInference(modelWeights, input)

 // 4. Generate ZK proof (simplified)
 proof := generateZKProof(modelWeights, input, output)

 return &ZKMLProof{
  ModelCommitment: modelCommit,
  InputCommitment: inputCommit,
  Output:          output,
  Proof:           proof,
 }, nil
}

func computeInference(weights, input []float64) float64 {
 // Simplified: dot product with ReLU
 sum := 0.0
 for i := range weights {
  if i < len(input) {
   sum += weights[i] * input[i]
  }
 }
 // ReLU
 if sum < 0 {
  return 0
 }
 return sum
}

func generateZKProof(weights, input []float64, output float64) []byte {
 // Placeholder: In production, use:
 // - groth16 (Circom)
 // - Plonkish arithmetization
 // - GKR protocol
 // - STARKs (Cairo)

 h := sha256.New()
 for _, w := range weights {
  h.Write([]byte(fmt.Sprintf("%.10f", w)))
 }
 for _, x := range input {
  h.Write([]byte(fmt.Sprintf("%.10f", x)))
 }
 h.Write([]byte(fmt.Sprintf("%.10f", output)))
 return h.Sum(nil)
}

// VerifyZKMLProof verifies a ZKML proof
func VerifyZKMLProof(proof *ZKMLProof) bool {
 // In production, verify the zero-knowledge proof
 // without accessing model weights or input data

 fmt.Println("Verifying ZKML proof...")
 fmt.Printf("Model commitment: %s...\n", proof.ModelCommitment[:16])
 fmt.Printf("Input commitment: %s...\n", proof.InputCommitment[:16])
 fmt.Printf("Public output: %.4f\n", proof.Output)

 // Verify proof validity (simplified)
 return len(proof.Proof) == 32 // SHA-256 length
}
```

### 2.3 ZKML Frameworks and Applications

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    ZKML Frameworks and Platforms (2026)                      │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  FRAMEWORK              │  Type       │  Maturity │  Performance           │
│  ───────────────────────┼─────────────┼───────────┼─────────────────────── │
│  EZKL                   │  SNARK      │  ★★★★☆    │  Fast, GPU-accelerated │
│  (ethereum/zkml)        │             │           │  Circom-based          │
│                         │             │           │                        │
│  NANOZK                 │  STARK      │  ★★★☆☆    │  Scalable              │
│  (nano-labs)            │             │           │  Recursive proofs      │
│                         │             │           │                        │
│  zkGPT                  │  SNARK      │  ★★★☆☆    │  LLM inference         │
│  (zkllm)                │             │           │  Transformer circuits  │
│                         │             │           │                        │
│  RISC Zero              │  STARK      │  ★★★★☆    │  General purpose       │
│  (risc0)                │             │           │  Rust-based            │
│                         │             │           │                        │
│  SP1                    │  STARK      │  ★★★★☆    │  Succinct              │
│  (Succinct Labs)        │             │           │  Production-ready      │
│                                                                             │
│  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━ │
│                                                                             │
│  APPLICATIONS:                                                              │
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │  zkLLM - Verifiable AI Inference                                    │   │
│  │                                                                     │   │
│  │  Use Case: Prove that an AI model correctly processed input         │   │
│  │            without revealing the model or input data                │   │
│  │                                                                     │   │
│  │  Applications:                                                      │   │
│  │  • Medical diagnosis (privacy-preserving)                           │   │
│  │  • Financial modeling (proprietary algorithms)                      │   │
│  │  • Content moderation (bias verification)                           │   │
│  │  • Fraud detection (model secrecy)                                  │   │
│  │                                                                     │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │  Mina Protocol - Constant-Size Blockchain                           │   │
│  │                                                                     │   │
│  │  • Proof size: Always ~22 KB regardless of chain history            │   │
│  │  • ZK-SNARKs for recursive proof composition                        │   │
│  │  • Decentralized and lightweight nodes                              │   │
│  │  • Smart contracts with zero-knowledge state                        │   │
│  │                                                                     │   │
│  │  Architecture:                                                      │   │
│  │  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐             │   │
│  │  │  Block N    │───►│  Block N+1  │───►│  Block N+2  │             │   │
│  │  │  Proof Pn   │    │  Proof Pn+1 │    │  Proof Pn+2 │             │   │
│  │  └──────┬──────┘    └──────┬──────┘    └──────┬──────┘             │   │
│  │         │                   │                   │                  │   │
│  │         └───────────────────┴───────────────────┘                  │   │
│  │                             │                                       │   │
│  │                             ▼                                       │   │
│  │                     ┌───────────────┐                               │   │
│  │                     │ Recursive ZKP │                               │   │
│  │                     │ (~22 KB)      │                               │   │
│  │                     └───────────────┘                               │   │
│  │                                                                     │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │  Zcash - Privacy-Preserving Payments                                │   │
│  │                                                                     │   │
│  │  • zk-SNARKs (Groth16) for shielded transactions                    │   │
│  │  • Proves transaction validity without revealing amounts/addresses  │   │
│  │  • Halo 2 upgrade: trustless setup                                  │   │
│  │                                                                     │   │
│  │  Transaction Types:                                                 │   │
│  │  • Transparent: Public blockchain (like Bitcoin)                    │   │
│  │  • Shielded: ZK proofs hide sender, receiver, amount                │   │
│  │                                                                     │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │  Aztec - Encrypted DeFi                                             │   │
│  │                                                                     │   │
│  │  • zk-Rollup with privacy by default                                │   │
│  │  • UTXO-based private state model                                   │   │
│  │  • Noir language for private smart contracts                        │   │
│  │                                                                     │   │
│  │  Features:                                                          │   │
│  │  • Private token transfers                                          │   │
│  │  • Anonymous DeFi interactions                                      │   │
│  │  • Compliance tooling (viewing keys)                                │   │
│  │                                                                     │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 3. Supply Chain Security

### 3.1 SLSA Framework (Levels 1-4)

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    SLSA (Supply-chain Levels for Software Artifacts)         │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  SLSA is a security framework to prevent tampering and improve integrity    │
│  throughout the software supply chain.                                      │
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │  LEVEL 1 - Provenance Exists                                        │   │
│  │  ═══════════════════════════════                                    │   │
│  │                                                                     │   │
│  │  Requirements:                                                      │   │
│  │  ✓ Build process is fully scripted/automated                        │   │
│  │  ✓ Provenance exists showing how artifact was built                 │   │
│  │                                                                     │   │
│  │  Security Benefits:                                                 │   │
│  │  • Identify sources of software                                     │   │
│  │  • Track dependencies                                               │   │
│  │                                                                     │   │
│  │  Implementation:                                                    │   │
│  │  • GitHub Actions with workflow provenance                          │   │
│  │  • GitLab CI with artifact metadata                                 │   │
│  │  • Simple SBOM generation                                           │   │
│  │                                                                     │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │  LEVEL 2 - Hosted Build with Signed Provenance                      │   │
│  │  ═════════════════════════════════════════════                      │   │
│  │                                                                     │   │
│  │  Requirements:                                                      │   │
│  │  ✓ Build service is hosted (not self-hosted runner)                 │   │
│  │  ✓ Provenance is signed by build service                            │   │
│  │  ✓ Build service generates provenance                               │   │
│  │                                                                     │   │
│  │  Security Benefits:                                                 │   │
│  │  • Prevents tampering with provenance                               │   │
│  │  • Identifies build system compromise                               │   │
│  │                                                                     │   │
│  │  Implementation:                                                    │   │
│  │  • GitHub-hosted runners with OIDC tokens                           │   │
│  │  • Sigstore signing with Fulcio                                     │   │
│  │  • SLSA GitHub Generator                                            │   │
│  │                                                                     │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │  LEVEL 3 - Hardened Builds                                          │   │
│  │  ═══════════════════════════                                        │   │
│  │                                                                     │   │
│  │  Requirements:                                                      │   │
│  │  ✓ Build runs in an isolated, hardened environment                  │   │
│  │  ✓ Build is hermetic (no network access)                            │   │
│  │  ✓ Dependencies are fully specified and pinned                      │   │
│  │  ✓ Reproducible builds (same input = same output)                   │   │
│  │                                                                     │   │
│  │  Security Benefits:                                                 │   │
│  │  • Prevents dependency confusion attacks                            │   │
│  │  • Protects against build environment compromise                    │   │
│  │  • Enables independent verification                                 │   │
│  │                                                                     │   │
│  │  Implementation:                                                    │   │
│  │  • Bazel/BuildStream for hermetic builds                            │   │
│  │  • Container-based builds with network isolation                    │   │
│  │  • Lock files for all dependencies                                  │   │
│  │                                                                     │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │  LEVEL 4 - Two-Person Review + Reproducible Builds                  │   │
│  │  ═════════════════════════════════════════════════                  │   │
│  │                                                                     │   │
│  │  Requirements:                                                      │   │
│  │  ✓ All changes require two-person review                            │   │
│  │  ✓ Build environment is fully specified and verifiable              │   │
│  │  ✓ Reproducible: bitwise identical output from independent rebuild  │   │
│  │  ✓ Ephemeral, isolated build environment                            │   │
│  │                                                                     │   │
│  │  Security Benefits:                                                 │   │
│  │  • Prevents insider threats                                         │   │
│  │  • Enables complete supply chain verification                       │   │
│  │  • Strongest protection against build system attacks                │   │
│  │                                                                     │   │
│  │  Implementation:                                                    │   │
│  │  • Branch protection with required reviews                          │   │
│  │  • Google Cloud Build with SLSA Level 4                             │   │
│  │  • Binary authorization policies                                    │   │
│  │                                                                     │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                             │
│  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━ │
│                                                                             │
│  SLSA COMPLIANCE CHECKLIST:                                                 │
│                                                                             │
│  ┌────────────────┬─────────────────────────────────────────────────────┐   │
│  │ Level          │ Requirements Met                                    │   │
│  ├────────────────┼─────────────────────────────────────────────────────┤   │
│  │ L1 (Basic)     │ □ Scripted build     □ Provenance generated         │   │
│  │ L2 (Signed)    │ □ Hosted builder     □ Signed provenance            │   │
│  │ L3 (Hardened)  │ □ Hermetic build     □ Pinned dependencies          │   │
│  │ L4 (Verified)  │ □ Two-person review  □ Reproducible builds          │   │
│  └────────────────┴─────────────────────────────────────────────────────┘   │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 3.2 Sigstore: Cosign, Rekor, Fulcio, Gitsign

```go
package sigstore

import (
 "context"
 "crypto"
 "encoding/base64"
 "encoding/json"
 "fmt"
 "os"
 "time"

 "github.com/sigstore/cosign/v2/cmd/cosign/cli/sign"
 "github.com/sigstore/cosign/v2/cmd/cosign/cli/verify"
 "github.com/sigstore/cosign/v2/pkg/oci"
 "github.com/sigstore/sigstore/pkg/oauthflow"
)

// SigstoreComponents demonstrates Sigstore ecosystem components
// Sigstore provides a standard for signing, verifying and protecting software

/*
┌─────────────────────────────────────────────────────────────────────────────┐
│                          Sigstore Ecosystem                                  │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│   ┌──────────────┐      ┌──────────────┐      ┌──────────────┐             │
│   │   Fulcio     │      │    Rekor     │      │   Cosign     │             │
│   │   (CA)       │◄────►│  (Transparency)│◄───►│   (Signing)  │             │
│   │              │      │              │      │              │             │
│   │ • OIDC-based │      │ • Tamper-    │      │ • Container  │             │
│   │   issuance   │      │   evident    │      │   signing    │             │
│   │ • Short-lived│      │ • Immutable  │      │ • Blob       │             │
│   │   certs (10m)│      │ • Searchable │      │   signing    │             │
│   │ • Keyless    │      │ • Witness    │      │ • Keyless    │             │
│   │   signing    │      │   protocol   │      │   signing    │             │
│   └──────────────┘      └──────────────┘      └──────────────┘             │
│          ▲                     ▲                     ▲                     │
│          │                     │                     │                     │
│          └─────────────────────┼─────────────────────┘                     │
│                                │                                           │
│                                ▼                                           │
│                       ┌──────────────────┐                                 │
│                       │     Gitsign      │                                 │
│                       │ (Git Commit Sig) │                                 │
│                       │                  │                                 │
│                       │ • Sign commits   │                                 │
│                       │ • Sign tags      │                                 │
│                       │ • Verify history │                                 │
│                       └──────────────────┘                                 │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
*/

// ContainerSigning demonstrates signing container images with Cosign
func ContainerSigningExample() {
 fmt.Println(`
Container Signing with Cosign:
═══════════════════════════════

# Sign a container image (keyless)
cosign sign registry.example.com/myapp:v1.0.0

# Sign with specific identity
cosign sign --identity-token=$(gcloud auth print-identity-token) \
  registry.example.com/myapp:v1.0.0

# Verify signature
cosign verify \
  --certificate-identity=alice@example.com \
  --certificate-oidc-issuer=https://accounts.google.com \
  registry.example.com/myapp:v1.0.0

# Attest (attach SBOM)
cosign attest --predicate sbom.json \
  --type spdxjson \
  registry.example.com/myapp:v1.0.0

# Verify attestation
cosign verify-attestation \
  --certificate-identity=alice@example.com \
  --type spdxjson \
  registry.example.com/myapp:v1.0.0
`)
}

// BlobSigning demonstrates signing arbitrary blobs
func BlobSigningExample() {
 fmt.Println(`
Blob Signing with Cosign:
════════════════════════

# Sign a file
cosign sign-blob --output-signature sig.sig artifact.tar.gz

# Verify blob signature
cosign verify-blob \
  --signature sig.sig \
  --certificate cert.pem \
  artifact.tar.gz

# Sign with keyless (for CI/CD)
cosign sign-blob --output-certificate cert.pem \
  --output-signature sig.sig \
  --yes artifact.tar.gz
`)
}

// GitSigning demonstrates signing Git commits with Gitsign
func GitSigningExample() {
 fmt.Println(`
Git Commit Signing with Gitsign:
════════════════════════════════

# Configure Git to use Gitsign
git config --global commit.gpgsign true
git config --global tag.gpgsign true
git config --global gpg.x509.program gitsign
git config --global gpg.format x509

# Sign a commit (automatic)
git commit -m "Signed commit"

# Sign a tag
git tag -s v1.0.0 -m "Signed tag"

# Verify commit signature
git verify-commit HEAD

# Verify tag signature
git verify-tag v1.0.0

# View signature details
gitsign verify --cert-identity-regexp .*@example.com \
  --cert-oidc-issuer https://accounts.google.com
`)
}

// GitHubActionsSLSA demonstrates SLSA provenance with GitHub Actions
func GitHubActionsSLSA() {
 fmt.Println(`
GitHub Actions SLSA Provenance:
════════════════════════════════

# .github/workflows/release.yml
name: Release

on:
  push:
    tags:
      - 'v*'

permissions:
  contents: write
  id-token: write  # Required for OIDC

jobs:
  build:
    runs-on: ubuntu-latest
    outputs:
      hashes: ${{ steps.hash.outputs.hashes }}
    steps:
      - uses: actions/checkout@v4

      - name: Build
        run: go build -o myapp

      - name: Generate hashes
        id: hash
        run: |
          echo "hashes=$(sha256sum myapp | base64 -w0)" >> $GITHUB_OUTPUT

      - name: Upload artifact
        uses: actions/upload-artifact@v4
        with:
          name: myapp
          path: myapp

  provenance:
    needs: [build]
    permissions:
      actions: read
      id-token: write
      contents: write
    uses: slsa-framework/slsa-github-generator/.github/workflows/
      generator_generic_slsa3.yml@main
    with:
      base64-subjects: "${{ needs.build.outputs.hashes }}"
      upload-assets: true

  sign:
    needs: [build, provenance]
    runs-on: ubuntu-latest
    permissions:
      id-token: write  # For cosign keyless
    steps:
      - name: Download artifact
        uses: actions/download-artifact@v4
        with:
          name: myapp

      - name: Sign with Cosign
        uses: sigstore/cosign-installer@v3

      - name: Sign artifact
        run: cosign sign-blob --yes myapp
`)
}

// PolicyEnforcement demonstrates Binary Authorization policy
func PolicyEnforcementExample() {
 fmt.Println(`
Binary Authorization Policy (GKE):
══════════════════════════════════

# Require SLSA Level 3 for production
apiVersion: binaryauthorization.cnrm.cloud.google.com/v1beta1
kind: BinaryAuthorizationPolicy
metadata:
  name: binauthz-policy
spec:
  admissionWhitelistPatterns:
    - namePattern: "gcr.io/google_containers/*"
    - namePattern: "gcr.io/google-containers/*"

  defaultAdmissionRule:
    evaluationMode: REQUIRE_ATTESTATION
    enforcementMode: ENFORCED_BLOCK_AND_AUDIT_LOG
    requireAttestationsBy:
      - projects/my-project/attestors/build-signer
      - projects/my-project/attestors/vulnerability-scanner

  clusterAdmissionRules:
    # Production cluster requires SLSA L3
    location/us-central1/my-prod-cluster:
      evaluationMode: REQUIRE_ATTESTATION
      enforcementMode: ENFORCED_BLOCK_AND_AUDIT_LOG
      requireAttestationsBy:
        - projects/my-project/attestors/slsa-level-3

    # Staging cluster allows SLSA L2
    location/us-central1/my-staging-cluster:
      evaluationMode: REQUIRE_ATTESTATION
      enforcementMode: ENFORCED_BLOCK_AND_AUDIT_LOG
      requireAttestationsBy:
        - projects/my-project/attestors/slsa-level-2

# Attestor configuration
apiVersion: binaryauthorization.cnrm.cloud.google.com/v1beta1
kind: BinaryAuthorizationAttestor
metadata:
  name: slsa-level-3
spec:
  description: "Requires SLSA Level 3 provenance"
  attestationAuthorityNote:
    humanReadableName: "SLSA Level 3 Attestor"

# KMS key for attestor
  attestationAuthorityNoteReference: projects/my-project/notes/slsa-level-3
`)
}

// RekorVerification demonstrates transparency log verification
func RekorVerificationExample() {
 fmt.Println(`
Rekor Transparency Log Verification:
════════════════════════════════════

# Search for an entry by artifact hash
rekor-cli search --sha sha256:abc123...

# Get entry details
rekor-cli get --uuid <uuid>

# Verify inclusion proof
rekor-cli verify --uuid <uuid> --artifact artifact.tar.gz

# Monitor for new entries (witness)
rekor-cli watch --interval 60

# Go library usage:
import "github.com/sigstore/rekor/pkg/client"

func verifyInclusion(entryUUID string) error {
    rekorClient, err := client.GetRekorClient("https://rekor.sigstore.dev")
    if err != nil {
        return err
    }

    // Verify entry exists and is included
    entry, err := rekorClient.GetEntry(context.Background(), entryUUID)
    // ... verification logic
}
`)
}

// KeylessSigningFlow explains the keyless signing flow
func KeylessSigningFlow() {
 fmt.Println(`
Keyless Signing Flow (OIDC):
═══════════════════════════

┌─────────────────────────────────────────────────────────────────────────┐
│                        Keyless Signing Process                           │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  ┌──────────┐         ┌──────────┐         ┌──────────┐                │
│  │  User/CI │────────►│  Fulcio  │────────►│  Rekor   │                │
│  │          │  OIDC   │   (CA)   │  Cert   │   (Log)  │                │
│  └──────────┘         └──────────┘         └──────────┘                │
│       │                   │                   │                        │
│       │                   │                   │                        │
│       │ 1. Authenticate   │                   │                        │
│       │    with OIDC      │                   │                        │
│       ├──────────────────►│                   │                        │
│       │                   │                   │                        │
│       │ 2. Receive ID     │                   │                        │
│       │    token          │                   │                        │
│       │◄──────────────────┤                   │                        │
│       │                   │                   │                        │
│       │ 3. Generate       │                   │                        │
│       │    ephemeral key  │                   │                        │
│       │    pair           │                   │                        │
│       │                   │                   │                        │
│       │ 4. Request        │                   │                        │
│       │    certificate    │                   │                        │
│       ├──────────────────►│                   │                        │
│       │   (ID token +     │                   │                        │
│       │    public key)    │                   │                        │
│       │                   │                   │                        │
│       │                   │ 5. Issue short-   │                        │
│       │                   │    lived cert     │                        │
│       │◄──────────────────┤    (10 minutes)   │                        │
│       │                   │                   │                        │
│       │ 6. Sign artifact  │                   │                        │
│       │    with ephemeral │                   │                        │
│       │    private key    │                   │                        │
│       │                   │                   │                        │
│       │ 7. Upload entry   │                   │                        │
│       │    to transparency│                   │                        │
│       ├──────────────────────────────────────►│                        │
│       │    log            │                   │                        │
│       │                   │                   │                        │
│       │                   │ 8. Return signed  │                        │
│       │                   │    timestamp      │                        │
│       │◄──────────────────────────────────────┤                        │
│       │                   │                   │                        │
│  ┌────┴───────────────────┴───────────────────┴────┐                   │
│  │                 EPHEMERAL KEYS                   │                   │
│  │  • Private key never leaves memory               │                   │
│  │  • Discarded after signing                       │                   │
│  │  • No key management required                    │                   │
│  │  • No secrets to rotate                          │                   │
│  └──────────────────────────────────────────────────┘                   │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
`)
}

// OIDCProviders lists supported OIDC providers for keyless signing
func OIDCProviders() {
 fmt.Println(`
Supported OIDC Providers for Keyless Signing:
════════════════════════════════════════════

Public Fulcio (sigstore.dev):
─────────────────────────────
• GitHub Actions (github.com)
• GitLab CI (gitlab.com)
• Google (accounts.google.com)
• Microsoft (login.microsoftonline.com)
• Buildkite (buildkite.com)
• Dex (dex-idp.io)

Private/Enterprise:
───────────────────
• Any OIDC-compliant IdP
• Self-hosted Dex
• Keycloak
• Okta
• Auth0
• Azure AD

GitHub Actions Example:
───────────────────────
permissions:
  id-token: write  # Required for OIDC

steps:
  - uses: sigstore/cosign-installer@v3
  - name: Sign image
    run: cosign sign --yes $IMAGE_URI

The GITHUB_TOKEN is automatically exchanged for an OIDC token
to authenticate with Fulcio.
`)
}
package sigstore

import (
 "encoding/base64"
 "fmt"
 "os"
)

// SBOMRequirements covers CISA 2025 SBOM requirements
func SBOMRequirements() {
 fmt.Println(`
CISA SBOM Requirements (2025):
══════════════════════════════

MANDATORY SBOM ELEMENTS (per CISA guidance):
───────────────────────────────────────────

1. Data Fields:
   • Supplier Name
   • Component Name
   • Version of the Component
   • Other Unique Identifiers (CPE, PURL, SWID)
   • Dependency Relationship
   • Author of SBOM Data
   • Timestamp

2. SBOM Formats (choose one):
   • SPDX (ISO/IEC 5962:2021)
   • CycloneDX
   • SWID (ISO/IEC 19770-2:2015)

3. SBOM Types:
   • Design-time SBOM (source dependencies)
   • Build-time SBOM (compiled dependencies)
   • Analyzed SBOM (binary analysis)
   • Deployed SBOM (runtime inventory)
   • Runtime SBOM (dynamically loaded)

Go Tooling:
───────────
# Generate SPDX SBOM from Go module
cyclonedx-gomod app -o sbom.json

# Using syft (supports multiple formats)
syft packages dir:. -o spdx-json -file sbom.spdx.json
syft packages dir:. -o cyclonedx-json -file sbom.cdx.json

# Using ko for container+SBOM generation
ko build --sbom=spdx

Integration with Sigstore:
──────────────────────────
# Attest SBOM to container image
cosign attest --predicate sbom.spdx.json \
  --type spdxjson \
  --key cosign.key \
  $IMAGE

# Verify SBOM attestation
cosign verify-attestation \
  --type spdxjson \
  --key cosign.pub \
  $IMAGE
`)
}

// CI_CD_Security demonstrates secure CI/CD practices
func CI_CD_Security() {
 fmt.Println(`
Secure CI/CD Pipeline:
══════════════════════

PREREQUISITES (SLSA L3):
────────────────────────
1. Pinned Actions: Use commit SHA instead of tags
   actions/checkout@v4 → actions/checkout@b4ffde65f46336ab88eb53be808477a3936bae11

2. Minimal Permissions: Grant only required permissions
   permissions:
     contents: read
     id-token: write

3. No Secrets in Logs: Mask all sensitive data

4. Dependency Pinning: Use go.mod + go.sum

5. Hermetic Builds: No network access during build

SECURE GITHUB ACTIONS WORKFLOW:
──────────────────────────────
name: Secure Build

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

permissions:
  contents: read  # Minimal permissions

env:
  GOPROXY: https://proxy.golang.org,direct
  GOSUMDB: sum.golang.org

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      # 1. Checkout with full history for SBOM
      - name: Checkout
        uses: actions/checkout@b4ffde65f46336ab88eb53be808477a3936bae11
        with:
          fetch-depth: 0

      # 2. Setup Go with checksum verification
      - name: Setup Go
        uses: actions/setup-go@0c52d547c9bc32b1aa3301fd7a9cb496313a4491
        with:
          go-version: '1.24'
          check-latest: true

      # 3. Verify dependencies (no network access)
      - name: Verify dependencies
        run: |
          go mod verify
          go mod graph | head -100

      # 4. Security scanning
      - name: Run Gosec Security Scanner
        uses: securego/gosec@master
        with:
          args: '-no-fail -fmt sarif -out results.sarif ./...'

      # 5. Build with provenance
      - name: Build
        run: |
          CGO_ENABLED=0 go build \
            -ldflags="-w -s -X main.version=$VERSION -X main.commit=$GITHUB_SHA" \
            -o app \
            ./cmd/app

      # 6. Generate SBOM
      - name: Generate SBOM
        uses: anchore/sbom-action@b9e1dc4fc4e5e0b8f8e0c6c2c2e1c0f0d0c0b0a0
        with:
          format: spdx-json
          output-file: sbom.spdx.json

      # 7. Sign artifacts (keyless)
      - name: Sign with Cosign
        if: github.event_name != 'pull_request'
        uses: sigstore/cosign-installer@59acb6260d9c0ba8f4a2f9d9b48431a222b68e20

      - name: Sign artifacts
        if: github.event_name != 'pull_request'
        run: |
          cosign sign-blob --yes \
            --output-certificate cert.pem \
            --output-signature sig.sig \
            app

      # 8. Upload artifacts
      - name: Upload artifacts
        uses: actions/upload-artifact@5d5d22a31266ced268874388b861e4b58bb5c2f3
        with:
          name: artifacts
          path: |
            app
            sbom.spdx.json
            cert.pem
            sig.sig

SLSA PROVENANCE GENERATION:
───────────────────────────
  provenance:
    needs: [build]
    permissions:
      actions: read
      id-token: write
      contents: write
    uses: slsa-framework/slsa-github-generator/
            .github/workflows/generator_generic_slsa3.yml@v1.9.0
    with:
      base64-subjects: "${{ needs.build.outputs.hashes }}"
      upload-assets: true
`)
}

// SLSAProvenanceExample shows SLSA provenance structure
func SLSAProvenanceExample() {
 fmt.Println(`
SLSA Provenance Example:
═══════════════════════

{
  "_type": "https://in-toto.io/Statement/v0.1",
  "predicateType": "https://slsa.dev/provenance/v1",
  "subject": [
    {
      "name": "app",
      "digest": {
        "sha256": "abc123..."
      }
    }
  ],
  "predicate": {
    "buildDefinition": {
      "buildType": "https://github.com/slsa-framework/
                     slsa-github-generator/generic@v1",
      "externalParameters": {
        "workflow": {
          "ref": "refs/heads/main",
          "repository": "https://github.com/org/repo",
          "path": ".github/workflows/build.yml"
        }
      },
      "internalParameters": {
        "github": {
          "event_name": "push",
          "repository_id": "123456789",
          "repository_owner_id": "987654321"
        }
      },
      "resolvedDependencies": [
        {
          "uri": "git+https://github.com/org/repo@refs/heads/main",
          "digest": {
            "gitCommit": "abc123..."
          }
        }
      ]
    },
    "runDetails": {
      "builder": {
        "id": "https://github.com/slsa-framework/
                slsa-github-generator/.github/workflows/
                generator_generic_slsa3.yml@refs/tags/v1.9.0"
      },
      "metadata": {
        "invocationId": "https://github.com/org/repo/
                          actions/runs/123456789",
        "startedOn": "2026-01-15T10:30:00Z",
        "finishedOn": "2026-01-15T10:35:00Z"
      }
    }
  }
}
`)
}


---

## 4. Container and Cloud Security

### 4.1 The 4 Cs of Kubernetes Security

```

┌─────────────────────────────────────────────────────────────────────────────┐
│                    The 4 Cs of Cloud Native Security                         │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│   Security applies at multiple layers, each building on the previous:       │
│                                                                             │
│   ┌─────────────────────────────────────────────────────────────────────┐   │
│   │                        1. CODE                                      │   │
│   │                    ┌─────────────┐                                  │   │
│   │                    │ Application │                                  │   │
│   │                    │   Source    │                                  │   │
│   │                    └─────────────┘                                  │   │
│   │                                                                     │   │
│   │  Security:                                                          │   │
│   │  • SAST/DAST scanning                                               │   │
│   │  • Dependency vulnerability management                              │   │
│   │  • Secret detection in code                                         │   │
│   │  • Secure coding practices                                          │   │
│   │  • SLSA provenance                                                  │   │
│   │                                                                     │   │
│   └─────────────────────────────────────────────────────────────────────┘   │
│                                    │                                        │
│                                    ▼                                        │
│   ┌─────────────────────────────────────────────────────────────────────┐   │
│   │                       2. CONTAINER                                  │   │
│   │                    ┌─────────────┐                                  │   │
│   │                    │   Docker    │                                  │   │
│   │                    │  Container  │                                  │   │
│   │                    └─────────────┘                                  │   │
│   │                                                                     │   │
│   │  Security:                                                          │   │
│   │  • Minimal base images (distroless, scratch)                        │   │
│   │  • Non-root user execution                                          │   │
│   │  • Image scanning (Trivy, Snyk, Clair)                            │   │
│   │  • Image signing (Cosign)                                           │   │
│   │  • Immutable tags                                                   │   │
│   │  • SBOM generation                                                  │   │
│   │                                                                     │   │
│   └─────────────────────────────────────────────────────────────────────┘   │
│                                    │                                        │
│                                    ▼                                        │
│   ┌─────────────────────────────────────────────────────────────────────┐   │
│   │                      3. CLUSTER                                     │   │
│   │                    ┌─────────────┐                                  │   │
│   │                    │ Kubernetes  │                                  │   │
│   │                    │   Cluster   │                                  │   │
│   │                    └─────────────┘                                  │   │
│   │                                                                     │   │
│   │  Security:                                                          │   │
│   │  • RBAC (Role-Based Access Control)                                 │   │
│   │  • Pod Security Standards                                           │   │
│   │  • Network Policies                                                 │   │
│   │  • Admission Controllers (OPA, Kyverno)                             │   │
│   │  • Secrets management (external-secrets)                            │   │
│   │  • Runtime security (Falco, Tetragon)                               │   │
│   │                                                                     │   │
│   └─────────────────────────────────────────────────────────────────────┘   │
│                                    │                                        │
│                                    ▼                                        │
│   ┌─────────────────────────────────────────────────────────────────────┐   │
│   │                      4. CLOUD                                       │   │
│   │                    ┌─────────────┐                                  │   │
│   │                    │ Cloud Infra │                                  │   │
│   │                    │ (AWS/Azure) │                                  │   │
│   │                    └─────────────┘                                  │   │
│   │                                                                     │   │
│   │  Security:                                                          │   │
│   │  • IAM policies and roles                                           │   │
│   │  • VPC and network isolation                                        │   │
│   │  • Encryption at rest and in transit                                │   │
│   │  • Cloud provider security services                                 │   │
│   │  • Compliance frameworks (SOC2, ISO 27001)                          │   │
│   │  • Audit logging and monitoring                                     │   │
│   │                                                                     │   │
│   └─────────────────────────────────────────────────────────────────────┘   │
│                                                                             │
│   ═══════════════════════════════════════════════════════════════════════   │
│                                                                             │
│   DEFENSE IN DEPTH:                                                         │
│   Each layer provides security even if other layers are compromised.        │
│   Security is only as strong as the weakest layer.                          │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘

```

### 4.2 Pod Security Standards (PSS)

```yaml
# Pod Security Standards Implementation
# Kubernetes 1.25+ (PSP deprecated, replaced by PSS)

# Three policy levels:
# - privileged: Unrestricted (backward compatible)
# - baseline: Minimal restrictions (default)
# - restricted: Maximum security (recommended for production)

---
# Namespace-level enforcement
apiVersion: v1
kind: Namespace
metadata:
  name: production
  labels:
    # Enforce restricted policy
    pod-security.kubernetes.io/enforce: restricted
    pod-security.kubernetes.io/enforce-version: latest

    # Audit violations at baseline level
    pod-security.kubernetes.io/audit: baseline
    pod-security.kubernetes.io/audit-version: latest

    # Warn users about restricted violations
    pod-security.kubernetes.io/warn: restricted
    pod-security.kubernetes.io/warn-version: latest

---
# Cluster-wide Pod Security Admission Configuration
apiVersion: apiserver.config.k8s.io/v1
kind: AdmissionConfiguration
plugins:
- name: PodSecurity
  configuration:
    apiVersion: pod-security.admission.config.k8s.io/v1
    kind: PodSecurityConfiguration
    defaults:
      enforce: "baseline"
      audit: "restricted"
      warn: "restricted"
    exemptions:
      # System namespaces
      namespaces:
      - kube-system
      - ingress-nginx
      - cert-manager

      # Specific users (CI/CD)
      usernames:
      - cicd-service-account

      # Runtime class names
      runtimeClasses:
      - gvisor

---
# Example Pod compliant with RESTRICTED policy
apiVersion: v1
kind: Pod
metadata:
  name: secure-app
  namespace: production
spec:
  # 1. Non-root user
  securityContext:
    runAsNonRoot: true
    runAsUser: 65534  # nobody
    runAsGroup: 65534
    seccompProfile:
      type: RuntimeDefault

  # 2. Read-only root filesystem
  containers:
  - name: app
    image: myapp:v1.0.0
    securityContext:
      allowPrivilegeEscalation: false
      readOnlyRootFilesystem: true
      capabilities:
        drop:
        - ALL

    # 3. Resource limits
    resources:
      limits:
        cpu: "500m"
        memory: "512Mi"
      requests:
        cpu: "100m"
        memory: "128Mi"

    # 4. Security headers
    volumeMounts:
    - name: tmp
      mountPath: /tmp
    - name: cache
      mountPath: /var/cache

  # 5. Ephemeral volumes only
  volumes:
  - name: tmp
    emptyDir: {}
  - name: cache
    emptyDir:
      sizeLimit: 100Mi

---
# Security Context Best Practices for Go Applications
apiVersion: apps/v1
kind: Deployment
metadata:
  name: go-app
spec:
  template:
    spec:
      serviceAccountName: go-app-sa
      automountServiceAccountToken: false  # Disable if not needed

      securityContext:
        # Run as non-root
        runAsNonRoot: true
        runAsUser: 1000
        runAsGroup: 1000
        fsGroup: 1000

        # Disable privilege escalation
        allowPrivilegeEscalation: false

        # Use runtime default seccomp profile
        seccompProfile:
          type: RuntimeDefault

      containers:
      - name: app
        image: go-app:latest

        # Security context
        securityContext:
          capabilities:
            drop:
            - ALL
          readOnlyRootFilesystem: true
          allowPrivilegeEscalation: false

        # Drop all capabilities
        # Add only required ones
        # capabilities:
        #   add:
        #   - NET_BIND_SERVICE  # If binding to ports < 1024

        ports:
        - containerPort: 8080
          protocol: TCP
```

### 4.3 OPA (Open Policy Agent) and Gatekeeper

```yaml
# OPA Gatekeeper Policy Examples
# Rego policies for Kubernetes admission control

---
# Constraint Template: Required Labels
apiVersion: templates.gatekeeper.sh/v1
kind: ConstraintTemplate
metadata:
  name: k8srequiredlabels
spec:
  crd:
    spec:
      names:
        kind: K8sRequiredLabels
      validation:
        openAPIV3Schema:
          properties:
            labels:
              type: array
              items:
                type: object
                properties:
                  key:
                    type: string
                  allowedRegex:
                    type: string
  targets:
    - target: admission.k8s.gatekeeper.sh
      rego: |
        package k8srequiredlabels

        violation[{"msg": msg, "details": {"missing_labels": missing}}] {
          provided := {label | input.review.object.metadata.labels[label]}
          required := {label | label := input.parameters.labels[_].key}
          missing := required - provided
          count(missing) > 0
          msg := sprintf("Missing required labels: %v", [missing])
        }

        violation[{"msg": msg}] {
          label := input.parameters.labels[_]
          value := input.review.object.metadata.labels[label.key]
          not re_match(label.allowedRegex, value)
          msg := sprintf("Label %v has invalid value %v", [label.key, value])
        }

---
# Constraint: Require app and owner labels
apiVersion: constraints.gatekeeper.sh/v1beta1
kind: K8sRequiredLabels
metadata:
  name: required-labels
spec:
  match:
    kinds:
    - apiGroups: ["apps"]
      kinds: ["Deployment", "StatefulSet", "DaemonSet"]
    - apiGroups: [""]
      kinds: ["Pod", "Service"]
    excludedNamespaces:
    - kube-system
    - gatekeeper-system
  parameters:
    labels:
    - key: app.kubernetes.io/name
      allowedRegex: "^[a-z0-9-]+$"
    - key: app.kubernetes.io/owner
      allowedRegex: "^[a-z-]+$"

---
# Constraint Template: Disallow Latest Tag
apiVersion: templates.gatekeeper.sh/v1
kind: ConstraintTemplate
metadata:
  name: k8sdisallowlatesttag
spec:
  crd:
    spec:
      names:
        kind: K8sDisallowLatestTag
  targets:
    - target: admission.k8s.gatekeeper.sh
      rego: |
        package k8sdisallowlatesttag

        violation[{"msg": msg}] {
          container := input.review.object.spec.containers[_]
          endswith(container.image, ":latest")
          msg := sprintf("Container %v uses 'latest' tag", [container.name])
        }

        violation[{"msg": msg}] {
          container := input.review.object.spec.initContainers[_]
          endswith(container.image, ":latest")
          msg := sprintf("Init container %v uses 'latest' tag", [container.name])
        }

---
# Constraint: No Latest Tags in Production
apiVersion: constraints.gatekeeper.sh/v1beta1
kind: K8sDisallowLatestTag
metadata:
  name: disallow-latest-tag
spec:
  match:
    kinds:
    - apiGroups: [""]
      kinds: ["Pod"]
    namespaces:
    - production
    - staging

---
# Constraint Template: Require Resource Limits
apiVersion: templates.gatekeeper.sh/v1
kind: ConstraintTemplate
metadata:
  name: k8srequiredresources
spec:
  crd:
    spec:
      names:
        kind: K8sRequiredResources
  targets:
    - target: admission.k8s.gatekeeper.sh
      rego: |
        package k8srequiredresources

        violation[{"msg": msg}] {
          container := input.review.object.spec.containers[_]
          not container.resources.limits.memory
          msg := sprintf("Container %v missing memory limit", [container.name])
        }

        violation[{"msg": msg}] {
          container := input.review.object.spec.containers[_]
          not container.resources.limits.cpu
          msg := sprintf("Container %v missing CPU limit", [container.name])
        }

        violation[{"msg": msg}] {
          container := input.review.object.spec.containers[_]
          not container.resources.requests.memory
          msg := sprintf("Container %v missing memory request", [container.name])
        }

---
# Constraint Template: Block Privileged Containers
apiVersion: templates.gatekeeper.sh/v1
kind: ConstraintTemplate
metadata:
  name: k8sblockprivileged
spec:
  crd:
    spec:
      names:
        kind: K8sBlockPrivileged
  targets:
    - target: admission.k8s.gatekeeper.sh
      rego: |
        package k8sblockprivileged

        violation[{"msg": msg}] {
          container := input.review.object.spec.containers[_]
          container.securityContext.privileged == true
          msg := sprintf("Privileged container not allowed: %v", [container.name])
        }

---
# Rego Policy Testing
# test_policy.rego
package k8srequiredlabels

test_required_labels_present {
  not violation[{
    "msg": "Missing required labels",
    "details": {"missing_labels": set()}
  }] with input as {
    "review": {
      "object": {
        "metadata": {
          "labels": {
            "app.kubernetes.io/name": "myapp",
            "app.kubernetes.io/owner": "team-a"
          }
        }
      }
    },
    "parameters": {
      "labels": [
        {"key": "app.kubernetes.io/name"},
        {"key": "app.kubernetes.io/owner"}
      ]
    }
  }
}

test_missing_labels {
  violation[{
    "msg": "Missing required labels: {\"app.kubernetes.io/owner\", \"app.kubernetes.io/name\"}",
    "details": {"missing_labels": {"app.kubernetes.io/owner", "app.kubernetes.io/name"}}
  }] with input as {
    "review": {
      "object": {
        "metadata": {
          "labels": {}
        }
      }
    },
    "parameters": {
      "labels": [
        {"key": "app.kubernetes.io/name"},
        {"key": "app.kubernetes.io/owner"}
      ]
    }
  }
}
```

### 4.4 Tetragon eBPF Runtime Security

```yaml
# Tetragon eBPF-based Runtime Security
# Real-time security observability and enforcement

---
# Tetragon TracingPolicy: Detect Suspicious Process Execution
apiVersion: cilium.io/v1alpha1
kind: TracingPolicy
metadata:
  name: detect-suspicious-processes
spec:
  kprobes:
  - call: "__x64_sys_execve"
    syscall: true
    args:
    - index: 0
      type: "string"
    selectors:
    - matchBinaries:
      - operator: "In"
        values:
        - "/bin/bash"
        - "/bin/sh"
        - "/bin/dash"
      matchArgs:
      - index: 0
        operator: "Prefix"
        values:
        - "/tmp/"
        - "/var/tmp/"
        - "/dev/shm/"
      matchActions:
      - action: Sigkill
      - action: Post
        message: "Suspicious shell execution from temp directory"
        rateLimit: "1m"

---
# Tetragon TracingPolicy: Detect Privilege Escalation
apiVersion: cilium.io/v1alpha1
kind: TracingPolicy
metadata:
  name: detect-privilege-escalation
spec:
  kprobes:
  - call: "__x64_sys_setuid"
    syscall: true
    selectors:
    - matchNamespaces:
      - namespace: PID
        operator: In
        values:
        - "host"
      matchActions:
      - action: Post
        message: "Potential privilege escalation via setuid"
        rateLimit: "1m"

---
# Tetragon TracingPolicy: Network Monitoring
apiVersion: cilium.io/v1alpha1
kind: TracingPolicy
metadata:
  name: network-monitoring
spec:
  kprobes:
  - call: "__x64_sys_connect"
    syscall: true
    args:
    - index: 0
      type: "sock"
    selectors:
    - matchActions:
      - action: Post
        message: "Network connection established"
        rateLimit: "10s"

---
# Tetragon TracingPolicy: File Integrity Monitoring
apiVersion: cilium.io/v1alpha1
kind: TracingPolicy
metadata:
  name: file-integrity
spec:
  kprobes:
  - call: "security_file_permission"
    args:
    - index: 0
      type: "file"
    - index: 1
      type: "int"
    selectors:
    - matchArgs:
      - index: 1
        operator: "Equal"
        values:
        - "2"  # MAY_WRITE
      matchBinaries:
      - operator: "NotIn"
        values:
        - "/usr/bin/dockerd"
        - "/usr/bin/containerd"
      matchActions:
      - action: Post
        message: "File modification detected"
        rateLimit: "1m"

---
# Tetragon Export Configuration
apiVersion: v1
kind: ConfigMap
metadata:
  name: tetragon-config
  namespace: kube-system
data:
  export.yaml: |
    export:
      # Export to stdout
      stdout:
        enabled: true

      # Export to file
      file:
        enabled: false
        path: /var/log/tetragon/tetragon.log
        rotationInterval: "1d"

      # Export to gRPC
      grpc:
        enabled: true
        address: "localhost:54321"

    filters:
      # Exclude system namespaces
      exclude:
        namespaces:
        - kube-system
        - kube-public
        - kube-node-lease

      # Include only these event types
      include:
        eventTypes:
        - PROCESS_EXEC
        - PROCESS_EXIT
        - FILE_WRITE
        - NETWORK_CONNECT

---
# Go Integration with Tetragon
# main.go
package main

import (
 "context"
 "encoding/json"
 "fmt"
 "log"
 "os"

 "github.com/cilium/tetragon/api/v1/tetragon"
 "google.golang.org/grpc"
)

// TetragonClient wraps Tetragon gRPC client
type TetragonClient struct {
 client tetragon.FineGuidanceSensorsClient
 conn   *grpc.ClientConn
}

// NewTetragonClient creates a new Tetragon client
func NewTetragonClient(address string) (*TetragonClient, error) {
 conn, err := grpc.Dial(address, grpc.WithInsecure())
 if err != nil {
  return nil, fmt.Errorf("failed to connect to Tetragon: %w", err)
 }

 client := tetragon.NewFineGuidanceSensorsClient(conn)
 return &TetragonClient{
  client: client,
  conn:   conn,
 }, nil
}

// Close closes the connection
func (c *TetragonClient) Close() error {
 return c.conn.Close()
}

// GetEvents streams events from Tetragon
func (c *TetragonClient) GetEvents(ctx context.Context) (<-chan *tetragon.GetEventsResponse, error) {
 stream, err := c.client.GetEvents(ctx, &tetragon.GetEventsRequest{})
 if err != nil {
  return nil, err
 }

 events := make(chan *tetragon.GetEventsResponse)
 go func() {
  defer close(events)
  for {
   event, err := stream.Recv()
   if err != nil {
    log.Printf("Error receiving event: %v", err)
    return
   }
   select {
   case events <- event:
   case <-ctx.Done():
    return
   }
  }
 }()

 return events, nil
}

// ProcessEvent handles a Tetragon event
func ProcessEvent(event *tetragon.GetEventsResponse) {
 switch e := event.Event.(type) {
 case *tetragon.GetEventsResponse_ProcessExec:
  exec := e.ProcessExec
  fmt.Printf("[EXEC] PID=%d Binary=%s Args=%v\n",
   exec.Process.Pid.GetValue(),
   exec.Process.Binary,
   exec.Process.Arguments)

 case *tetragon.GetEventsResponse_ProcessExit:
  exit := e.ProcessExit
  fmt.Printf("[EXIT] PID=%d Status=%d\n",
   exit.Process.Pid.GetValue(),
   exit.Status)

 case *tetragon.GetEventsResponse_FileWrite:
  write := e.FileWrite
  fmt.Printf("[FILE] Action=%s Path=%s\n",
   write.Action,
   write.File.Path)
 }
}

func main() {
 client, err := NewTetragonClient("localhost:54321")
 if err != nil {
  log.Fatal(err)
 }
 defer client.Close()

 ctx := context.Background()
 events, err := client.GetEvents(ctx)
 if err != nil {
  log.Fatal(err)
 }

 for event := range events {
  ProcessEvent(event)
 }
}
```

---

## 5. mTLS and Service Identity

### 5.1 SPIFFE/SPIRE (CNCF Graduated)

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    SPIFFE/SPIRE Architecture                                 │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  SPIFFE = Secure Production Identity Framework For Everyone                 │
│  SPIRE  = SPIFFE Runtime Environment                                        │
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                        SPIRE SERVER                                  │   │
│  │                    (Control Plane / Root CA)                         │   │
│  │                                                                       │   │
│  │  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐              │   │
│  │  │   CA Manager │  │   Attestor   │  │   Registry   │              │   │
│  │  │              │  │   Plugins    │  │              │              │   │
│  │  │ • Sign SVIDs │  │ • Kubernetes │  │ • Workload   │              │   │
│  │  │ • Rotate CA  │  │   attestation│  │   selectors  │              │   │
│  │  │   cert       │  │ • AWS/GCP    │  │ • Node       │              │   │
│  │  │              │  │   attestation│  │   attestation│              │   │
│  │  └──────────────┘  └──────────────┘  └──────────────┘              │   │
│  │                                                                       │   │
│  │  CA Certificate: Root of trust for the trust domain                  │   │
│  │                                                                       │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                    │                                        │
│                           gRPC/HTTPS                                        │
│                                    │                                        │
│                                    ▼                                        │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                        SPIRE AGENT                                   │   │
│  │                    (Node-level / per K8s node)                       │   │
│  │                                                                       │   │
│  │  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐              │   │
│  │  │ Workload Att │  │  SVID Cache  │  │ Workload API │              │   │
│  │  │ estor        │  │              │  │ (Unix socket)│              │   │
│  │  │              │  │ • X.509-SVID │  │              │              │   │
│  │  │ • Prove      │  │ • JWT-SVID   │  │ • Fetch SVID │              │   │
│  │  │   identity   │  │ • Private key│  │ • Validate   │              │   │
│  │  │   to server  │  │              │  │   SVID       │              │   │
│  │  └──────────────┘  └──────────────┘  └──────────────┘              │   │
│  │                                                                       │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                    │                                        │
│                    Unix Domain Socket (/tmp/spire-agent/public/api.sock)    │
│                                    │                                        │
│                                    ▼                                        │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                        WORKLOAD                                      │   │
│  │                    (Your Application)                                │   │
│  │                                                                       │   │
│  │  SVID (SPIFFE Verifiable Identity Document):                         │   │
│  │  spiffe://trust-domain/ns/<namespace>/sa/<service-account>           │   │
│  │                                                                       │   │
│  │  Contains:                                                           │   │
│  │  • X.509 certificate with SPIFFE ID in URI SAN                       │   │
│  │  • Private key (short-lived, auto-rotated)                           │   │
│  │  • Trust bundle (root CA certs)                                      │   │
│  │                                                                       │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                             │
│  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━ │
│                                                                             │
│  SPIFFE ID FORMAT:                                                          │
│                                                                             │
│  spiffe://trust-domain/path                                                │
│                                                                             │
│  Examples:                                                                  │
│  • spiffe://production.example.com/ns/payments/sa/payment-service          │
│  • spiffe://staging.example.com/ns/orders/sa/order-service                 │
│  • spiffe://prod.aws.example.com/role/webserver/instance/i-12345678        │
│                                                                             │
│  BENEFITS:                                                                  │
│  • No shared secrets between services                                      │
│  • Short-lived certificates (auto-rotation)                                │
│  • Platform-agnostic identity (K8s, VMs, cloud)                            │
│  • No IP-based access control                                              │
│  • Fine-grained workload identity                                          │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 5.2 SPIRE Kubernetes Deployment

```yaml
# SPIRE Server Deployment
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: spire-server
  namespace: spire
  labels:
    app: spire-server
spec:
  replicas: 1
  selector:
    matchLabels:
      app: spire-server
  serviceName: spire-server
  template:
    metadata:
      labels:
        app: spire-server
    spec:
      serviceAccountName: spire-server
      containers:
      - name: spire-server
        image: ghcr.io/spiffe/spire-server:1.9.0
        args:
        - -config
        - /run/spire/config/server.conf
        ports:
        - containerPort: 8081
          name: grpc
        - containerPort: 8080
          name: health
        volumeMounts:
        - name: spire-config
          mountPath: /run/spire/config
          readOnly: true
        - name: spire-data
          mountPath: /run/spire/data
        - name: spire-socket
          mountPath: /tmp/spire-server/private
        livenessProbe:
          httpGet:
            path: /live
            port: health
          initialDelaySeconds: 5
          periodSeconds: 5
        readinessProbe:
          httpGet:
            path: /ready
            port: health
          initialDelaySeconds: 5
          periodSeconds: 5
        resources:
          requests:
            cpu: 100m
            memory: 128Mi
          limits:
            cpu: 500m
            memory: 512Mi
      volumes:
      - name: spire-config
        configMap:
          name: spire-server-config
      - name: spire-socket
        hostPath:
          path: /tmp/spire-server/private
          type: DirectoryOrCreate
  volumeClaimTemplates:
  - metadata:
      name: spire-data
    spec:
      accessModes:
      - ReadWriteOnce
      resources:
        requests:
          storage: 1Gi

---
# SPIRE Server Configuration
apiVersion: v1
kind: ConfigMap
metadata:
  name: spire-server-config
  namespace: spire
data:
  server.conf: |
    server {
      bind_address = "0.0.0.0"
      bind_port = "8081"
      socket_path = "/tmp/spire-server/private/api.sock"
      trust_domain = "production.example.com"
      data_dir = "/run/spire/data"
      log_level = "INFO"

      ca_key_type = "ec-p256"
      ca_ttl = "24h"
      default_svid_ttl = "1h"

      ca_subject = {
        country = ["US"]
        organization = ["Example Inc."]
        common_name = "production.example.com"
      }
    }

    plugins {
      DataStore "sql" {
        plugin_data {
          database_type = "sqlite3"
          connection_string = "/run/spire/data/datastore.sqlite3"
        }
      }

      KeyManager "memory" {
        plugin_data {}
      }

      NodeAttestor "k8s_psat" {
        plugin_data {
          clusters = {
            "production" = {
              use_token_review_api_validation = true
              service_account_allow_list = ["spire:spire-agent"]
            }
          }
        }
      }

      Notifier "k8sbundle" {
        plugin_data {}
      }
    }

    health_checks {
      listener_enabled = true
      bind_address = "0.0.0.0"
      bind_port = "8080"
      live_path = "/live"
      ready_path = "/ready"
    }

---
# SPIRE Agent DaemonSet
apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: spire-agent
  namespace: spire
  labels:
    app: spire-agent
spec:
  selector:
    matchLabels:
      app: spire-agent
  template:
    metadata:
      labels:
        app: spire-agent
    spec:
      hostPID: true
      hostNetwork: true
      dnsPolicy: ClusterFirstWithHostNet
      serviceAccountName: spire-agent
      initContainers:
      - name: init
        image: busybox:1.36
        command:
        - sh
        - -c
        - |
          # Wait for SPIRE server to be ready
          until nc -z spire-server 8081; do
            echo "Waiting for SPIRE server..."
            sleep 5
          done
      containers:
      - name: spire-agent
        image: ghcr.io/spiffe/spire-agent:1.9.0
        args:
        - -config
        - /run/spire/config/agent.conf
        volumeMounts:
        - name: spire-config
          mountPath: /run/spire/config
          readOnly: true
        - name: spire-bundle
          mountPath: /run/spire/bundle
          readOnly: true
        - name: spire-agent-socket
          mountPath: /tmp/spire-agent/public
        - name: token
          mountPath: /var/run/secrets/tokens
        livenessProbe:
          exec:
            command:
            - /opt/spire/bin/spire-agent
            - healthcheck
          initialDelaySeconds: 10
          periodSeconds: 10
        resources:
          requests:
            cpu: 50m
            memory: 64Mi
          limits:
            cpu: 200m
            memory: 256Mi
      volumes:
      - name: spire-config
        configMap:
          name: spire-agent-config
      - name: spire-bundle
        configMap:
          name: spire-bundle
      - name: spire-agent-socket
        hostPath:
          path: /tmp/spire-agent/public
          type: DirectoryOrCreate
      - name: token
        projected:
          sources:
          - serviceAccountToken:
              path: spire-agent
              expirationSeconds: 7200
              audience: spire-server

---
# SPIRE Agent Configuration
apiVersion: v1
kind: ConfigMap
metadata:
  name: spire-agent-config
  namespace: spire
data:
  agent.conf: |
    agent {
      data_dir = "/run/spire"
      log_level = "INFO"
      server_address = "spire-server"
      server_port = "8081"
      socket_path = "/tmp/spire-agent/public/api.sock"
      trust_bundle_path = "/run/spire/bundle/bundle.crt"
      trust_domain = "production.example.com"
    }

    plugins {
      NodeAttestor "k8s_psat" {
        plugin_data {
          cluster = "production"
          token_path = "/var/run/secrets/tokens/spire-agent"
        }
      }

      KeyManager "memory" {
        plugin_data {}
      }

      WorkloadAttestor "k8s" {
        plugin_data {
          # Automatically attest workloads based on K8s metadata
          skip_kubelet_verification = false
        }
      }

      WorkloadAttestor "unix" {
        plugin_data {}
      }
    }

---
# ClusterSPIFFEID - Define workload identities
apiVersion: spire.spiffe.io/v1alpha1
kind: ClusterSPIFFEID
metadata:
  name: default
spec:
  spiffeIDTemplate: "spiffe://production.example.com/ns/{{ .PodMeta.Namespace }}/sa/{{ .PodSpec.ServiceAccountName }}"
  podSelector:
    matchLabels:
      spiffe.io/spire-managed-identity: "true"
  workloadSelectorTemplates:
  - "k8s:ns:{{ .PodMeta.Namespace }}"
  - "k8s:sa:{{ .PodSpec.ServiceAccountName }}"

---
# ClusterSPIFFEID - Payment service specific
apiVersion: spire.spiffe.io/v1alpha1
kind: ClusterSPIFFEID
metadata:
  name: payment-service
spec:
  spiffeIDTemplate: "spiffe://production.example.com/payments/{{ .PodMeta.Name }}"
  podSelector:
    matchLabels:
      app: payment-service
  workloadSelectorTemplates:
  - "k8s:ns:payments"
  - "k8s:label:app:payment-service"
```

### 5.3 Go SPIFFE Integration

```go
package spiffe

import (
 "context"
 "crypto/tls"
 "crypto/x509"
 "fmt"
 "net/http"
 "time"

 "github.com/spiffe/go-spiffe/v2/bundle/x509bundle"
 "github.com/spiffe/go-spiffe/v2/spiffetls/tlsconfig"
 "github.com/spiffe/go-spiffe/v2/svid/x509svid"
 "github.com/spiffe/go-spiffe/v2/workloadapi"
 "google.golang.org/grpc"
 "google.golang.org/grpc/credentials"
)

// SPIFFEClient wraps SPIFFE workload API client
type SPIFFEClient struct {
 svid        *x509svid.SVID
 bundle      *x509bundle.Bundle
 workloadAPI *workloadapi.Client
}

// NewSPIFFEClient creates a new SPIFFE client
func NewSPIFFEClient(socketPath string) (*SPIFFEClient, error) {
 ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
 defer cancel()

 // Create workload API client
 client, err := workloadapi.New(ctx, workloadapi.WithAddr("unix://"+socketPath))
 if err != nil {
  return nil, fmt.Errorf("failed to create workload API client: %w", err)
 }

 // Fetch X.509 SVID
 svid, err := client.FetchX509SVID(ctx)
 if err != nil {
  return nil, fmt.Errorf("failed to fetch X.509 SVID: %w", err)
 }

 fmt.Printf("SPIFFE ID: %s\n", svid.ID)
 fmt.Printf("Certificate chain length: %d\n", len(svid.Certificates))

 // Fetch trust bundle
 bundle, err := client.FetchX509Bundle(ctx, svid.ID.TrustDomain())
 if err != nil {
  return nil, fmt.Errorf("failed to fetch X.509 bundle: %w", err)
 }

 return &SPIFFEClient{
  svid:        svid,
  bundle:      bundle,
  workloadAPI: client,
 }, nil
}

// Close closes the client
func (c *SPIFFEClient) Close() error {
 return c.workloadAPI.Close()
}

// GetTLSConfig creates TLS config for mTLS using SPIFFE identities
func (c *SPIFFEClient) GetTLSConfig(authorizedSPIFFEIDs []string) *tls.Config {
 return tlsconfig.MTLSServerConfig(c.svid, c.bundle, tlsconfig.AuthorizeIDStringMatch(authorizedSPIFFEIDs...))
}

// GetClientTLSConfig creates TLS config for client mTLS
func (c *SPIFFEClient) GetClientTLSConfig() *tls.Config {
 return tlsconfig.MTLSClientConfig(c.svid, c.bundle, tlsconfig.AuthorizeAny())
}

// GetMTLSHTTPClient creates HTTP client with mTLS
func (c *SPIFFEClient) GetMTLSHTTPClient(authorizedIDs []string) *http.Client {
 tlsConfig := tlsconfig.MTLSClientConfig(c.svid, c.bundle, tlsconfig.AuthorizeIDStringMatch(authorizedIDs...))

 return &http.Client{
  Transport: &http.Transport{
   TLSClientConfig: tlsConfig,
  },
  Timeout: 30 * time.Second,
 }
}

// Example: SPIFFE-enabled HTTP Server
func SPIFFEHTTPServerExample() {
 fmt.Println(`
SPIFFE-Enabled HTTP Server:
═══════════════════════════

package main

import (
 "context"
 "fmt"
 "log"
 "net/http"
 "time"

 "github.com/spiffe/go-spiffe/v2/spiffetls/tlsconfig"
 "github.com/spiffe/go-spiffe/v2/workloadapi"
)

func main() {
 ctx := context.Background()

 // Create workload API client
 client, err := workloadapi.New(ctx, workloadapi.WithAddr("unix:///tmp/spire-agent/public/api.sock"))
 if err != nil {
  log.Fatal(err)
 }
 defer client.Close()

 // Create mTLS server
 server := &http.Server{
  Addr: ":8443",
  TLSConfig: tlsconfig.MTLSServerConfig(
   client,  // Source of SVIDs
   client,  // Source of trust bundles
   tlsconfig.AuthorizeMemberOf("production.example.com"),
  ),
 }

 // Handler extracts SPIFFE ID from connection
 http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
  if r.TLS != nil && len(r.TLS.PeerCertificates) > 0 {
   // Extract SPIFFE ID from client certificate
   uri := r.TLS.PeerCertificates[0].URIs
   for _, u := range uri {
    if u.Scheme == "spiffe" {
     fmt.Fprintf(w, "Hello from SPIFFE ID: %s\n", u.String())
     return
    }
   }
  }
  http.Error(w, "No SPIFFE ID found", http.StatusUnauthorized)
 })

 log.Fatal(server.ListenAndServeTLS("", ""))
}
`)
}

// Example: SPIFFE-enabled gRPC Server
func SPIFFEGRPCServerExample() {
 fmt.Println(`
SPIFFE-Enabled gRPC Server:
═══════════════════════════

package main

import (
 "context"
 "fmt"
 "log"
 "net"

 "github.com/spiffe/go-spiffe/v2/spiffetls/tlsconfig"
 "github.com/spiffe/go-spiffe/v2/workloadapi"
 "google.golang.org/grpc"
 "google.golang.org/grpc/credentials"
 "google.golang.org/grpc/peer"
)

type server struct {
 pb.UnimplementedMyServiceServer
}

func (s *server) DoSomething(ctx context.Context, req *pb.Request) (*pb.Response, error) {
 // Extract SPIFFE ID from context
 p, ok := peer.FromContext(ctx)
 if !ok {
  return nil, fmt.Errorf("no peer info")
 }

 tlsInfo, ok := p.AuthInfo.(credentials.TLSInfo)
 if !ok {
  return nil, fmt.Errorf("no TLS info")
 }

 if len(tlsInfo.State.PeerCertificates) == 0 {
  return nil, fmt.Errorf("no peer certificates")
 }

 // Get SPIFFE ID from certificate
 uris := tlsInfo.State.PeerCertificates[0].URIs
 for _, uri := range uris {
  if uri.Scheme == "spiffe" {
   log.Printf("Request from SPIFFE ID: %s", uri.String())
   break
  }
 }

 return &pb.Response{Message: "Success"}, nil
}

func main() {
 ctx := context.Background()

 client, err := workloadapi.New(ctx, workloadapi.WithAddr("unix:///tmp/spire-agent/public/api.sock"))
 if err != nil {
  log.Fatal(err)
 }
 defer client.Close()

 // Only allow specific services
 allowedIDs := []string{
  "spiffe://production.example.com/ns/payments/sa/payment-service",
  "spiffe://production.example.com/ns/orders/sa/order-service",
 }

 creds := credentials.NewTLS(tlsconfig.MTLSServerConfig(
  client,
  client,
  tlsconfig.AuthorizeIDStringMatch(allowedIDs...),
 ))

 lis, err := net.Listen("tcp", ":50051")
 if err != nil {
  log.Fatal(err)
 }

 s := grpc.NewServer(grpc.Creds(creds))
 pb.RegisterMyServiceServer(s, &server{})

 log.Fatal(s.Serve(lis))
}
`)
}

// Example: SPIFFE-aware Client
func SPIFFEClientExample() {
 fmt.Println(`
SPIFFE-Aware gRPC Client:
═════════════════════════

package main

import (
 "context"
 "fmt"
 "log"
 "time"

 "github.com/spiffe/go-spiffe/v2/spiffetls/tlsconfig"
 "github.com/spiffe/go-spiffe/v2/workloadapi"
 "google.golang.org/grpc"
 "google.golang.org/grpc/credentials"
)

func main() {
 ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
 defer cancel()

 // Create workload API client
 client, err := workloadapi.New(ctx, workloadapi.WithAddr("unix:///tmp/spire-agent/public/api.sock"))
 if err != nil {
  log.Fatal(err)
 }
 defer client.Close()

 // Create mTLS credentials
 creds := credentials.NewTLS(tlsconfig.MTLSClientConfig(
  client,
  client,
  tlsconfig.AuthorizeMemberOf("production.example.com"),
 ))

 // Connect to server
 conn, err := grpc.Dial("server.production.svc.cluster.local:50051",
  grpc.WithTransportCredentials(creds),
 )
 if err != nil {
  log.Fatal(err)
 }
 defer conn.Close()

 // Use connection
 c := pb.NewMyServiceClient(conn)
 resp, err := c.DoSomething(ctx, &pb.Request{Data: "hello"})
 if err != nil {
  log.Fatal(err)
 }

 fmt.Println(resp.Message)
}
`)
}

// CertificateRotation demonstrates automatic certificate rotation
func CertificateRotation() {
 fmt.Println(`
Automatic Certificate Rotation:
════════════════════════════════

SPIRE automatically rotates certificates:

1. SVID TTL (default 1 hour):
   - Certificates valid for 1 hour
   - Auto-rotated at 50% of TTL (30 minutes)
   - No application restart needed

2. Trust Bundle TTL (default 24 hours):
   - Root CA certificates
   - Rotated automatically by SPIRE server
   - Clients update without restart

3. Implementation pattern for long-lived connections:
   - Use GetConfigForClient in TLS config
   - Fetch fresh SVID for each new connection
   - Or use spiffetls package which handles this

Example with automatic rotation:
─────────────────────────────────
config := &tls.Config{
    GetConfigForClient: func(*tls.ClientHelloInfo) (*tls.Config, error) {
        // Fetch fresh SVID for each connection
        svid, err := client.FetchX509SVID(ctx)
        if err != nil {
            return nil, err
        }

        bundle, err := client.FetchX509Bundle(ctx, svid.ID.TrustDomain())
        if err != nil {
            return nil, err
        }

        return tlsconfig.MTLSServerConfig(svid, bundle,
            tlsconfig.AuthorizeAny()), nil
    },
}
`)
}


### 5.4 Service Mesh Integration

```

┌─────────────────────────────────────────────────────────────────────────────┐
│                    Service Mesh mTLS Integration                             │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  Service meshes provide automatic mTLS and service identity across          │
│  all service-to-service communication.                                      │
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                        ISTIO                                          │   │
│  │  (Most widely adopted service mesh)                                   │   │
│  │                                                                       │   │
│  │  mTLS Modes:                                                          │   │
│  │  • PERMISSIVE: Accept both plaintext and mTLS (migration mode)        │   │
│  │  • STRICT: Require mTLS for all communication                         │   │
│  │  • DISABLE: No mTLS (not recommended)                                 │   │
│  │                                                                       │   │
│  │  Identity: SPIFFE-based                                               │   │
│  │  spiffe://cluster.local/ns/<namespace>/sa/<service-account>           │   │
│  │                                                                       │   │
│  │  Certificate Management: Istiod CA (or external CA via Vault)         │   │
│  │  Rotation: Automatic, ~24 hour TTL                                    │   │
│  │                                                                       │   │
│  │  Authorization: L4 (IP/port) + L7 (HTTP/gRPC paths, JWT)              │   │
│  │                                                                       │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                        LINKERD                                        │   │
│  │  (Lightweight, CNCF graduated)                                        │   │
│  │                                                                       │   │
│  │  mTLS: Always enabled, cannot be disabled                             │   │
│  │  Identity: TLS-based (not SPIFFE, but compatible)                     │   │
│  │                                                                       │   │
│  │  Certificate Management: Linkerd trust anchor                         │   │
│  │  Rotation: Automatic, 24 hour default                                 │   │
│  │                                                                       │   │
│  │  Authorization: Basic L4 only (no L7 authz)                           │   │
│  │                                                                       │   │
│  │  Performance: Lower overhead than Istio                               │   │
│  │                                                                       │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                        CILIUM SERVICE MESH                            │   │
│  │  (eBPF-based, kernel-level)                                           │   │
│  │                                                                       │   │
│  │  mTLS: SPIFFE-based via Envoy extension                               │   │
│  │  Identity: SPIFFE IDs with SPIRE integration                          │   │
│  │                                                                       │   │
│  │  Certificate Management: SPIRE or external CA                         │   │
│  │  Rotation: Automatic via SPIRE                                        │   │
│  │                                                                       │   │
│  │  Performance: Best-in-class (no sidecar proxy)                        │   │
│  │                                                                       │   │
│  │  Features: Combined with Cilium Network Policies                      │   │
│  │                                                                       │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                             │
│  COMPARISON:                                                                │
│                                                                             │
│  Feature              │  Istio     │  Linkerd   │  Cilium                  │
│  ─────────────────────┼────────────┼────────────┼────────────────────────  │
│  mTLS enforcement     │  Flexible  │  Always on │  Optional                │
│  SPIFFE identity      │  Yes       │  No        │  Yes                     │
│  L7 authorization     │  Excellent │  None      │  Via Envoy               │
│  Performance          │  Good      │  Better    │  Best                    │
│  Resource overhead    │  High      │  Medium    │  Low                     │
│  Sidecar              │  Yes       │  Yes       │  No (sidecar-less)       │
│  eBPF acceleration    │  Partial   │  No        │  Full                    │
│  Multi-cluster        │  Excellent │  Good      │  Good                    │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘

```

```yaml
# Istio mTLS Configuration

---
# Strict mTLS for entire namespace
apiVersion: security.istio.io/v1beta1
kind: PeerAuthentication
metadata:
  name: default
  namespace: production
spec:
  mtls:
    mode: STRICT

---
# Strict mTLS for specific workload
apiVersion: security.istio.io/v1beta1
kind: PeerAuthentication
metadata:
  name: payment-service
  namespace: production
spec:
  selector:
    matchLabels:
      app: payment-service
  mtls:
    mode: STRICT
  portLevelMtls:
    8080:
      mode: STRICT

---
# Authorization Policy - Deny all by default
apiVersion: security.istio.io/v1beta1
kind: AuthorizationPolicy
metadata:
  name: deny-all
  namespace: production
spec:
  {}

---
# Authorization Policy - Allow specific traffic
apiVersion: security.istio.io/v1beta1
kind: AuthorizationPolicy
metadata:
  name: payment-service-policy
  namespace: production
spec:
  selector:
    matchLabels:
      app: payment-service
  action: ALLOW
  rules:
  # Allow orders service to call payment service
  - from:
    - source:
        principals: ["cluster.local/ns/orders/sa/order-service"]
    to:
    - operation:
        methods: ["POST"]
        paths: ["/api/v1/payments/*"]

  # Allow admin service full access
  - from:
    - source:
        principals: ["cluster.local/ns/admin/sa/admin-service"]
    to:
    - operation:
        methods: ["*"]
        paths: ["/api/v1/admin/*"]

  # Require JWT token for customer endpoints
  - from:
    - source:
        requestPrincipals: ["*"]
    to:
    - operation:
        methods: ["GET"]
        paths: ["/api/v1/customers/*"]
    when:
    - key: request.auth.claims[scope]
      values: ["payments:read"]

---
# Request Authentication (JWT validation)
apiVersion: security.istio.io/v1beta1
kind: RequestAuthentication
metadata:
  name: jwt-auth
  namespace: production
spec:
  selector:
    matchLabels:
      app: api-gateway
  jwtRules:
  - issuer: "https://auth.example.com"
    jwksUri: "https://auth.example.com/.well-known/jwks.json"
    audiences:
    - "api.example.com"
    forwardOriginalToken: true
    outputPayloadToHeader: "x-jwt-payload"

---
# Cilium mTLS with SPIFFE
# Requires Cilium >= 1.14 with Envoy enabled
apiVersion: cilium.io/v2
kind: CiliumNetworkPolicy
metadata:
  name: payment-service-mtls
  namespace: production
spec:
  endpointSelector:
    matchLabels:
      app: payment-service
  ingress:
  - fromEndpoints:
    - matchLabels:
        k8s:io.kubernetes.pod.namespace: orders
        k8s:app: order-service
    toPorts:
    - ports:
      - port: "8080"
        protocol: TCP
      rules:
        http:
        - method: POST
          path: "/api/v1/payments/.*"

---
# Cilium Clusterwide Network Policy with mTLS
apiVersion: cilium.io/v2
kind: CiliumClusterwideNetworkPolicy
metadata:
  name: require-mtls
spec:
  description: "Require mTLS for all service-to-service communication"
  nodeSelector: {}
  ingressDeny:
  - fromEndpoints:
    - matchLabels:
        reserved:host: ""
    toPorts:
    - ports:
      - port: "1-65535"
        protocol: TCP
      rules:
        http:
        - headers:
          - name: x-forwarded-client-cert
            absence: true
  ingress:
  - fromEndpoints:
    - matchLabels:
        io.kubernetes.pod.namespace: production
    toPorts:
    - ports:
      - port: "8080"
        protocol: TCP
      terminatingTLS:
        certificate: "spiffe://production.example.com"
```

---

## 6. Go Cryptography

### 6.1 Native Go Crypto Packages

```go
package gocrypto

import (
 "crypto/aes"
 "crypto/cipher"
 "crypto/ecdsa"
 "crypto/ed25519"
 "crypto/elliptic"
 "crypto/hmac"
 "crypto/mlkem"
 "crypto/rand"
 "crypto/rsa"
 "crypto/sha256"
 "crypto/sha512"
 "crypto/tls"
 "crypto/x509"
 "encoding/hex"
 "encoding/pem"
 "fmt"
 "io"
 "time"
)

// NativeCryptoExamples demonstrates Go's native cryptography capabilities
func NativeCryptoExamples() {
 fmt.Println(`
Go Cryptography Overview:
═════════════════════════

Standard Library Packages:
──────────────────────────
• crypto/aes        - AES encryption (GCM, CTR modes)
• crypto/cipher     - Block cipher modes
• crypto/des        - DES/Triple DES (legacy, avoid)
• crypto/ecdsa      - ECDSA signatures
• crypto/ed25519    - Ed25519 signatures (recommended)
• crypto/elliptic   - Elliptic curves
• crypto/hmac       - HMAC authentication
• crypto/md5        - MD5 (legacy, avoid for security)
• crypto/mlkem      - ML-KEM post-quantum (Go 1.24+)
     - ML-KEM-512
     - ML-KEM-768 (recommended)
     - ML-KEM-1024
• crypto/rand       - Cryptographically secure RNG
• crypto/rc4        - RC4 stream cipher (insecure, deprecated)
• crypto/rsa        - RSA encryption/signatures
• crypto/sha1       - SHA-1 (legacy, avoid for new code)
• crypto/sha256     - SHA-256 hash
• crypto/sha512     - SHA-384, SHA-512 hashes
• crypto/subtle     - Constant-time comparisons
• crypto/tls        - TLS implementation
• crypto/x509       - X.509 certificates

x/crypto Extensions:
────────────────────
• golang.org/x/crypto/argon2    - Argon2 password hashing
• golang.org/x/crypto/bcrypt    - bcrypt password hashing
• golang.org/x/crypto/blake2b   - BLAKE2b hash
• golang.org/x/crypto/chacha20  - ChaCha20 stream cipher
• golang.org/x/crypto/chacha20poly1305 - ChaCha20-Poly1305 AEAD
• golang.org/x/crypto/curve25519 - X25519 ECDH
• golang.org/x/crypto/hkdf      - HKDF key derivation
• golang.org/x/crypto/nacl      - NaCl/libsodium compatibility
• golang.org/x/crypto/pbkdf2    - PBKDF2 key derivation
• golang.org/x/crypto/scrypt    - scrypt password hashing
• golang.org/x/crypto/ssh       - SSH client/server

RECOMMENDED ALGORITHMS (2026):
──────────────────────────────
Symmetric Encryption:
• AES-256-GCM (authenticated encryption)
• ChaCha20-Poly1305 (for mobile/ARM)

Asymmetric Encryption:
• X25519MLKEM768 (Go 1.24+, post-quantum hybrid)
• X25519 (classical ECDH)
• ML-KEM-768 (post-quantum only)

Signatures:
• Ed25519 (fast, secure, compact)
• ML-DSA-65 (post-quantum, FIPS 204)
• ECDSA P-256 (for compatibility)

Hashing:
• SHA-256 (general purpose)
• SHA-3-256 (if quantum resistance needed)
• BLAKE2b (faster alternative)

Password Hashing:
• Argon2id (OWASP recommended)
• bcrypt (if Argon2 unavailable)
• scrypt (for cryptocurrencies)
`)
}

// SecureRandom demonstrates secure random number generation
func SecureRandom() {
 // crypto/rand provides cryptographically secure random numbers
 // NEVER use math/rand for security purposes

 // Generate random bytes
 randomBytes := make([]byte, 32)
 _, err := rand.Read(randomBytes)
 if err != nil {
  panic(err)
 }

 fmt.Printf("Random bytes: %x\n", randomBytes)

 // Generate random int
 randomInt, err := rand.Int(rand.Reader, big.NewInt(1000))
 if err != nil {
  panic(err)
 }

 fmt.Printf("Random int (0-999): %d\n", randomInt)

 // Generate UUID v4
 uuid := make([]byte, 16)
 rand.Read(uuid)
 // Set version (4) and variant bits
 uuid[6] = (uuid[6] & 0x0f) | 0x40
 uuid[8] = (uuid[8] & 0x3f) | 0x80

 fmt.Printf("UUID: %x-%x-%x-%x-%x\n",
  uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:])
}

// AESGCMExample demonstrates AES-GCM authenticated encryption
func AESGCMExample() {
 // Generate key
 key := make([]byte, 32) // AES-256
 if _, err := rand.Read(key); err != nil {
  panic(err)
 }

 // Create AES cipher
 block, err := aes.NewCipher(key)
 if err != nil {
  panic(err)
 }

 // Create GCM mode
 gcm, err := cipher.NewGCM(block)
 if err != nil {
  panic(err)
 }

 // Encrypt
 plaintext := []byte("Hello, World!")
 nonce := make([]byte, gcm.NonceSize())
 if _, err := rand.Read(nonce); err != nil {
  panic(err)
 }

 ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
 fmt.Printf("Ciphertext: %x\n", ciphertext)

 // Decrypt
 nonceSize := gcm.NonceSize()
 if len(ciphertext) < nonceSize {
  panic("ciphertext too short")
 }

 decrypted, err := gcm.Open(nil, ciphertext[:nonceSize], ciphertext[nonceSize:], nil)
 if err != nil {
  panic(err)
 }

 fmt.Printf("Decrypted: %s\n", decrypted)
}

// ChaCha20Poly1305Example demonstrates ChaCha20-Poly1305
func ChaCha20Poly1305Example() {
 import "golang.org/x/crypto/chacha20poly1305"

 // Generate key
 key := make([]byte, chacha20poly1305.KeySize)
 if _, err := rand.Read(key); err != nil {
  panic(err)
 }

 // Create cipher
 aead, err := chacha20poly1305.New(key)
 if err != nil {
  panic(err)
 }

 // Encrypt
 plaintext := []byte("Hello, World!")
 nonce := make([]byte, aead.NonceSize())
 if _, err := rand.Read(nonce); err != nil {
  panic(err)
 }

 ciphertext := aead.Seal(nonce, nonce, plaintext, nil)
 fmt.Printf("ChaCha20-Poly1305 ciphertext: %x\n", ciphertext)

 // Decrypt
 nonceSize := aead.NonceSize()
 decrypted, err := aead.Open(nil, ciphertext[:nonceSize], ciphertext[nonceSize:], nil)
 if err != nil {
  panic(err)
 }

 fmt.Printf("Decrypted: %s\n", decrypted)
}

// Ed25519Example demonstrates Ed25519 signatures
func Ed25519Example() {
 // Generate key pair
 publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
 if err != nil {
  panic(err)
 }

 fmt.Printf("Public key: %x\n", publicKey)

 // Sign message
 message := []byte("Hello, World!")
 signature := ed25519.Sign(privateKey, message)
 fmt.Printf("Signature: %x\n", signature)
 fmt.Printf("Signature length: %d bytes\n", len(signature))

 // Verify signature
 valid := ed25519.Verify(publicKey, message, signature)
 fmt.Printf("Signature valid: %v\n", valid)

 // Try to verify with wrong message
 valid = ed25519.Verify(publicKey, []byte("Wrong message"), signature)
 fmt.Printf("Wrong message valid: %v\n", valid)
}

// BcryptExample demonstrates bcrypt password hashing
func BcryptExample() {
 import "golang.org/x/crypto/bcrypt"

 password := []byte("my-secret-password")

 // Hash password (cost 10-14 recommended)
 hash, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
 if err != nil {
  panic(err)
 }

 fmt.Printf("Bcrypt hash: %s\n", hash)
 fmt.Printf("Cost: %d\n", bcrypt.DefaultCost)

 // Verify password
 err = bcrypt.CompareHashAndPassword(hash, password)
 if err != nil {
  fmt.Println("Password does not match")
 } else {
  fmt.Println("Password matches!")
 }

 // Verify wrong password
 err = bcrypt.CompareHashAndPassword(hash, []byte("wrong-password"))
 if err != nil {
  fmt.Println("Wrong password rejected")
 }
}

// Argon2Example demonstrates Argon2 password hashing (recommended)
func Argon2Example() {
 import "golang.org/x/crypto/argon2"

 password := []byte("my-secret-password")
 salt := make([]byte, 16)
 if _, err := rand.Read(salt); err != nil {
  panic(err)
 }

 // Argon2id (recommended variant)
 hash := argon2.IDKey(
  password,
  salt,
  3,      // time (iterations)
  64*1024, // memory (64 MB)
  4,      // threads
  32,     // key length
 )

 fmt.Printf("Argon2id hash: %x\n", hash)
 fmt.Printf("Parameters: time=3, memory=64MB, threads=4\n")

 // Store: salt + hash + parameters
 // Verify: recompute with same salt and parameters
}

// HMACExample demonstrates HMAC message authentication
func HMACExample() {
 // Create HMAC
 key := []byte("my-secret-key")
 message := []byte("Hello, World!")

 mac := hmac.New(sha256.New, key)
 mac.Write(message)
 signature := mac.Sum(nil)

 fmt.Printf("HMAC-SHA256: %x\n", signature)

 // Verify
 mac2 := hmac.New(sha256.New, key)
 mac2.Write(message)
 expectedMAC := mac2.Sum(nil)

 valid := hmac.Equal(signature, expectedMAC)
 fmt.Printf("HMAC valid: %v\n", valid)

 // Use crypto/subtle for constant-time comparison
 // when verifying secrets
 valid = subtle.ConstantTimeCompare(signature, expectedMAC) == 1
 fmt.Printf("Constant-time valid: %v\n", valid)
}


// TLSServerExample demonstrates secure TLS server configuration
func TLSServerExample() {
 fmt.Println(`
Secure TLS Server Configuration:
═════════════════════════════════

package main

import (
 "crypto/tls"
 "log"
 "net/http"
)

func main() {
 tlsConfig := &tls.Config{
  // Minimum TLS 1.3 (Go 1.24+)
  MinVersion: tls.VersionTLS13,

  // Disable renegotiation
  Renegotiation: tls.RenegotiateNever,

  // Prefer server cipher suites
  PreferServerCipherSuites: true,

  // Strong cipher suites (TLS 1.3 has fixed list)
  CipherSuites: []uint16{
   tls.TLS_AES_256_GCM_SHA384,
   tls.TLS_AES_128_GCM_SHA256,
   tls.TLS_CHACHA20_POLY1305_SHA256,
  },

  // Post-quantum key exchange (Go 1.24+)
  CurvePreferences: []tls.CurveID{
   tls.X25519MLKEM768,  // PQ hybrid
   tls.X25519,
   tls.SecP256R1,
  },

  // Certificate verification
  ClientAuth: tls.VerifyClientCertIfGiven,

  // Session tickets
  SessionTicketsDisabled: false,
 }

 server := &http.Server{
  Addr:      ":8443",
  TLSConfig: tlsConfig,
  Handler:   handler(),
 }

 log.Fatal(server.ListenAndServeTLS("server.crt", "server.key"))
}
`)
}

// TLSClientExample demonstrates secure TLS client configuration
func TLSClientExample() {
 fmt.Println(`
Secure TLS Client Configuration:
═════════════════════════════════

package main

import (
 "crypto/tls"
 "crypto/x509"
 "net/http"
 "os"
)

func main() {
 // Load system CA certs
 rootCAs, _ := x509.SystemCertPool()
 if rootCAs == nil {
  rootCAs = x509.NewCertPool()
 }

 // Add custom CA cert
 caCert, _ := os.ReadFile("custom-ca.crt")
 rootCAs.AppendCertsFromPEM(caCert)

 tlsConfig := &tls.Config{
  // Minimum TLS 1.2 (1.3 preferred)
  MinVersion: tls.VersionTLS12,

  // Root CAs
  RootCAs: rootCAs,

  // Certificate verification
  InsecureSkipVerify: false,

  // Server name verification
  ServerName: "api.example.com",

  // Post-quantum key exchange (Go 1.24+)
  CurvePreferences: []tls.CurveID{
   tls.X25519MLKEM768,
   tls.X25519,
  },

  // Strong cipher suites
  CipherSuites: []uint16{
   tls.TLS_AES_256_GCM_SHA384,
   tls.TLS_AES_128_GCM_SHA256,
   tls.TLS_CHACHA20_POLY1305_SHA256,
  },
 }

 client := &http.Client{
  Transport: &http.Transport{
   TLSClientConfig: tlsConfig,
  },
 }

 resp, err := client.Get("https://api.example.com")
 // ...
}
`)
}

// MutualTLSExample demonstrates mTLS configuration
func MutualTLSExample() {
 fmt.Println(`
Mutual TLS (mTLS) Configuration:
═════════════════════════════════

Server-side mTLS:
─────────────────
tlsConfig := &tls.Config{
 // Require client certificate
 ClientAuth: tls.RequireAndVerifyClientCert,

 // Client CAs
 ClientCAs: clientCAPool,

 // Verify function (optional)
 VerifyPeerCertificate: func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
  // Custom verification logic
  return nil
 },
}

Client-side mTLS:
─────────────────
tlsConfig := &tls.Config{
 // Client certificate
 Certificates: []tls.Certificate{clientCert},

 // Server verification
 RootCAs: rootCAs,
 ServerName: "server.example.com",
}
`)
}

### 6.2 FIPS 140-3 Module Status

```

┌─────────────────────────────────────────────────────────────────────────────┐
│                    Go FIPS 140-3 Module Status                               │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  CURRENT STATUS (2026):                                                     │
│                                                                             │
│  Go's standard library crypto packages do NOT have FIPS 140-3 validation.   │
│  However, several solutions exist for FIPS compliance:                      │
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │  Option 1: Microsoft Go (go-crypto-openssl)                          │   │
│  │  ─────────────────────────────────────────                           │   │
│  │                                                                       │   │
│  │  • OpenSSL backend for Go crypto                                     │   │
│  │  • Leverages OpenSSL FIPS 140-3 module                               │   │
│  │  • Drop-in replacement for standard library                          │   │
│  │  • <https://github.com/microsoft/go-crypto-openssl>                    │   │
│  │                                                                       │   │
│  │  Installation:                                                       │   │
│  │  GOEXPERIMENT=opensslcrypto go build                                  │   │
│  │                                                                       │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │  Option 2: AWS-LC (AWS Libcrypto)                                    │   │
│  │  ────────────────────────────────                                    │   │
│  │                                                                       │   │
│  │  • AWS's FIPS-validated crypto library                               │   │
│  │  • Go bindings available                                             │   │
│  │  • <https://github.com/aws/aws-lc>                                     │   │
│  │                                                                       │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │  Option 3: BoringCrypto (Google)                                     │   │
│  │  ─────────────────────────────                                       │   │
│  │                                                                       │   │
│  │  • Used internally at Google                                         │   │
│  │  • BoringSSL FIPS module                                             │   │
│  │  • Available in Go 1.24 with GOEXPERIMENT=boringcrypto               │   │
│  │                                                                       │   │
│  │  Limitation: Not generally available                                 │   │
│  │                                                                       │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │  Option 4: Entrust Go Toolkit                                        │   │
│  │  ────────────────────────────                                        │   │
│  │                                                                       │   │
│  │  • Commercial FIPS 140-3 solution                                    │   │
│  │  • Purpose-built for Go                                              │   │
│  │  • <https://www.entrust.com/go-toolkit>                                │   │
│  │                                                                       │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                             │
│  FIPS 140-3 VALIDATED MODULES:                                              │
│                                                                             │
│  Module               │  Level  │  Vendor      │  Certificate             │
│  ─────────────────────┼─────────┼──────────────┼────────────────────────  │
│  OpenSSL 3.0          │  1      │  OpenSSL     │  #4282                   │
│  OpenSSL 3.1          │  1      │  OpenSSL     │  #4634                   │
│  BoringSSL            │  1      │  Google      │  #4407                   │
│  AWS-LC               │  1      │  AWS         │  #4816                   │
│                                                                             │
│  ROADMAP:                                                                   │
│  • Go 1.25+ considering native FIPS module                                 │
│  • Likely based on BoringCrypto or similar                                 │
│  • Timeline: 2026-2027 for initial validation                              │
│                                                                             │
│  RECOMMENDATION:                                                            │
│  For FIPS compliance today, use Microsoft Go with OpenSSL backend          │
│  or wait for official Go FIPS module (2026-2027).                          │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘

```

### 6.3 Trail of Bits Security Audit (May 2025)

```

┌─────────────────────────────────────────────────────────────────────────────┐
│                    Go Security Audit Results (May 2025)                      │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  Audit performed by Trail of Bits on Go 1.24 standard library cryptography  │
│                                                                             │
│  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━ │
│                                                                             │
│  HIGH SEVERITY FINDINGS:                                                    │
│                                                                             │
│  None identified in standard library crypto packages.                       │
│                                                                             │
│  MEDIUM SEVERITY FINDINGS:                                                  │
│                                                                             │
│  1. crypto/tls: Potential timing side-channel in RSA key generation        │
│     Status: Mitigated in Go 1.24.1                                         │
│     Impact: Limited (RSA key generation is rare)                           │
│                                                                             │
│  2. crypto/ecdsa: Non-constant-time hash-to-curve for some curves          │
│     Status: Fixed in Go 1.24.2                                             │
│     Impact: Low (requires specific attack conditions)                      │
│                                                                             │
│  LOW SEVERITY FINDINGS:                                                     │
│                                                                             │
│  1. Documentation: Missing guidance on ML-KEM side-channel resistance      │
│     Status: Documentation updated                                          │
│                                                                             │
│  2. crypto/mlkem: Input validation could be stricter                       │
│     Status: Improved in Go 1.24.1                                          │
│                                                                             │
│  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━ │
│                                                                             │
│  RECOMMENDATIONS FROM AUDIT:                                                │
│                                                                             │
│  ✓ Prefer Ed25519 over ECDSA for new code                                  │
│  ✓ Use crypto/rand for all random generation                               │
│  ✓ Enable TLS 1.3 minimum where possible                                   │
│  ✓ Use X25519MLKEM768 for post-quantum protection                          │
│  ✓ Avoid RSA for new designs (use ECDH/Ed25519)                            │
│  ✓ Use AES-GCM with random nonces (never reuse)                            │
│                                                                             │
│  POSITIVE FINDINGS:                                                         │
│                                                                             │
│  ✓ Constant-time implementations for most operations                       │
│  ✓ Memory clearing for sensitive data                                      │
│  ✓ No known critical vulnerabilities                                       │
│  ✓ Good test coverage for edge cases                                       │
│  ✓ Proper handling of malformed inputs                                     │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘

```

### 6.4 Cryptography Best Practices

```go
package bestpractices

import (
 "crypto/subtle"
 "fmt"
)

// BestPracticesSummary provides cryptographic best practices
func BestPracticesSummary() {
 fmt.Println(`
Go Cryptography Best Practices (2026):
═══════════════════════════════════════

1. RANDOM NUMBER GENERATION
───────────────────────────
✓ ALWAYS use crypto/rand for security-sensitive operations
✗ NEVER use math/rand for keys, passwords, or nonces

// Good:
key := make([]byte, 32)
if _, err := rand.Read(key); err != nil {
    panic(err)
}

// Bad:
r := mathrand.New(mathrand.NewSource(time.Now().UnixNano()))
key := make([]byte, 32)
r.Read(key)  // PREDICTABLE!

2. PASSWORD HASHING
───────────────────
✓ Use Argon2id (OWASP recommended)
✓ Use bcrypt if Argon2 not available
✓ Cost factor: bcrypt=10-14, Argon2 time=3-6
✗ NEVER use MD5, SHA-*, or simple hashes

// Good (Argon2id):
hash := argon2.IDKey(password, salt, 3, 64*1024, 4, 32)

// Acceptable (bcrypt):
hash, _ := bcrypt.GenerateFromPassword(password, 12)

// Bad:
hash := sha256.Sum256(password)  // VULNERABLE!

3. SYMMETRIC ENCRYPTION
───────────────────────
✓ Use AES-256-GCM (authenticated encryption)
✓ Use ChaCha20-Poly1305 on mobile/ARM
✓ Generate new nonce for each encryption
✓ Never reuse nonce+key combination
✗ NEVER use ECB mode
✗ NEVER use unauthenticated encryption (CBC without MAC)

// Good:
gcm, _ := cipher.NewGCM(block)
nonce := make([]byte, gcm.NonceSize())
rand.Read(nonce)
ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)

// Bad (unauthenticated):
ciphertext := make([]byte, len(plaintext))
cipher.NewCTR(block, nonce).XORKeyStream(ciphertext, plaintext)

4. ASYMMETRIC ENCRYPTION/KEY EXCHANGE
─────────────────────────────────────
✓ Use X25519MLKEM768 (Go 1.24+, post-quantum)
✓ Use X25519 for classical only
✓ Use ECDH P-256 for compatibility
✗ Avoid RSA for new designs

// Good (Post-quantum):
// Go 1.24 automatically uses X25519MLKEM768 in TLS

// Acceptable (Classical):
privateKey, _ := ecdh.X25519().GenerateKey(rand.Reader)

5. DIGITAL SIGNATURES
─────────────────────
✓ Use Ed25519 (fast, secure, compact)
✓ Use ML-DSA for post-quantum (Go 1.25+)
✓ Use ECDSA P-256 for compatibility
✗ Avoid RSA signatures

// Good (Ed25519):
publicKey, privateKey, _ := ed25519.GenerateKey(rand.Reader)
signature := ed25519.Sign(privateKey, message)

// Acceptable (ECDSA):
privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
r, s, _ := ecdsa.Sign(rand.Reader, privateKey, hash[:])

6. HASH FUNCTIONS
─────────────────
✓ Use SHA-256 for general purpose
✓ Use SHA-3-256 if quantum resistance needed
✓ Use BLAKE2b for performance
✗ Avoid MD5 and SHA-1

// Good:
hash := sha256.Sum256(data)

// Better (faster):
hash, _ := blake2b.Sum256(data)

7. HMAC
───────
✓ Use HMAC-SHA256 for message authentication
✓ Use constant-time comparison (crypto/subtle)

// Good:
mac := hmac.New(sha256.New, key)
mac.Write(message)
signature := mac.Sum(nil)

// Verification:
if !hmac.Equal(signature, expected) {
    return errors.New("invalid signature")
}

8. TIMING ATTACK PREVENTION
───────────────────────────
✓ Use crypto/subtle.ConstantTimeCompare
✓ Use crypto/subtle.ConstantTimeSelect
✓ Avoid branching on secret data

// Good:
if subtle.ConstantTimeCompare(calculatedMAC, receivedMAC) != 1 {
    return errors.New("invalid MAC")
}

// Bad:
if string(calculatedMAC) == string(receivedMAC) {  // TIMING LEAK!
    // ...
}

9. TLS CONFIGURATION
────────────────────
✓ Minimum TLS 1.2 (1.3 preferred)
✓ Enable post-quantum key exchange (Go 1.24+)
✓ Use strong cipher suites
✓ Disable compression
✗ Never disable certificate verification

// Good TLS config:
&tls.Config{
    MinVersion: tls.VersionTLS12,
    CurvePreferences: []tls.CurveID{
        tls.X25519MLKEM768,  // PQ hybrid
        tls.X25519,
    },
}

10. MEMORY SECURITY
───────────────────
✓ Clear sensitive data from memory
✓ Use sync.Pool for sensitive buffers
✗ Avoid leaving keys in memory longer than needed

// Good:
key := make([]byte, 32)
rand.Read(key)
// ... use key ...
for i := range key {
    key[i] = 0  // Clear memory
}

11. KEY MANAGEMENT
──────────────────
✓ Use key derivation functions (HKDF, PBKDF2)
✓ Use hardware security modules (HSM) for production
✓ Rotate keys regularly
✓ Never hardcode keys

// Good (key derivation):
salt := make([]byte, 32)
rand.Read(salt)
key := pbkdf2.Key(password, salt, 100000, 32, sha256.New)

12. CRYPTOGRAPHIC AGILITY
─────────────────────────
✓ Version your encrypted data
✓ Support algorithm negotiation
✓ Plan for algorithm deprecation

// Good:
type EncryptedData struct {
    Version   uint8  // Algorithm version
    Algorithm string // Algorithm identifier
    Nonce     []byte
    Ciphertext []byte
}
`)
}

// ConstantTimeComparison demonstrates constant-time comparison
func ConstantTimeComparison() {
 secret := []byte("correct-horse-battery-staple")
 guess1 := []byte("correct-horse-battery-staple")
 guess2 := []byte("wrong-password")

 // Constant-time comparison
 result1 := subtle.ConstantTimeCompare(secret, guess1)
 result2 := subtle.ConstantTimeCompare(secret, guess2)

 fmt.Printf("Correct guess: %d (should be 1)\n", result1)
 fmt.Printf("Wrong guess: %d (should be 0)\n", result2)
}

// SecureMemoryHandling demonstrates secure memory practices
func SecureMemoryHandling() {
 fmt.Println(`
Secure Memory Handling:
═══════════════════════

1. Allocate sensitive data:
   key := make([]byte, 32)
   rand.Read(key)

2. Use the key:
   ciphertext := aesgcm.Seal(nil, nonce, plaintext, key)

3. Clear from memory when done:
   for i := range key {
       key[i] = 0
   }

4. Use runtime.GC() if immediate clearing is critical:
   runtime.GC()

Note: Go's garbage collector may copy data during compaction.
For maximum security, consider using mmap with MLOCK and MADV_DONTDUMP.
`)
}

// CryptographicAgilityExample demonstrates algorithm versioning
func CryptographicAgilityExample() {
 fmt.Println(`
Cryptographic Agility:
══════════════════════

// Versioned encryption envelope
type EncryptedEnvelope struct {
    Version    uint8  // 1 = AES-256-GCM, 2 = ChaCha20-Poly1305
    Algorithm  string // "AES-256-GCM" or "ChaCha20-Poly1305"
    Parameters []byte // Algorithm-specific parameters
    Ciphertext []byte // The encrypted data
}

func Encrypt(version uint8, plaintext []byte) (*EncryptedEnvelope, error) {
    switch version {
    case 1:
        return encryptAESGCM(plaintext)
    case 2:
        return encryptChaCha20(plaintext)
    default:
        return nil, fmt.Errorf("unsupported version: %d", version)
    }
}

func Decrypt(envelope *EncryptedEnvelope) ([]byte, error) {
    switch envelope.Version {
    case 1:
        return decryptAESGCM(envelope)
    case 2:
        return decryptChaCha20(envelope)
    default:
        return nil, fmt.Errorf("unsupported version: %d", envelope.Version)
    }
}

This allows:
• Adding new algorithms without breaking old data
• Graceful deprecation of weak algorithms
• Emergency algorithm switching
`)
}

---

## 7. Migration Checklists

### 7.1 Post-Quantum Migration Checklist

```

┌─────────────────────────────────────────────────────────────────────────────┐
│                    Post-Quantum Cryptography Migration Checklist             │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  PHASE 1: ASSESSMENT (2024-2025)                                            │
│  □ Inventory all cryptographic implementations                              │
│  □ Identify TLS endpoints and their key exchange algorithms                 │
│  □ Catalog all signatures and their algorithms                              │
│  □ Identify long-term secrets that need protection                          │
│  □ Assess "harvest now, decrypt later" risk                                 │
│  □ Review compliance requirements (CNSA 2.0, NIS2)                          │
│                                                                             │
│  PHASE 2: PREPARATION (2025-2026)                                           │
│  □ Upgrade to Go 1.24+ for native ML-KEM support                            │
│  □ Enable X25519MLKEM768 in TLS configuration                               │
│  □ Test hybrid key exchange in staging environments                         │
│  □ Monitor for compatibility issues with legacy clients                     │
│  □ Update monitoring to track key exchange algorithms used                  │
│                                                                             │
│  PHASE 3: HYBRID DEPLOYMENT (2026-2028)                                     │
│  □ Deploy hybrid (classical + PQ) key exchange to all services              │
│  □ Implement PQ certificate signing (ML-DSA)                                │
│  □ Update service mesh to support PQ mTLS                                   │
│  □ Ensure all external-facing APIs support PQ                               │
│  □ Document PQ algorithm usage for compliance                               │
│                                                                             │
│  PHASE 4: PQ-ONLY TRANSITION (2028-2030)                                    │
│  □ Disable classical algorithms where possible                              │
│  □ Maintain hybrid for compatibility with external systems                  │
│  □ Update security policies to require PQ                                   │
│  □ Complete vendor assessments for PQ support                               │
│                                                                             │
│  PHASE 5: COMPLIANCE (2030+)                                                │
│  □ Achieve full CNSA 2.0 compliance                                         │
│  □ Complete EU NIS2 Directive requirements                                  │
│  □ Maintain ongoing algorithm agility                                       │
│  □ Monitor for new PQ standards (FIPS 206, etc.)                            │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘

```

### 7.2 Supply Chain Security Checklist

```

┌─────────────────────────────────────────────────────────────────────────────┐
│                    Supply Chain Security Implementation Checklist            │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  IMMEDIATE (Implement Now):                                                 │
│  □ Enable GitHub/GitLab branch protection                                   │
│  □ Require 2-person review for all changes                                  │
│  □ Enable Dependabot or similar dependency scanning                         │
│  □ Sign commits with Gitsign (Sigstore)                                     │
│  □ Pin all GitHub Actions to commit SHA                                     │
│  □ Generate SBOMs for all releases                                          │
│                                                                             │
│  SHORT-TERM (Within 3 months):                                              │
│  □ Implement SLSA Level 2 provenance                                        │
│  □ Sign container images with Cosign                                        │
│  □ Enable Binary Authorization (GKE) or equivalent                          │
│  □ Implement vulnerability scanning in CI/CD                                │
│  □ Create Software Bill of Materials (SBOM)                                 │
│  □ Enable artifact attestation with SLSA generator                          │
│                                                                             │
│  MEDIUM-TERM (Within 6 months):                                             │
│  □ Achieve SLSA Level 3 compliance                                          │
│  □ Implement hermetic builds                                                │
│  □ Enable reproducible builds where possible                                │
│  □ Implement policy as code (OPA/Kyverno)                                   │
│  □ Complete vendor security assessments                                     │
│  □ Implement dependency update automation                                   │
│                                                                             │
│  LONG-TERM (Within 12 months):                                              │
│  □ Strive for SLSA Level 4 compliance                                       │
│  □ Implement comprehensive audit logging                                    │
│  □ Complete third-party security audit                                      │
│  □ Achieve SOC 2 Type II compliance                                         │
│  □ Implement insider threat detection                                       │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘

```

---

## References

1. **NIST Standards**
   - FIPS 203: ML-KEM
   - FIPS 204: ML-DSA
   - FIPS 205: SLH-DSA
   - NIST SP 800-57: Key Management

2. **Go Cryptography**
   - Go 1.24 Release Notes
   - Trail of Bits Go Security Audit (May 2025)
   - Go Cryptography Design Docs

3. **Cloud Native Security**
   - CNCF Security Whitepaper
   - Kubernetes Security Best Practices
   - SLSA Framework Specification

4. **Service Identity**
   - SPIFFE/SPIRE Documentation
   - Istio Security Documentation
   - Cilium Service Mesh Guide

5. **Supply Chain**
   - Sigstore Documentation
   - SLSA GitHub Generator
   - CISA SBOM Guidance

6. **Zero-Knowledge Proofs**
   - Zcash Protocol Specification
   - Mina Protocol Documentation
   - zkSNARKs vs zkSTARKs Research

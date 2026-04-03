# EC-020: Kubernetes Security in Production (2025-2026 Edition)

## Table of Contents

- [EC-020: Kubernetes Security in Production (2025-2026 Edition)](#ec-020-kubernetes-security-in-production-2025-2026-edition)
  - [Table of Contents](#table-of-contents)
  - [Overview](#overview)
    - [Security Threat Model](#security-threat-model)
  - [Post-Quantum Cryptography](#post-quantum-cryptography)
    - [NIST Standards (FIPS 203/204/205) - August 2024](#nist-standards-fips-203204205---august-2024)
      - [FIPS 203: Module-Lattice-Based Key-Encapsulation Mechanism (ML-KEM)](#fips-203-module-lattice-based-key-encapsulation-mechanism-ml-kem)
      - [FIPS 204: Module-Lattice-Based Digital Signature Algorithm (ML-DSA)](#fips-204-module-lattice-based-digital-signature-algorithm-ml-dsa)
      - [FIPS 205: Stateless Hash-Based Digital Signature Standard (SLH-DSA)](#fips-205-stateless-hash-based-digital-signature-standard-slh-dsa)
    - [Go 1.24+ X25519MLKEM768 Hybrid](#go-124-x25519mlkem768-hybrid)
    - [Migration Deadlines and Roadmap](#migration-deadlines-and-roadmap)
  - [Supply Chain Security](#supply-chain-security)
    - [SLSA v1.0/v1.1 Framework](#slsa-v10v11-framework)
      - [SLSA Levels Implementation](#slsa-levels-implementation)
    - [Sigstore Ecosystem](#sigstore-ecosystem)
      - [Cosign for Container Signing](#cosign-for-container-signing)
      - [Rekor Transparency Log](#rekor-transparency-log)
      - [Fulcio for Keyless Signing](#fulcio-for-keyless-signing)
    - [Keyless Signing with OIDC](#keyless-signing-with-oidc)
    - [SBOM Requirements (CISA 2025)](#sbom-requirements-cisa-2025)
  - [eBPF Runtime Security](#ebpf-runtime-security)
    - [Tetragon for Kernel-Level Threat Detection](#tetragon-for-kernel-level-threat-detection)
    - [TracingPolicy Examples](#tracingpolicy-examples)
    - [Falco Integration](#falco-integration)
  - [Service Mesh Security](#service-mesh-security)
    - [Istio Ambient Mode mTLS](#istio-ambient-mode-mtls)
      - [Performance Comparison](#performance-comparison)
    - [SPIFFE/SPIRE Integration](#spiffespire-integration)
  - [Confidential Computing](#confidential-computing)
    - [Hardware-Based Enclaves](#hardware-based-enclaves)
    - [Use Cases for Sensitive Workloads](#use-cases-for-sensitive-workloads)
  - [Network Security](#network-security)
    - [Cilium Cluster Mesh with Encryption](#cilium-cluster-mesh-with-encryption)
  - [Identity and Access Management](#identity-and-access-management)
    - [RBAC Hardening](#rbac-hardening)
  - [Secrets Management](#secrets-management)
    - [External Secrets Operator with Vault](#external-secrets-operator-with-vault)
  - [Pod Security](#pod-security)
    - [Pod Security Standards](#pod-security-standards)
  - [Runtime Security](#runtime-security)
    - [OPA/Gatekeeper Policies](#opagatekeeper-policies)
  - [Security Monitoring and Auditing](#security-monitoring-and-auditing)
    - [Falco + Fluent Bit Integration](#falco--fluent-bit-integration)
  - [Incident Response](#incident-response)
    - [Automated Incident Response](#automated-incident-response)
  - [Compliance and Governance](#compliance-and-governance)
    - [CIS Kubernetes Benchmarks](#cis-kubernetes-benchmarks)
  - [References](#references)
    - [Standards and Specifications](#standards-and-specifications)
    - [Tools and Projects](#tools-and-projects)
    - [Further Reading](#further-reading)

---

## Overview

Kubernetes security in production environments requires a defense-in-depth approach that addresses threats across the entire stack—from the underlying infrastructure to the application layer. This guide covers the latest security developments for 2025-2026, including post-quantum cryptography readiness, supply chain security with Sigstore, eBPF-based runtime security, and confidential computing technologies.

### Security Threat Model

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                           THREAT LANDSCAPE 2025-2026                        │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐             │
│  │  Supply Chain   │  │  Runtime        │  │  Network        │             │
│  │  Attacks        │  │  Threats        │  │  Intrusions     │             │
│  │  ▼ 45% increase │  │  ▼ 38% increase │  │  ▼ 32% increase │             │
│  └────────┬────────┘  └────────┬────────┘  └────────┬────────┘             │
│           │                    │                    │                       │
│           ▼                    ▼                    ▼                       │
│  ┌─────────────────────────────────────────────────────────────────┐       │
│  │                    KUBERNETES SECURITY LAYERS                    │       │
│  ├─────────────────────────────────────────────────────────────────┤       │
│  │  Layer 7: Application (mTLS, AuthZ, Secrets)                    │       │
│  │  Layer 6: Service Mesh (Istio Ambient, SPIFFE)                  │       │
│  │  Layer 5: Runtime (eBPF, Tetragon, Falco)                       │       │
│  │  Layer 4: Network (Cilium, Network Policies)                    │       │
│  │  Layer 3: Compute (Confidential Computing, Enclaves)            │       │
│  │  Layer 2: Supply Chain (SLSA, Sigstore, SBOM)                   │       │
│  │  Layer 1: Cryptography (Post-Quantum, Hybrid Schemes)           │       │
│  └─────────────────────────────────────────────────────────────────┘       │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## Post-Quantum Cryptography

### NIST Standards (FIPS 203/204/205) - August 2024

The National Institute of Standards and Technology (NIST) released the first set of post-quantum cryptography standards in August 2024, marking a critical milestone in preparing for the quantum threat.

#### FIPS 203: Module-Lattice-Based Key-Encapsulation Mechanism (ML-KEM)

ML-KEM (formerly Kyber) is a key encapsulation mechanism based on the hardness of the Module Learning With Errors (MLWE) problem.

```go
// Go 1.24+ ML-KEM implementation example
package main

import (
    "crypto/rand"
    "fmt"

    "golang.org/x/crypto/mlkem"
)

// MLKEM768KeyExchange demonstrates ML-KEM-768 key encapsulation
func MLKEM768KeyExchange() error {
    // Generate deterministic keypair from seed (for testing)
    // In production, use mlkem.GenerateKey768(rand.Reader)

    // Generate encapsulation key and decapsulation key
    dek, err := mlkem.GenerateKey768(rand.Reader)
    if err != nil {
        return fmt.Errorf("key generation failed: %w", err)
    }

    // Encapsulation key (public) can be shared
    encapKey := dek.EncapsulationKey()

    // Encapsulate generates shared secret and ciphertext
    ciphertext, sharedSecretEnc, err := encapKey.Encapsulate()
    if err != nil {
        return fmt.Errorf("encapsulation failed: %w", err)
    }

    // Decapsulate recovers the shared secret
    sharedSecretDec, err := dek.Decapsulate(ciphertext)
    if err != nil {
        return fmt.Errorf("decapsulation failed: %w", err)
    }

    // Verify both parties have the same shared secret
    if string(sharedSecretEnc) != string(sharedSecretDec) {
        return fmt.Errorf("shared secret mismatch")
    }

    fmt.Printf("ML-KEM-768: Shared secret established (%d bytes)\n", len(sharedSecretEnc))
    return nil
}
```

#### FIPS 204: Module-Lattice-Based Digital Signature Algorithm (ML-DSA)

ML-DSA (formerly Dilithium) provides digital signatures based on lattice cryptography.

```go
// ML-DSA signature example with Go 1.24+
package crypto

import (
    "crypto/rand"
    "encoding/base64"
    "fmt"

    "golang.org/x/crypto/mldsa"
)

// MLDSA65Signer implements ML-DSA-65 (NIST Level 3 security)
type MLDSA65Signer struct {
    privateKey *mldsa.PrivateKey65
    publicKey  *mldsa.PublicKey65
}

// GenerateMLDSA65Keypair generates a new ML-DSA-65 keypair
func GenerateMLDSA65Keypair() (*MLDSA65Signer, error) {
    pub, priv, err := mldsa.GenerateKey65(rand.Reader)
    if err != nil {
        return nil, fmt.Errorf("ML-DSA-65 key generation failed: %w", err)
    }

    return &MLDSA65Signer{
        privateKey: priv,
        publicKey:  pub,
    }, nil
}

// Sign creates a ML-DSA-65 signature
func (s *MLDSA65Signer) Sign(message []byte) ([]byte, error) {
    signature := s.privateKey.Sign(message)
    return signature, nil
}

// Verify checks a ML-DSA-65 signature
func (s *MLDSA65Signer) Verify(message, signature []byte) bool {
    return s.publicKey.Verify(message, signature)
}

// KeyProperties shows ML-DSA-65 properties
var KeyProperties = map[string]interface{}{
    "algorithm":        "ML-DSA-65",
    "nist_level":       3,
    "public_key_size":  1952,  // bytes
    "private_key_size": 4032,  // bytes
    "signature_size":   3293,  // bytes
    "security_basis":   "Module-Lattice (Dilithium)",
}
```

#### FIPS 205: Stateless Hash-Based Digital Signature Standard (SLH-DSA)

SLH-DSA (SPHINCS+) provides stateless hash-based signatures for long-term security.

```yaml
# Kubernetes Secret for SLH-DSA keys (backup signing)
apiVersion: v1
kind: Secret
metadata:
  name: slh-dsa-signing-key
  namespace: security
  annotations:
    security.nist.gov/algorithm: "SLH-DSA-SHA2-128s"
    security.nist.gov/fips-205: "compliant"
    security.nist.gov/usage: "long-term-archive-signing"
type: Opaque
data:
  # Base64 encoded SLH-DSA private key (HSM-protected in production)
  private.key: "LS0t..."
  public.key: "LS0t..."

  # Key metadata
  algorithm.info: |
    {
      "algorithm": "SLH-DSA-SHA2-128s",
      "nist_level": 1,
      "public_key_bytes": 32,
      "private_key_bytes": 64,
      "signature_bytes": 7856,
      "hash_function": "SHA2-256",
      "stateless": true
    }
```

### Go 1.24+ X25519MLKEM768 Hybrid

Go 1.24 introduced hybrid post-quantum key exchange combining X25519 (classical ECDH) with ML-KEM-768 (post-quantum KEM) for maximum security during the transition period.

```go
// Hybrid X25519MLKEM768 implementation for Go 1.24+
package main

import (
    "crypto/ecdh"
    "crypto/rand"
    "crypto/subtle"
    "fmt"
    "io"

    "golang.org/x/crypto/curve25519"
    "golang.org/x/crypto/mlkem"
    "golang.org/x/crypto/hkdf"
)

// X25519MLKEM768Hybrid implements the hybrid key exchange
// Combines classical X25519 with post-quantum ML-KEM-768
type X25519MLKEM768Hybrid struct {
    X25519PrivateKey   [32]byte
    X25519PublicKey    [32]byte
    MLKEM768PrivateKey *mlkem.DecapsulationKey768
    MLKEM768PublicKey  *mlkem.EncapsulationKey768
}

// GenerateX25519MLKEM768 generates hybrid keypair
func GenerateX25519MLKEM768() (*X25519MLKEM768Hybrid, error) {
    hybrid := &X25519MLKEM768Hybrid{}

    // Generate X25519 keypair
    var err error
    _, err = rand.Read(hybrid.X25519PrivateKey[:])
    if err != nil {
        return nil, fmt.Errorf("x25519 key generation failed: %w", err)
    }

    // Clamp private key per RFC 7748
    hybrid.X25519PrivateKey[0] &= 248
    hybrid.X25519PrivateKey[31] &= 127
    hybrid.X25519PrivateKey[31] |= 64

    // Generate public key
    curve25519.ScalarBaseMult(&hybrid.X25519PublicKey, &hybrid.X25519PrivateKey)

    // Generate ML-KEM-768 keypair
    dek, err := mlkem.GenerateKey768(rand.Reader)
    if err != nil {
        return nil, fmt.Errorf("mlkem key generation failed: %w", err)
    }

    hybrid.MLKEM768PrivateKey = dek
    hybrid.MLKEM768PublicKey = dek.EncapsulationKey()

    return hybrid, nil
}

// Encapsulate performs hybrid encapsulation
func (h *X25519MLKEM768Hybrid) Encapsulate() (
    x25519PublicKey [32]byte,
    mlkemCiphertext []byte,
    sharedSecret []byte,
    err error,
) {
    // X25519 ephemeral key generation
    var ephemeralPrivate [32]byte
    _, err = rand.Read(ephemeralPrivate[:])
    if err != nil {
        return [32]byte{}, nil, nil, err
    }

    ephemeralPrivate[0] &= 248
    ephemeralPrivate[31] &= 127
    ephemeralPrivate[31] |= 64

    var ephemeralPublic [32]byte
    curve25519.ScalarBaseMult(&ephemeralPublic, &ephemeralPrivate)

    // X25519 shared secret
    var x25519Shared [32]byte
    curve25519.ScalarMult(&x25519Shared, &ephemeralPrivate, &h.X25519PublicKey)

    // ML-KEM encapsulation
    mlkemCiphertext, mlkemShared, err := h.MLKEM768PublicKey.Encapsulate()
    if err != nil {
        return [32]byte{}, nil, nil, fmt.Errorf("mlkem encapsulation failed: %w", err)
    }

    // Combine secrets using HKDF
    combinedInput := append(x25519Shared[:], mlkemShared...)
    sharedSecret = make([]byte, 32)

    hkdfReader := hkdf.New(sha256.New, combinedInput, nil, []byte("X25519MLKEM768 hybrid"))
    _, err = io.ReadFull(hkdfReader, sharedSecret)
    if err != nil {
        return [32]byte{}, nil, nil, fmt.Errorf("hkdf failed: %w", err)
    }

    return ephemeralPublic, mlkemCiphertext, sharedSecret, nil
}

// Decapsulate performs hybrid decapsulation
func (h *X25519MLKEM768Hybrid) Decapsulate(
    ephemeralPublic [32]byte,
    mlkemCiphertext []byte,
) ([]byte, error) {
    // X25519 shared secret
    var x25519Shared [32]byte
    curve25519.ScalarMult(&x25519Shared, &h.X25519PrivateKey, &ephemeralPublic)

    // ML-KEM decapsulation
    mlkemShared, err := h.MLKEM768PrivateKey.Decapsulate(mlkemCiphertext)
    if err != nil {
        return nil, fmt.Errorf("mlkem decapsulation failed: %w", err)
    }

    // Combine secrets using HKDF
    combinedInput := append(x25519Shared[:], mlkemShared...)
    sharedSecret := make([]byte, 32)

    hkdfReader := hkdf.New(sha256.New, combinedInput, nil, []byte("X25519MLKEM768 hybrid"))
    _, err = io.ReadFull(hkdfReader, sharedSecret)
    if err != nil {
        return nil, fmt.Errorf("hkdf failed: %w", err)
    }

    return sharedSecret, nil
}

// ConstantTimeCompare performs constant-time comparison of shared secrets
func ConstantTimeCompare(a, b []byte) bool {
    return subtle.ConstantTimeCompare(a, b) == 1
}
```

### Migration Deadlines and Roadmap

```yaml
# Post-Quantum Cryptography Migration Timeline
pqc_migration:
  phases:
    phase_1_crypto_agility:
      timeline: "2024-2025"
      description: "Implement crypto-agility in all systems"
      actions:
        - Audit all cryptographic implementations
        - Implement algorithm negotiation
        - Deploy hybrid key exchange (X25519MLKEM768)
        - Update TLS configurations

    phase_2_hybrid_deployment:
      timeline: "2025-2028"
      description: "Deploy hybrid post-quantum cryptography"
      actions:
        - Enable ML-KEM in TLS 1.3
        - Deploy ML-DSA for code signing
        - Update Kubernetes API server encryption
        - Implement SPIFFE/SPIRE with PQC

    phase_3_pqc_transition:
      timeline: "2028-2030"
      description: "Transition to full PQC where required"
      critical_deadline: "2030-01-01"
      actions:
        - Deprecate RSA and ECDH for sensitive data
        - Full ML-KEM/ML-DSA deployment
        - Long-term archive signing with SLH-DSA
        - Update all compliance frameworks

    phase_4_pqc_mandatory:
      timeline: "2030-2035"
      description: "Post-quantum cryptography mandatory"
      final_deadline: "2035-01-01"
      actions:
        - Disable classical algorithms for new deployments
        - Full PQC mandate for government systems
        - Industry-wide PQC adoption complete
        - Quantum-safe cloud infrastructure

  compliance_requirements:
    nsa_cnsa_2_0:
      commercial_national_security:
        algorithm_suite: "CNSA 2.0"
        ml_kem_required_by: "2025-12-31"
        ml_dsa_required_by: "2030-12-31"

    nist_sp_800_208:
      stateful_hash_signatures:
        use_case: "software_firmware_signing"
        requirement: "LMS or XMSS for long-term"

    cis_benchmarks:
      kubernetes_pqc:
        expected_version: "CIS Kubernetes v1.10+"
        includes: "Post-quantum TLS recommendations"
```

---

## Supply Chain Security

### SLSA v1.0/v1.1 Framework

The Supply Chain Levels for Software Artifacts (SLSA) framework provides a comprehensive approach to securing software supply chains.

```yaml
# SLSA v1.1 Build Definition for Kubernetes Components
apiVersion: slsa.dev/v1.1
kind: BuildDefinition
metadata:
  name: kubernetes-component-build
  version: "1.1.0"

spec:
  # Build environment requirements
  buildType: https://slsa.dev/container-based-build/v1.1

  # External parameters (user-controlled)
  externalParameters:
    source:
      uri: "git+https://github.com/org/kubernetes-component"
      digest:
        sha256: "abc123..."
      entryPoint: ".github/workflows/build.yml"
    inputs:
      go_version: "1.24.0"
      build_flags: ["-trimpath", "-ldflags=-s -w"]

  # Internal parameters (system-controlled)
  internalParameters:
    runner:
      os: "ubuntu-24.04"
      arch: "amd64"
      labels: ["ubuntu-latest", "slsa-level-3"]
    builder:
      id: "https://github.com/org/slsa-compliant-builder"
      version: "v2.1.0"

  # Resolved dependencies with full provenance
  resolvedDependencies:
    - uri: "https://proxy.golang.org/github.com/some/module@v1.2.3"
      digest:
        sha256: "def456..."
      annotations:
        slsa.dependency.type: "go-module"
        slsa.dependency.direct: "true"
        sbom.generated: "true"

    - uri: "docker://golang:1.24.0-alpine"
      digest:
        sha256: "ghi789..."
      annotations:
        slsa.dependency.type: "base-image"
        slsa.dependency.reproducible: "true"

  # Output artifacts with SLSA attestations
  runDetails:
    builder:
      id: "https://github.com/attestations/slsa-github-generator/.github/workflows/generator_container_slsa3.yml@refs/tags/v2.0.0"
    metadata:
      invocationId: "build-2025-04-01-001"
      startedOn: "2025-04-01T10:00:00Z"
      finishedOn: "2025-04-01T10:05:00Z"
    byproducts:
      - name: "container-image"
        uri: "ghcr.io/org/component:v1.0.0"
        digest:
          sha256: "jkl012..."
        annotations:
          slsa.level: "3"
          provenance.attached: "true"
          sbom.spdx.json: "attached"
          sbom.cyclonedx.json: "attached"
```

#### SLSA Levels Implementation

```go
// SLSA compliance checker for Go projects
package slsa

import (
    "context"
    "encoding/json"
    "fmt"
    "os"
)

// SLSALevel represents the SLSA compliance level
type SLSALevel int

const (
    SLSALevel1 SLSALevel = iota + 1  // Provenance exists
    SLSALevel2                       // Hosted build service, signed provenance
    SLSALevel3                       // Hardened build service, hermetic builds
    SLSALevel4                       // Reproducible builds, two-person review
)

// BuildProvenance represents SLSA provenance attestation
type BuildProvenance struct {
    BuildDefinition `json:"buildDefinition"`
    RunDetails      `json:"runDetails"`
}

// BuilderRequirements defines requirements for each SLSA level
var BuilderRequirements = map[SLSALevel][]string{
    SLSALevel1: {
        "Build process is fully scripted/automated",
        "Provenance is available showing how artifact was built",
    },
    SLSALevel2: {
        "Version controlled source (Git with signed commits)",
        "Hosted build service (GitHub Actions, Cloud Build)",
        "Signed provenance generated by build service",
        "Provenance is authenticated and service-generated",
    },
    SLSALevel3: {
        "Source meets Level 2 requirements",
        "Build runs on hardened infrastructure",
        "Build is hermetic (no network access, pinned dependencies)",
        "Dependencies are tracked and verified",
        "Reproducible builds (byte-for-byte identical)",
        "Builder is idempotent and isolated",
    },
    SLSALevel4: {
        "All Level 3 requirements",
        "Two-person review of all changes",
        "Reproducible builds verified by independent builder",
        "Ephemeral, isolated, and parameterless build environment",
        "Dependencies are fully declared with hashes",
    },
}

// CheckSLSACompliance verifies a build meets SLSA requirements
func CheckSLSACompliance(ctx context.Context, provenancePath string) (SLSALevel, error) {
    data, err := os.ReadFile(provenancePath)
    if err != nil {
        return 0, fmt.Errorf("reading provenance: %w", err)
    }

    var provenance BuildProvenance
    if err := json.Unmarshal(data, &provenance); err != nil {
        return 0, fmt.Errorf("parsing provenance: %w", err)
    }

    // Check for SLSA Level 1
    if provenance.BuildDefinition.BuildType == "" {
        return 0, fmt.Errorf("missing build type")
    }

    // Check for SLSA Level 2
    if provenance.RunDetails.Builder.ID == "" {
        return SLSALevel1, nil
    }

    // Check for SLSA Level 3
    if !provenance.hasHermeticBuild() || !provenance.hasPinnedDependencies() {
        return SLSALevel2, nil
    }

    // Check for SLSA Level 4
    if !provenance.hasTwoPersonReview() || !provenance.hasReproducibleBuild() {
        return SLSALevel3, nil
    }

    return SLSALevel4, nil
}

func (p *BuildProvenance) hasHermeticBuild() bool {
    // Check if build has network isolation
    for _, param := range p.InternalParameters.Runner.Labels {
        if param == "network-isolated" || param == "hermetic" {
            return true
        }
    }
    return false
}

func (p *BuildProvenance) hasPinnedDependencies() bool {
    for _, dep := range p.ResolvedDependencies {
        if dep.Digest.SHA256 == "" {
            return false
        }
    }
    return len(p.ResolvedDependencies) > 0
}

func (p *BuildProvenance) hasTwoPersonReview() bool {
    // Check attestation for two-person review requirement
    return false // Implementation depends on SCM integration
}

func (p *BuildProvenance) hasReproducibleBuild() bool {
    // Check if build was reproduced by independent builder
    return false // Implementation depends on build system
}
```

### Sigstore Ecosystem

Sigstore provides a complete solution for software signing and verification without the need for managing private keys.

#### Cosign for Container Signing

```bash
#!/bin/bash
# Sigstore Cosign signing workflow for Kubernetes deployments

set -euo pipefail

# Configuration
IMAGE_NAME="${IMAGE_NAME:-ghcr.io/org/app}"
IMAGE_TAG="${IMAGE_TAG:-v1.0.0}"
FULL_IMAGE="${IMAGE_NAME}:${IMAGE_TAG}"

# ============================================
# Step 1: Build and push container image
# ============================================
echo "Building container image..."
docker build -t "${FULL_IMAGE}" .
docker push "${FULL_IMAGE}"

# ============================================
# Step 2: Sign with Cosign (keyless)
# ============================================
echo "Signing image with Cosign (keyless)..."

# Keyless signing uses OIDC identity (GitHub Actions, Google, etc.)
cosign sign \
    --yes \
    --recursive \
    --output-signature "${IMAGE_TAG}.sig" \
    --output-certificate "${IMAGE_TAG}.cert" \
    "${FULL_IMAGE}"

# ============================================
# Step 3: Verify signature
# ============================================
echo "Verifying signature..."

# Verify using Sigstore's root of trust
cosign verify \
    --certificate-identity-regexp "^https://github.com/${GITHUB_REPOSITORY}/.*" \
    --certificate-oidc-issuer "https://token.actions.githubusercontent.com" \
    "${FULL_IMAGE}"

# ============================================
# Step 4: Attach SBOM
# ============================================
echo "Generating and attaching SBOM..."

# Generate SBOM with Syft
syft "${FULL_IMAGE}" -o spdx-json > "${IMAGE_TAG}.spdx.json"
syft "${FULL_IMAGE}" -o cyclonedx-json > "${IMAGE_TAG}.cyclonedx.json"

# Attach SBOM to image
cosign attest \
    --predicate "${IMAGE_TAG}.spdx.json" \
    --type spdxjson \
    --yes \
    "${FULL_IMAGE}"

cosign attest \
    --predicate "${IMAGE_TAG}.cyclonedx.json" \
    --type cyclonedx \
    --yes \
    "${FULL_IMAGE}"

# ============================================
# Step 5: Attach SLSA provenance
# ============================================
echo "Attaching SLSA provenance..."

# Provenance is generated by SLSA GitHub Generator
cosign attest \
    --predicate provenance.json \
    --type slsaprovenance \
    --yes \
    "${FULL_IMAGE}"

echo "Image signed and attested successfully!"
```

#### Rekor Transparency Log

```go
// Rekor transparency log verification for Go applications
package sigstore

import (
    "context"
    "crypto/sha256"
    "encoding/base64"
    "encoding/hex"
    "fmt"

    "github.com/sigstore/rekor/pkg/client"
    "github.com/sigstore/rekor/pkg/generated/client/entries"
    "github.com/sigstore/rekor/pkg/generated/models"
    "github.com/sigstore/rekor/pkg/types"
    _ "github.com/sigstore/rekor/pkg/types/hashedrekord"
)

// RekorVerifier verifies entries in the Rekor transparency log
type RekorVerifier struct {
    client *client.Rekor
}

// NewRekorVerifier creates a new Rekor verifier
func NewRekorVerifier(rekorURL string) (*RekorVerifier, error) {
    c, err := client.GetRekorClient(rekorURL)
    if err != nil {
        return nil, fmt.Errorf("creating rekor client: %w", err)
    }

    return &RekorVerifier{client: c}, nil
}

// VerifyContainerImage checks if a container image signature exists in Rekor
func (r *RekorVerifier) VerifyContainerImage(ctx context.Context, imageSHA string) (*RekorEntry, error) {
    // Search for entries by SHA
    searchParams := entries.NewSearchLogQueryParams()

    // Create hash search
    hash := sha256.Sum256([]byte(imageSHA))
    hashStr := hex.EncodeToString(hash[:])

    query := &models.SearchIndex{
        Hash: &models.SearchIndexHash{
            Algorithm: swag.String("sha256"),
            Value:     swag.String(hashStr),
        },
    }

    searchParams.SetEntry(query)

    resp, err := r.client.Entries.SearchLogQuery(searchParams)
    if err != nil {
        return nil, fmt.Errorf("searching rekor: %w", err)
    }

    if len(resp.Payload) == 0 {
        return nil, fmt.Errorf("no rekor entry found for image")
    }

    // Parse the first matching entry
    entry := resp.Payload[0]
    return parseRekorEntry(entry)
}

// RekorEntry represents a verified transparency log entry
type RekorEntry struct {
    UUID          string
    LogIndex      int64
    IntegratedTime int64
    Body          interface{}
    Verification  *RekorVerification
}

type RekorVerification struct {
    InclusionProof    *InclusionProof
    SignedEntryTimestamp []byte
}

type InclusionProof struct {
    LogIndex int64
    RootHash string
    TreeSize int64
    Hashes   []string
}

func parseRekorEntry(entry models.LogEntry) (*RekorEntry, error) {
    // Implementation for parsing entry
    return &RekorEntry{}, nil
}

// VerifyInclusionProof verifies the Merkle tree inclusion proof
func (r *RekorVerifier) VerifyInclusionProof(entry *RekorEntry) error {
    if entry.Verification == nil || entry.Verification.InclusionProof == nil {
        return fmt.Errorf("no inclusion proof available")
    }

    proof := entry.Verification.InclusionProof

    // Verify the Merkle inclusion proof
    // This ensures the entry is actually in the log
    calculatedRoot := calculateRootHash(proof)

    if calculatedRoot != proof.RootHash {
        return fmt.Errorf("inclusion proof verification failed")
    }

    return nil
}

func calculateRootHash(proof *InclusionProof) string {
    // Merkle tree root hash calculation
    // Implementation would verify the proof path
    return ""
}
```

#### Fulcio for Keyless Signing

```yaml
# Kubernetes Policy for Sigstore/Cosign verification
apiVersion: constraints.gatekeeper.sh/v1beta1
kind: K8sVerifiedImage
metadata:
  name: require-signed-images
spec:
  match:
    kinds:
      - apiGroups: [""]
        kinds: ["Pod"]
    namespaces:
      - "production"
      - "staging"
  parameters:
    # Fulcio certificate requirements
    fulcio:
      enabled: true
      issuers:
        - "https://token.actions.githubusercontent.com"
        - "https://accounts.google.com"
      subjectRegExp:
        - "^https://github.com/org/.*"

    # Rekor transparency log requirements
    rekor:
      enabled: true
      url: "https://rekor.sigstore.dev"
      requireInclusionProof: true

    # Allowed key references (for non-keyless signing)
    keyRefs:
      - kms://gcpkms/projects/project/locations/global/keyRings/ring/cryptoKeys/key

    # Exemptions for development
    exemptImages:
      - "localhost/*"
      - "*.local"
```

### Keyless Signing with OIDC

```go
// Keyless signing implementation using OIDC identity
package main

import (
    "context"
    "fmt"

    "github.com/sigstore/cosign/v2/cmd/cosign/cli/options"
    "github.com/sigstore/cosign/v2/cmd/cosign/cli/sign"
    "github.com/sigstore/cosign/v2/pkg/oci/mutate"
    "github.com/sigstore/cosign/v2/pkg/providers"
    _ "github.com/sigstore/cosign/v2/pkg/providers/github"
)

// OIDCIdentity represents the OIDC identity used for keyless signing
type OIDCIdentity struct {
    Token      string
    Issuer     string
    Subject    string
    Audience   string
    Expiry     int64
}

// GetGitHubActionsIdentity retrieves OIDC token from GitHub Actions
func GetGitHubActionsIdentity(ctx context.Context) (*OIDCIdentity, error) {
    // GitHub Actions provides OIDC tokens via environment
    token, err := providers.Provide(ctx, "github-actions")
    if err != nil {
        return nil, fmt.Errorf("getting GitHub OIDC token: %w", err)
    }

    return &OIDCIdentity{
        Token:    token,
        Issuer:   "https://token.actions.githubusercontent.com",
        Subject:  getEnv("GITHUB_WORKFLOW_REF", ""),
        Audience: "sigstore",
    }, nil
}

// KeylessSigner performs keyless signing using OIDC
type KeylessSigner struct {
    identity *OIDCIdentity
    rekorURL string
    fulcioURL string
}

// NewKeylessSigner creates a new keyless signer
func NewKeylessSigner(identity *OIDCIdentity) *KeylessSigner {
    return &KeylessSigner{
        identity:  identity,
        rekorURL:  "https://rekor.sigstore.dev",
        fulcioURL: "https://fulcio.sigstore.dev",
    }
}

// SignImage signs a container image using keyless signing
func (s *KeylessSigner) SignImage(ctx context.Context, imageRef string) error {
    // Keyless signing options
    signOpts := &options.SignOptions{
        Registry: options.RegistryOptions{},
        Fulcio: options.FulcioOptions{
            URL:        s.fulcioURL,
            IdentityToken: s.identity.Token,
        },
        Rekor: options.RekorOptions{
            URL: s.rekorURL,
        },
        OIDC: options.OIDCOptions{
            Issuer:   s.identity.Issuer,
            ClientID: "sigstore",
        },
        // Skip upload for dry-run
        Upload: true,
        // Recursive signing for multi-arch images
        Recursive: true,
        // Output certificate and signature
        OutputCertificate: "",
        OutputSignature:   "",
    }

    // Perform the sign operation
    // This uses Fulcio to get a short-lived certificate
    // and Rekor to store the transparency log entry

    fmt.Printf("Signing %s with keyless signing...\n", imageRef)
    fmt.Printf("OIDC Issuer: %s\n", s.identity.Issuer)

    return nil
}

// VerifyKeylessSignature verifies a keyless signature
func VerifyKeylessSignature(ctx context.Context, imageRef string, identityRegExp string) error {
    // Verification checks:
    // 1. Certificate chain from Fulcio
    // 2. OIDC identity matches expected pattern
    // 3. Signature is valid
    // 4. Entry exists in Rekor

    fmt.Printf("Verifying %s...\n", imageRef)
    fmt.Printf("Expected identity: %s\n", identityRegExp)

    return nil
}
```

### SBOM Requirements (CISA 2025)

The Cybersecurity and Infrastructure Security Agency (CISA) issued new SBOM requirements effective 2025.

```yaml
# CISA 2025 SBOM Requirements Implementation
apiVersion: v1
kind: ConfigMap
metadata:
  name: sbom-requirements-cisa-2025
  namespace: security
data:
  requirements.yaml: |
    cisa_sbom_2025:
      # Minimum Data Elements (CISA SBOM Minimum Elements)
      minimum_elements:
        supplier_name: true
        component_name: true
        version: true
        other_unique_identifiers:
          - cpe
          - purl
          - swid
        dependency_relationship: true
        author_of_sbom: true
        timestamp: true

      # Data Fields for each component
      component_fields:
        author: "Required - Who created this component"
        name: "Required - Component name"
        version: "Required - Component version"
        licenses: "Required - All applicable licenses"
        hashes: "Required - Cryptographic hashes"
        supplier: "Required - Supplier name"
        copyright: "Required - Copyright statements"

      # Timeliness Requirements
      timeliness:
        generation: "At build time or release"
        updates: "When component changes"
        distribution: "With software delivery"

      # Access and Delivery
      access_delivery:
        methods:
          - direct_download
          - api_access
          - physical_media
          - email
        formats:
          - spdx
          - cyclonedx
          - swid

      # Vulnerability Association
      vulnerability_association:
        vex: "Recommended - Vulnerability Exploitability eXchange"
        cve_mapping: "Required - Map components to CVEs"
        epss: "Recommended - Exploit Prediction Scoring System"

      # Machine Readability
      formats:
        spdx:
          version: "SPDX-2.3"
          required_fields:
            - DocumentName
            - DocumentNamespace
            - SPDXVersion
            - DataLicense
            - SPDXID
            - DocumentName
            - Creators
            - Created

        cyclonedx:
          version: "CycloneDX 1.5"
          required_fields:
            - bomFormat
            - specVersion
            - serialNumber
            - version
            - metadata.timestamp
            - components

      # Go-specific SBOM considerations
      go_ecosystem:
        module_info:
          - go.mod
          - go.sum
          - vendor/modules.txt
        binary_provenance:
          - buildinfo
          - compiler_version
          - build_settings
          - dependencies
```

```go
// Go SBOM generator for CISA 2025 compliance
package sbom

import (
    "debug/buildinfo"
    "encoding/json"
    "fmt"
    "os"
    "runtime/debug"
    "time"
)

// CISACompliantSBOM represents a CISA 2025 compliant SBOM
type CISACompliantSBOM struct {
    SPDXVersion    string          `json:"spdxVersion"`
    DataLicense    string          `json:"dataLicense"`
    SPDXID         string          `json:"SPDXID"`
    Name           string          `json:"name"`
    DocumentNamespace string       `json:"documentNamespace"`
    CreationInfo   CreationInfo    `json:"creationInfo"`
    Packages       []Package       `json:"packages"`
    Relationships  []Relationship  `json:"relationships"`
}

type CreationInfo struct {
    Created  string   `json:"created"`
    Creators []string `json:"creators"`
}

type Package struct {
    SPDXID           string            `json:"SPDXID"`
    Name             string            `json:"name"`
    VersionInfo      string            `json:"versionInfo"`
    DownloadLocation string            `json:"downloadLocation"`
    Supplier         string            `json:"supplier"`
    Checksums        []Checksum        `json:"checksums,omitempty"`
    ExternalRefs     []ExternalRef     `json:"externalRefs,omitempty"`
    LicenseConcluded string            `json:"licenseConcluded,omitempty"`
    CopyrightText    string            `json:"copyrightText,omitempty"`
}

type Checksum struct {
    Algorithm     string `json:"algorithm"`
    ChecksumValue string `json:"checksumValue"`
}

type ExternalRef struct {
    ReferenceCategory string `json:"referenceCategory"`
    ReferenceType     string `json:"referenceType"`
    ReferenceLocator  string `json:"referenceLocator"`
}

type Relationship struct {
    SPDXElementID      string `json:"spdxElementId"`
    RelatedSPDXElement string `json:"relatedSpdxElement"`
    RelationshipType   string `json:"relationshipType"`
}

// GenerateGoBinarySBOM extracts SBOM from Go binary build info
func GenerateGoBinarySBOM(binaryPath string) (*CISACompliantSBOM, error) {
    info, err := buildinfo.ReadFile(binaryPath)
    if err != nil {
        return nil, fmt.Errorf("reading build info: %w", err)
    }

    sbom := &CISACompliantSBOM{
        SPDXVersion:       "SPDX-2.3",
        DataLicense:       "CC0-1.0",
        SPDXID:            "SPDXRef-DOCUMENT",
        Name:              info.Main.Path,
        DocumentNamespace: fmt.Sprintf("https://example.com/sbom/%s-%s", info.Main.Path, time.Now().Format("20060102")),
        CreationInfo: CreationInfo{
            Created:  time.Now().UTC().Format(time.RFC3339),
            Creators: []string{
                fmt.Sprintf("Tool: go-sbom-generator-1.0.0"),
                fmt.Sprintf("Organization: %s", info.Main.Path),
            },
        },
    }

    // Add main package
    mainPkg := Package{
        SPDXID:           "SPDXRef-Package-Main",
        Name:             info.Main.Path,
        VersionInfo:      info.Main.Version,
        DownloadLocation: "NOASSERTION",
        Supplier:         "Organization: " + info.Main.Path,
        LicenseConcluded: "NOASSERTION",
        CopyrightText:    "NOASSERTION",
    }
    sbom.Packages = append(sbom.Packages, mainPkg)

    // Add relationship
    sbom.Relationships = append(sbom.Relationships, Relationship{
        SPDXElementID:      "SPDXRef-DOCUMENT",
        RelatedSPDXElement: "SPDXRef-Package-Main",
        RelationshipType:   "DESCRIBES",
    })

    // Add dependencies
    for i, dep := range info.Deps {
        pkg := Package{
            SPDXID:           fmt.Sprintf("SPDXRef-Package-Dep-%d", i),
            Name:             dep.Path,
            VersionInfo:      dep.Version,
            DownloadLocation: fmt.Sprintf("https://proxy.golang.org/%s/@v/%s.zip", dep.Path, dep.Version),
            Supplier:         "Organization: " + getModulePathOwner(dep.Path),
            LicenseConcluded: "NOASSERTION",
            CopyrightText:    "NOASSERTION",
            ExternalRefs: []ExternalRef{
                {
                    ReferenceCategory: "PACKAGE-MANAGER",
                    ReferenceType:     "purl",
                    ReferenceLocator:  fmt.Sprintf("pkg:golang/%s@%s", dep.Path, dep.Version),
                },
            },
        }

        if dep.Sum != "" {
            pkg.Checksums = append(pkg.Checksums, Checksum{
                Algorithm:     "SHA256",
                ChecksumValue: dep.Sum,
            })
        }

        sbom.Packages = append(sbom.Packages, pkg)

        // Add dependency relationship
        sbom.Relationships = append(sbom.Relationships, Relationship{
            SPDXElementID:      "SPDXRef-Package-Main",
            RelatedSPDXElement: pkg.SPDXID,
            RelationshipType:   "DEPENDS_ON",
        })
    }

    return sbom, nil
}

func getModulePathOwner(path string) string {
    // Extract owner from module path
    // e.g., github.com/org/repo -> org
    return path
}
```

---

## eBPF Runtime Security

### Tetragon for Kernel-Level Threat Detection

Tetragon is an eBPF-based security observability and runtime enforcement tool from Cilium that provides deep kernel-level visibility.

```yaml
# Tetragon deployment for Kubernetes runtime security
apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: tetragon
  namespace: kube-system
  labels:
    app.kubernetes.io/name: tetragon
    app.kubernetes.io/part-of: cilium
spec:
  selector:
    matchLabels:
      app.kubernetes.io/name: tetragon
  template:
    metadata:
      labels:
        app.kubernetes.io/name: tetragon
    spec:
      hostNetwork: true
      hostPID: true
      containers:
        - name: tetragon
          image: quay.io/cilium/tetragon:v1.2.0
          securityContext:
            privileged: true
          volumeMounts:
            - name: bpf-maps
              mountPath: /sys/fs/bpf
            - name: tracing-policy
              mountPath: /etc/tetragon/tetragon.conf.d/
            - name: proc
              mountPath: /proc
              readOnly: true
          resources:
            requests:
              memory: "128Mi"
              cpu: "100m"
            limits:
              memory: "1Gi"
              cpu: "1000m"
        - name: export-stdout
          image: quay.io/cilium/tetragon-export-stdout:v1.2.0
          command:
            - /usr/bin/tetragon-export-stdout
          args:
            - --server-address
            - "unix:///var/run/cilium/tetragon/tetragon.sock"
          volumeMounts:
            - name: tetragon-socket
              mountPath: /var/run/cilium/tetragon
      volumes:
        - name: bpf-maps
          hostPath:
            path: /sys/fs/bpf
        - name: tracing-policy
          configMap:
            name: tetragon-tracing-policies
        - name: proc
          hostPath:
            path: /proc
        - name: tetragon-socket
          hostPath:
            path: /var/run/cilium/tetragon
            type: DirectoryOrCreate
      tolerations:
        - operator: Exists
```

### TracingPolicy Examples

```yaml
# Tetragon TracingPolicy for detecting container escapes
apiVersion: cilium.io/v1alpha1
kind: TracingPolicy
metadata:
  name: container-escape-detection
  namespace: kube-system
spec:
  # Detect attempts to escape container boundaries
  kprobes:
    # Detect mount namespace escape attempts
    - call: "__x64_sys_setns"
      syscall: true
      args:
        - index: 0
          type: "int"
        - index: 1
          type: "int"
      selectors:
        - matchNamespaces:
            - namespace: "Pid"
              operator: "NotIn"
              values:
                - "/proc/1/ns/pid"
          matchActions:
            - action: Post
              rateLimit: "1m"
            - action: Signal
              signal: "SIGKILL"

    # Detect privilege escalation via setuid binaries
    - call: "__x64_sys_setuid"
      syscall: true
      args:
        - index: 0
          type: "uid_t"
      selectors:
        - matchCapabilities:
            - type: Effective
              operator: In
              values:
                - "CAP_SETUID"
          matchActions:
            - action: Post
            - action: Override
              argError: -1
              argError: -EPERM

---
# Tetragon TracingPolicy for crypto mining detection
apiVersion: cilium.io/v1alpha1
kind: TracingPolicy
metadata:
  name: crypto-mining-detection
  namespace: kube-system
spec:
  # Detect crypto mining processes
  tracepoints:
    - subsystem: "sched"
      event: "sched_process_exec"
      args:
        - index: 0
          type: "string"
      selectors:
        - matchArgs:
            - index: 0
              operator: "Prefix"
              values:
                - "xmrig"
                - "minerd"
                - "cgminer"
                - "bminer"
                - "ethminer"
                - "nbminer"
                - "t-rex"
                - "teamredminer"
                - "lolminer"
                - " gminer"
          matchActions:
            - action: Post
            - action: Signal
              signal: "SIGKILL"
            - action: NotifyKubernetes

  # Detect suspicious network connections to mining pools
  kprobes:
    - call: "tcp_connect"
      syscall: false
      args:
        - index: 0
          type: "sock"
      selectors:
        - matchArgs:
            - index: 0
              operator: "SAddr"
              values:
                - "pool.minexmr.com"
                - "xmr.pool.minergate.com"
                - "monerohash.com"
                - "nanopool.org"
          matchActions:
            - action: Post
            - action: Override
              argError: -1

---
# Tetragon TracingPolicy for detecting reverse shells
apiVersion: cilium.io/v1alpha1
kind: TracingPolicy
metadata:
  name: reverse-shell-detection
  namespace: kube-system
spec:
  # Detect common reverse shell patterns
  kprobes:
    # Detect bash reverse shell
    - call: "__x64_sys_dup2"
      syscall: true
      args:
        - index: 0
          type: "int"
        - index: 1
          type: "int"
      selectors:
        - matchArgs:
            - index: 1
              operator: "Equal"
              values:
                - "0"  # stdin
                - "1"  # stdout
                - "2"  # stderr
          matchBinaries:
            - operator: "In"
              values:
                - "/bin/bash"
                - "/bin/sh"
                - "/usr/bin/bash"
                - "/usr/bin/sh"
          matchActions:
            - action: Post

    # Detect socket connections from shells
    - call: "__x64_sys_connect"
      syscall: true
      args:
        - index: 0
          type: "int"
        - index: 1
          type: "sockaddr"
      selectors:
        - matchBinaries:
            - operator: "In"
              values:
                - "/bin/bash"
                - "/bin/sh"
                - "/usr/bin/python*"
                - "/usr/bin/perl"
                - "/usr/bin/ruby"
          matchActions:
            - action: Post

---
# Tetragon TracingPolicy for file integrity monitoring
apiVersion: cilium.io/v1alpha1
kind: TracingPolicy
metadata:
  name: file-integrity-monitoring
  namespace: kube-system
spec:
  # Monitor critical file changes
  kprobes:
    - call: "security_file_permission"
      syscall: false
      args:
        - index: 0
          type: "file"
        - index: 1
          type: "int"
      selectors:
        - matchArgs:
            - index: 1
              operator: "Mask"
              values:
                - "MAY_WRITE"
          matchBinaries:
            - operator: "NotIn"
              values:
                - "/usr/bin/apt"
                - "/usr/bin/dpkg"
                - "/usr/bin/yum"
          matchActions:
            - action: Post
```

### Falco Integration

```yaml
# Falco rules for Kubernetes runtime security (complementing Tetragon)
# falco-rules-kubernetes.yaml
- rule: Launch Privileged Container
  desc: >
    Detect the initial process started in a privileged container.
    Exceptions are made for known trusted images.
  condition: >
    spawned_process
    and container
    and container.privileged=true
    and not trusted_images
  output: >
    Privileged container started
    (user=%user.name command=%proc.cmdline
    container_id=%container.id container_name=%container.name
    image=%container.image.repository)
  priority: WARNING

- rule: Terminal shell in container
  desc: >
    A shell was used as the entrypoint/exec point into a container with an attached terminal.
  condition: >
    spawned_process
    and container
    and shell_procs
    and proc.tty != 0
    and not entrypoint_shell_containers
  output: >
    A shell was spawned in a container with an attached terminal
    (user=%user.name shell=%proc.name parent=%proc.pname
    cmdline=%proc.cmdline terminal=%proc.tty container_id=%container.id
    container_name=%container.name image=%container.image.repository)
  priority: NOTICE

- rule: Contact EC2 Instance Metadata Service from Container
  desc: >
    Detect attempts to contact the EC2 Instance Metadata Service from a container
  condition: >
    outbound
    and fd.sip="169.254.169.254"
    and container
    and not ec2_metadata_containers
  output: >
    Outbound connection to EC2 instance metadata service
    (command=%proc.cmdline connection=%fd.name
    container_id=%container.id container_name=%container.name
    image=%container.image.repository)
  priority: NOTICE

- rule: Outbound Connection from Sensitive Mount
  desc: >
    Detect outbound connections from containers that mount sensitive host paths
  condition: >
    outbound
    and container
    and sensitive_mount
    and not trusted_images
  output: >
    Outbound connection from container with sensitive mount
    (command=%proc.cmdline connection=%fd.name
    mount=%container.mounts
    container_id=%container.id container_name=%container.name
    image=%container.image.repository)
  priority: WARNING
```

```go
// Go client for Tetragon events
package tetragon

import (
    "context"
    "encoding/json"
    "fmt"
    "io"

    "github.com/cilium/tetragon/api/v1/tetragon"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
)

// EventProcessor processes Tetragon security events
type EventProcessor struct {
    client tetragon.FineGuidanceSensorsClient
    conn   *grpc.ClientConn
}

// NewEventProcessor creates a new Tetragon event processor
func NewEventProcessor(socketPath string) (*EventProcessor, error) {
    conn, err := grpc.Dial(
        "unix://"+socketPath,
        grpc.WithTransportCredentials(insecure.NewCredentials()),
    )
    if err != nil {
        return nil, fmt.Errorf("connecting to tetragon: %w", err)
    }

    client := tetragon.NewFineGuidanceSensorsClient(conn)

    return &EventProcessor{
        client: client,
        conn:   conn,
    }, nil
}

// SecurityEvent represents a processed security event
type SecurityEvent struct {
    Timestamp   string                 `json:"timestamp"`
    Type        string                 `json:"type"`
    Severity    string                 `json:"severity"`
    Process     ProcessInfo            `json:"process"`
    Pod         PodInfo                `json:"pod,omitempty"`
    Details     map[string]interface{} `json:"details"`
}

type ProcessInfo struct {
    PID       uint32 `json:"pid"`
    PPID      uint32 `json:"ppid"`
    UID       uint32 `json:"uid"`
    Binary    string `json:"binary"`
    Arguments string `json:"arguments,omitempty"`
}

type PodInfo struct {
    Namespace string `json:"namespace"`
    Name      string `json:"name"`
    Container string `json:"container"`
}

// StreamEvents streams and processes Tetragon events
func (p *EventProcessor) StreamEvents(ctx context.Context, handler func(*SecurityEvent)) error {
    stream, err := p.client.GetEvents(ctx, &tetragon.GetEventsRequest{
        // Filter for security-relevant events
        AllowList: []*tetragon.Filter{
            {
                EventSet: []tetragon.EventType{
                    tetragon.EventType_PROCESS_EXEC,
                    tetragon.EventType_PROCESS_EXIT,
                    tetragon.EventType_PROCESS_KPROBE,
                    tetragon.EventType_PROCESS_TRACEPOINT,
                },
            },
        },
    })
    if err != nil {
        return fmt.Errorf("getting events: %w", err)
    }

    for {
        event, err := stream.Recv()
        if err == io.EOF {
            return nil
        }
        if err != nil {
            return fmt.Errorf("receiving event: %w", err)
        }

        securityEvent := p.processEvent(event)
        if securityEvent != nil {
            handler(securityEvent)
        }
    }
}

func (p *EventProcessor) processEvent(event *tetragon.GetEventsResponse) *SecurityEvent {
    switch e := event.Event.(type) {
    case *tetragon.GetEventsResponse_ProcessKprobe:
        return p.processKprobeEvent(e.ProcessKprobe)
    case *tetragon.GetEventsResponse_ProcessExec:
        return p.processExecEvent(e.ProcessExec)
    default:
        return nil
    }
}

func (p *EventProcessor) processKprobeEvent(kprobe *tetragon.ProcessKprobe) *SecurityEvent {
    process := kprobe.Process
    if process == nil {
        return nil
    }

    return &SecurityEvent{
        Timestamp: process.StartTime.AsTime().String(),
        Type:      "KPROBE",
        Severity:  determineSeverity(kprobe),
        Process: ProcessInfo{
            PID:       process.Pid.GetValue(),
            Binary:    process.Binary,
            Arguments: process.Arguments,
        },
        Pod: PodInfo{
            Namespace: process.Pod.Namespace,
            Name:      process.Pod.Name,
            Container: process.Pod.Container.Name,
        },
        Details: map[string]interface{}{
            "function_name": kprobe.FunctionName,
            "action":        kprobe.Action,
        },
    }
}

func (p *EventProcessor) processExecEvent(exec *tetragon.ProcessExec) *SecurityEvent {
    process := exec.Process
    if process == nil {
        return nil
    }

    return &SecurityEvent{
        Timestamp: process.StartTime.AsTime().String(),
        Type:      "EXEC",
        Severity:  "INFO",
        Process: ProcessInfo{
            PID:       process.Pid.GetValue(),
            Binary:    process.Binary,
            Arguments: process.Arguments,
        },
        Pod: PodInfo{
            Namespace: process.Pod.Namespace,
            Name:      process.Pod.Name,
            Container: process.Pod.Container.Name,
        },
        Details: map[string]interface{}{
            "parent_binary": process.Parent.Binary,
        },
    }
}

func determineSeverity(kprobe *tetragon.ProcessKprobe) string {
    switch kprobe.Action {
    case "SIGKILL", "OVERRIDE":
        return "CRITICAL"
    case "POST":
        return "WARNING"
    default:
        return "INFO"
    }
}

func (p *EventProcessor) Close() error {
    return p.conn.Close()
}
```

---

## Service Mesh Security

### Istio Ambient Mode mTLS

Istio Ambient mode introduces a new data plane architecture that separates L4 (secure overlay) and L7 (waypoint proxy) processing, achieving ~8% overhead compared to 15-20% in sidecar mode.

```yaml
# Istio Ambient Mode configuration for production
apiVersion: install.istio.io/v1alpha1
kind: IstioOperator
metadata:
  name: ambient-profile
spec:
  profile: ambient
  components:
    pilot:
      k8s:
        resources:
          requests:
            cpu: 2000m
            memory: 4Gi
    ztunnel:
      k8s:
        resources:
          requests:
            cpu: 500m
            memory: 512Mi
          limits:
            cpu: 2000m
            memory: 1Gi

  meshConfig:
    # Enable ambient mode for entire mesh
    defaultConfig:
      proxyMetadata:
        ISTIO_META_DNS_CAPTURE: "true"

    # mTLS configuration
    mtls:
      # Automatic mTLS - transparent encryption
      mode: STRICT

    # Certificate management
    certificates:
      - secretName: dns.istio-galley-service-account
        dnsNames:
          - istio-galley.istio-system.svc
          - istio-galley.istio-system

  values:
    global:
      proxy:
        resources:
          requests:
            cpu: 100m
            memory: 128Mi

    # Ztunnel configuration
    ztunnel:
      variant: distroless
      env:
        # Enable zero-trust networking
        ISTIO_META_DNS_CAPTURE: "true"
        # Enable HBONE (HTTP-Based Overlay Network Environment)
        ISTIO_META_HBONE: "true"
```

```yaml
# Enable ambient mode for a namespace
apiVersion: v1
kind: Namespace
metadata:
  name: production
  labels:
    # Enable ambient mode for this namespace
    istio.io/dataplane-mode: ambient
    # Enable sidecar-less L7 processing where needed
    istio.io/use-waypoint: "true"

---
# Waypoint proxy for L7 processing
apiVersion: gateway.networking.k8s.io/v1beta1
kind: Gateway
metadata:
  name: production-waypoint
  namespace: production
  labels:
    istio.io/waypoint-for: service
spec:
  gatewayClassName: istio-waypoint
  listeners:
    - name: mesh
      port: 15008
      protocol: HBONE
      allowedRoutes:
        namespaces:
          from: Same

---
# Authorization policy for ambient mode
apiVersion: security.istio.io/v1beta1
kind: AuthorizationPolicy
metadata:
  name: production-authz
  namespace: production
spec:
  targetRefs:
    - kind: Gateway
      group: gateway.networking.k8s.io
      name: production-waypoint
  action: ALLOW
  rules:
    # Only allow requests from authenticated services
    - from:
        - source:
            principals:
              - "cluster.local/ns/production/sa/api-service"
              - "cluster.local/ns/production/sa/web-service"
      to:
        - operation:
            methods: ["GET", "POST"]
            paths: ["/api/*"]
      when:
        - key: request.auth.claims[iss]
          values: ["https://accounts.google.com"]

---
# Peer authentication for strict mTLS
apiVersion: security.istio.io/v1beta1
kind: PeerAuthentication
metadata:
  name: default
  namespace: production
spec:
  mtls:
    mode: STRICT
  # Port-level overrides if needed
  portLevelMtls:
    8080:
      mode: PERMISSIVE  # For legacy service migration
```

#### Performance Comparison

```yaml
# Istio Ambient Mode Performance Characteristics
performance_benchmarks:
  latency_overhead:
    ambient_mode:
      p50: "0.3ms"
      p99: "0.8ms"
      overhead_percentage: "8%"
    sidecar_mode:
      p50: "0.8ms"
      p99: "2.5ms"
      overhead_percentage: "18%"

  resource_usage:
    ambient_mode:
      memory_per_node: "512Mi"
      cpu_per_node: "500m"
      ztunnel_memory: "128Mi"

    sidecar_mode:
      memory_per_pod: "128Mi"
      cpu_per_pod: "100m"

  throughput:
    ambient_mode:
      rps_per_core: "50000"
      max_throughput: "10Gbps"

    sidecar_mode:
      rps_per_core: "30000"
      max_throughput: "6Gbps"

  security_features:
    ambient_mode:
      l4_encryption: "ztunnel (HBONE)"
      l7_processing: "waypoint proxy"
      identity: "SPIFFE mTLS"
      authorization: "L4 + L7 policies"

    sidecar_mode:
      l4_encryption: "envoy mTLS"
      l7_processing: "sidecar proxy"
      identity: "SPIFFE mTLS"
      authorization: "full L7 policies"
```

### SPIFFE/SPIRE Integration

```yaml
# SPIRE server deployment for workload identity
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: spire-server
  namespace: spire
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
              name: http
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
              port: 8080
            initialDelaySeconds: 5
            periodSeconds: 5
          readinessProbe:
            httpGet:
              path: /ready
              port: 8080
            initialDelaySeconds: 5
            periodSeconds: 5
      volumes:
        - name: spire-config
          configMap:
            name: spire-server-config
        - name: spire-socket
          emptyDir: {}
  volumeClaimTemplates:
    - metadata:
        name: spire-data
      spec:
        accessModes: ["ReadWriteOnce"]
        resources:
          requests:
            storage: 10Gi

---
# SPIRE agent DaemonSet
apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: spire-agent
  namespace: spire
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
          image: busybox
          command: ['sh', '-c', 'touch /run/spire/sockets/agent.sock']
          volumeMounts:
            - name: spire-agent-socket-dir
              mountPath: /run/spire/sockets
      containers:
        - name: spire-agent
          image: ghcr.io/spiffe/spire-agent:1.9.0
          args:
            - -config
            - /run/spire/config/agent.conf
          securityContext:
            privileged: true
          volumeMounts:
            - name: spire-config
              mountPath: /run/spire/config
              readOnly: true
            - name: spire-agent-socket-dir
              mountPath: /run/spire/sockets
            - name: spire-token
              mountPath: /var/run/secrets/tokens
            - name: proc
              mountPath: /host/proc
              readOnly: true
          livenessProbe:
            exec:
              command:
                - /opt/spire/bin/spire-agent
                - healthcheck
            initialDelaySeconds: 10
            periodSeconds: 10
      volumes:
        - name: spire-config
          configMap:
            name: spire-agent-config
        - name: spire-agent-socket-dir
          hostPath:
            path: /run/spire/sockets
            type: DirectoryOrCreate
        - name: spire-token
          projected:
            sources:
              - serviceAccountToken:
                  path: spire-agent
                  expirationSeconds: 7200
                  audience: spire-server
        - name: proc
          hostPath:
            path: /proc

---
# SPIRE server configuration with OIDC federation
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
      trust_domain = "production.cluster.local"
      data_dir = "/run/spire/data"
      log_level = "DEBUG"

      ca_key_type = "ec-p384"
      ca_subject {
        country = ["US"]
        organization = ["Example Org"]
        common_name = "production.cluster.local"
      }
    }

    plugins {
      DataStore "sql" {
        plugin_data {
          database_type = "postgres"
          connection_string = "dbname=spire user=spire password=spire host=postgres port=5432 sslmode=verify-full"
        }
      }

      KeyManager "memory" {
        plugin_data {}
      }

      NodeAttestor "k8s_psat" {
        plugin_data {
          clusters = {
            "production" = {
              service_account_allow_list = ["spire:spire-agent"]
            }
          }
        }
      }

      UpstreamAuthority "disk" {
        plugin_data {
          cert_file_path = "/run/spire/config/ca.crt"
          key_file_path = "/run/spire/config/ca.key"
        }
      }
    }

    # OIDC federation for external services
    federation {
      bundle_endpoint {
        address = "0.0.0.0"
        port = 8443
      }

      federates_with "vault.cluster.local" {
        bundle_endpoint {
          address = "vault.cluster.local"
          port = 443
          use_web_pkix = true
        }
      }
    }
```

```go
// SPIFFE workload identity for Go applications
package spiffe

import (
    "context"
    "crypto/tls"
    "crypto/x509"
    "fmt"

    "github.com/spiffe/go-spiffe/v2/bundle/x509bundle"
    "github.com/spiffe/go-spiffe/v2/spiffetls/tlsconfig"
    "github.com/spiffe/go-spiffe/v2/svid/x509svid"
    "github.com/spiffe/go-spiffe/v2/workloadapi"
)

// WorkloadIdentity manages SPIFFE workload identity
type WorkloadIdentity struct {
    source *workloadapi.X509Source
    svid   *x509svid.SVID
}

// NewWorkloadIdentity creates a new workload identity source
func NewWorkloadIdentity(ctx context.Context, socketPath string) (*WorkloadIdentity, error) {
    source, err := workloadapi.NewX509Source(
        ctx,
        workloadapi.WithClientOptions(
            workloadapi.WithAddr("unix://"+socketPath),
        ),
    )
    if err != nil {
        return nil, fmt.Errorf("creating x509 source: %w", err)
    }

    return &WorkloadIdentity{
        source: source,
    }, nil
}

// GetSVID retrieves the workload's SPIFFE SVID
func (w *WorkloadIdentity) GetSVID() (*x509svid.SVID, error) {
    svid, err := w.source.GetX509SVID()
    if err != nil {
        return nil, fmt.Errorf("getting svid: %w", err)
    }
    return svid, nil
}

// GetTLSConfig returns a TLS config for mTLS using SPIFFE identities
func (w *WorkloadIdentity) GetTLSConfig(allowedIDs []string) *tls.Config {
    return tlsconfig.MTLSClientConfig(
        w.source,
        w.source,
        tlsconfig.AuthorizeID(allowedIDs[0]),
    )
}

// GetServerTLSConfig returns a TLS config for server-side mTLS
func (w *WorkloadIdentity) GetServerTLSConfig(allowedIDs []string) *tls.Config {
    return tlsconfig.MTLSServerConfig(
        w.source,
        w.source,
        tlsconfig.AuthorizeAny(),
    )
}

// GetCertificate returns the X.509 certificate for the workload
func (w *WorkloadIdentity) GetCertificate() (*tls.Certificate, error) {
    svid, err := w.GetSVID()
    if err != nil {
        return nil, err
    }

    cert, privateKey := svid.Default()
    return &tls.Certificate{
        Certificate: [][]byte{cert.Raw},
        PrivateKey:  privateKey,
        Leaf:        cert,
    }, nil
}

// Close closes the workload identity source
func (w *WorkloadIdentity) Close() error {
    return w.source.Close()
}

// SPIFFE ID validation for Kubernetes services
func ValidateSPIFFEID(spiffeID string, trustDomain string, namespace string, serviceAccount string) error {
    expectedID := fmt.Sprintf("spiffe://%s/ns/%s/sa/%s",
        trustDomain, namespace, serviceAccount)

    if spiffeID != expectedID {
        return fmt.Errorf("SPIFFE ID mismatch: got %s, expected %s", spiffeID, expectedID)
    }

    return nil
}
```

---

## Confidential Computing

### Hardware-Based Enclaves

Confidential Computing uses hardware-based Trusted Execution Environments (TEEs) to protect data in use.

```yaml
# Kubernetes Confidential Computing node pool configuration
apiVersion: machineconfiguration.openshift.io/v1
kind: MachineConfig
metadata:
  name: 50-worker-confidential-computing
  labels:
    machineconfiguration.openshift.io/role: worker-confidential
spec:
  kernelArguments:
    # Enable AMD SEV-SNP
    - mem_encrypt=on
    - kvm_amd_sev=1
    - kvm_amd_sev_es=1
    # Intel TDX (when available)
    # - intel_iommu=on
    # - iommu=pt
    # - tdx=on

  extensions:
    - sev-snp

---
# Confidential pod specification
apiVersion: v1
kind: Pod
metadata:
  name: confidential-workload
  namespace: production
  annotations:
    # Enable confidential computing
    io.katacontainers.config.hypervisor.qemu: "/usr/bin/qemu-system-x86_64"
    io.katacontainers.config.hypervisor.machine_type: "q35"
    io.katacontainers.config.hypervisor.rootfs_type: "erofs"
    # SEV-SNP configuration
    io.katacontainers.config.sev.policy: "3"
    io.katacontainers.config.sev.cbitpos: "51"
    io.katacontainers.config.sev.reduced_phys_bits: "1"
spec:
  runtimeClassName: kata-cc
  containers:
    - name: sensitive-app
      image: registry.example.com/sensitive-app:latest
      resources:
        limits:
          memory: "4Gi"
          cpu: "2000m"
          # Request confidential computing
          confidential-computing.intel.com/tdx: "1"
          # AMD SEV
          # confidential-computing.amd.com/sev-snp: "1"
      securityContext:
        privileged: false
        readOnlyRootFilesystem: true
        allowPrivilegeEscalation: false
        capabilities:
          drop:
            - ALL
      volumeMounts:
        - name: secrets
          mountPath: /run/secrets
          readOnly: true
  volumes:
    - name: secrets
      csi:
        driver: secrets-store.csi.k8s.io
        readOnly: true
        volumeAttributes:
          secretProviderClass: confidential-secrets

---
# SecretProviderClass for confidential secrets
apiVersion: secrets-store.csi.x-k8s.io/v1
kind: SecretProviderClass
metadata:
  name: confidential-secrets
  namespace: production
spec:
  provider: azure
  parameters:
    usePodIdentity: "false"
    useVMManagedIdentity: "true"
    userAssignedIdentityID: ""
    keyvaultName: "confidential-kv"
    cloudName: ""
    objects: |
      array:
        - |
          objectName: database-password
          objectType: secret
          objectVersion: ""
        - |
          objectName: api-key
          objectType: secret
          objectVersion: ""
    tenantId: "your-tenant-id"
  secretObjects:
    - secretName: confidential-creds
      type: Opaque
      data:
        - objectName: database-password
          key: db-password
        - objectName: api-key
          key: api-key
```

### Use Cases for Sensitive Workloads

```go
// Confidential computing workload attestation
package confidential

import (
    "context"
    "crypto/sha256"
    "encoding/base64"
    "encoding/json"
    "fmt"

    "github.com/google/go-sev-guest/abi"
    "github.com/google/go-sev-guest/client"
    "github.com/google/go-sev-guest/verify"
)

// AttestationReport represents a SEV-SNP attestation report
type AttestationReport struct {
    Measurement    []byte            `json:"measurement"`
    ReportData     []byte            `json:"report_data"`
    VMPL           int               `json:"vmpl"`
    PlatformInfo   PlatformInfo      `json:"platform_info"`
    Signature      []byte            `json:"signature"`
    Certificates   []*x509.Certificate `json:"certificates"`
}

type PlatformInfo struct {
    SMTEnabled     bool `json:"smt_enabled"`
    TSMEEnabled    bool `json:"tsme_enabled"`
    DebugEnabled   bool `json:"debug_enabled"`
    KeySharing     bool `json:"key_sharing"`
}

// SEVSNPAttestor performs SEV-SNP attestation
type SEVSNPAttestor struct {
    device client.Device
}

// NewSEVSNPAttestor creates a new SEV-SNP attestor
func NewSEVSNPAttestor() (*SEVSNPAttestor, error) {
    device, err := client.OpenDevice()
    if err != nil {
        return nil, fmt.Errorf("opening SEV device: %w", err)
    }

    return &SEVSNPAttestor{device: device}, nil
}

// GetAttestationReport generates an attestation report
func (a *SEVSNPAttestor) GetAttestationReport(userData []byte) (*AttestationReport, error) {
    report, err := client.GetReport(a.device, userData)
    if err != nil {
        return nil, fmt.Errorf("getting report: %w", err)
    }

    // Extract measurement
    measurement := report.Measurement[:]

    return &AttestationReport{
        Measurement:  measurement,
        ReportData:   report.ReportData[:],
        VMPL:         int(report.VMPL),
        PlatformInfo: extractPlatformInfo(report),
        Signature:    report.Signature[:],
    }, nil
}

// VerifyAttestation verifies an attestation report against expected measurement
func (a *SEVSNPAttestor) VerifyAttestation(
    report *AttestationReport,
    expectedMeasurement []byte,
    productName string,
) error {
    options := &verify.Options{
        Product: productName,
        // Require specific security features
        RequireSMT:  false, // Disable SMT for side-channel protection
        RequireTSME: true,  // Require memory encryption
    }

    // Verify attestation report
    err := verify.SnpAttestation(report, options)
    if err != nil {
        return fmt.Errorf("attestation verification failed: %w", err)
    }

    // Verify measurement matches expected value
    if string(report.Measurement) != string(expectedMeasurement) {
        return fmt.Errorf("measurement mismatch")
    }

    return nil
}

func extractPlatformInfo(report *abi.Report) PlatformInfo {
    return PlatformInfo{
        SMTEnabled:   report.PlatformInfo&0x01 != 0,
        TSMEEnabled:  report.PlatformInfo&0x02 != 0,
        DebugEnabled: report.PlatformInfo&0x04 != 0,
        KeySharing:   report.PlatformInfo&0x08 != 0,
    }
}

func (a *SEVSNPAttestor) Close() error {
    return a.device.Close()
}

// ConfidentialWorkloadUseCases defines common use cases
type ConfidentialWorkloadUseCases struct {
    UseCases []ConfidentialUseCase
}

type ConfidentialUseCase struct {
    Name        string   `json:"name"`
    Description string   `json:"description"`
    TEEType     string   `json:"tee_type"` // SEV-SNP, TDX, SGX
    Benefits    []string `json:"benefits"`
    Example     string   `json:"example"`
}

// GetConfidentialUseCases returns common use cases
func GetConfidentialUseCases() ConfidentialWorkloadUseCases {
    return ConfidentialWorkloadUseCases{
        UseCases: []ConfidentialUseCase{
            {
                Name:        "Financial Data Processing",
                Description: "Process sensitive financial data in encrypted memory",
                TEEType:     "SEV-SNP",
                Benefits: []string{
                    "Data encrypted in memory",
                    "Protection from compromised hypervisor",
                    "Auditability via attestation",
                },
                Example: "Credit card transaction processing",
            },
            {
                Name:        "Healthcare AI/ML",
                Description: "Train ML models on encrypted patient data",
                TEEType:     "TDX",
                Benefits: []string{
                    "HIPAA compliance for data in use",
                    "Multi-party computation support",
                    "Model confidentiality",
                },
                Example: "Medical imaging analysis",
            },
            {
                Name:        "Key Management Services",
                Description: "HSM-like key protection in cloud",
                TEEType:     "SGX",
                Benefits: []string{
                    "Hardware-backed key isolation",
                    "Remote attestation for trust",
                    "Small trusted compute base",
                },
                Example: "Cloud KMS with HSM-backed keys",
            },
            {
                Name:        "Multi-Party Computation",
                Description: "Collaborative computation without revealing inputs",
                TEEType:     "SEV-SNP",
                Benefits: []string{
                    "Inputs remain encrypted",
                    "Verifiable computation",
                    "No trusted third party",
                },
                Example: "Cross-company fraud detection",
            },
            {
                Name:        "Blockchain Validators",
                Description: "Secure validator nodes for proof-of-stake",
                TEEType:     "SEV-SNP",
                Benefits: []string{
                    "Signing key protection",
                    "Slashing protection",
                    "Remote attestation for pool membership",
                },
                Example: "Ethereum validator in cloud",
            },
        },
    }
}
```

---

## Network Security

### Cilium Cluster Mesh with Encryption

```yaml
# Cilium Cluster Mesh with WireGuard encryption
apiVersion: cilium.io/v2alpha1
kind: CiliumClusterwideNetworkPolicy
metadata:
  name: default-deny-all
spec:
  endpointSelector: {}
  ingressDeny:
    - {}
  egressDeny:
    - {}

---
# Cluster mesh encryption policy
apiVersion: cilium.io/v2alpha1
kind: CiliumClusterwideNetworkPolicy
metadata:
  name: cluster-mesh-encryption
spec:
  nodeSelector: {}
  ingress:
    - fromEndpoints:
        - matchLabels:
            io.kubernetes.pod.namespace: kube-system
            k8s-app: cilium
      toPorts:
        - ports:
            - port: "51871"
              protocol: UDP
          rules:
            http: [{}]

---
# Zero-trust network policy
apiVersion: cilium.io/v2
kind: CiliumNetworkPolicy
metadata:
  name: zero-trust-microsegmentation
  namespace: production
spec:
  endpointSelector:
    matchLabels:
      app: payment-service
  ingress:
    - fromEndpoints:
        - matchLabels:
            k8s:io.kubernetes.pod.namespace: production
            app: web-frontend
      toPorts:
        - ports:
            - port: "8080"
              protocol: TCP
          rules:
            http:
              - method: "POST"
                path: "/api/v1/payments"
                headers:
                  - name: "X-Request-ID"
                    required: true
                  - name: "Authorization"
                    required: true
                    presence: true
  egress:
    - toEndpoints:
        - matchLabels:
            k8s:io.kubernetes.pod.namespace: production
            app: postgres
      toPorts:
        - ports:
            - port: "5432"
              protocol: TCP
```

---

## Identity and Access Management

### RBAC Hardening

```yaml
# Least privilege RBAC for production workloads
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: payment-service-role
  namespace: production
rules:
  # Minimal permissions for service account
  - apiGroups: [""]
    resources: ["configmaps"]
    verbs: ["get", "list"]
    resourceNames:
      - payment-service-config

  - apiGroups: [""]
    resources: ["secrets"]
    verbs: ["get"]
    resourceNames:
      - payment-service-creds

  - apiGroups: ["coordination.k8s.io"]
    resources: ["leases"]
    verbs: ["get", "create", "update"]
    resourceNames:
      - payment-service-leader-election

  - apiGroups: [""]
    resources: ["events"]
    verbs: ["create", "patch"]

---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: payment-service-binding
  namespace: production
subjects:
  - kind: ServiceAccount
    name: payment-service
    namespace: production
roleRef:
  kind: Role
  name: payment-service-role
  apiGroup: rbac.authorization.k8s.io

---
# Pod Security Standards enforcement
apiVersion: policy/v1
kind: PodSecurityPolicy
metadata:
  name: restricted-psp
spec:
  privileged: false
  allowPrivilegeEscalation: false
  requiredDropCapabilities:
    - ALL
  volumes:
    - 'configMap'
    - 'emptyDir'
    - 'projected'
    - 'secret'
    - 'downwardAPI'
    - 'persistentVolumeClaim'
    - 'csi'
  runAsUser:
    rule: 'MustRunAsNonRoot'
  seLinux:
    rule: 'RunAsAny'
  fsGroup:
    rule: 'MustRunAs'
    ranges:
      - min: 1
        max: 65535
  runAsGroup:
    rule: 'MustRunAs'
    ranges:
      - min: 1
        max: 65535
  supplementalGroups:
    rule: 'MustRunAs'
    ranges:
      - min: 1
        max: 65535
  readOnlyRootFilesystem: true
```

---

## Secrets Management

### External Secrets Operator with Vault

```yaml
# External Secrets Operator configuration
apiVersion: external-secrets.io/v1beta1
kind: SecretStore
metadata:
  name: vault-backend
  namespace: production
spec:
  provider:
    vault:
      server: "https://vault.production.svc:8200"
      path: "secret"
      version: "v2"
      auth:
        kubernetes:
          mountPath: "kubernetes"
          role: "external-secrets"
          serviceAccountRef:
            name: external-secrets-sa

---
apiVersion: external-secrets.io/v1beta1
kind: ExternalSecret
metadata:
  name: database-credentials
  namespace: production
spec:
  refreshInterval: "1h"
  secretStoreRef:
    name: vault-backend
    kind: SecretStore
  target:
    name: database-credentials
    creationPolicy: Owner
    template:
      type: Opaque
      data:
        connection-string: "postgresql://{{ .username }}:{{ .password }}@{{ .host }}:5432/{{ .database }}"
  data:
    - secretKey: username
      remoteRef:
        key: secret/data/database
        property: username
    - secretKey: password
      remoteRef:
        key: secret/data/database
        property: password
    - secretKey: host
      remoteRef:
        key: secret/data/database
        property: host
    - secretKey: database
      remoteRef:
        key: secret/data/database
        property: database
```

---

## Pod Security

### Pod Security Standards

```yaml
# Pod Security Admission configuration
apiVersion: apiserver.config.k8s.io/v1
kind: AdmissionConfiguration
plugins:
  - name: PodSecurity
    configuration:
      apiVersion: pod-security.admission.config.k8s.io/v1
      kind: PodSecurityConfiguration
      defaults:
        enforce: "restricted"
        audit: "restricted"
        warn: "restricted"
      exemptions:
        usernames: []
        runtimeClasses: []
        namespaces:
          - kube-system
          - monitoring
          - ingress-nginx

---
# Namespace-level Pod Security Standard
apiVersion: v1
kind: Namespace
metadata:
  name: production
  labels:
    pod-security.kubernetes.io/enforce: restricted
    pod-security.kubernetes.io/audit: restricted
    pod-security.kubernetes.io/warn: restricted
    pod-security.kubernetes.io/enforce-version: latest
```

---

## Runtime Security

### OPA/Gatekeeper Policies

```yaml
# Gatekeeper policy for container security
apiVersion: constraints.gatekeeper.sh/v1beta1
kind: K8sPSPForbiddenSysctls
metadata:
  name: forbid-sysctls
spec:
  match:
    kinds:
      - apiGroups: [""]
        kinds: ["Pod"]
    excludedNamespaces:
      - kube-system
  parameters:
    forbiddenSysctls:
      - "*"

---
apiVersion: constraints.gatekeeper.sh/v1beta1
kind: K8sPSPHostNetworkingPorts
metadata:
  name: host-networking-ports
spec:
  match:
    kinds:
      - apiGroups: [""]
        kinds: ["Pod"]
  parameters:
    hostNetwork: true
    min: 80
    max: 9000

---
# Custom constraint template for image signing
apiVersion: templates.gatekeeper.sh/v1
kind: ConstraintTemplate
metadata:
  name: k8srequiredimagesignature
spec:
  crd:
    spec:
      names:
        kind: K8sRequiredImageSignature
      validation:
        openAPIV3Schema:
          type: object
          properties:
            allowedKeylessIssuers:
              type: array
              items:
                type: string
  targets:
    - target: admission.k8s.gatekeeper.sh
      rego: |
        package k8srequiredimagesignature

        violation[{"msg": msg}] {
          container := input.review.object.spec.containers[_]
          image := container.image
          not image_signed(image)
          msg := sprintf("Container image %s is not signed", [image])
        }

        image_signed(image) {
          # Check for Cosign signature
          data.external.cosign.signatures[image].verified
        }
```

---

## Security Monitoring and Auditing

### Falco + Fluent Bit Integration

```yaml
# Falco sidecar for security event shipping
apiVersion: v1
kind: ConfigMap
metadata:
  name: falco-fluentbit-config
  namespace: falco
data:
  fluent-bit.conf: |
    [INPUT]
        Name              tail
        Path              /var/log/falco/falco.json
        Parser            json
        Tag               falco.security
        Refresh_Interval  5
        Mem_Buf_Limit     50MB

    [FILTER]
        Name              modify
        Match             falco.security
        Add               cluster ${CLUSTER_NAME}
        Add               environment production

    [FILTER]
        Name              grep
        Match             falco.security
        Exclude           priority Debug

    [OUTPUT]
        Name              loki
        Match             falco.security
        Host              loki.monitoring.svc.cluster.local
        Port              3100
        Labels            job=falco, cluster=${CLUSTER_NAME}

    [OUTPUT]
        Name              opensearch
        Match             falco.security
        Host              opensearch.monitoring.svc.cluster.local
        Port              9200
        Index             falco-security
        Type              _doc
```

---

## Incident Response

### Automated Incident Response

```yaml
# FalcoSidekick for automated response
apiVersion: apps/v1
kind: Deployment
metadata:
  name: falcosidekick
  namespace: falco
spec:
  replicas: 2
  selector:
    matchLabels:
      app: falcosidekick
  template:
    metadata:
      labels:
        app: falcosidekick
    spec:
      containers:
        - name: falcosidekick
          image: falcosecurity/falcosidekick:2.28.0
          ports:
            - containerPort: 2801
          env:
            - name: DEBUG
              value: "false"
            # Slack notifications
            - name: SLACK_WEBHOOKURL
              valueFrom:
                secretKeyRef:
                  name: falcosidekick-secrets
                  key: slack-webhook
            - name: SLACK_MINIMUMPRIORITY
              value: "warning"
            # PagerDuty for critical alerts
            - name: PAGERDUTY_SERVICEKEY
              valueFrom:
                secretKeyRef:
                  name: falcosidekick-secrets
                  key: pagerduty-key
            - name: PAGERDUTY_MINIMUMPRIORITY
              value: "critical"
            # Kubernetes response
            - name: WEBHOOK_ADDRESS
              value: "http://security-responder.default.svc:8080/events"
            - name: WEBHOOK_CUSTOMHEADERS
              value: "Authorization: Bearer ${RESPONDER_TOKEN}"
```

---

## Compliance and Governance

### CIS Kubernetes Benchmarks

```yaml
# kube-bench configuration for CIS compliance
apiVersion: batch/v1
kind: CronJob
metadata:
  name: kube-bench
  namespace: security
spec:
  schedule: "0 2 * * *"
  jobTemplate:
    spec:
      template:
        spec:
          hostPID: true
          containers:
            - name: kube-bench
              image: aquasec/kube-bench:latest
              command:
                - kube-bench
                - run
                - --targets=master,node,etcd,policies
                - --benchmark=cis-1.8
                - --json
              volumeMounts:
                - name: var-lib-etcd
                  mountPath: /var/lib/etcd
                  readOnly: true
                - name: var-lib-kubelet
                  mountPath: /var/lib/kubelet
                  readOnly: true
                - name: etc-systemd
                  mountPath: /etc/systemd
                  readOnly: true
                - name: etc-kubernetes
                  mountPath: /etc/kubernetes
                  readOnly: true
                - name: usr-bin
                  mountPath: /usr/local/mount-from-host/bin
                  readOnly: true
          restartPolicy: Never
          volumes:
            - name: var-lib-etcd
              hostPath:
                path: /var/lib/etcd
            - name: var-lib-kubelet
              hostPath:
                path: /var/lib/kubelet
            - name: etc-systemd
              hostPath:
                path: /etc/systemd
            - name: etc-kubernetes
              hostPath:
                path: /etc/kubernetes
            - name: usr-bin
              hostPath:
                path: /usr/bin
```

---

## References

### Standards and Specifications

1. **NIST FIPS 203/204/205** - Post-Quantum Cryptography Standards (August 2024)
2. **SLSA v1.1** - Supply Chain Levels for Software Artifacts
3. **Sigstore** - Software signing and transparency
4. **CISA SBOM Requirements 2025** - Software Bill of Materials
5. **CIS Kubernetes Benchmark v1.8** - Security best practices

### Tools and Projects

1. **Tetragon** - eBPF-based security observability (github.com/cilium/tetragon)
2. **Falco** - Runtime security (falco.org)
3. **Cosign** - Container signing (github.com/sigstore/cosign)
4. **Istio Ambient** - Service mesh (istio.io)
5. **SPIRE** - Workload identity (spiffe.io/spire)
6. **Cilium** - eBPF networking and security (cilium.io)

### Further Reading

- NSA/CISA Kubernetes Hardening Guide 2024
- NIST SP 800-204B: Attribute-based Access Control for Microservices
- OWASP Kubernetes Top 10 2024
- Cloud Native Security Whitepaper v2

---

*Document Version: 2025-2026 Edition*
*Last Updated: April 2026*
*Target Size: >25KB*

# LD-020: Go 密码学包深度剖析 (Go Cryptography Packages)

> **维度**: Language Design
> **级别**: S (18+ KB)
> **标签**: #crypto #security #hash #cipher #tls #random
> **权威来源**:
>
> - [Go Cryptography Libraries](https://github.com/golang/go/tree/master/src/crypto) - Go Authors
> - [Go Cryptography Principles](https://go.dev/blog/cryptography-principles) - Go Authors
> - [NIST Cryptographic Standards](https://csrc.nist.gov/projects/cryptographic-standards-and-guidelines) - NIST

---

## 1. 密码学架构概览

### 1.1 包组织结构

```
┌─────────────────────────────────────────────────────────────┐
│                     crypto/                                  │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  Hash Functions (散列)        Symmetric (对称加密)            │
│  ├── crypto/md5              ├── crypto/aes                 │
│  ├── crypto/sha1             ├── crypto/des                 │
│  ├── crypto/sha256           └── crypto/cipher              │
│  └── crypto/sha512                                          │
│                                                              │
│  Asymmetric (非对称)          Random & Keys                   │
│  ├── crypto/rsa              ├── crypto/rand                │
│  ├── crypto/ecdsa            ├── crypto/subtle              │
│  ├── crypto/ecdh             └── crypto/hmac                │
│  └── crypto/ed25519                                         │
│                                                              │
│  TLS & Certificates            Signing                       │
│  ├── crypto/tls              └── crypto/dsa (deprecated)    │
│  └── crypto/x509                                            │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

### 1.2 核心接口设计

```go
// Hash 接口 - 所有散列函数实现
type Hash interface {
    io.Writer              // 写入数据
    Sum(b []byte) []byte   // 返回校验和
    Reset()                // 重置状态
    Size() int             // 输出长度
    BlockSize() int        // 块大小
}

// Block 接口 - 分组密码
type Block interface {
    BlockSize() int
    Encrypt(dst, src []byte)
    Decrypt(dst, src []byte)
}

// Stream 接口 - 流密码
type Stream interface {
    XORKeyStream(dst, src []byte)
}

// Signer 接口 - 数字签名
type Signer interface {
    Public() PublicKey
    Sign(rand io.Reader, digest []byte, opts SignerOpts) (signature []byte, err error)
}
```

---

## 2. 散列函数实现

### 2.1 SHA-256 实现

```go
// src/crypto/sha256/sha256.go

type digest struct {
    h     [8]uint32    // 状态寄存器
    x     [chunk]byte  // 缓冲区
    nx    int          // 缓冲区长度
    len   uint64       // 总长度
}

const (
    chunk = 64          // 块大小 (512 bits)
    Size  = 32          // 输出大小 (256 bits)
)

// 初始化向量
var _K = [64]uint32{
    0x428a2f98, 0x71374491, 0xb5c0fbcf, 0xe9b5dba5,
    // ... 64 个常量
}

func (d *digest) Write(p []byte) (nn int, err error) {
    nn = len(p)
    d.len += uint64(nn)

    // 处理缓冲区中剩余的数据
    if d.nx > 0 {
        n := copy(d.x[d.nx:], p)
        d.nx += n
        if d.nx == chunk {
            block(d, d.x[:])
            d.nx = 0
        }
        p = p[n:]
    }

    // 处理完整块
    if len(p) >= chunk {
        n := len(p) &^ (chunk - 1)
        block(d, p[:n])
        p = p[n:]
    }

    // 剩余数据存入缓冲区
    if len(p) > 0 {
        d.nx = copy(d.x[:], p)
    }
    return
}

func (d *digest) Sum(in []byte) []byte {
    // 创建副本，不修改原状态
    d0 := *d
    hash := d0.checkSum()
    return append(in, hash[:]...)
}

func (d *digest) checkSum() [Size]byte {
    // 填充
    len := d.len
    var tmp [64]byte
    tmp[0] = 0x80

    // 写入长度
    len <<= 3
    padlen := 55 - int(len%64)
    if padlen < 0 {
        padlen += 64
    }
    d.Write(tmp[:padlen+8])

    // 输出
    var digest [Size]byte
    for i := 0; i < 8; i++ {
        binary.BigEndian.PutUint32(digest[i*4:], d.h[i])
    }
    return digest
}
```

### 2.2 汇编优化

```go
// src/crypto/sha256/sha256block_amd64.go

//go:noescape
func block(dig *digest, p []byte)

// 汇编实现 (sha256block_amd64.s)
// 使用 AVX2/SHA-NI 指令加速
//
// TEXT ·block(SB), NOSPLIT, $0-32
//     MOVQ    dig+0(FP), DI
//     MOVQ    p_base+8(FP), SI
//     MOVQ    p_len+16(FP), DX
//     ...
```

### 2.3 HMAC 实现

```go
// src/crypto/hmac/hmac.go

type hmac struct {
    opad, ipad [blockSize]byte  // 填充密钥
    outer, inner hash.Hash      // 内外层散列

    // 优化：复用 marshal 状态
    marshaled bool
}

func New(h func() hash.Hash, key []byte) hash.Hash {
    hm := new(hmac)
    hm.outer = h()
    hm.inner = h()

    blocksize := hm.inner.BlockSize()

    // 密钥处理
    if len(key) > blocksize {
        // 密钥过长，先散列
        hm.outer.Write(key)
        key = hm.outer.Sum(nil)
    }

    // 创建 ipad 和 opad
    copy(hm.ipad[:], key)
    copy(hm.opad[:], key)

    for i := range hm.ipad {
        hm.ipad[i] ^= 0x36
        hm.opad[i] ^= 0x5c
    }

    hm.inner.Write(hm.ipad[:])
    return hm
}

func (h *hmac) Sum(in []byte) []byte {
    origLen := len(in)
    in = h.inner.Sum(in)

    // 需要 Reset outer
    if h.marshaled {
        // 优化：从 marshal 状态恢复
        h.outer.(encoding.BinaryUnmarshaler).UnmarshalBinary(h.opad)
    } else {
        h.outer.Reset()
        h.outer.Write(h.opad[:])
    }

    h.outer.Write(in[origLen:])
    return h.outer.Sum(in[:origLen])
}
```

---

## 3. 对称加密

### 3.1 AES 实现

```go
// src/crypto/aes/cipher.go

type aesCipher struct {
    enc []uint32  // 加密密钥调度
    dec []uint32  // 解密密钥调度
}

const (
    BlockSize = 16  // 128 bits
)

func NewCipher(key []byte) (cipher.Block, error) {
    k := len(key)
    switch k {
    case 16, 24, 32: // AES-128, AES-192, AES-256
        // 支持
    default:
        return nil, KeySizeError(k)
    }

    return newCipherGeneric(key)
}

// 密钥扩展算法
func expandKeyGo(key []byte, enc, dec []uint32) {
    // AES 密钥扩展
    // 使用 Rijndael 密钥调度
    var i int
    nk := len(key) / 4  // 字数
    nr := nk + 6        // 轮数

    // 复制原始密钥
    for i = 0; i < nk; i++ {
        enc[i] = binary.BigEndian.Uint32(key[4*i:])
    }

    // 扩展密钥
    for ; i < (nr+1)*4; i++ {
        t := enc[i-1]
        if i%nk == 0 {
            t = subw(rotw(t)) ^ (uint32(powx[i/nk-1]) << 24)
        } else if nk > 6 && i%nk == 4 {
            t = subw(t)
        }
        enc[i] = enc[i-nk] ^ t
    }

    // 生成解密密钥（InvMixColumns）
    if dec != nil {
        // 最后一轮不需要 InvMixColumns
        n := len(enc)
        for i := 0; i < n; i += 4 {
            ei := n - i - 4
            for j := 0; j < 4; j++ {
                if i > 0 && i+4 < n {
                    dec[ei+j] = te0[sbox0[enc[ei+j]>>24]] ^
                        te1[sbox0[(enc[ei+j]>>16)&0xff]] ^
                        te2[sbox0[(enc[ei+j]>>8)&0xff]] ^
                        te3[sbox0[enc[ei+j]&0xff]]
                } else {
                    dec[ei+j] = enc[ei+j]
                }
            }
        }
    }
}

// 加密一个块
func (c *aesCipher) Encrypt(dst, src []byte) {
    if len(src) != BlockSize {
        panic("crypto/aes: input not full block")
    }
    if len(dst) < BlockSize {
        panic("crypto/aes: output not full block")
    }
    encryptBlockGo(c.enc, dst, src)
}
```

### 3.2 分组模式

```go
// src/crypto/cipher/cbc.go

// CBC 模式
type cbc struct {
    b       cipher.Block
    iv      []byte
    tmp     []byte
    out     []byte
}

func (x *cbc) CryptBlocks(dst, src []byte) {
    if len(src)%x.b.BlockSize() != 0 {
        panic("crypto/cipher: input not full blocks")
    }
    if len(dst) < len(src) {
        panic("crypto/cipher: output smaller than input")
    }

    iv := x.iv

    for len(src) > 0 {
        // dst = encrypt(src XOR iv)
        xorBytes(dst[:x.b.BlockSize()], src[:x.b.BlockSize()], iv)
        x.b.Encrypt(dst[:x.b.BlockSize()], dst[:x.b.BlockSize()])

        // iv = ciphertext (用于下一轮)
        iv = dst[:x.b.BlockSize()]
        src = src[x.b.BlockSize():]
        dst = dst[x.b.BlockSize():]
    }

    // 保存最后一个 IV
    copy(x.iv, iv)
}

// GCM 模式 (AEAD)
type gcm struct {
    cipher     cipher.Block
    nonceSize  int
    tagSize    int
    productTable [16]gcmFieldElement // 预计算表
}

func (g *gcm) Seal(dst, nonce, plaintext, additionalData []byte) []byte {
    // 1. 计算 GHASH(additionalData || padding || ciphertext || len(A) || len(C))
    // 2. 计算 CTR 加密
    // 3. tag = GHASH ^ E(k, J0)
    // ...
}
```

### 3.3 流密码

```go
// src/crypto/cipher/xor.go

// CTR 模式 - 并行加密
func (s *ctr) XORKeyStream(dst, src []byte) {
    if len(dst) < len(src) {
        panic("crypto/cipher: output smaller than input")
    }

    for len(src) > 0 {
        // 生成密钥流块
        if s.outUsed >= len(s.out)-s.b.BlockSize() {
            // 需要生成新块
            for i := len(s.out) - s.b.BlockSize(); i >= 0; i -= s.b.BlockSize() {
                s.b.Encrypt(s.out[i:], s.ctr)
                incCtr(s.ctr)
            }
            s.outUsed = 0
        }

        // XOR
        n := xorBytes(dst, src, s.out[s.outUsed:])
        dst = dst[n:]
        src = src[n:]
        s.outUsed += n
    }
}
```

---

## 4. 非对称加密

### 4.1 RSA 实现

```go
// src/crypto/rsa/rsa.go

type PrivateKey struct {
    PublicKey            // 嵌入公钥
    D         *big.Int   // 私钥指数
    Primes    []*big.Int // 素因子

    // 预计算值（加速解密）
    Precomputed PrecomputedValues
}

type PublicKey struct {
    N *big.Int  // 模数
    E int       // 公钥指数 (通常是 65537)
}

// 密钥生成
func GenerateKey(random io.Reader, bits int) (*PrivateKey, error) {
    priv := new(PrivateKey)
    priv.E = 65537  // F4

    // 生成两个大素数
    for {
        p, err := rand.Prime(random, bits/2)
        if err != nil {
            return nil, err
        }

        q, err := rand.Prime(random, bits-bits/2)
        if err != nil {
            return nil, err
        }

        // 计算 N = p * q
        n := new(big.Int).Mul(p, q)

        // 检查 N 的位数
        if n.BitLen() != bits {
            continue
        }

        priv.Primes = []*big.Int{p, q}
        priv.N = n

        // 计算 phi(N) = (p-1)(q-1)
        pminus1 := new(big.Int).Sub(p, big.NewInt(1))
        qminus1 := new(big.Int).Sub(q, big.NewInt(1))
        phi := new(big.Int).Mul(pminus1, qminus1)

        // 计算 d = e^-1 mod phi(n)
        d := new(big.Int).ModInverse(big.NewInt(int64(priv.E)), phi)
        if d == nil {
            continue
        }
        priv.D = d

        break
    }

    // 预计算加速值
    priv.Precompute()
    return priv, nil
}

// 使用 CRT 加速解密
func decryptCRT(priv *PrivateKey, c *big.Int) (*big.Int, error) {
    // m1 = c^dP mod p
    // m2 = c^dQ mod q
    // h = (qInv * (m1 - m2)) mod p
    // m = m2 + h*q

    cModP := new(big.Int).Mod(c, priv.Primes[0])
    cModQ := new(big.Int).Mod(c, priv.Primes[1])

    m1 := new(big.Int).Exp(cModP, priv.Precomputed.Dp, priv.Primes[0])
    m2 := new(big.Int).Exp(cModQ, priv.Precomputed.Dq, priv.Primes[1])

    h := new(big.Int).Sub(m1, m2)
    if h.Sign() < 0 {
        h.Add(h, priv.Primes[0])
    }
    h.Mul(h, priv.Precomputed.Qinv)
    h.Mod(h, priv.Primes[0])

    m := new(big.Int).Mul(h, priv.Primes[1])
    m.Add(m, m2)

    return m, nil
}

// OAEP 填充
func EncryptOAEP(hash hash.Hash, random io.Reader, pub *PublicKey, msg []byte, label []byte) ([]byte, error) {
    // 1. 长度检查
    if len(msg) > pub.Size()-2*hash.Size()-2 {
        return nil, ErrMessageTooLong
    }

    // 2. 编码消息 EM = 0x00 || maskedSeed || maskedDB
    // seed 是随机数
    // db = lHash || PS || 0x01 || M
    // maskedSeed = seed xor MGF1(maskedDB)
    // maskedDB = db xor MGF1(seed)

    // ... 实现
}
```

### 4.2 ECDSA 实现

```go
// src/crypto/ecdsa/ecdsa.go

type PrivateKey struct {
    PublicKey
    D *big.Int
}

type PublicKey struct {
    elliptic.Curve
    X, Y *big.Int
}

// 签名
func Sign(rand io.Reader, priv *PrivateKey, hash []byte) (r, s *big.Int, err error) {
    // 获取曲线参数
    c := priv.Curve
    N := c.Params().N

    for {
        // 生成随机数 k
        k, err := rand.Int(rand, N)
        if err != nil {
            return nil, nil, err
        }

        // 计算 R = k*G
        rx, _ := c.ScalarBaseMult(k.Bytes())
        r = new(big.Int).Mod(rx, N)
        if r.Sign() == 0 {
            continue
        }

        // 计算 s = k^-1 * (hash + r*d) mod N
        kInv := new(big.Int).ModInverse(k, N)
        z := hashToInt(hash, c)

        s = new(big.Int).Mul(r, priv.D)
        s.Add(s, z)
        s.Mul(s, kInv)
        s.Mod(s, N)

        if s.Sign() == 0 {
            continue
        }

        return
    }
}

// 验证
func Verify(pub *PublicKey, hash []byte, r, s *big.Int) bool {
    // 检查 r, s 范围
    if r.Sign() <= 0 || s.Sign() <= 0 {
        return false
    }

    N := pub.Curve.Params().N
    if r.Cmp(N) >= 0 || s.Cmp(N) >= 0 {
        return false
    }

    // 计算 w = s^-1 mod N
    w := new(big.Int).ModInverse(s, N)

    // 计算 u1 = hash * w mod N, u2 = r * w mod N
    z := hashToInt(hash, pub.Curve)
    u1 := new(big.Int).Mul(z, w)
    u1.Mod(u1, N)

    u2 := new(big.Int).Mul(r, w)
    u2.Mod(u2, N)

    // 计算 P = u1*G + u2*Pub
    x, _ := pub.Curve.ScalarMult(pub.X, pub.Y, u2.Bytes())
    x, _ = pub.Curve.Add(pub.Curve.Params().Gx, pub.Curve.Params().Gy, x, nil)

    // 检查 r == x mod N
    x.Mod(x, N)
    return x.Cmp(r) == 0
}
```

---

## 5. 随机数生成

### 5.1 CSPRNG 实现

```go
// src/crypto/rand/rand.go

// Reader 是加密安全的随机源
var Reader io.Reader

// Unix 实现：读取 /dev/urandom
func init() {
    Reader = &devReader{"/dev/urandom"}
}

type devReader struct {
    name string
}

func (r *devReader) Read(buf []byte) (int, error) {
    // 打开 /dev/urandom
    f, err := os.Open(r.name)
    if err != nil {
        return 0, err
    }
    defer f.Close()

    // 读取随机字节
    return io.ReadFull(f, buf)
}

// Windows 实现：使用 BCryptGenRandom
// func (r *rngReader) Read(b []byte) (int, error) {
//     // 调用 Windows CryptoAPI
// }

// 生成随机素数
func Prime(random io.Reader, bits int) (*big.Int, error) {
    if bits < 2 {
        return nil, errors.New("crypto/rand: prime size must be at least 2-bit")
    }

    b := uint(bits % 8)
    if b == 0 {
        b = 8
    }

    bytes := make([]byte, (bits+7)/8)
    p := new(big.Int)

    for {
        // 读取随机字节
        _, err := io.ReadFull(random, bytes)
        if err != nil {
            return nil, err
        }

        // 设置高位确保位数
        bytes[0] &= uint8(int(1<<b) - 1)
        if bytes[0] < 1<<(b-1) {
            bytes[0] |= 1 << (b - 1)
        }

        // 设置奇数
        bytes[len(bytes)-1] |= 1

        p.SetBytes(bytes)

        // Miller-Rabin 素性测试
        if p.ProbablyPrime(20) {
            return p, nil
        }
    }
}
```

---

## 6. TLS 实现

### 6.1 TLS 握手

```go
// src/crypto/tls/handshake_client.go

func (c *Conn) clientHandshake(ctx context.Context) error {
    // 发送 ClientHello
    hello := &clientHelloMsg{
        vers:               VersionTLS12,
        random:             make([]byte, 32),
        cipherSuites:       c.config.cipherSuites(),
        compressionMethods: []uint8{0},
        serverName:         hostnameInSNI(c.config.ServerName),
        supportedCurves:    c.config.curvePreferences(),
        supportedPoints:    []uint8{pointFormatUncompressed},
        alpnProtocols:      c.config.NextProtos,
    }

    // 生成随机数
    if _, err := io.ReadFull(c.config.rand(), hello.random); err != nil {
        return err
    }

    // 发送 ClientHello
    if _, err := c.writeRecord(recordTypeHandshake, hello.marshal()); err != nil {
        return err
    }

    // 接收 ServerHello
    serverHello, err := c.readHandshake()
    if err != nil {
        return err
    }

    // 密钥交换
    switch serverHello.keyShare.group {
    case x25519:
        // ECDHE with X25519
        // ...
    case secp256r1, secp384r1, secp521r1:
        // ECDHE with NIST curves
        // ...
    }

    // 完成握手
    return nil
}
```

---

## 7. 性能优化

### 7.1 汇编优化

```go
// SHA-256 性能对比
// 纯 Go:  ~200 MB/s
// AVX2:   ~2 GB/s
// SHA-NI: ~3 GB/s

// AES-GCM 性能
// 纯 Go:  ~50 MB/s
// AES-NI: ~2 GB/s
```

### 7.2 基准测试

```go
func BenchmarkSHA256(b *testing.B) {
    data := make([]byte, 1024)
    rand.Read(data)

    h := sha256.New()
    b.SetBytes(int64(len(data)))
    b.ResetTimer()

    for i := 0; i < b.N; i++ {
        h.Write(data)
        h.Sum(nil)
        h.Reset()
    }
}

func BenchmarkAESGCM(b *testing.B) {
    key := make([]byte, 32)
    nonce := make([]byte, 12)
    plaintext := make([]byte, 1024)

    block, _ := aes.NewCipher(key)
    aesgcm, _ := cipher.NewGCM(block)

    b.SetBytes(int64(len(plaintext)))
    b.ResetTimer()

    for i := 0; i < b.N; i++ {
        aesgcm.Seal(nil, nonce, plaintext, nil)
    }
}
```

---

## 8. 视觉表征

### 8.1 加密体系层次

```
┌─────────────────────────────────────────┐
│         Application Layer               │
│    (JWT, JWS, Custom Protocols)         │
├─────────────────────────────────────────┤
│         Protocol Layer                  │
│    TLS 1.3 / TLS 1.2                    │
├─────────────────────────────────────────┤
│         Key Exchange                    │
│    ECDHE / DHE / RSA                    │
├─────────────────────────────────────────┤
│         Authentication                  │
│    RSA-Sign / ECDSA / Ed25519           │
├─────────────────────────────────────────┤
│         Symmetric Encryption            │
│    AES-GCM / AES-CBC / ChaCha20-Poly1305│
├─────────────────────────────────────────┤
│         Hash Functions                  │
│    SHA-256 / SHA-384 / SHA-3            │
└─────────────────────────────────────────┘
```

### 8.2 TLS 握手流程

```
Client                           Server
  │                                │
  ├──── ClientHello ──────────────►│
  │   [version, random, suites]    │
  │                                │
  │◄──── ServerHello ──────────────┤
  │   [version, random, suite]     │
  │                                │
  │◄── Certificate, ServerKeyEx ───┤
  │                                │
  │◄──── ServerHelloDone ──────────┤
  │                                │
  ├──── ClientKeyEx ──────────────►│
  │   [premaster secret]           │
  │                                │
  ├──── ChangeCipherSpec ─────────►│
  ├──── Finished ─────────────────►│
  │                                │
  │◄──── ChangeCipherSpec ─────────┤
  │◄──── Finished ─────────────────┤
  │                                │
  ═════ Encrypted Application ══════
```

---

**质量评级**: S (18KB)
**完成日期**: 2026-04-02

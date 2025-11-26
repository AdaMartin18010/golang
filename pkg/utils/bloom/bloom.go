package bloom

import (
	"hash"
	"hash/fnv"
	"sync"
)

// BloomFilter 布隆过滤器
type BloomFilter struct {
	bits      []bool        // 位数组
	size      uint64        // 位数组大小
	hashFuncs []hash.Hash64 // 哈希函数
	numHashes uint          // 哈希函数数量
	mu        sync.RWMutex  // 互斥锁
}

// NewBloomFilter 创建布隆过滤器
func NewBloomFilter(size uint64, numHashes uint) *BloomFilter {
	if numHashes == 0 {
		numHashes = 3 // 默认3个哈希函数
	}
	
	hashFuncs := make([]hash.Hash64, numHashes)
	for i := uint(0); i < numHashes; i++ {
		hashFuncs[i] = fnv.New64a()
	}
	
	return &BloomFilter{
		bits:      make([]bool, size),
		size:      size,
		hashFuncs: hashFuncs,
		numHashes: numHashes,
	}
}

// Add 添加元素
func (bf *BloomFilter) Add(item []byte) {
	bf.mu.Lock()
	defer bf.mu.Unlock()
	
	for i := uint(0); i < bf.numHashes; i++ {
		hashFunc := fnv.New64a()
		hashFunc.Write(item)
		hashFunc.Write([]byte{byte(i)}) // 添加种子以生成不同的哈希值
		hash := hashFunc.Sum64()
		index := hash % bf.size
		bf.bits[index] = true
	}
}

// AddString 添加字符串元素
func (bf *BloomFilter) AddString(item string) {
	bf.Add([]byte(item))
}

// Contains 检查元素是否存在
func (bf *BloomFilter) Contains(item []byte) bool {
	bf.mu.RLock()
	defer bf.mu.RUnlock()
	
	for i := uint(0); i < bf.numHashes; i++ {
		hashFunc := fnv.New64a()
		hashFunc.Write(item)
		hashFunc.Write([]byte{byte(i)})
		hash := hashFunc.Sum64()
		index := hash % bf.size
		if !bf.bits[index] {
			return false
		}
	}
	return true
}

// ContainsString 检查字符串元素是否存在
func (bf *BloomFilter) ContainsString(item string) bool {
	return bf.Contains([]byte(item))
}

// Clear 清空布隆过滤器
func (bf *BloomFilter) Clear() {
	bf.mu.Lock()
	defer bf.mu.Unlock()
	for i := range bf.bits {
		bf.bits[i] = false
	}
}

// Size 获取位数组大小
func (bf *BloomFilter) Size() uint64 {
	return bf.size
}

// Count 估算元素数量（近似值）
func (bf *BloomFilter) Count() uint64 {
	bf.mu.RLock()
	defer bf.mu.RUnlock()
	
	setBits := uint64(0)
	for _, bit := range bf.bits {
		if bit {
			setBits++
		}
	}
	
	// 使用公式估算：m * ln(1 - X/m) / k
	// 其中 m 是位数组大小，X 是设置的位数，k 是哈希函数数量
	if setBits == 0 {
		return 0
	}
	
	ratio := float64(setBits) / float64(bf.size)
	if ratio >= 1.0 {
		return bf.size // 已满，无法估算
	}
	
	// 简化估算公式
	estimated := -float64(bf.size) * log(1.0-ratio) / float64(bf.numHashes)
	if estimated < 0 {
		return 0
	}
	return uint64(estimated)
}

// FalsePositiveRate 计算假阳性率
func (bf *BloomFilter) FalsePositiveRate(numItems uint64) float64 {
	if numItems == 0 {
		return 0.0
	}
	
	// 公式: (1 - e^(-k*n/m))^k
	// 其中 k 是哈希函数数量，n 是元素数量，m 是位数组大小
	ratio := float64(bf.numHashes) * float64(numItems) / float64(bf.size)
	prob := 1.0 - exp(-ratio)
	return pow(prob, float64(bf.numHashes))
}

// OptimalSize 计算最优位数组大小
func OptimalSize(numItems uint64, falsePositiveRate float64) uint64 {
	if falsePositiveRate <= 0 || falsePositiveRate >= 1 {
		return 0
	}
	// 公式: m = -n * ln(p) / (ln(2)^2)
	// 其中 n 是元素数量，p 是假阳性率
	return uint64(-float64(numItems) * log(falsePositiveRate) / (log(2.0) * log(2.0)))
}

// OptimalHashCount 计算最优哈希函数数量
func OptimalHashCount(numItems, size uint64) uint {
	if size == 0 {
		return 0
	}
	// 公式: k = (m/n) * ln(2)
	// 其中 m 是位数组大小，n 是元素数量
	optimal := float64(size) / float64(numItems) * log(2.0)
	if optimal < 1 {
		return 1
	}
	return uint(optimal)
}

// 辅助函数
func log(x float64) float64 {
	// 简化的自然对数实现
	if x <= 0 {
		return 0
	}
	// 使用泰勒级数近似
	result := 0.0
	term := (x - 1.0) / (x + 1.0)
	term2 := term * term
	power := term
	for i := 1; i < 100; i += 2 {
		result += power / float64(i)
		power *= term2
	}
	return 2.0 * result
}

func exp(x float64) float64 {
	// 简化的指数函数实现
	result := 1.0
	term := 1.0
	for i := 1; i < 100; i++ {
		term *= x / float64(i)
		result += term
		if term < 1e-10 {
			break
		}
	}
	return result
}

func pow(x, y float64) float64 {
	// 简化的幂函数实现
	if y == 0 {
		return 1.0
	}
	if y == 1 {
		return x
	}
	result := 1.0
	for i := 0; i < int(y); i++ {
		result *= x
	}
	return result
}


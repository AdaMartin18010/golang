package ring

import (
	"errors"
	"fmt"
	"hash/fnv"
	"sort"
	"sync"
)

var (
	ErrNodeExists   = errors.New("node already exists")
	ErrNodeNotFound = errors.New("node not found")
	ErrEmptyRing    = errors.New("ring is empty")
)

// Node represents a physical cache node
type Node struct {
	ID              string
	Address         string
	Weight          int
	VirtualHashes   []uint32
}

// Ring implements consistent hashing with virtual nodes
type Ring struct {
	mu       sync.RWMutex
	nodes    map[string]*Node
	ring     map[uint32]string
	hashes   []uint32
	vnodes   int
}

// New creates a new consistent hash ring
func New(vnodes int) *Ring {
	if vnodes <= 0 {
		vnodes = 150
	}
	
	return &Ring{
		nodes:  make(map[string]*Node),
		ring:   make(map[uint32]string),
		vnodes: vnodes,
	}
}

// AddNode adds a node to the ring
func (r *Ring) AddNode(node *Node) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	if _, exists := r.nodes[node.ID]; exists {
		return ErrNodeExists
	}
	
	if node.Weight <= 0 {
		node.Weight = 1
	}
	
	// Calculate virtual node hashes
	vnodes := r.vnodes * node.Weight
	for i := 0; i < vnodes; i++ {
		hash := r.hash(fmt.Sprintf("%s-%d", node.ID, i))
		r.ring[hash] = node.ID
		node.VirtualHashes = append(node.VirtualHashes, hash)
	}
	
	r.nodes[node.ID] = node
	r.sortHashes()
	
	return nil
}

// RemoveNode removes a node from the ring
func (r *Ring) RemoveNode(nodeID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	node, exists := r.nodes[nodeID]
	if !exists {
		return ErrNodeNotFound
	}
	
	// Remove virtual nodes
	for _, hash := range node.VirtualHashes {
		delete(r.ring, hash)
	}
	
	delete(r.nodes, nodeID)
	r.sortHashes()
	
	return nil
}

// GetNode returns the node responsible for a key
func (r *Ring) GetNode(key string) (*Node, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	if len(r.hashes) == 0 {
		return nil, ErrEmptyRing
	}
	
	hash := r.hash(key)
	
	// Binary search for the first hash >= key hash
	idx := sort.Search(len(r.hashes), func(i int) bool {
		return r.hashes[i] >= hash
	})
	
	if idx == len(r.hashes) {
		idx = 0
	}
	
	nodeID := r.ring[r.hashes[idx]]
	return r.nodes[nodeID], nil
}

// GetNodes returns N distinct nodes for replication
func (r *Ring) GetNodes(key string, n int) ([]*Node, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	if len(r.hashes) == 0 {
		return nil, ErrEmptyRing
	}
	
	if n > len(r.nodes) {
		n = len(r.nodes)
	}
	
	hash := r.hash(key)
	idx := sort.Search(len(r.hashes), func(i int) bool {
		return r.hashes[i] >= hash
	})
	
	var nodes []*Node
	seen := make(map[string]bool)
	
	for len(nodes) < n {
		if idx >= len(r.hashes) {
			idx = 0
		}
		
		nodeID := r.ring[r.hashes[idx]]
		if !seen[nodeID] {
			nodes = append(nodes, r.nodes[nodeID])
			seen[nodeID] = true
		}
		idx++
	}
	
	return nodes, nil
}

// GetAllNodes returns all nodes in the ring
func (r *Ring) GetAllNodes() []*Node {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	nodes := make([]*Node, 0, len(r.nodes))
	for _, node := range r.nodes {
		nodes = append(nodes, node)
	}
	return nodes
}

// NodeCount returns the number of physical nodes
func (r *Ring) NodeCount() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.nodes)
}

// sortHashes sorts the hash values for binary search
func (r *Ring) sortHashes() {
	r.hashes = make([]uint32, 0, len(r.ring))
	for hash := range r.ring {
		r.hashes = append(r.hashes, hash)
	}
	sort.Slice(r.hashes, func(i, j int) bool {
		return r.hashes[i] < r.hashes[j]
	})
}

// hash generates a 32-bit hash using FNV-1a
func (r *Ring) hash(key string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(key))
	return h.Sum32()
}

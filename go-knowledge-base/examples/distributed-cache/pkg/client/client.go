package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
	
	"distributed-cache/internal/ring"
)

var (
	ErrKeyNotFound = errors.New("key not found")
	ErrServerError = errors.New("server error")
	ErrTimeout     = errors.New("request timeout")
)

// Config holds client configuration
type Config struct {
	Servers        []string
	Timeout        time.Duration
	MaxRetries     int
	PoolSize       int
	ConsistentHash bool
}

// Client provides a high-level interface to the distributed cache
type Client struct {
	config     *Config
	ring       *ring.Ring
	httpClient *http.Client
	pool       sync.Pool
}

// New creates a new cache client
func New(config *Config) (*Client, error) {
	if len(config.Servers) == 0 {
		return nil, errors.New("no servers configured")
	}
	if config.Timeout <= 0 {
		config.Timeout = 5 * time.Second
	}
	if config.MaxRetries <= 0 {
		config.MaxRetries = 3
	}
	
	client := &Client{
		config: config,
		httpClient: &http.Client{
			Timeout: config.Timeout,
		},
	}
	
	if config.ConsistentHash {
		client.ring = ring.New(150)
		for _, server := range config.Servers {
			node := &ring.Node{
				ID:      server,
				Address: server,
				Weight:  1,
			}
			if err := client.ring.AddNode(node); err != nil && !errors.Is(err, ring.ErrNodeExists) {
				return nil, err
			}
		}
	}
	
	return client, nil
}

// Get retrieves a value from cache
func (c *Client) Get(ctx context.Context, key string) ([]byte, error) {
	server, err := c.getServer(key)
	if err != nil {
		return nil, err
	}
	
	url := fmt.Sprintf("http://%s/cache/%s", server, key)
	
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	switch resp.StatusCode {
	case http.StatusOK:
		return io.ReadAll(resp.Body)
	case http.StatusNotFound:
		return nil, ErrKeyNotFound
	default:
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("%w: %s", ErrServerError, string(body))
	}
}

// Set stores a value in cache
func (c *Client) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	server, err := c.getServer(key)
	if err != nil {
		return err
	}
	
	url := fmt.Sprintf("http://%s/cache/%s", server, key)
	if ttl > 0 {
		url = fmt.Sprintf("%s?ttl=%d", url, int(ttl.Seconds()))
	}
	
	req, err := http.NewRequestWithContext(ctx, "PUT", url, bytes.NewReader(value))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/octet-stream")
	
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("%w: %s", ErrServerError, string(body))
	}
	
	return nil
}

// Delete removes a value from cache
func (c *Client) Delete(ctx context.Context, key string) error {
	server, err := c.getServer(key)
	if err != nil {
		return err
	}
	
	url := fmt.Sprintf("http://%s/cache/%s", server, key)
	
	req, err := http.NewRequestWithContext(ctx, "DELETE", url, nil)
	if err != nil {
		return err
	}
	
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("%w: %s", ErrServerError, string(body))
	}
	
	return nil
}

// MGet retrieves multiple values from cache
func (c *Client) MGet(ctx context.Context, keys []string) (map[string][]byte, error) {
	result := make(map[string][]byte)
	var mu sync.Mutex
	var wg sync.WaitGroup
	
	for _, key := range keys {
		wg.Add(1)
		go func(k string) {
			defer wg.Done()
			
			value, err := c.Get(ctx, k)
			if err != nil {
				return
			}
			
			mu.Lock()
			result[k] = value
			mu.Unlock()
		}(key)
	}
	
	wg.Wait()
	return result, nil
}

// MSet stores multiple values in cache
func (c *Client) MSet(ctx context.Context, items map[string][]byte, ttl time.Duration) error {
	var wg sync.WaitGroup
	errChan := make(chan error, len(items))
	
	for key, value := range items {
		wg.Add(1)
		go func(k string, v []byte) {
			defer wg.Done()
			if err := c.Set(ctx, k, v, ttl); err != nil {
				errChan <- err
			}
		}(key, value)
	}
	
	wg.Wait()
	close(errChan)
	
	for err := range errChan {
		if err != nil {
			return err
		}
	}
	
	return nil
}

// Stats returns cache statistics from a server
func (c *Client) Stats(ctx context.Context) (*Stats, error) {
	if len(c.config.Servers) == 0 {
		return nil, errors.New("no servers configured")
	}
	
	url := fmt.Sprintf("http://%s/stats", c.config.Servers[0])
	
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, ErrServerError
	}
	
	var stats Stats
	if err := json.NewDecoder(resp.Body).Decode(&stats); err != nil {
		return nil, err
	}
	
	return &stats, nil
}

// Close closes the client
func (c *Client) Close() error {
	return nil
}

// getServer returns the server address for a key
func (c *Client) getServer(key string) (string, error) {
	if c.ring != nil {
		node, err := c.ring.GetNode(key)
		if err != nil {
			return "", err
		}
		return node.Address, nil
	}
	
	// Simple round-robin if consistent hashing is disabled
	return c.config.Servers[0], nil
}

// Stats holds cache statistics
type Stats struct {
	Size      int64   `json:"size"`
	MaxSize   int64   `json:"max_size"`
	Items     int64   `json:"items"`
	Hits      int64   `json:"hits"`
	Misses    int64   `json:"misses"`
	HitRatio  float64 `json:"hit_ratio"`
}

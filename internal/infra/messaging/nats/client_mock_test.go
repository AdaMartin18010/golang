// Package nats provides mock implementations for unit testing.
package nats

import (
	"errors"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
)

// mockNatsConn implements natsConn interface for testing
type mockNatsConn struct {
	mu              sync.RWMutex
	connected       bool
	published       []mockPublishedMsg
	subscriptions   map[string][]*nats.Subscription
	stats           nats.Statistics
	publishErr      error
	subscribeErr    error
	queueSubErr     error
	requestErr      error
	requestResponse *nats.Msg
	lastSubject     string
	lastData        []byte
}

type mockPublishedMsg struct {
	Subject string
	Data    []byte
}

func newMockNatsConn() *mockNatsConn {
	return &mockNatsConn{
		connected:     true,
		subscriptions: make(map[string][]*nats.Subscription),
		stats:         nats.Statistics{},
		published:     make([]mockPublishedMsg, 0),
	}
}

func (m *mockNatsConn) Publish(subj string, data []byte) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	m.lastSubject = subj
	m.lastData = data
	
	if m.publishErr != nil {
		return m.publishErr
	}
	
	m.published = append(m.published, mockPublishedMsg{
		Subject: subj,
		Data:    data,
	})
	m.stats.OutMsgs++
	m.stats.OutBytes += uint64(len(data))
	return nil
}

func (m *mockNatsConn) Subscribe(subj string, cb nats.MsgHandler) (*nats.Subscription, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	if m.subscribeErr != nil {
		return nil, m.subscribeErr
	}
	
	sub := &nats.Subscription{
		Subject: subj,
	}
	m.subscriptions[subj] = append(m.subscriptions[subj], sub)
	return sub, nil
}

func (m *mockNatsConn) QueueSubscribe(subj, queue string, cb nats.MsgHandler) (*nats.Subscription, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	if m.queueSubErr != nil {
		return nil, m.queueSubErr
	}
	
	sub := &nats.Subscription{
		Subject: subj,
		Queue:   queue,
	}
	m.subscriptions[subj] = append(m.subscriptions[subj], sub)
	return sub, nil
}

func (m *mockNatsConn) Request(subj string, data []byte, timeout time.Duration) (*nats.Msg, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	if m.requestErr != nil {
		return nil, m.requestErr
	}
	
	if m.requestResponse != nil {
		return m.requestResponse, nil
	}
	
	return &nats.Msg{
		Subject: subj,
		Data:    []byte(`{"status":"ok"}`),
		Header:  nats.Header{},
	}, nil
}

func (m *mockNatsConn) IsConnected() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.connected
}

func (m *mockNatsConn) Close() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.connected = false
}

func (m *mockNatsConn) Stats() nats.Statistics {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.stats
}

func (m *mockNatsConn) ConnectedUrl() string {
	return "nats://mock:4222"
}

// Helper methods
func (m *mockNatsConn) setPublishErr(err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.publishErr = err
}

func (m *mockNatsConn) setSubscribeErr(err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.subscribeErr = err
}

func (m *mockNatsConn) setQueueSubErr(err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.queueSubErr = err
}

func (m *mockNatsConn) setRequestErr(err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.requestErr = err
}

func (m *mockNatsConn) setRequestResponse(msg *nats.Msg) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.requestResponse = msg
}

func (m *mockNatsConn) getPublishedMessages() []mockPublishedMsg {
	m.mu.RLock()
	defer m.mu.RUnlock()
	result := make([]mockPublishedMsg, len(m.published))
	copy(result, m.published)
	return result
}

func (m *mockNatsConn) disconnect() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.connected = false
}

// mockClientWithErrors creates a mock client that returns errors
func mockClientWithErrors(pubErr, subErr, reqErr error) *Client {
	mock := &mockNatsConn{
		connected:     true,
		subscriptions: make(map[string][]*nats.Subscription),
		stats:         nats.Statistics{},
		published:     make([]mockPublishedMsg, 0),
		publishErr:    pubErr,
		subscribeErr:  subErr,
		requestErr:    reqErr,
	}
	return newClientWithConn(mock)
}

// testErrors
var (
	errPublishFailed  = errors.New("publish failed")
	errSubscribeFailed = errors.New("subscribe failed")
	errRequestFailed  = errors.New("request failed")
)

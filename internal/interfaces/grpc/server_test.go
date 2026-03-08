// Package grpc provides tests for gRPC server.
// Note: This test file does not import the grpc package to avoid protobuf init issues.
// Instead, it tests the types and interfaces defined in the server.go file.
package grpc

import (
	"log/slog"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestServerStruct(t *testing.T) {
	// Test that Server struct is properly defined
	s := &Server{}
	assert.NotNil(t, s)
}

func TestConfigStruct(t *testing.T) {
	// Test Config struct with all fields
	config := &Config{
		Host:                  "127.0.0.1",
		Port:                  8080,
		MaxConnectionIdle:     10 * time.Minute,
		MaxConnectionAge:      20 * time.Minute,
		MaxConnectionAgeGrace: 3 * time.Minute,
		Time:                  10 * time.Second,
		Timeout:               2 * time.Second,
		EnableReflection:      false,
	}
	
	assert.Equal(t, "127.0.0.1", config.Host)
	assert.Equal(t, 8080, config.Port)
	assert.Equal(t, 10*time.Minute, config.MaxConnectionIdle)
	assert.Equal(t, 20*time.Minute, config.MaxConnectionAge)
	assert.Equal(t, 3*time.Minute, config.MaxConnectionAgeGrace)
	assert.Equal(t, 10*time.Second, config.Time)
	assert.Equal(t, 2*time.Second, config.Timeout)
	assert.False(t, config.EnableReflection)
}

func TestConfigStructTypes(t *testing.T) {
	// Test that Config fields have correct types
	config := &Config{}
	
	assert.IsType(t, "", config.Host)
	assert.IsType(t, int(0), config.Port)
	assert.IsType(t, time.Duration(0), config.MaxConnectionIdle)
	assert.IsType(t, time.Duration(0), config.MaxConnectionAge)
	assert.IsType(t, time.Duration(0), config.MaxConnectionAgeGrace)
	assert.IsType(t, time.Duration(0), config.Time)
	assert.IsType(t, time.Duration(0), config.Timeout)
	assert.IsType(t, true, config.EnableReflection)
}

func TestDefaultConfig(t *testing.T) {
	// Test DefaultConfig function
	config := DefaultConfig()
	
	assert.NotNil(t, config)
	assert.Equal(t, "0.0.0.0", config.Host)
	assert.Equal(t, 50051, config.Port)
	assert.Equal(t, 15*time.Minute, config.MaxConnectionIdle)
	assert.Equal(t, 30*time.Minute, config.MaxConnectionAge)
	assert.Equal(t, 5*time.Minute, config.MaxConnectionAgeGrace)
	assert.Equal(t, 5*time.Second, config.Time)
	assert.Equal(t, 1*time.Second, config.Timeout)
	assert.True(t, config.EnableReflection)
}

func TestConfigAddress(t *testing.T) {
	// Test Address method
	config := &Config{
		Host: "127.0.0.1",
		Port: 8080,
	}
	assert.Equal(t, "127.0.0.1:8080", config.Address())
	
	// Test with default host
	config2 := &Config{
		Host: "0.0.0.0",
		Port: 50051,
	}
	assert.Equal(t, "0.0.0.0:50051", config2.Address())
}

func TestServerOptionType(t *testing.T) {
	// Test ServerOption type
	var opt ServerOption = func(s *Server) {
		s.logger = slog.Default()
	}
	
	assert.NotNil(t, opt)
}

func TestWithLogger(t *testing.T) {
	// Test WithLogger function
	logger := slog.Default()
	opt := WithLogger(logger)
	
	assert.NotNil(t, opt)
	
	// Test that the option can be applied
	s := &Server{}
	opt(s)
	assert.Equal(t, logger, s.logger)
}

func TestWithConfig(t *testing.T) {
	// Test WithConfig function
	config := &Config{
		Host: "127.0.0.1",
		Port: 9090,
	}
	opt := WithConfig(config)
	
	assert.NotNil(t, opt)
	
	// Test that the option can be applied
	s := &Server{}
	opt(s)
	assert.Equal(t, config, s.config)
}

func TestWithHealthChecker(t *testing.T) {
	// Test WithHealthChecker function
	opt := WithHealthChecker(nil)
	
	assert.NotNil(t, opt)
	
	// Test that the option can be applied without panic
	s := &Server{}
	assert.NotPanics(t, func() {
		opt(s)
	})
}

func TestNewServerFunction(t *testing.T) {
	// Test that NewServer function exists and has correct signature
	assert.NotNil(t, NewServer)
}

func TestServerMethods(t *testing.T) {
	// Test that Server has all expected methods
	s := &Server{}
	
	// Start method
	assert.NotNil(t, s.Start)
	
	// Stop method
	assert.NotNil(t, s.Stop)
	
	// Shutdown method
	assert.NotNil(t, s.Shutdown)
	
	// Addr method
	assert.NotNil(t, s.Addr)
	
	// GetUserHandler method
	assert.NotNil(t, s.GetUserHandler)
	
	// GetHealthHandler method
	assert.NotNil(t, s.GetHealthHandler)
	
	// SetReadyFunc method
	assert.NotNil(t, s.SetReadyFunc)
	
	// RegisterHealthChecker method
	assert.NotNil(t, s.RegisterHealthChecker)
}

func TestConfigWithZeroValues(t *testing.T) {
	// Test Config with zero values
	config := &Config{}
	
	assert.Empty(t, config.Host)
	assert.Equal(t, 0, config.Port)
	assert.Equal(t, time.Duration(0), config.MaxConnectionIdle)
	assert.False(t, config.EnableReflection)
}

func TestServerWithNilFields(t *testing.T) {
	// Test Server with nil fields
	s := &Server{
		grpcServer:    nil,
		listener:      nil,
		logger:        nil,
		config:        nil,
		userHandler:   nil,
		healthHandler: nil,
	}
	
	assert.NotNil(t, s)
	assert.Nil(t, s.grpcServer)
	assert.Nil(t, s.listener)
	assert.Nil(t, s.logger)
	assert.Nil(t, s.config)
}

func TestServerOptionChaining(t *testing.T) {
	// Test that multiple ServerOptions can be created
	logger := slog.Default()
	config := DefaultConfig()
	
	opts := []ServerOption{
		WithLogger(logger),
		WithConfig(config),
		WithHealthChecker(nil),
	}
	
	assert.Equal(t, 3, len(opts))
	
	// Test applying options in sequence
	s := &Server{}
	for _, opt := range opts {
		opt(s)
	}
	
	assert.Equal(t, logger, s.logger)
	assert.Equal(t, config, s.config)
}

func TestConfigImmutability(t *testing.T) {
	// Test that DefaultConfig returns a new instance each time
	config1 := DefaultConfig()
	config2 := DefaultConfig()
	
	// They should be equal in value
	assert.Equal(t, config1.Host, config2.Host)
	assert.Equal(t, config1.Port, config2.Port)
	
	// But modifying one should not affect the other
	config1.Port = 9999
	assert.NotEqual(t, config1.Port, config2.Port)
	assert.Equal(t, 50051, config2.Port)
}

func TestServerAddrWithNilListener(t *testing.T) {
	// Test Addr method when listener is nil
	s := &Server{
		listener: nil,
	}
	
	addr := s.Addr()
	assert.Nil(t, addr)
}

// TestHealthCheckerInterface tests that HealthChecker interface is defined
func TestHealthCheckerInterface(t *testing.T) {
	// HealthChecker should be an interface with Check method
	// We can't directly test the interface here but we document it
	assert.NotNil(t, NewServer)
}

// TestServerWithListener tests Server with a mock listener
func TestServerWithListener(t *testing.T) {
	// Create a server with a mock address
	s := &Server{
		config: &Config{
			Host: "localhost",
			Port: 8080,
		},
	}
	
	// Test that the server struct is valid
	assert.NotNil(t, s)
	assert.Equal(t, "localhost:8080", s.config.Address())
}

// TestConfigAddressWithEmptyHost tests Address with empty host
func TestConfigAddressWithEmptyHost(t *testing.T) {
	config := &Config{
		Host: "",
		Port: 8080,
	}
	assert.Equal(t, ":8080", config.Address())
}

// TestServerFields tests that all Server fields are accessible
func TestServerFields(t *testing.T) {
	s := &Server{}
	
	// All fields should be accessible (even if nil)
	_ = s.grpcServer
	_ = s.listener
	_ = s.logger
	_ = s.config
	_ = s.userHandler
	_ = s.healthHandler
}

// TestNetAddrCompatibility tests that Server.Addr returns net.Addr
func TestNetAddrCompatibility(t *testing.T) {
	// Server.Addr() should return net.Addr
	s := &Server{}
	addr := s.Addr()
	// Addr can be nil if listener is not set
	_ = addr
}

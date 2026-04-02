package server

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"
	
	"github.com/gorilla/mux"
	
	"distributed-cache/internal/cache"
	"distributed-cache/internal/ring"
)

// Server represents the HTTP server
type Server struct {
	cache *cache.Cache
	ring  *ring.Ring
	nodeID string
}

// New creates a new server
func New(cache *cache.Cache, ring *ring.Ring, nodeID string) *Server {
	return &Server{
		cache:  cache,
		ring:   ring,
		nodeID: nodeID,
	}
}

// Router returns the HTTP router
func (s *Server) Router() http.Handler {
	r := mux.NewRouter()
	
	// Health checks
	r.HandleFunc("/health", s.handleHealth).Methods("GET")
	r.HandleFunc("/ready", s.handleReady).Methods("GET")
	
	// Cache operations
	r.HandleFunc("/cache/{key}", s.handleGet).Methods("GET")
	r.HandleFunc("/cache/{key}", s.handleSet).Methods("PUT", "POST")
	r.HandleFunc("/cache/{key}", s.handleDelete).Methods("DELETE")
	
	// Stats
	r.HandleFunc("/stats", s.handleStats).Methods("GET")
	r.HandleFunc("/stats/reset", s.handleResetStats).Methods("POST")
	
	// Cluster operations
	r.HandleFunc("/cluster/nodes", s.handleListNodes).Methods("GET")
	r.HandleFunc("/cluster/join", s.handleJoin).Methods("POST")
	r.HandleFunc("/cluster/leave", s.handleLeave).Methods("POST")
	
	return r
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{
		"status": "healthy",
		"node_id": s.nodeID,
	})
}

func (s *Server) handleReady(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ready",
		"node_id": s.nodeID,
	})
}

func (s *Server) handleGet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	
	value, err := s.cache.Get(key)
	if err != nil {
		if err == cache.ErrKeyNotFound {
			http.Error(w, "Key not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(value)
}

func (s *Server) handleSet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	
	value, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	
	// Parse TTL from query parameter
	var ttl time.Duration
	if ttlStr := r.URL.Query().Get("ttl"); ttlStr != "" {
		if ttlSecs, err := strconv.Atoi(ttlStr); err == nil {
			ttl = time.Duration(ttlSecs) * time.Second
		}
	}
	
	if err := s.cache.Set(key, value, ttl); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
		"key":    key,
	})
}

func (s *Server) handleDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	
	if s.cache.Delete(key) {
		w.WriteHeader(http.StatusNoContent)
	} else {
		http.Error(w, "Key not found", http.StatusNotFound)
	}
}

func (s *Server) handleStats(w http.ResponseWriter, r *http.Request) {
	stats := s.cache.Stats()
	
	response := map[string]interface{}{
		"node_id":    s.nodeID,
		"size":       stats.Size,
		"max_size":   stats.MaxSize,
		"items":      stats.Items,
		"hits":       stats.Hits,
		"misses":     stats.Misses,
		"hit_ratio":  stats.HitRatio(),
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (s *Server) handleResetStats(w http.ResponseWriter, r *http.Request) {
	// In a real implementation, this would reset the stats
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (s *Server) handleListNodes(w http.ResponseWriter, r *http.Request) {
	nodes := s.ring.GetAllNodes()
	
	response := make([]map[string]string, 0, len(nodes))
	for _, node := range nodes {
		response = append(response, map[string]string{
			"id":      node.ID,
			"address": node.Address,
		})
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (s *Server) handleJoin(w http.ResponseWriter, r *http.Request) {
	// In a real implementation, this would join the cluster
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "joined"})
}

func (s *Server) handleLeave(w http.ResponseWriter, r *http.Request) {
	// In a real implementation, this would leave the cluster
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "left"})
}

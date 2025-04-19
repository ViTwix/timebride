package loadbalancer

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"time"
)

// LoadBalancer interface defines the methods that all load balancer implementations must provide
type LoadBalancer interface {
	// AddServer adds a server to the load balancer
	AddServer(server *url.URL) error

	// RemoveServer removes a server from the load balancer
	RemoveServer(server *url.URL) error

	// GetNextServer returns the next server to handle a request
	GetNextServer() (*url.URL, error)

	// GetServers returns all servers in the load balancer
	GetServers() []*url.URL

	// HealthCheck performs a health check on all servers
	HealthCheck() map[*url.URL]bool

	// ServeHTTP handles HTTP requests
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

// BaseLoadBalancer provides common functionality for all load balancer implementations
type BaseLoadBalancer struct {
	servers []*url.URL
	mu      sync.RWMutex
	health  map[*url.URL]bool
	client  *http.Client
}

// NewBaseLoadBalancer creates a new base load balancer
func NewBaseLoadBalancer() *BaseLoadBalancer {
	return &BaseLoadBalancer{
		servers: make([]*url.URL, 0),
		health:  make(map[*url.URL]bool),
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

// AddServer adds a server to the load balancer
func (lb *BaseLoadBalancer) AddServer(server *url.URL) error {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	// Check if server already exists
	for _, s := range lb.servers {
		if s.String() == server.String() {
			return fmt.Errorf("server already exists: %s", server.String())
		}
	}

	lb.servers = append(lb.servers, server)
	lb.health[server] = true

	return nil
}

// RemoveServer removes a server from the load balancer
func (lb *BaseLoadBalancer) RemoveServer(server *url.URL) error {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	for i, s := range lb.servers {
		if s.String() == server.String() {
			lb.servers = append(lb.servers[:i], lb.servers[i+1:]...)
			delete(lb.health, server)
			return nil
		}
	}

	return fmt.Errorf("server not found: %s", server.String())
}

// GetServers returns all servers in the load balancer
func (lb *BaseLoadBalancer) GetServers() []*url.URL {
	lb.mu.RLock()
	defer lb.mu.RUnlock()

	servers := make([]*url.URL, len(lb.servers))
	copy(servers, lb.servers)

	return servers
}

// HealthCheck performs a health check on all servers
func (lb *BaseLoadBalancer) HealthCheck() map[*url.URL]bool {
	lb.mu.RLock()
	servers := make([]*url.URL, len(lb.servers))
	copy(servers, lb.servers)
	lb.mu.RUnlock()

	results := make(map[*url.URL]bool)

	for _, server := range servers {
		health := lb.checkServerHealth(server)
		results[server] = health

		lb.mu.Lock()
		lb.health[server] = health
		lb.mu.Unlock()
	}

	return results
}

// checkServerHealth checks if a server is healthy
func (lb *BaseLoadBalancer) checkServerHealth(server *url.URL) bool {
	healthURL := fmt.Sprintf("%s/health", server.String())

	resp, err := lb.client.Get(healthURL)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}

// RoundRobinLoadBalancer implements the LoadBalancer interface using a round-robin strategy
type RoundRobinLoadBalancer struct {
	*BaseLoadBalancer
	current int
}

// NewRoundRobinLoadBalancer creates a new round-robin load balancer
func NewRoundRobinLoadBalancer() *RoundRobinLoadBalancer {
	return &RoundRobinLoadBalancer{
		BaseLoadBalancer: NewBaseLoadBalancer(),
		current:          0,
	}
}

// GetNextServer returns the next server using round-robin strategy
func (lb *RoundRobinLoadBalancer) GetNextServer() (*url.URL, error) {
	lb.mu.RLock()
	defer lb.mu.RUnlock()

	if len(lb.servers) == 0 {
		return nil, fmt.Errorf("no servers available")
	}

	// Find the next healthy server
	startIndex := lb.current
	for i := 0; i < len(lb.servers); i++ {
		index := (startIndex + i) % len(lb.servers)
		server := lb.servers[index]

		if lb.health[server] {
			lb.current = (index + 1) % len(lb.servers)
			return server, nil
		}
	}

	return nil, fmt.Errorf("no healthy servers available")
}

// ServeHTTP handles HTTP requests using round-robin load balancing
func (lb *RoundRobinLoadBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	server, err := lb.GetNextServer()
	if err != nil {
		http.Error(w, "Service Unavailable", http.StatusServiceUnavailable)
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(server)
	proxy.ServeHTTP(w, r)
}

// LeastConnectionsLoadBalancer implements the LoadBalancer interface using a least connections strategy
type LeastConnectionsLoadBalancer struct {
	*BaseLoadBalancer
	connections map[*url.URL]int
}

// NewLeastConnectionsLoadBalancer creates a new least connections load balancer
func NewLeastConnectionsLoadBalancer() *LeastConnectionsLoadBalancer {
	return &LeastConnectionsLoadBalancer{
		BaseLoadBalancer: NewBaseLoadBalancer(),
		connections:      make(map[*url.URL]int),
	}
}

// GetNextServer returns the next server using least connections strategy
func (lb *LeastConnectionsLoadBalancer) GetNextServer() (*url.URL, error) {
	lb.mu.RLock()
	defer lb.mu.RUnlock()

	if len(lb.servers) == 0 {
		return nil, fmt.Errorf("no servers available")
	}

	var selectedServer *url.URL
	minConnections := -1

	for _, server := range lb.servers {
		if !lb.health[server] {
			continue
		}

		connections := lb.connections[server]
		if minConnections == -1 || connections < minConnections {
			minConnections = connections
			selectedServer = server
		}
	}

	if selectedServer == nil {
		return nil, fmt.Errorf("no healthy servers available")
	}

	return selectedServer, nil
}

// IncrementConnections increments the connection count for a server
func (lb *LeastConnectionsLoadBalancer) IncrementConnections(server *url.URL) {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	lb.connections[server]++
}

// DecrementConnections decrements the connection count for a server
func (lb *LeastConnectionsLoadBalancer) DecrementConnections(server *url.URL) {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	if lb.connections[server] > 0 {
		lb.connections[server]--
	}
}

// ServeHTTP handles HTTP requests using least connections load balancing
func (lb *LeastConnectionsLoadBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	server, err := lb.GetNextServer()
	if err != nil {
		http.Error(w, "Service Unavailable", http.StatusServiceUnavailable)
		return
	}

	lb.IncrementConnections(server)

	// Create a custom response writer to track when the request is complete
	rw := &responseWriter{
		ResponseWriter: w,
		onClose: func() {
			lb.DecrementConnections(server)
		},
	}

	proxy := httputil.NewSingleHostReverseProxy(server)
	proxy.ServeHTTP(rw, r)
}

// responseWriter is a custom response writer that calls onClose when the response is complete
type responseWriter struct {
	http.ResponseWriter
	onClose func()
}

// WriteHeader overrides the WriteHeader method to call onClose when the response is complete
func (rw *responseWriter) WriteHeader(code int) {
	rw.ResponseWriter.WriteHeader(code)
}

// Write overrides the Write method to call onClose when the response is complete
func (rw *responseWriter) Write(b []byte) (int, error) {
	n, err := rw.ResponseWriter.Write(b)
	if err == nil {
		rw.onClose()
	}
	return n, err
}

// IPHashLoadBalancer implements the LoadBalancer interface using an IP hash strategy
type IPHashLoadBalancer struct {
	*BaseLoadBalancer
}

// NewIPHashLoadBalancer creates a new IP hash load balancer
func NewIPHashLoadBalancer() *IPHashLoadBalancer {
	return &IPHashLoadBalancer{
		BaseLoadBalancer: NewBaseLoadBalancer(),
	}
}

// GetNextServer returns the next server using IP hash strategy
func (lb *IPHashLoadBalancer) GetNextServer() (*url.URL, error) {
	lb.mu.RLock()
	defer lb.mu.RUnlock()

	if len(lb.servers) == 0 {
		return nil, fmt.Errorf("no servers available")
	}

	// This is a placeholder. In a real implementation, you would hash the client IP
	// and use that to select a server. For simplicity, we'll just use the first server.
	return lb.servers[0], nil
}

// ServeHTTP handles HTTP requests using IP hash load balancing
func (lb *IPHashLoadBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	server, err := lb.GetNextServer()
	if err != nil {
		http.Error(w, "Service Unavailable", http.StatusServiceUnavailable)
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(server)
	proxy.ServeHTTP(w, r)
}

// LoadBalancerFactory creates a load balancer based on the provided configuration
func LoadBalancerFactory(balancerType string) (LoadBalancer, error) {
	switch balancerType {
	case "round_robin":
		return NewRoundRobinLoadBalancer(), nil
	case "least_connections":
		return NewLeastConnectionsLoadBalancer(), nil
	case "ip_hash":
		return NewIPHashLoadBalancer(), nil
	default:
		return nil, fmt.Errorf("unsupported load balancer type: %s", balancerType)
	}
}

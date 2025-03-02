package proxy

import (
	"fmt"
	"io"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/BeautyyuYanli/tcp_multi/pkg/config"
	"github.com/BeautyyuYanli/tcp_multi/pkg/logger"
)

// Proxy represents a TCP proxy
type Proxy struct {
	config             *config.Config
	logger             *logger.Logger
	listener           net.Listener
	backends           []string
	nextBackend        uint64
	backendConnections map[string]int
	backendMutex       sync.Mutex
}

// New creates a new TCP proxy
func New(cfg *config.Config, log *logger.Logger) *Proxy {
	return &Proxy{
		config:             cfg,
		logger:             log,
		backends:           cfg.Proxy.Backends,
		backendConnections: make(map[string]int),
	}
}

// Start starts the TCP proxy
func (p *Proxy) Start() error {
	addr := fmt.Sprintf("%s:%d", p.config.Server.Host, p.config.Server.Port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to start listener: %w", err)
	}
	p.listener = listener
	p.logger.Info("TCP proxy started on %s", addr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			p.logger.Error("Failed to accept connection: %v", err)
			continue
		}

		if p.config.Connection.KeepAlive {
			tcpConn, ok := conn.(*net.TCPConn)
			if ok {
				tcpConn.SetKeepAlive(true)
				tcpConn.SetKeepAlivePeriod(p.config.Connection.KeepAlivePeriod)
			}
		}

		go p.handleConnection(conn)
	}
}

// Stop stops the TCP proxy
func (p *Proxy) Stop() error {
	if p.listener != nil {
		return p.listener.Close()
	}
	return nil
}

// handleConnection handles a client connection
func (p *Proxy) handleConnection(clientConn net.Conn) {
	defer clientConn.Close()

	// Try each backend until one succeeds or all fail
	backendConn, backend, err := p.connectToBackend()
	if err != nil {
		p.logger.Error("Failed to connect to any backend: %v", err)
		return
	}
	defer func() {
		backendConn.Close()
		p.decrementBackendConnections(backend)
	}()

	p.logger.Debug("Connected to backend %s", backend)

	// Set connection parameters
	if tcpConn, ok := backendConn.(*net.TCPConn); ok {
		tcpConn.SetKeepAlive(p.config.Connection.KeepAlive)
		if p.config.Connection.KeepAlive {
			tcpConn.SetKeepAlivePeriod(p.config.Connection.KeepAlivePeriod)
		}
	}

	// Create a wait group to wait for both goroutines to finish
	var wg sync.WaitGroup
	wg.Add(2)

	// Copy data from client to backend
	go func() {
		defer wg.Done()
		buf := make([]byte, p.config.Connection.BufferSize)
		if _, err := io.CopyBuffer(backendConn, clientConn, buf); err != nil {
			if err != io.EOF {
				p.logger.Debug("Error copying client -> backend: %v", err)
			}
		}
		// Signal the backend connection to close after writing is done
		backendConn.(*net.TCPConn).CloseWrite()
	}()

	// Copy data from backend to client
	go func() {
		defer wg.Done()
		buf := make([]byte, p.config.Connection.BufferSize)
		if _, err := io.CopyBuffer(clientConn, backendConn, buf); err != nil {
			if err != io.EOF {
				p.logger.Debug("Error copying backend -> client: %v", err)
			}
		}
		// Signal the client connection to close after writing is done
		clientConn.(*net.TCPConn).CloseWrite()
	}()

	// Wait for both copy operations to complete
	wg.Wait()
	p.logger.Debug("Connection closed")
}

// connectToBackend attempts to connect to backends in sequence until one succeeds
func (p *Proxy) connectToBackend() (net.Conn, string, error) {
	if len(p.backends) == 0 {
		return nil, "", fmt.Errorf("no backends available")
	}

	// Get the starting backend index
	startIdx := int(atomic.AddUint64(&p.nextBackend, 1) - 1) % len(p.backends)
	
	// Try each backend starting from the next one
	for i := 0; i < len(p.backends); i++ {
		idx := (startIdx + i) % len(p.backends)
		backend := p.backends[idx]
		
		// Check if this backend has reached the maximum connections limit
		if !p.canConnectToBackend(backend) {
			continue
		}
		
		p.logger.Debug("Attempting to connect to backend %s", backend)
		backendConn, err := net.DialTimeout("tcp", backend, 1*time.Second)
		if err != nil {
			p.logger.Error("Failed to connect to backend %s: %v", backend, err)
			continue
		}
		
		p.incrementBackendConnections(backend)
		return backendConn, backend, nil
	}
	
	return nil, "", fmt.Errorf("all backends are unavailable or at maximum capacity")
}

// canConnectToBackend checks if a backend has capacity for more connections
func (p *Proxy) canConnectToBackend(backend string) bool {
	p.backendMutex.Lock()
	defer p.backendMutex.Unlock()
	
	// Remove the maximum connection check since MaxConnsPerBackend is not defined
	// Always return true to allow connections
	return true
}

// incrementBackendConnections increments the connection count for a backend
func (p *Proxy) incrementBackendConnections(backend string) {
	p.backendMutex.Lock()
	defer p.backendMutex.Unlock()
	
	p.backendConnections[backend]++
	p.logger.Debug("New connection to backend %s, active connections: %d", 
		backend, p.backendConnections[backend])
}

// decrementBackendConnections decrements the connection count for a backend
func (p *Proxy) decrementBackendConnections(backend string) {
	p.backendMutex.Lock()
	defer p.backendMutex.Unlock()
	
	p.backendConnections[backend]--
	p.logger.Debug("Closed connection to backend %s, active connections: %d", 
		backend, p.backendConnections[backend])
} 
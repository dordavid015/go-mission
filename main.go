package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

// getClientIP extracts the real client IP from the request
func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header first (for proxies/load balancers)
	xForwardedFor := r.Header.Get("X-Forwarded-For")
	if xForwardedFor != "" {
		// Take the first IP if there are multiple
		ips := strings.Split(xForwardedFor, ",")
		return strings.TrimSpace(ips[0])
	}

	// Check X-Real-IP header
	xRealIP := r.Header.Get("X-Real-IP")
	if xRealIP != "" {
		return xRealIP
	}

	// Fall back to RemoteAddr
	ip := r.RemoteAddr
	// Remove port if present
	if strings.Contains(ip, ":") {
		ip = strings.Split(ip, ":")[0]
	}
	
	return ip
}

// handler function that logs client IP and responds
func handler(w http.ResponseWriter, r *http.Request) {
	clientIP := getClientIP(r)
	
	// Log the request with client IP
	log.Printf("Request from client IP: %s - Method: %s - Path: %s - User-Agent: %s", 
		clientIP, r.Method, r.URL.Path, r.Header.Get("User-Agent"))
	
	// Send response
	response := fmt.Sprintf("Hello! Your IP address is: %s\n", clientIP)
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
}

// health check endpoint
func healthHandler(w http.ResponseWriter, r *http.Request) {
	clientIP := getClientIP(r)
	log.Printf("Health check from client IP: %s", clientIP)
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "healthy"}`))
}

func main() {
	// Set up routes
	http.HandleFunc("/", handler)
	http.HandleFunc("/health", healthHandler)
	
	port := ":8081"
	log.Printf("Starting server on port %s", port)
	log.Printf("Server ready to accept connections...")
	
	// Start the server
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

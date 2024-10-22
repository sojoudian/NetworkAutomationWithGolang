package main

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"time"
)

type Pinger struct {
	IPAddress string
}

func NewPinger(ipAddress string) *Pinger {
	if ipAddress == "" {
		ipAddress = "8.8.4.4" // Default IP if none provided
	}
	return &Pinger{IPAddress: ipAddress}
}

func (p *Pinger) ping() (string, error) {
	// Create context with timeout for better memory safety and security
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Prepare the command based on OS (example assumes Unix-like systems)
	// Replace `-c` with `-n` for Windows systems.
	cmd := exec.CommandContext(ctx, "ping", "-c", "4", p.IPAddress)

	// Capture the output
	var out bytes.Buffer
	cmd.Stdout = &out

	// Execute the command
	if err := cmd.Run(); err != nil {
		return "", err
	}

	return out.String(), nil
}

func pingConcurrently(ips []string) {
	results := make(chan string, len(ips)) // Buffer for results

	// Concurrently ping each IP using goroutines
	for _, ip := range ips {
		go func(ip string) {
			pinger := NewPinger(ip)
			result, err := pinger.ping()
			if err != nil {
				results <- fmt.Sprintf("Ping failed for %s: %v", ip, err)
			} else {
				results <- fmt.Sprintf("Ping result for %s:\n%s", ip, result)
			}
		}(ip)
	}

	// Collect and print results
	for range ips {
		fmt.Println(<-results)
	}
}

func main() {
	// List of IPs to ping
	ips := []string{"8.8.4.4", "8.8.8.8"}

	// Call the concurrent ping function
	pingConcurrently(ips)
}

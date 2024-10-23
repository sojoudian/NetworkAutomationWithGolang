package main

import (
	"fmt"
	"io"
	"log"
	"time"

	"golang.org/x/crypto/ssh"
)

func main() {
	// SSH client configuration
	config := &ssh.ClientConfig{
		User: "your_username",
		Auth: []ssh.AuthMethod{
			ssh.Password("your_password"),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // Automatically accept any host key (similar to AutoAddPolicy)
	}

	// Connect to the remote server
	conn, err := ssh.Dial("tcp", "your_remote_host:22", config)
	if err != nil {
		log.Fatalf("Failed to dial: %s", err)
	}
	defer conn.Close()

	// Start an interactive session
	session, err := conn.NewSession()
	if err != nil {
		log.Fatalf("Failed to create session: %s", err)
	}
	defer session.Close()

	// Allocate a terminal for the session
	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // Disable echoing
		ssh.TTY_OP_ISPEED: 14400, // Set input speed (baud rate)
		ssh.TTY_OP_OSPEED: 14400, // Set output speed (baud rate)
	}

	// Request a pseudo-terminal
	if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
		log.Fatalf("Request for pseudo terminal failed: %s", err)
	}

	// Create pipes for input and output
	stdin, err := session.StdinPipe()
	if err != nil {
		log.Fatalf("Unable to setup stdin: %s", err)
	}

	stdout, err := session.StdoutPipe()
	if err != nil {
		log.Fatalf("Unable to setup stdout: %s", err)
	}

	// Start the session
	if err := session.Shell(); err != nil {
		log.Fatalf("Failed to start shell: %s", err)
	}

	// Send a command to the remote shell
	_, err = stdin.Write([]byte("show version\n"))
	if err != nil {
		log.Fatalf("Failed to send command: %s", err)
	}

	// Pause execution for 1 second to allow the command to be executed
	time.Sleep(1 * time.Second)

	// Read and print the output
	buf := make([]byte, 1000)
	n, err := stdout.Read(buf)
	if err != nil && err != io.EOF {
		log.Fatalf("Failed to read stdout: %s", err)
	}

	// Print the output
	fmt.Println(string(buf[:n]))

	// Close the SSH connection
	if err := session.Close(); err != nil {
		log.Fatalf("Failed to close session: %s", err)
	}
}

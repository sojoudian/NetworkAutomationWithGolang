package main

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"io"
	"log"
	"time"
)

func saveRunningToStartup() {
	// Define the SSH connection details
	config := &ssh.ClientConfig{
		User: "u1",
		Auth: []ssh.AuthMethod{
			ssh.Password("cisco"),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // Insecure for testing, you can replace it with a proper callback
	}

	// Connect to the device (replace with your actual device details)
	conn, err := ssh.Dial("tcp", "192.168.122.100:22", config)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	// Create a new session
	session, err := conn.NewSession()
	if err != nil {
		log.Fatalf("Failed to create session: %v", err)
	}
	defer session.Close()

	// Request a pseudo-terminal for interactive input/output
	modes := ssh.TerminalModes{
		ssh.ECHO:  0,     // Disable echo
		ssh.TTY_OP_ISPEED: 14400, // Input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // Output speed = 14.4kbaud
	}
	if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
		log.Fatalf("request for pseudo terminal failed: %v", err)
	}

	// Start a shell
	stdin, err := session.StdinPipe()
	if err != nil {
		log.Fatalf("Unable to setup stdin for session: %v", err)
	}
	stdout, err := session.StdoutPipe()
	if err != nil {
		log.Fatalf("Unable to setup stdout for session: %v", err)
	}
	if err := session.Shell(); err != nil {
		log.Fatalf("Failed to start shell: %v", err)
	}

	// Send the commands to the shell
	commands := []string{
		"enable\n",
		"cisco\n", // Assuming 'cisco' is the enable password
		"copy running-config startup-config\n",
		"\n", // Confirm the copy action
		"show running-config | include clock\n", // Verify the running config
	}

	// Write commands to the session
	for _, cmd := range commands {
		io.WriteString(stdin, cmd)
		time.Sleep(2 * time.Second) // Wait for the device to process the command
	}

	// Read and print the output
	buf := make([]byte, 10000)
	n, err := stdout.Read(buf)
	if err != nil && err != io.EOF {
		log.Fatalf("Failed to read stdout: %v", err)
	}
	fmt.Println("Output:")
	fmt.Println(string(buf[:n]))

	// Close the session
	session.Close()
}

func main() {
	saveRunningToStartup()
}

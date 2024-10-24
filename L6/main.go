package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"golang.org/x/crypto/ssh"
)

const sourcesListContent = `deb http://deb.debian.org/debian/ bookworm main contrib non-free non-free-firmware
deb-src http://deb.debian.org/debian/ bookworm main contrib non-free non-free-firmware

deb http://deb.debian.org/debian/ bookworm-updates main contrib non-free non-free-firmware
deb-src http://deb.debian.org/debian/ bookworm-updates main contrib non-free non-free-firmware

deb http://deb.debian.org/debian-security bookworm-security main contrib non-free non-free-firmware
deb-src http://deb.debian.org/debian-security bookworm-security main contrib non-free non-free-firmware

deb http://deb.debian.org/debian/ bookworm-backports main contrib non-free non-free-firmware
deb-src http://deb.debian.org/debian/ bookworm-backports main contrib non-free non-free-firmware`

func main() {
	// Define command line flags
	host := flag.String("host", "", "SSH server hostname or IP")
	user := flag.String("user", "", "SSH username")
	pass := flag.String("pass", "", "SSH password")
	rootPass := flag.String("rootpass", "", "Root password")
	flag.Parse()

	// Validate input
	if *host == "" || *user == "" || *pass == "" || *rootPass == "" {
		log.Fatal("Please provide all required flags: -host, -user, -pass, and -rootpass")
	}

	// Configure SSH client
	config := &ssh.ClientConfig{
		User: *user,
		Auth: []ssh.AuthMethod{
			ssh.Password(*pass),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         30 * time.Second,
	}

	// Connect to the remote server
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:22", *host), config)
	if err != nil {
		log.Fatalf("Failed to dial: %s", err)
	}
	defer client.Close()

	// First, update sources.list to remove CD-ROM entry and add HTTP sources
	updateSourcesScript := fmt.Sprintf(`#!/bin/bash
echo '%s' | su - -c '
cp /etc/apt/sources.list /etc/apt/sources.list.backup
echo "%s" > /etc/apt/sources.list
'
`, *rootPass, sourcesListContent)

	// Update sources.list first
	session, err := client.NewSession()
	if err != nil {
		log.Fatalf("Failed to create session: %s", err)
	}

	fmt.Println("Updating sources.list...")
	if err := runCommand(session, updateSourcesScript); err != nil {
		log.Fatalf("Failed to update sources.list: %s", err)
	}
	session.Close()

	// Now install expect
	installScript := fmt.Sprintf(`#!/bin/bash
echo '%s' | su - -c 'apt-get update && apt-get install -y expect'
`, *rootPass)

	session, err = client.NewSession()
	if err != nil {
		log.Fatalf("Failed to create session: %s", err)
	}

	fmt.Println("Installing expect...")
	if err := runCommand(session, installScript); err != nil {
		log.Fatalf("Failed to install expect: %s", err)
	}
	session.Close()

	// Create system update script
	updateScript := fmt.Sprintf(`expect << 'EOF'
spawn su -
expect "Password: "
send "%s\r"
expect "# "
send "apt update\r"
expect "# "
send "apt upgrade -y\r"
expect "# "
send "exit\r"
expect eof
EOF`, *rootPass)

	// Create new session for the system update
	session, err = client.NewSession()
	if err != nil {
		log.Fatalf("Failed to create session: %s", err)
	}
	defer session.Close()

	// Run the system update
	fmt.Println("Executing system update process...")
	if err := runCommand(session, updateScript); err != nil {
		log.Fatalf("Failed to execute update process: %s", err)
	}

	fmt.Println("Update process completed successfully!")
}

func runCommand(session *ssh.Session, command string) error {
	// Create pipes for capturing output
	stdout, err := session.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdout pipe: %s", err)
	}

	stderr, err := session.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to create stderr pipe: %s", err)
	}

	// Start command
	if err := session.Start(command); err != nil {
		return fmt.Errorf("failed to start command: %s", err)
	}

	// Print stdout in real-time
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
	}()

	// Print stderr in real-time
	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			fmt.Fprintf(os.Stderr, "Error: %s\n", scanner.Text())
		}
	}()

	// Wait for command to complete
	if err := session.Wait(); err != nil {
		return fmt.Errorf("command failed: %s", err)
	}

	return nil
}

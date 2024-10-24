package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"time"

	"golang.org/x/crypto/ssh"
)

const sourcesListContent = `# Debian 12 (Bookworm) main repositories
deb http://deb.debian.org/debian/ bookworm main contrib non-free non-free-firmware
deb-src http://deb.debian.org/debian/ bookworm main contrib non-free non-free-firmware

# Debian 12 (Bookworm) updates
deb http://deb.debian.org/debian/ bookworm-updates main contrib non-free non-free-firmware
deb-src http://deb.debian.org/debian/ bookworm-updates main contrib non-free non-free-firmware

# Security updates
deb http://deb.debian.org/debian-security bookworm-security main contrib non-free non-free-firmware
deb-src http://deb.debian.org/debian-security bookworm-security main contrib non-free non-free-firmware

# Backports (optional, if you want newer versions of some packages)
deb http://deb.debian.org/debian/ bookworm-backports main contrib non-free non-free-firmware
deb-src http://deb.debian.org/debian/ bookworm-backports main contrib non-free non-free-firmware`

func main() {
	// Define command line flags
	host := flag.String("host", "", "SSH server hostname or IP")
	user := flag.String("user", "", "SSH username")
	pass := flag.String("pass", "", "SSH password")
	flag.Parse()

	// Validate input
	if *host == "" || *user == "" || *pass == "" {
		log.Fatal("Please provide all required flags: -host, -user, and -pass")
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

	// Create a new SSH session
	session, err := client.NewSession()
	if err != nil {
		log.Fatalf("Failed to create session: %s", err)
	}
	defer session.Close()

	// Backup existing sources.list
	backupCmd := "cp /etc/apt/sources.list /etc/apt/sources.list.backup"
	if err := runCommand(session, backupCmd); err != nil {
		log.Fatalf("Failed to backup sources.list: %s", err)
	}
	session.Close()

	// Create new session for writing sources.list
	session, err = client.NewSession()
	if err != nil {
		log.Fatalf("Failed to create session: %s", err)
	}

	// Write new sources.list
	writeCmd := fmt.Sprintf("echo '%s' | sudo tee /etc/apt/sources.list", sourcesListContent)
	if err := runCommand(session, writeCmd); err != nil {
		log.Fatalf("Failed to write sources.list: %s", err)
	}
	session.Close()

	// Run apt update
	session, err = client.NewSession()
	if err != nil {
		log.Fatalf("Failed to create session: %s", err)
	}
	fmt.Println("Running apt update...")
	if err := runCommand(session, "sudo apt update"); err != nil {
		log.Fatalf("Failed to run apt update: %s", err)
	}
	session.Close()

	// Run apt upgrade
	session, err = client.NewSession()
	if err != nil {
		log.Fatalf("Failed to create session: %s", err)
	}
	fmt.Println("Running apt upgrade...")
	if err := runCommand(session, "sudo apt upgrade -y"); err != nil {
		log.Fatalf("Failed to run apt upgrade: %s", err)
	}
	session.Close()

	fmt.Println("Update process completed successfully!")
}

func runCommand(session *ssh.Session, command string) error {
	// Create pipe for capturing output
	stdout, err := session.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdout pipe: %s", err)
	}

	// Start command
	if err := session.Start(command); err != nil {
		return fmt.Errorf("failed to start command: %s", err)
	}

	// Print output in real-time
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	// Wait for command to complete
	if err := session.Wait(); err != nil {
		return fmt.Errorf("command failed: %s", err)
	}

	return nil
}

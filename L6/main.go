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

	// Create expect-like script for su authentication
	suScript := fmt.Sprintf(`expect << 'EOF'
spawn su -
expect "Password: "
send "%s\r"
expect "# "
send "cp /etc/apt/sources.list /etc/apt/sources.list.backup\r"
expect "# "
send "cat > /etc/apt/sources.list << 'EOSOURCES'\r"
send "%s\r"
send "EOSOURCES\r"
expect "# "
send "apt update\r"
expect "# "
send "apt upgrade -y\r"
expect "# "
send "exit\r"
expect eof
EOF`, *rootPass, sourcesListContent)

	// First ensure expect is installed
	session, err := client.NewSession()
	if err != nil {
		log.Fatalf("Failed to create session: %s", err)
	}

	fmt.Println("Checking if expect is installed...")
	if err := runCommand(session, "which expect || (apt-get update && apt-get install -y expect)"); err != nil {
		log.Fatalf("Failed to verify/install expect: %s", err)
	}
	session.Close()

	// Create new session for the main script
	session, err = client.NewSession()
	if err != nil {
		log.Fatalf("Failed to create session: %s", err)
	}
	defer session.Close()

	// Run the expect script
	fmt.Println("Executing system update process...")
	if err := runCommand(session, suScript); err != nil {
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

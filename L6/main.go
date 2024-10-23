package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"

	"github.com/howeyc/gopass"
)

func main() {
	// Take user inputs for the server IP and username
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter the Debian server IP address: ")
	serverIP, _ := reader.ReadString('\n')
	serverIP = strings.TrimSpace(serverIP) // Trim newline

	fmt.Print("Enter the username: ")
	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username) // Trim newline

	// Prompt for password securely (masked with '*')
	fmt.Print("Enter the password: ")
	passwordBytes, err := gopass.GetPasswdMasked() // Password input is masked with '*'
	if err != nil {
		log.Fatalf("Failed to get password: %v", err)
	}
	password := string(passwordBytes)

	// SSH configuration
	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // Avoid strict host key checking for this demo
		Timeout:         5 * time.Second,             // Set timeout for connection
	}

	// Connect to the remote server
	conn, err := ssh.Dial("tcp", serverIP+":22", config)
	if err != nil {
		log.Fatalf("Failed to dial: %s", err)
	}
	defer conn.Close()

	// Create a session
	session, err := conn.NewSession()
	if err != nil {
		log.Fatalf("Failed to create session: %s", err)
	}
	defer session.Close()

	// Get a pseudo-terminal (PTY) for interaction with `su -`
	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // Disable echoing
		ssh.TTY_OP_ISPEED: 14400, // Input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // Output speed = 14.4kbaud
	}

	// Request a PTY to handle interactive commands
	err = session.RequestPty("xterm", 80, 40, modes)
	if err != nil {
		log.Fatalf("Failed to request PTY: %v", err)
	}

	// Start the shell
	stdin, err := session.StdinPipe()
	if err != nil {
		log.Fatalf("Failed to get stdin pipe: %v", err)
	}
	stdout, err := session.StdoutPipe()
	if err != nil {
		log.Fatalf("Failed to get stdout pipe: %v", err)
	}
	session.Stderr = os.Stderr

	// Start the shell session
	err = session.Shell()
	if err != nil {
		log.Fatalf("Failed to start shell: %v", err)
	}

	// Switch to root user with `su -` and provide password
	go func() {
		// Send `su -` command
		fmt.Fprintln(stdin, "su -")

		// Simulate typing the password for `su -`
		time.Sleep(1 * time.Second) // Wait for the prompt
		fmt.Fprintln(stdin, password)

		// Run the necessary commands after switching to root
		fmt.Fprintln(stdin, "
echo "# Debian 12 (Bookworm) main repositories" >> /etc/apt/sources.list
echo "deb http://deb.debian.org/debian/ bookworm main contrib non-free non-free-firmware" >> /etc/apt/sources.list
echo "deb-src http://deb.debian.org/debian/ bookworm main contrib non-free non-free-firmware" >> /etc/apt/sources.list
echo "# Debian 12 (Bookworm) updates" >> /etc/apt/sources.list
echo "deb http://deb.debian.org/debian/ bookworm-updates main contrib non-free non-free-firmware" >> /etc/apt/sources.list
echo "deb-src http://deb.debian.org/debian/ bookworm-updates main contrib non-free non-free-firmware" >> /etc/apt/sources.list
echo "# Security updates" >> /etc/apt/sources.list
echo "deb http://deb.debian.org/debian-security bookworm-security main contrib non-free non-free-firmware" >> /etc/apt/sources.list
echo "deb-src http://deb.debian.org/debian-security bookworm-security main contrib non-free non-free-firmware" >> /etc/apt/sources.list
echo "# Backports (optional, if you want newer versions of some packages)" >> /etc/apt/sources.list
echo "deb http://deb.debian.org/debian/ bookworm-backports main contrib non-free non-free-firmware" >> /etc/apt/sources.list
echo "deb-src http://deb.debian.org/debian/ bookworm-backports main contrib non-free non-free-firmware" >> /etc/apt/sources.list
sudo apt update
sudo apt upgrade -y
")

		// Exit the session after commands are done
		fmt.Fprintln(stdin, "exit")
	}()

	// Capture the session output and print to console
	go io.Copy(os.Stdout, stdout)

	// Wait for the session to complete
	err = session.Wait()
	if err != nil {
		log.Fatalf("Session finished with error: %v", err)
	}

	fmt.Println("Commands executed successfully on the remote server after switching to root.")
}

package main

import (
	"bufio"
	"fmt"
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

	// Request a pseudo-terminal (PTY) to handle interactive commands
	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // Disable echoing
		ssh.TTY_OP_ISPEED: 14400, // Input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // Output speed = 14.4kbaud
	}
	err = session.RequestPty("xterm", 80, 40, modes)
	if err != nil {
		log.Fatalf("Failed to request PTY: %v", err)
	}

	// Create pipes for session input/output
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

	// Capture the stdout and look for the password prompt
	go func() {
		buf := make([]byte, 1024)
		for {
			n, _ := stdout.Read(buf)
			output := string(buf[:n])
			fmt.Print(output) // Print output to terminal for debugging

			// Look for the password prompt after `su -`
			if strings.Contains(output, "Password:") {
				fmt.Fprintln(stdin, password) // Send the password for `su -`
				break
			}
		}
	}()

	// Send `su -` command
	fmt.Fprintln(stdin, "su -")

	// Wait for the session to become root, then send commands
	go func() {
		// Wait for a short while before sending commands
		time.Sleep(2 * time.Second)

		// Send individual commands to update sources.list
		fmt.Fprintln(stdin, "echo '# Debian 12 (Bookworm) main repositories' | sudo tee -a /etc/apt/sources.list")
		fmt.Fprintln(stdin, "echo 'deb http://deb.debian.org/debian/ bookworm main contrib non-free non-free-firmware' | sudo tee -a /etc/apt/sources.list")
		fmt.Fprintln(stdin, "echo 'deb-src http://deb.debian.org/debian/ bookworm main contrib non-free non-free-firmware' | sudo tee -a /etc/apt/sources.list")

		fmt.Fprintln(stdin, "echo '# Debian 12 (Bookworm) updates' | sudo tee -a /etc/apt/sources.list")
		fmt.Fprintln(stdin, "echo 'deb http://deb.debian.org/debian/ bookworm-updates main contrib non-free non-free-firmware' | sudo tee -a /etc/apt/sources.list")
		fmt.Fprintln(stdin, "echo 'deb-src http://deb.debian.org/debian/ bookworm-updates main contrib non-free non-free-firmware' | sudo tee -a /etc/apt/sources.list")

		fmt.Fprintln(stdin, "echo '# Security updates' | sudo tee -a /etc/apt/sources.list")
		fmt.Fprintln(stdin, "echo 'deb http://deb.debian.org/debian-security bookworm-security main contrib non-free non-free-firmware' | sudo tee -a /etc/apt/sources.list")
		fmt.Fprintln(stdin, "echo 'deb-src http://deb.debian.org/debian-security bookworm-security main contrib non-free non-free-firmware' | sudo tee -a /etc/apt/sources.list")

		fmt.Fprintln(stdin, "echo '# Backports (optional, if you want newer versions of some packages)' | sudo tee -a /etc/apt/sources.list")
		fmt.Fprintln(stdin, "echo 'deb http://deb.debian.org/debian/ bookworm-backports main contrib non-free non-free-firmware' | sudo tee -a /etc/apt/sources.list")
		fmt.Fprintln(stdin, "echo 'deb-src http://deb.debian.org/debian/ bookworm-backports main contrib non-free non-free-firmware' | sudo tee -a /etc/apt/sources.list")

		// Run apt update and upgrade
		fmt.Fprintln(stdin, "sudo apt update")
		fmt.Fprintln(stdin, "sudo apt upgrade -y")

		// Exit the session after commands are done
		fmt.Fprintln(stdin, "exit")
	}()

	// Wait for the session to complete
	err = session.Wait()
	if err != nil {
		log.Fatalf("Session finished with error: %v", err)
	}

	fmt.Println("Commands executed successfully on the remote server after switching to root.")
}

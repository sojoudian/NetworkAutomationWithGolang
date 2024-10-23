package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"

	"golang.org/x/crypto/ssh"
)

func main() {
	// Take user inputs for the server IP, username, and password
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter the Debian server IP address: ")
	serverIP, _ := reader.ReadString('\n')
	serverIP = serverIP[:len(serverIP)-1] // Trim newline

	fmt.Print("Enter the username: ")
	username, _ := reader.ReadString('\n')
	username = username[:len(username)-1] // Trim newline

	fmt.Print("Enter the password: ")
	password, _ := reader.ReadString('\n')
	password = password[:len(password)-1] // Trim newline

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

	// // Capture both stdout and stderr from the session
	// var stdoutBuf, stderrBuf io.Writer
	// session.Stdout = os.Stdout
	// session.Stderr = os.Stderr

	// Run the necessary commands on the remote server
	commands := `
	echo "# Debian 12 (Bookworm) main repositories" | sudo tee -a /etc/apt/sources.list
	echo "deb http://deb.debian.org/debian/ bookworm main contrib non-free non-free-firmware" | sudo tee -a /etc/apt/sources.list
	echo "deb-src http://deb.debian.org/debian/ bookworm main contrib non-free non-free-firmware" | sudo tee -a /etc/apt/sources.list
	echo "# Debian 12 (Bookworm) updates" | sudo tee -a /etc/apt/sources.list
	echo "deb http://deb.debian.org/debian/ bookworm-updates main contrib non-free non-free-firmware" | sudo tee -a /etc/apt/sources.list
	echo "deb-src http://deb.debian.org/debian/ bookworm-updates main contrib non-free non-free-firmware" | sudo tee -a /etc/apt/sources.list
	echo "# Security updates" | sudo tee -a /etc/apt/sources.list
	echo "deb http://deb.debian.org/debian-security bookworm-security main contrib non-free non-free-firmware" | sudo tee -a /etc/apt/sources.list
	echo "deb-src http://deb.debian.org/debian-security bookworm-security main contrib non-free non-free-firmware" | sudo tee -a /etc/apt/sources.list
	echo "# Backports (optional, if you want newer versions of some packages)" | sudo tee -a /etc/apt/sources.list
	echo "deb http://deb.debian.org/debian/ bookworm-backports main contrib non-free non-free-firmware" | sudo tee -a /etc/apt/sources.list
	echo "deb-src http://deb.debian.org/debian/ bookworm-backports main contrib non-free non-free-firmware" | sudo tee -a /etc/apt/sources.list
	sudo apt update
	sudo apt upgrade -y
	`

	// Run the commands on the remote server
	err = session.Run(commands)
	if err != nil {
		log.Fatalf("Failed to run commands on remote server: %v", err)
	}

	fmt.Println("Commands executed successfully on the remote server.")
}

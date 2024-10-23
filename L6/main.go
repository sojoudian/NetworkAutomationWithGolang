package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

func main() {
	// Define the repositories to be added to sources.list
	repoLines := `
# Debian 12 (Bookworm) main repositories
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
deb-src http://deb.debian.org/debian/ bookworm-backports main contrib non-free non-free-firmware
`

	// Path to sources.list
	sourcesList := "/etc/apt/sources.list"

	// Step 1: Open the sources.list file in append mode
	file, err := os.OpenFile(sourcesList, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open %s: %v", sourcesList, err)
	}
	defer file.Close()

	// Step 2: Append the repository lines to sources.list
	_, err = file.WriteString(repoLines)
	if err != nil {
		log.Fatalf("Failed to write to %s: %v", sourcesList, err)
	}

	fmt.Println("Successfully added Debian 12 repositories to sources.list")

	// Step 3: Run "apt update"
	fmt.Println("Running apt update...")
	cmdUpdate := exec.Command("sudo", "apt", "update")
	cmdUpdate.Stdout = os.Stdout
	cmdUpdate.Stderr = os.Stderr

	err = cmdUpdate.Run()
	if err != nil {
		log.Fatalf("Failed to run apt update: %v", err)
	}
	fmt.Println("apt update completed successfully.")

	// Step 4: Run "apt upgrade -y"
	fmt.Println("Running apt upgrade -y...")
	cmdUpgrade := exec.Command("sudo", "apt", "upgrade", "-y")
	cmdUpgrade.Stdout = os.Stdout
	cmdUpgrade.Stderr = os.Stderr

	err = cmdUpgrade.Run()
	if err != nil {
		log.Fatalf("Failed to run apt upgrade: %v", err)
	}
	fmt.Println("apt upgrade completed successfully.")
}

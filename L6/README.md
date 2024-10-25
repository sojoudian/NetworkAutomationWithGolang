# Debian Sources Update Automation

This Go application automates the process of updating the sources.list file on a Debian server and running system updates. It connects to a remote Debian server via SSH, updates the sources.list file with Debian 12 (Bookworm) repositories, and performs system updates.

## Features

- Remote SSH connection to Debian server
- Automatic backup of existing sources.list
- Updates sources.list with Debian 12 (Bookworm) repositories
- Installation of 'expect' package for automation
- Executes apt update and upgrade commands
- Real-time command output display
- Error handling and logging

## Prerequisites

- Go 1.16 or higher
- Access to a Debian server with:
  - SSH access enabled
  - Root password access
  - Password authentication enabled

## Installation

1. Clone the repository:
```bash
git clone https://github.com/sojoudian/NetworkAutomationWithGolang.git
cd NetworkAutomationWithGolang/L6/
```

2. Install the required dependency:
```bash
go get golang.org/x/crypto/ssh
```

3. Build the application:
```bash
go build -o debian-updater
```

## Usage

Run the application with the following command-line flags:

```bash
./debian-updater -host=SERVER_IP -user=USERNAME -pass=PASSWORD -rootpass=ROOT_PASSWORD
```

### Required Flags:
- `-host`: The IP address or hostname of your Debian server
- `-user`: SSH username
- `-pass`: SSH password
- `-rootpass`: Root password for system updates

### Example:
```bash
./debian-updater -host=192.168.1.100 -user=admin -pass=your_password -rootpass=root_password
```

## What the Application Does

1. Connects to the specified Debian server via SSH
2. Creates a backup of the existing /etc/apt/sources.list file
3. Updates sources.list with the following repositories:
   - Debian 12 (Bookworm) main repositories
   - Debian 12 updates
   - Security updates
   - Backports
4. Installs the 'expect' package for automated password entry
5. Runs `apt update` to refresh package lists
6. Runs `apt upgrade -y` to upgrade all packages

## Output

The application provides real-time output of all operations, including:
- SSH connection status
- File backup confirmation
- Sources list update status
- Expect package installation progress
- apt update and upgrade progress

## Error Handling

The application includes error handling for:
- SSH connection failures
- Authentication issues
- Command execution errors
- File operation failures
- Package installation errors

## Security Considerations

- The passwords are passed as command-line arguments, which might be visible in process listings
- For production use, consider implementing SSH key authentication
- The application requires root password access
- Host key verification is disabled for simplicity (uses InsecureIgnoreHostKey)

## Backup

A backup of the original sources.list is automatically created at:
```
/etc/apt/sources.list.backup
```

## Troubleshooting

### Common Issues:

1. SSH Connection Failed:
   - Verify the server IP address is correct
   - Ensure SSH service is running on the server
   - Check if port 22 is open and accessible

2. Authentication Failed:
   - Verify username and password are correct
   - Ensure password authentication is enabled in SSH config

3. Root Access Issues:
   - Verify the root password is correct
   - Check if root login is enabled on the system

4. Package Installation Failures:
   - Ensure the system has internet connectivity
   - Verify the sources.list is properly updated
   - Check for sufficient disk space

## Contributing

Feel free to submit issues, fork the repository, and create pull requests for any improvements.

## License
MIT

## Note

This tool is designed for Debian 12 (Bookworm). If you're using a different version of Debian, modify the sources.list content in the code accordingly.


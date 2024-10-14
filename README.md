
# Network Automation with Go

## Overview
This project demonstrates the use of **Go** for automating network tasks, specifically interacting with Cisco devices via SSH. The application connects to a network device, enters enable mode, and saves the running configuration to the startup configuration, ensuring that changes made to the device are persistent.

## Features
- Establishes an SSH connection to a network device.
- Automates sending commands like entering enable mode and saving configurations.
- Verifies configuration changes.
- Uses Go's powerful concurrency and networking capabilities.

## Requirements
To run this project, you need:
- Go 1.18 or later installed.
- A Cisco device accessible over SSH.
- The `golang.org/x/crypto/ssh` package installed.

## Installation

1. **Install Go**:
   If you haven't installed Go yet, download it from [here](https://golang.org/dl/).

2. **Clone the repository**:
   ```bash
   git clone https://github.com/yourusername/network-automation-go.git
   cd network-automation-go
   ```

3. **Install dependencies**:
   Install the necessary Go package by running:
   ```bash
   go get golang.org/x/crypto/ssh
   ```

## Usage

1. Open the `main.go` file and replace the placeholder values for the device's IP address, username, and password with your actual network device credentials.

2. Run the Go application:
   ```bash
   go run main.go
   ```

3. The application will connect to the Cisco device, enter enable mode, and save the running configuration to the startup configuration.

## Customization

You can easily modify this code to execute other network automation tasks such as:
- Backing up configurations.
- Pushing configuration changes.
- Gathering operational data from network devices.

## License
This project is licensed under the MIT License.

## Contribution
Contributions are welcome! Please fork the repository and submit a pull request for any enhancements or bug fixes.

## Contact
For any inquiries, please contact [your.email@example.com](mailto:your.email@example.com).


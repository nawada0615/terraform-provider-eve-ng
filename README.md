# Terraform Provider for EVE-NG

A Terraform provider for managing EVE-NG (Emulated Virtual Environment Next Generation) resources with enhanced reliability and comprehensive error handling.

## Features

This provider supports the following EVE-NG resources with robust error handling and fallback mechanisms:

### Resources
- **eve_folder** - Manage EVE-NG folders
- **eve_lab** - Manage EVE-NG labs
- **eve_network** - Manage lab networks
- **eve_node** - Manage lab nodes (QEMU, Dynamips, IOL, Docker, VPCS)
- **eve_interface_attachment** - Manage node interface connections

### Data Sources
- **eve_templates** - Get available node templates
- **eve_network_types** - Get available network types
- **eve_icons** - Get available icons
- **eve_status** - Get system status

## Recent Improvements

### ✅ Enhanced Reliability
- **Fallback Mechanisms**: Network individual read failures automatically fallback to network list
- **Comprehensive Error Handling**: Detailed error messages and proper status code validation
- **Type Safety**: Robust handling of API response type variations
- **Debug Logging**: Extensive logging for troubleshooting and monitoring

### 🛡️ Robust Error Handling
- **API Response Validation**: Proper validation of all API responses
- **State Management**: Correct handling of resource not found scenarios
- **Authentication**: Enhanced login error handling with detailed messages
- **Network Management**: Resilient network operations with fallback support

## Installation

### Local Installation

1. Clone this repository
2. Build the provider:
   ```bash
   make build
   ```
3. Install locally:
   ```bash
   make install
   ```

### Usage

```hcl
terraform {
  required_providers {
    eve-ng = {
      source = "local/eve-ng"
      version = "1.0.0"
    }
  }
}

provider "eve-ng" {
  endpoint = "http://192.168.0.101"
  username = "admin"
  password = "eve"
  insecure_skip_verify = true
}

# Create a lab
resource "eve_lab" "test_lab" {
  path        = "/"
  name        = "test-lab"
  author      = "Terraform"
  description = "Test lab created by Terraform"
  version     = "1"
}

# Create a network
resource "eve_network" "net1" {
  lab_file = eve_lab.test_lab.file
  name     = "net1"
  type     = "bridge"
}

# Create a QEMU node with custom options
resource "eve_node" "linux_node" {
  lab_file       = eve_lab.test_lab.file
  name           = "linux-node"
  type           = "qemu"
  template       = "linux"
  cpu            = 2
  ram            = 2048
  qemu_version   = "2.4.0"
  qemu_arch      = "x86_64"
  qemu_nic       = "virtio-net-pci"
  qemu_options   = "-enable-kvm"
  cpulimit       = true
  desired_state  = "started"
}

# Connect node to network
resource "eve_interface_attachment" "net1_connection" {
  lab_file        = eve_lab.test_lab.file
  node_id         = eve_node.linux_node.id
  interface_index = 0
  target          = "network:${eve_network.net1.id}"
}

# Create a second node
resource "eve_node" "linux_node2" {
  lab_file       = eve_lab.test_lab.file
  name           = "linux-node2"
  type           = "qemu"
  template       = "linux"
  cpu            = 1
  ram            = 1024
  qemu_version   = "2.4.0"
  qemu_arch      = "x86_64"
  qemu_nic       = "virtio-net-pci"
  qemu_options   = "-enable-kvm"
  cpulimit       = true
  desired_state  = "started"
}

# Create a second network for node-to-node connection
resource "eve_network" "connection_net" {
  lab_file   = eve_lab.test_lab.file
  name       = "connection-net"
  type       = "bridge"
  visibility = "0"  # Hidden network
}

# Connect both nodes to the same network for direct communication
resource "eve_interface_attachment" "node1_to_net" {
  lab_file        = eve_lab.test_lab.file
  node_id         = eve_node.linux_node.id
  interface_index = 0
  target          = "network:${eve_network.connection_net.id}"
}

resource "eve_interface_attachment" "node2_to_net" {
  lab_file        = eve_lab.test_lab.file
  node_id         = eve_node.linux_node2.id
  interface_index = 0
  target          = "network:${eve_network.connection_net.id}"
}
```

## Development

### Prerequisites
- Go 1.21 or later
- Terraform 1.0 or later

### Building
```bash
make build
```

### Testing
```bash
# Run all tests
make test

# Run tests with race detection
make test-race

# Run specific test suites
make test-unit
make test-client
make test-mock
```

### Code Quality
```bash
# Format code
make fmt

# Run linter
make lint

# Run vet
make vet
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for your changes
5. Run the test suite
6. Submit a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

- EVE-NG community for the excellent network emulation platform
- HashiCorp for the Terraform plugin SDK

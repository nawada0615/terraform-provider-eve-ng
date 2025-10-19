# EVE-NG API Documentation

## Overview
This document summarizes the REST API specifications for EVE-NG (EVE-NG Community Edition 6.2.0-4). It is based on verification results from actual hardware.

## Basic Information
- **Base URL**: `http://192.168.0.101/api/`
- **Authentication**: Cookie-based authentication (`unetlab_session`)
- **Content-Type**: `application/json`
- **API Version**: 6.2.0-4

## Authentication

### Login
```http
POST /api/auth/login
Content-Type: application/json

{
  "username": "admin",
  "password": "eve",
  "html5": 1
}
```

**Response**:
```json
{
  "code": 200,
  "status": "success",
  "message": "User logged in (90013)."
}
```

**Note**: The `html5` parameter is required.

### Logout
```http
GET /api/auth/logout
Cookie: unetlab_session=<session_id>
```

### Authentication Check
```http
GET /api/auth
Cookie: unetlab_session=<session_id>
```

## System Information

### Get Status
```http
GET /api/status
Cookie: unetlab_session=<session_id>
```

**Response**:
```json
{
  "code": 200,
  "status": "success",
  "message": "Fetched system status (60001).",
  "data": {
    "version": "6.2.0-4",
    "qemu_version": "2.4.0",
    "uksm": "unsupported",
    "ksm": "enabled",
    "cpulimit": "enabled",
    "cpu": 0,
    "disk": 35,
    "cached": 98,
    "mem": 3,
    "swap": 0,
    "iol": 0,
    "dynamips": 0,
    "qemu": 0,
    "docker": 0,
    "vpcs": 0
  }
}
```

## Folder Management

### Get Folder List
```http
GET /api/folders/
Cookie: unetlab_session=<session_id>
```

**Response**:
```json
{
  "code": 200,
  "status": "success",
  "message": "Successfully listed path (60007).",
  "data": {
    "folders": [
      {
        "name": "projects",
        "path": "/projects"
      }
    ],
    "labs": [
      {
        "file": "lab.unl",
        "path": "/lab.unl",
        "umtime": 1760274335,
        "mtime": "12 Oct 2025 15:05"
      }
    ]
  }
}
```

### Create Folder
```http
POST /api/folders
Content-Type: application/json
Cookie: unetlab_session=<session_id>

{
  "path": "/new-folder",
  "name": "new-folder"
}
```

## Lab Management

### Create Lab
```http
POST /api/labs
Content-Type: application/json
Cookie: unetlab_session=<session_id>

{
  "path": "/",
  "name": "test-lab",
  "author": "admin",
  "description": "Test lab",
  "version": "1"
}
```

**Response**:
```json
{
  "code": 201,
  "status": "success",
  "message": "Lab has been saved (60023).",
  "data": {
    "id": 1
  }
}
```

### Get Lab
```http
GET /api/labs/<lab_file>
Cookie: unetlab_session=<session_id>
```

### Update Lab
```http
PUT /api/labs/<lab_file>
Content-Type: application/json
Cookie: unetlab_session=<session_id>

{
  "id": 1,
  "name": "updated-lab",
  "author": "admin",
  "description": "Updated lab description",
  "version": "2"
}
```

### Delete Lab
```http
DELETE /api/labs/<lab_file>
Cookie: unetlab_session=<session_id>
```

## Node Management

### Create Node
```http
POST /api/labs/<lab_file>/nodes
Content-Type: application/json
Cookie: unetlab_session=<session_id>

{
  "type": "qemu",
  "template": "linux",
  "name": "test-node",
  "image": "linux-ubuntu-22.04",
  "icon": "Router-2D-Gen-White-S.svg",
  "top": 100,
  "left": 100,
  "cpu": 1,
  "ram": 1024,
  "ethernet": 4,
  "serial": 0,
  "delay": 0,
  "console": "telnet",
  "qemu_version": "2.4.0",
  "qemu_arch": "x86_64",
  "qemu_nic": "virtio-net-pci",
  "qemu_options": "-enable-kvm",
  "cpulimit": 1,
  "desired_state": "stopped"
}
```

**Response**:
```json
{
  "code": 201,
  "status": "success",
  "message": "Node has been saved (60025).",
  "data": {
    "id": 1
  }
}
```

### Get Node
```http
GET /api/labs/<lab_file>/nodes/<node_id>
Cookie: unetlab_session=<session_id>
```

**Response**:
```json
{
  "code": 200,
  "status": "success",
  "message": "Successfully listed node (60025).",
  "data": {
    "console": "telnet",
    "config": 0,
    "delay": 0,
    "left": 100,
    "icon": "Router-2D-Gen-White-S.svg",
    "image": "linux-ubuntu-22.04",
    "name": "test-node",
    "status": 0,
    "template": "linux",
    "type": "qemu",
    "top": 100,
    "url": "/html5/#/client/...",
    "cpulimit": 1,
    "cpu": 1,
    "ethernet": 4,
    "ram": 1024,
    "uuid": "9de1d69f-5f2d-4e47-bb3f-312125576dc8",
    "firstmac": "00:50:00:00:01:00",
    "qemu_options": "-machine type=pc,accel=kvm -vga std -usbdevice tablet -boot order=cd -cpu host",
    "qemu_version": "2.4.0",
    "qemu_arch": "x86_64",
    "qemu_nic": "virtio-net-pci"
  }
}
```

### Update Node
```http
PUT /api/labs/<lab_file>/nodes/<node_id>
Content-Type: application/json
Cookie: unetlab_session=<session_id>

{
  "id": 1,
  "name": "updated_node_name",
  "top": 150,
  "left": 150
}
```

### Delete Node
```http
DELETE /api/labs/<lab_file>/nodes/<node_id>
Cookie: unetlab_session=<session_id>
```

### Node Operations
```http
GET /api/labs/<lab_file>/nodes/<node_id>/start
GET /api/labs/<lab_file>/nodes/<node_id>/stop
GET /api/labs/<lab_file>/nodes/<node_id>/wipe
Cookie: unetlab_session=<session_id>
```

## Node Interface Management

### Get Node Interfaces
```http
GET /api/labs/<lab_file>/nodes/<node_id>/interfaces
Cookie: unetlab_session=<session_id>
```

**Response**:
```json
{
  "code": 200,
  "status": "success",
  "message": "Successfully listed node interfaces (60030).",
  "data": {
    "ethernet": [
      {
        "name": "eth0",
        "network_id": 1
      }
    ],
    "serial": [
      {
        "name": "serial0",
        "remote_id": null,
        "remote_if": null
      }
    ],
    "id": 1,
    "sort": "ethernet"
  }
}
```

### Update Node Interfaces
```http
PUT /api/labs/<lab_file>/nodes/<node_id>/interfaces
Content-Type: application/json
Cookie: unetlab_session=<session_id>

{
  "0": 1
}
```

## Network Management

### Create Network
```http
POST /api/labs/<lab_file>/networks
Content-Type: application/json
Cookie: unetlab_session=<session_id>

{
  "type": "bridge",
  "name": "test-network",
  "icon": "lan.png",
  "top": 200,
  "left": 200,
  "visibility": "1"
}
```

**Response**:
```json
{
  "code": 201,
  "status": "success",
  "message": "Network has been saved (60028).",
  "data": {
    "id": 1
  }
}
```

### Get Network
```http
GET /api/labs/<lab_file>/networks/<network_id>
Cookie: unetlab_session=<session_id>
```

### Update Network
```http
PUT /api/labs/<lab_file>/networks/<network_id>
Content-Type: application/json
Cookie: unetlab_session=<session_id>

{
  "id": 1,
  "name": "updated-network",
  "visibility": "0"
}
```

### Delete Network
```http
DELETE /api/labs/<lab_file>/networks/<network_id>
Cookie: unetlab_session=<session_id>
```

## Data Sources

### Get Templates
```http
GET /api/list/templates/
Cookie: unetlab_session=<session_id>
```

**Response**:
```json
{
  "code": 200,
  "status": "success",
  "message": "Successfully listed templates (60003).",
  "data": {
    "linux": {
      "name": "linux",
      "type": "qemu",
      "description": "Linux template",
      "icon": "Server-2D-Linux-S.svg"
    }
  }
}
```

### Get Network Types
```http
GET /api/list/networks
Cookie: unetlab_session=<session_id>
```

**Response**:
```json
{
  "code": 200,
  "status": "success",
  "message": "Successfully listed network types (60002).",
  "data": {
    "bridge": "bridge",
    "ovs": "ovs",
    "pnet0": "pnet0",
    "pnet1": "pnet1",
    "pnet2": "pnet2",
    "pnet3": "pnet3",
    "pnet4": "pnet4",
    "pnet5": "pnet5",
    "pnet6": "pnet6",
    "pnet7": "pnet7",
    "pnet8": "pnet8",
    "pnet9": "pnet9"
  },
  "icons": {
    "01-Cloud-Default.svg": "01-Cloud-Default.svg",
    "cloud.png": "cloud.png",
    "Router-2D-Gen-White-S.svg": "Router-2D-Gen-White-S.svg"
  }
}
```

## Error Responses

### General Error Format
```json
{
  "code": 404,
  "status": "fail",
  "message": "Resource not found (60038)."
}
```

### Authentication Error
```json
{
  "code": 412,
  "status": "unauthorized",
  "message": "User is not authenticated or session timed out (90001)."
}
```

### Permission Error
```json
{
  "code": 401,
  "status": "forbidden",
  "message": "User is not authenticated or session timed out (90001)."
}
```

## Parameter Specifications

### Lab Resource (eve_lab)

| Parameter | Required/Optional | Default Value | Description |
|-----------|------------------|---------------|-------------|
| path | Required | - | Lab path (must start with `/`) |
| name | Required | - | Lab name |
| author | Optional | "" | Author name |
| description | Optional | "" | Description |
| body | Optional | "" | Body text |
| version | Optional | "1" | Version |
| scripttimeout | Optional | 300 | Script timeout (seconds) |
| lock | Optional | false | Lock status |

### Node Resource (eve_node)

| Parameter | Required/Optional | Default Value | Description |
|-----------|------------------|---------------|-------------|
| lab_file | Required | - | Lab file path |
| name | Required | - | Node name |
| type | Required | - | Node type (qemu, dynamips, iol, docker, vpcs) |
| template | Required | - | Template name |
| image | Optional | "" | Image name |
| icon | Optional | "" | Icon |
| top | Optional | 0 | Y coordinate |
| left | Optional | 0 | X coordinate |
| cpu | Optional | 1 | CPU count |
| ram | Optional | 1024 | RAM (MB) |
| ethernet | Optional | 0 | Ethernet interface count |
| serial | Optional | 0 | Serial interface count |
| delay | Optional | 0 | Boot delay (seconds) |
| console | Optional | "telnet" | Console type |
| qemu_version | Optional | "2.4.0" | QEMU version |
| qemu_arch | Optional | "x86_64" | QEMU architecture |
| qemu_nic | Optional | "virtio-net-pci" | QEMU NIC type |
| qemu_options | Optional | "" | QEMU options |
| cpulimit | Optional | true | CPU limit |
| desired_state | Optional | "stopped" | Desired state |

### Network Resource (eve_network)

| Parameter | Required/Optional | Default Value | Description |
|-----------|------------------|---------------|-------------|
| lab_file | Required | - | Lab file path |
| name | Required | - | Network name |
| type | Required | - | Network type (bridge, ovs, pnet0-9) |
| icon | Optional | "" | Icon |
| top | Optional | 0 | Y coordinate |
| left | Optional | 0 | X coordinate |
| visibility | Optional | "1" | Visibility (0=hidden, 1=visible) |

## Notes

- All API endpoints require authentication via `unetlab_session` cookie
- The `html5` parameter is required for login requests
- Network visibility can be controlled using the `visibility` parameter (0=hidden, 1=visible)
- Node IDs and Network IDs are returned as strings in the format `<lab_file>:<type>:<id>`
- Interface attachments require numeric IDs extracted from the string format

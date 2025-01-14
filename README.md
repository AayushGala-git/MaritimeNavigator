# Project3

## Overview

This project simulates communication between satellites, vessels, and ground stations. It is designed with modular components in Go to handle satellite networking, vessel simulations, and ground station interactions.

---

## Project Structure

```
.
├── README.md               # Project documentation
├── api/
│   └── httpserver.go       # HTTP server implementation
├── cmd/
│   └── main.go             # Main entry point of the application
├── config.json             # Configuration file for the application
├── database.json           # Database configuration
├── go.mod                  # Go module configuration
├── pkg/
│   ├── common/
│   │   ├── config.go       # Configuration loading utilities
│   │   └── logger.go       # Logging utilities
│   ├── groundstation/
│   │   ├── database.go     # Ground station database operations
│   │   └── groundstation.go# Ground station logic
│   ├── protocol/
│   │   ├── handlers.go     # Protocol handler implementations
│   │   └── protocol.go     # Protocol definitions
│   ├── satellite/
│   │   ├── network.go      # Satellite network interactions
│   │   ├── satellite.go    # Satellite core functionality
│   │   └── simulator.go    # Satellite simulation logic
│   └── vessel/
│       ├── simulator.go    # Vessel simulation logic
│       └── vessel.go       # Vessel core functionality
└── run                     # Compiled binary (after build)
```

---

## Getting Started

### Prerequisites

- **Go 1.20+**
- Configuration files: `config.json` and `database.json`

### Steps to Run the Project

#### 1. Install Dependencies

Use Go modules to download and install required dependencies:

```bash
go mod tidy
```

#### 2. Run the Application

**Option 1: Direct Execution**

Run the main Go file directly:

```bash
go run cmd/main.go
```

**Option 2: Build and Execute**

Compile the project into an executable binary and run it:

```bash
go build -o run cmd/main.go
./run
```

---

## Configuration

### config.json

This file contains the application settings. Update it as needed for your environment. A sample structure might look like:

```json
{
   "ground_station_address": "127.0.0.1:8080",
   "satellites": [
      {
         "id": "Satellite-1",
         "port": 8001,
         "neighbors": [
            {
               "id": "Satellite-2",
               "latency": 50,
               "packet_loss": 0.1
            },
            {
               "id": "Satellite-3",
               "latency": 70,
               "packet_loss": 0.15
            }
         ]
      },
      {
         "id": "Satellite-2",
         "port": 8002,
         "neighbors": [
            {
               "id": "Satellite-1",
               "latency": 50,
               "packet_loss": 0.1
            },
            {
               "id": "Satellite-4",
               "latency": 40,
               "packet_loss": 0.05
            }
         ]
      },
      {
         "id": "Satellite-3",
         "port": 8003,
         "neighbors": [
            {
               "id": "Satellite-1",
               "latency": 70,
               "packet_loss": 0.15
            },
            {
               "id": "Satellite-5",
               "latency": 60,
               "packet_loss": 0.1
            }
         ]
      },
      {
         "id": "Satellite-4",
         "port": 8004,
         "neighbors": [
            {
               "id": "Satellite-2",
               "latency": 40,
               "packet_loss": 0.05
            },
            {
               "id": "Satellite-5",
               "latency": 30,
               "packet_loss": 0.2
            }
         ]
      },
      {
         "id": "Satellite-5",
         "port": 8005,
         "neighbors": [
            {
               "id": "Satellite-3",
               "latency": 60,
               "packet_loss": 0.1
            },
            {
               "id": "Satellite-4",
               "latency": 30,
               "packet_loss": 0.2
            }
         ]
      }
   ],
   "vessels": [
      { "id": "Vessel-1", "satellite": "Satellite-1" },
      { "id": "Vessel-2", "satellite": "Satellite-1" },
      { "id": "Vessel-3", "satellite": "Satellite-2" },
      { "id": "Vessel-4", "satellite": "Satellite-2" },
      { "id": "Vessel-5", "satellite": "Satellite-3" },
      { "id": "Vessel-6", "satellite": "Satellite-3" },
      { "id": "Vessel-7", "satellite": "Satellite-4" },
      { "id": "Vessel-8", "satellite": "Satellite-4" },
      { "id": "Vessel-9", "satellite": "Satellite-5" },
      { "id": "Vessel-10", "satellite": "Satellite-5" }
   ]
}
```

### database.json

Defines the database-related configurations for storing application data.

---

## Key Features

### HTTP Server (`api/httpserver.go`)

Provides HTTP endpoints for external interaction.

### Ground Station (`pkg/groundstation`)

Handles ground station database operations and communication logic.

### Satellite Module (`pkg/satellite`)

Simulates satellite operations, including networking and behavior.

### Vessel Module (`pkg/vessel`)

Implements vessel simulation and core functionalities.

### Protocol Module (`pkg/protocol`)

Manages communication protocols and their handlers.
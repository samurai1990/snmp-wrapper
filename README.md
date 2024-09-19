# SNMP-Wrapper for MikroTik REST API

This application is a wrapper for the SNMP protocol, written in **Go (Golang)**. It retrieves information from the **MikroTik REST API**, converts the data, and exposes it via **SNMP**.

## Features

- Retrieves device information via MikroTik's REST API.
- Converts the retrieved data into SNMP-compatible format **zabbix mikrotik template**.
- Exposes SNMP data for external monitoring tools.
- Lightweight, efficient, and built using Golang.

## Getting Started


### Prerequisites

- Go version 1.23.1 or later
- MikroTik device with REST API enabled
- SNMP monitoring tools (optional)

### Installation

1. Clone the repository:
2. Build the application:

    ```bash
    make build
    ```

3. Run the application:

    ```bash
    example: ./snmp-wrapper -l 127.0.0.1 -p 1611 -i 10 -RHost 192.168.1.1:8081 -U mikrotik -P mikrotik 
    ```

## Set Up as a Service
To set up the application as a systemd service:

```bash
make install-service
```

This will install the app as a system service, enabling it to start at boot and run in the background.

To uninstall the service:
```bash
make uninstall-service
```
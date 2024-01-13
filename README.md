# turbotanuki
![Project Logo](Logo.png)


### Your all in one http load tester


## Description

turbotanuki is a Go-based tool for load testing HTTP servers. It provides a simple command-line interface with various flags to customize your requests. Test The limits of your service, your ddos protections solution and much more.

## Features

- **Concurrent Requests:** Control the number of concurrent requests with the `-c` or `--cunnreq` flag (default: 1).

- **Total Number of Requests:** Set the total number of requests with the `-n` or `--numreq` flag (default:1).
- **Save Directives to File:** Save tanuki directives to a file with the `-s` or `--save` flag for reuse.
- **Specify URL:** Set the URL for the deed with the `-u` or `--url` flag.

- **File Input:** (pending) Utilize tanuki directives (commands) from a file with the `-f` or `--file` flag for more complex requests.

## Getting Started

### Prerequisites

- Go installed on your machine.

### Installation

Clone the repository:

```bash
git clone https://github.com/msalbrain/turbotanuki
```

Compile the project:

```bash
cd turbotanuki
go build -o tt
```

### Usage

```bash
./tt [flags]
```

Example:

```bash
./tt -u https://example.com -n 10
```

For more options, use the `-h` or `--help` flag:

```bash
./tt --help
```
 



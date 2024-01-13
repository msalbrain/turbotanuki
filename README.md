# turbotanuki
![Project Logo](Logo.png)


### Your all in one http load tester


## Description

turbotanuki is a Go-based tool for load testing HTTP servers. It provides a simple command-line interface with various flags to customize your requests. Test The limits of your service, your ddos protections solution and much more.

## Features

- **Concurrent Requests:** Control the number of concurrent requests with the `-c` or `--cunnreq` flag (default: 1).

- **Total Number of Requests:** Set the total number of requests with the `-n` or `--numreq` flag (default:1).

- **timeout**: Sets the timeout per request to 500 milliseconds. `-t` or `--timeout`

- **Save Directives to File:** Save tanuki directives to a file with the `-s` or `--save` flag for reuse.

- **Specify URL:** Set the URL for the deed with the `-u` or `--url` flags.

- **Method:** Set the method of the request to be made. `--method` or `-m`.

- **File Input:** (pending) Utilize tanuki directives (commands) from a file with the `-f` or `--file` flag for more complex requests.

- **Body:** Specifies the request body content as `{"key": "value"}` for POST requests. `-b` or `--body`

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
tt --url https://example.com/api --numreq 10 --cunnreq 5 --timeout 500 --file /path/to/directives.txt --method POST --header "Content-Type: application/json" --header "Authorization: Bearer token" --body '{"key": "value"}'
```

For more options, use the `-h` or `--help` flag:

```bash
tt --help
```
 



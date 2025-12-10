# Reverse-Scan

Perform reverse DNS lookups on huge network ranges

This utility uses the "Dispatcher/Workers" pattern discribed here :
- https://gobyexample.com/worker-pools
- https://nesv.github.io/golang/2014/02/25/worker-queues-in-go.html
- http://marcio.io/2015/07/handling-1-million-requests-per-minute-with-golang/

# Getting Started

Download the binary :

```bash
wget https://github.com/amine7536/reverse-scan/releases/download/v0.2.1/reverse-scan
chmod +x reverse-scan
```

Usage :

```bash
./reverse-scan --help
Reverse Scan

Usage:
  reverse-scan [flags]
  reverse-scan [command]

Available Commands:
  help        Help about any command
  version     Print the version number

Flags:
  -c, --cidr string     CIDR notation (e.g., 192.168.1.0/24)
  -e, --end string      ip range end
  -h, --help            help for reverse-scan
  -o, --output string   csv output file
  -s, --start string    ip range start
  -w, --workers int     number of workers (default 8)

Use "reverse-scan [command] --help" for more information about a command
```

Run with IP range:

```bash
./reverse-scan --start 37.160.0.0 --end 37.175.255.255 --output /tmp/out.csv -w 1024
2017/06/30 15:01:29 Resolving from 37.160.0.0 to 37.175.255.255
2017/06/30 15:01:29 Calculated CIDR is 37.160.0.0/12
2017/06/30 15:01:29 Number of IPs to scan: 1048576
2017/06/30 15:01:29 Starting 1024 Workers
   9s [======================================================>-------------]  81%
```

Or run with CIDR notation:

```bash
./reverse-scan --cidr 127.0.0.1/24 --output /tmp/out.csv -w 1024
2017/06/30 15:01:29 Resolving from 127.0.0.0 to 127.0.0.255
2017/06/30 15:01:29 Calculated CIDR is 127.0.0.0/24
2017/06/30 15:01:29 Number of IPs to scan: 256
2017/06/30 15:01:29 Starting 1024 Workers
   1s [===========================================================>------]  91%
```

You can specify either:
- IP range using `--start` and `--end` flags, or
- CIDR notation using `--cidr` flag

You specify the number of workers with the option `-w`, by default the utility starts with 8 workers.
You must also specify an output CSV file.

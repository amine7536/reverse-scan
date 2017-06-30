# Reverse-Scan

Perform reverse DNS lookups on huge network ranges

# Getting Started

Download the binary :

```bash
wget https://github.com/amine7536/reverse-scan/releases/download/v0.2.1/reverse-scan
chmod +x reverse-scan
```

Usage :

```bash
./reverse-scan --help
Revere Lookup

Usage:
  reverse-scan [flags]
  reverse-scan [command]

Available Commands:
  help        Help about any command
  version     Print the version number

Flags:
  -e, --end string      Range End
  -h, --help            help for reverse-scan
  -o, --output string   Output File
  -s, --start string    Range Start
  -w, --workers int     Number of Workers (default 16)

Use "reverse-scan [command] --help" for more information about a command.
```

Run :

```bash
./reverse-scan --start 37.160.0.0 --end 37.175.255.255 --output /tmp/out.csv -w 1024
2017/06/30 15:01:29 Resolving from 37.160.0.0 to 37.175.255.255
2017/06/30 15:01:29 Caluculated CIDR is 37.160.0.0/12
2017/06/30 15:01:29 Number of IPs to scan: 1048576
2017/06/30 15:01:29 Starting 1024 Workers
   9s [======================================================>-------------]  81%
```

You specify the number of workers with the option `-w`, by default the utility start with 16 workers.
You must also specify an output CSV file.
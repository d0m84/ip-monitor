# IP-Monitor

A basic Golang daemon which detects changes of the local IP or a DNS address.

If a change is detected the configured triggers (e.g. scripts) are executed with the old and the new IP address as parameters.

## Installation
1. Download and unpack appropriate release to e.g. `/usr/local/bin/ip-monitor`
2. Register and Enable as Systemd service. An example unit file can be found under `etc/` in this repository
3. Alternatively: (Enable as SysVinit service. An example init script can be found under `etc/` in this repository)
4. Configure as described below. An example configuration file can be found under `etc/` in this repository
5. Start and check journal logs (Systemd) or logfile (SysVinit)


## Configuration
The JSON config file format allows to configure several monitors with one or more triggers each.


### Main Section
| Key | Value | Comment |
| - | - | - |
| log_level | string(info\|debug) | Log level of daemon process |
| log_timestamps | bool | Timestamps in log output. Useful if running with SysVinit |
| interval | int() | Monitor interval in seconds. A too low value might cause a Service/DNS provider ban |
| http_provider | string() | URL of the HTTP IP provider to use. Must respond only with an IP address string |


### Monitors Section
| Key | Value | Comment |
| - | - | - |
| name | string() | Friendly monitor name for log output |
| domain | string() | DNS address to monitor. If empty HTTP mode is used = Local IP address |
| ip_version | string(ip4\|ip6) | Define if IPv4 or IPv6 address of domain should be monitored |
| triggers | list(string) | Path to the trigger scripts or executables. Can be one or many. |


## Modes of operation
As described above the daemon can monitor IP address changes via HTTP or DNS.

In the HTTP mode usually only local Internet IP address changes can be detected, while in DNS mode any system can be monitored which populates a DNS address.

It is expected that only **one** IP address is populated at a time. If multiple addresses can be found, a warning will be logged and only the first IP address will be taken into consideration. Depending on the DNS setup this might cause a frequent execution of the configured triggers.

### Use cases

#### HTTP mode
Useful to monitor local Internet IP address changes.
"If my local Internet IP address changes I would like to automatically perform action ..."

#### DNS mode
Useful to monitor remote IP changes.
"My router automatically updates a DynDNS address. If the address changes I would like to automatically perform action ..."


## Triggers
Triggers are executables or scripts which get executed once an IP address change is detected.
The daemon will execute those as follows:
```
/path/to/executable $OLD_IP $NEW_IP
```

This allows you to script what needs to be done on the system. Examples would be:
- iptables changes
- Configuration changes of service
- Service restarts
- Send an email
- ...

Trigger scripts **must** be executable and in case of scripts contain a Shebang.

**Example**
```bash
#!/bin/bash

echo "Old IP was: $1"
echo "New IP is: $2"
```

## Key concepts
- Allow to configure several monitors with several triggers each
- Avoid to use the local DNS cache (not supported for CNAME targets)
- Use HTTP or DNS resolvers
- Support IPv4 and IPv6 scenarios

## Disclaimer
The daemon is only tested on Debian and Ubuntu Linux.
It might work on Darwin, MacOS or other Linux distributions as well.

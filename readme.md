# Mini DNS Resolver
### Built using Go.

- read [docs](https://datatracker.ietf.org/doc/html/rfc1035) for implementation of the code.
- just put domain name in `main.go` file and run using `go run main.go` to see the output.

- for `queryString := "cses.fi"` :
```
Queueing for cses.fi to 192.203.230.10 
Queueing for cses.fi to 194.0.25.30 
Queueing for cses.fi to 194.146.106.26 
Queueing for cses.fi to 204.61.216.98 
Queueing for cses.fi to 194.0.1.14 
Queueing for cses.fi to 77.72.229.253 
Queueing for cses.fi to 194.0.11.104 
Queueing for cses.fi to 87.239.120.11 
Queueing for cses.fi to 193.166.4.1 
Querying for name server ip.
Queueing for ns1.digitalocean.com to 192.203.230.10
Queueing for ns1.digitalocean.com to 192.35.51.30 
Queueing for ns1.digitalocean.com to 173.245.59.148 
Queueing for ns2.digitalocean.com to 192.203.230.10 
Queueing for ns2.digitalocean.com to 192.35.51.30 
Queueing for ns2.digitalocean.com to 173.245.59.148 
Queueing for ns3.digitalocean.com to 192.203.230.10 
Queueing for ns3.digitalocean.com to 192.35.51.30 
Queueing for ns3.digitalocean.com to 173.245.59.148 
Queueing for cses.fi to 198.41.222.173
 
The resolved IP address is 188.166.104.231
```
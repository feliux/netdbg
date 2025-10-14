# netdbg – Feature Ideas and Roadmap

This document lists current and potential future features for the `netdbg` network debugging tool, especially focused on usability in restricted environments (containers, Kubernetes, distroless images, etc.).

---

## Core Features (Current & Planned)

### 1. Whois and Domain/IP Information
- Perform whois queries for domains and IP addresses.
- Display registration, ASN, country, and related metadata.

### 2. Traceroute (ICMP and TCP)
- Trace the network path to a destination using ICMP, UDP, or TCP.
- Show each hop, round-trip times, and optionally resolve hostnames.

### 3. Advanced DNS Query (dig-like)
- Query DNS records (A, AAAA, MX, TXT, CNAME, NS, PTR, SOA, etc.).
- Specify DNS server, protocol (UDP/TCP), and port.
- Show query/response details and timing.

### 4. Minimal HTTP(s) Client
- Perform HTTP/HTTPS requests (GET, POST, HEAD, etc.).
- Display status code, headers, response body, and timing.
- Support custom headers and request bodies.

### 5. Port Scanning
- Scan a range of TCP/UDP ports on a target host.
- Report open/closed ports and response times.

### 6. Proxy Connectivity Testing
- Test connectivity through HTTP or SOCKS proxies.
- Show connection success, handshake details, and errors.

### 7. Latency and Jitter Measurement
- Measure round-trip time (RTT) and jitter to a target (ICMP, TCP, UDP).
- Report statistics: min/avg/max RTT, packet loss, etc.

### 8. Local Network Discovery
- List local network interfaces and IP addresses.
- Optionally perform ARP scan or ping sweep on the local subnet.

### 9. Custom Traffic Generation
- Send custom TCP/UDP packets (payload, flags, etc.) for advanced testing.
- Useful for firewall, NAT, and DPI troubleshooting.

### 10. Structured Output Formats
- Support output in JSON, YAML, or table format.
- Facilitate parsing and automation in CI/CD or scripting.

### 11. Kubernetes Integration
- Display pod network information (IP, routes, DNS, etc.).
- Optionally run netdbg commands inside other pods (via `kubectl exec`).
    - **Option 1: Wrapper CLI**
      Add a subcommand (e.g., `netdbg kexec`) that wraps `kubectl exec` to run netdbg (or any command) inside a target pod.
      Example:
      ```
      netdbg kexec -n my-namespace -p mypod -- nc -a 8.8.8.8 -p 53
      ```
      This would execute:
      ```
      kubectl exec -n my-namespace mypod -- netdbg nc -a 8.8.8.8 -p 53
      ```
      The wrapper can accept namespace, pod, container, and arbitrary netdbg arguments.
    - **Option 2: Generate Command**
      Generate and print the appropriate `kubectl exec` command for the user to copy and run.
    - **Option 3: Advanced API Integration**
      Use the Kubernetes API (client-go) to upload files and execute commands in pods programmatically (requires kubeconfig and permissions).

### 12. MTU and Fragmentation Diagnostics
- Detect effective MTU to a destination (like `ping -M do -s ...`).
- Help diagnose fragmentation and path MTU issues.

---

## Prioritization Suggestions

- **Short-term:** Advanced DNS, traceroute, port scan, and HTTP client.
- **Mid-term:** Proxy testing, structured output, and local network discovery.
- **Long-term:** Kubernetes integration, custom traffic, and MTU diagnostics.

---

## Contribution

If you have ideas or want to contribute to any of these features, please open an issue or pull request!

---

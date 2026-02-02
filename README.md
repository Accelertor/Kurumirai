# Simple Chat Server

This project is a **minimal chat server** designed to be used with **netcat (`nc`) as the client**.

It is intentionally simple: no encryption, no authentication, no fancy protocols.  
The purpose is to understand how network connections, message broadcasting, and basic concurrency work at a low level.

---

##  Overview

This server allows multiple clients to connect using `netcat` and exchange messages in real time.

- Clients connect via TCP
- Messages are broadcast to all connected users
- The server handles multiple connections simultaneously
---

##  Goals

- Learn socket programming
- Understand TCP connections
- Practice concurrency
- Keep the system transparent and easy to reason about

This is a learning project, not a production chat system.

---

##  Features

- Plain TCP chat server
- Works with `netcat` as the client
- Broadcasts messages to all connected clients
- Minimal protocol (line-based text)
- No encryption, no authentication

---

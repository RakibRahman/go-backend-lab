# How backend works (request flow)

```
Browser → DNS → Server (EC2) → Firewall → Nginx → Backend → Database → Response back
```

- **DNS** — turns a domain name into an IP address, so the browser knows which machine to talk to.
- **EC2 (the server)** — the actual machine the backend and Nginx run on.
- **Firewall / Security Group** — a gatekeeper on the server. It only lets traffic through on allowed ports (80 = http, 443 = https). Everything else gets dropped before it reaches anything.
- **Nginx** — sits in front of the backend and receives the request first. Its job here is just to forward ("proxy") the request to wherever the actual backend app is listening, e.g. `localhost:3001`.
  - "localhost:3001" just means the backend process is running on the *same* machine as Nginx, listening on port 3001.
  - Nginx = receptionist. It doesn't do the actual work — it looks at the request and forwards it to the right app.
- **Backend** (Spring Boot / Node.js / Go, etc.) — this is where the real logic happens:
  - validates the request (is this input even valid?)
  - runs business logic (what should actually happen?)
  - talks to the database
  - keeps secrets safe — DB credentials, API keys, etc. never go to the browser; only the backend holds them
- **Database** — stores/retrieves the actual data the backend asked for.
- **Response** — travels back the same path in reverse: Database → Backend → Nginx → Browser.

**Why not let the browser talk to the database directly?**
Because the browser is untrusted and public — anyone can open dev tools and see/change anything sent to it. The backend is the trusted middle layer: it's the only thing allowed to hold secrets and decide what's actually allowed to happen to the data.


# Where TCP/UDP fit in

TCP/UDP sit one layer below HTTP — they're how two machines actually open a connection and move bytes, before HTTP messages enter the picture.

Mapped onto the flow above:
- **Browser → DNS** — **UDP**. One small query/response ("what's the IP for this domain?"). No need for a reliable connection; if it's lost, the browser just retries.
- **Browser → Nginx**, **Nginx → Backend**, **Backend → Database** — all **TCP**. A TCP connection (the "handshake") is opened first, then HTTP (or the Postgres wire protocol) rides on top of it. TCP guarantees bytes arrive complete, in order, and retransmits anything lost.

Rule of thumb:
- **TCP** when correctness/completeness matters more than speed (web pages, API calls, DB queries — one dropped byte and the response is garbage).
- **UDP** when speed matters more than occasionally losing a bit (DNS lookups, video/voice calls, live game state).

So in this flow, every arrow is TCP except the DNS one.

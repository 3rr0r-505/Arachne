<div align="center">

# Arachne - Concurrent Web Crawler

![net/http](https://img.shields.io/badge/net%2Fhttp-stdlib-00ADD8?logo=go&logoColor=00ADD8)
![Concurrency](https://img.shields.io/badge/Concurrency-Goroutines-00ADD8?logo=go&logoColor=00ADD8)
![x/net/html](https://img.shields.io/badge/HTML-x%2Fnet%2Fhtml-00ADD8?logo=go&logoColor=00ADD8)
![x/time/rate](https://img.shields.io/badge/RateLimit-x%2Ftime%2Frate-00ADD8?logo=go&logoColor=00ADD8)
![Output](https://img.shields.io/badge/Output-JSON-lightgrey?logo=json&logoColor=lightgrey)
![Interface](https://img.shields.io/badge/UI-Terminal-darkgreen?logo=gnubash)
![License](https://img.shields.io/badge/License-Apache%202.0-73e4bf?logo=opensourceinitiative&logoColor=73e4bf)

A concurrent web crawler built to solidify goroutines, channels, worker pools, and context-based cancellation — with a lightweight recon angle for HTB/offsec use.

</div>

---

## 🔐 Key Features

- Concurrent worker pool with configurable goroutine count
- Unbounded job queue — internal shuttler goroutine backed by a growable buffer, so a full frontier never deadlocks a worker mid-push
- Mutex-guarded visited-set — atomic check-and-mark, no duplicate crawls under concurrent load
- Domain-scope control — in-domain, subdomains, or full external following, decided explicitly per crawl
- Context-aware cancellation — Ctrl-C and `--ctx-timeout` abort in-flight requests, not just future ones
- Per-domain rate limiting via token bucket, independent budget per host
- Retry with exponential backoff on transport-level failures
- Proxy support for request interception (e.g. Caido, Burp) with scoped TLS bypass
- JS file and `<form>` extraction in a single HTML tokenizer pass
- Streamed, real-time results as pages are crawled — not batched at the end
- Self-documenting JSON export via `--out`, records the exact command that produced it
- `--max-pages` hard cap, enforced with a lock-free atomic counter

---

## 🛠️ Tech Stack

| Layer | Technology |
|---|---|
| Language | Go |
| HTTP Client | net/http (stdlib) |
| HTML Parsing | golang.org/x/net/html |
| Rate Limiting | golang.org/x/time/rate (per-domain token bucket) |
| Concurrency | goroutines, channels, sync.WaitGroup, sync.Mutex, atomic |
| Cancellation | context.Context (signal + timeout driven) |
| Output | encoding/json (stdlib) |

## 📦 Requirements

- Go 1.22+ (build only — the compiled binary has no runtime dependency)

---

## 🚀 Setup

### Option A — go install
```bash
go install github.com/3rr0r-505/arachne/cmd/arachne@latest
```
Installs the `arachne` binary to `$(go env GOPATH)/bin` — make sure that's on your `PATH`.

### Option B — Build from Source
```bash
git clone https://github.com/3rr0r-505/arachne.git
cd arachne
go build -o ./bin/arachne ./cmd/arachne
```

### Run a Crawl
```bash
./bin/arachne --url <seed-url> --workers <count> --depth <max-depth>
```

Example:

```bash
./bin/arachne --url https://example.com/ --workers 10 --depth 2 --subs
```

### 3️⃣ Save Results as JSON (Optional)
Add `--out` to save the full crawl as a self-documenting JSON file. The extension is always normalized to `.json`, and parent directories are created automatically if missing.

```bash
./bin/arachne --url https://example.com/ --workers 10 --subs --js --forms --out reports/scan1.json
```

---

## ⚙️ Flags

| Flag | Description | Default |
|---|---|---|
| `-url` | Seed URL to start crawling (required) | — |
| `-depth` | Max crawl depth (`-1` = unlimited) | `-1` |
| `-workers` | Number of concurrent worker goroutines | `10` |
| `-timeout` | Per-request timeout | `10s` |
| `-max-pages` | Hard cap on total pages crawled (`0` = unlimited) | `0` |
| `-subs` | Include subdomains as in-scope | `false` |
| `-external` | Follow external links too, not just report them | `false` |
| `-rate` | Requests/sec per domain (`0` = unlimited) | `0` |
| `-retries` | Max retry attempts on failed fetch | `2` |
| `-ctx-timeout` | Overall crawl timeout / cancellation deadline | none |
| `-forms` | Extract forms found on pages | `false` |
| `-js` | Extract JS file links | `false` |
| `-proxy` | Proxy URL to route requests through (e.g. Caido) | none |
| `-out` | Output file path — content is always JSON regardless of extension given | stdout only |

Cancel a running crawl anytime with `Ctrl+C` — in-flight requests abort immediately via `context.Context`, no hang, no second signal needed.

---

## 📄 License

This project is licensed under **Apache License 2.0**.
Free to use, modify, and learn from.

## ❤️ Support

If you like this project, consider giving it a ⭐ on GitHub!
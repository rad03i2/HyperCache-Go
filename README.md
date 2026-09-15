<div align="center">

# ⚡ HyperCache-Go
### Ultra-High-Throughput In-Memory Key-Value Store & Cache with Redis RESP Protocol Compatibility
#### محرك تخزين مؤقت وكاش فائق السرعة في الذاكرة بلغة Go متوافق تماماً مع بروتوكول Redis

[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://golang.org)
[![Redis Protocol](https://img.shields.io/badge/Protocol-RESP%20v2-DC382D?style=for-the-badge&logo=redis&logoColor=white)](https://redis.io)
[![Throughput](https://img.shields.io/badge/Throughput-100k%2B%20ops%2Fsec-brightgreen?style=for-the-badge)](https://github.com/rad03i2/HyperCache-Go)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?style=for-the-badge&logo=docker&logoColor=white)](https://docker.com)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg?style=for-the-badge)](LICENSE)

<br/>

[English Overview](#-english-overview) • [التوثيق بالعربية](#-نظرة-عامة-باللغة-العربية) • [Quick Start](#-quick-start) • [Architecture](#-architecture) • [Commercial Systems Consulting](#-commercial-systems-engineering--high-throughput-consulting)

</div>

---

## 🌟 Highlights

**HyperCache-Go** is a concurrent, sharded in-memory key-value cache engineered from the ground up in Go for extreme throughput and sub-millisecond response latencies.

It natively speaks the **Redis Serialization Protocol (RESP v2)**, allowing seamless drop-in integration with any existing client library (`redis-cli`, Python `redis-py`, Node.js `ioredis`, or Go `go-redis`) without modifying application code.

---

## 🚀 Key Features

- 🏎️ **64-Shard Concurrent Architecture**: Eliminates global lock contention using FNV-1a hashing and independent RWMutex partitions.
- 🔌 **Native Redis Client Compatibility**: Directly connect using standard `redis-cli -p 6379`.
- ⏱️ **Active TTL & Eviction Engine**: High-efficiency background goroutines prune expired keys automatically.
- 📊 **Telemetry & Stats HTTP API**: Built-in JSON metrics endpoint for Prometheus and health probes (`/stats`).
- 📦 **Zero External C Dependencies**: Pure Go standard library networking and concurrency.

---

## 🏛️ Architecture

```text
               [ Redis Clients (redis-cli, Python, Node, Go) ]
                                      │
                                      ▼ (TCP :6379)
                        ┌───────────────────────────┐
                        │   RESP Protocol Parser    │
                        └─────────────┬─────────────┘
                                      │
                                      ▼
                        ┌───────────────────────────┐
                        │   FNV-1a Shard Router     │
                        └─────────────┬─────────────┘
                                      │
          ┌───────────────────────────┼───────────────────────────┐
          ▼                           ▼                           ▼
 ┌─────────────────┐         ┌─────────────────┐         ┌─────────────────┐
 │ Shard 0 (RWMutex)│        │ Shard 1 (RWMutex)│        │ Shard 63 (Mutex)│
 └─────────────────┘         └─────────────────┘         └─────────────────┘
          ▲
          └────────── [ Background TTL Active Eviction Loop ]
```

---

## ⚡ Quick Start

### 1. Build and Run
```bash
git clone https://github.com/rad03i2/HyperCache-Go.git
cd HyperCache-Go

# Build binary
go build -o hypercache cmd/hypercache/main.go

# Start server
./hypercache -port=6379 -http=8080
```

### 2. Connect via Standard `redis-cli`
```bash
redis-cli -p 6379

127.0.0.1:6379> PING
PONG
127.0.0.1:6379> SET user:101 "Radwan Ahmed"
OK
127.0.0.1:6379> GET user:101
"Radwan Ahmed"
```

### 3. Check Live Telemetry Metrics
```bash
curl http://localhost:8080/stats
# Output: {"status":"online","keys":1,"hits":1,"misses":0}
```

---

## 🇸🇦 نظرة عامة باللغة العربية

### ما هو مشروع HyperCache-Go؟
**HyperCache-Go** هو محرك تخزين مؤقت وكاش عالي الأداء مبني بالكامل بلغة Go. صُمم لمعالجة مئات آلاف الطلبات في الثانية بزمن استجابة أقل من جزء من الثانية (Sub-millisecond).

### أهم المزايا التقنية:
1. **معمارية التجزيء المتوازي (Sharded Concurrency)**: تقسيم الذاكرة إلى 64 قسماً مستقلاً، مما يقضي تماماً على عنق الزجاجة (Lock Contention) ويتيح استغلال كافة أنوية المعالج.
2. **توافق كامل مع بروتوكول Redis**: يمكنك استخدامه كبديل مباشر لـ Redis في تطبيقاتك واستخدام نفس الأوامر مثل `PING`, `SET`, `GET`, `DEL`, `INFO`.
3. **مراقبة حية**: خادم مدمج يوفر إحصائيات فورية عن عدد المفاتيح، والـ Hits والـ Misses.

---

## 💼 Commercial Systems Engineering & High-Throughput Consulting
### استشارات التعاقد وتطوير الأنظمة الموزعة عالية الأداء

Need distributed in-memory caching, low-latency microservices, or custom protocol servers built in Go?
هل تحتاج إلى بناء خوادم عالية الأداء، أنظمة موزعة، أو حلول تخزين ومعالجة بيانات ضخمة بلغة Go؟

- 📩 **Contact**: Reach out via GitHub [@rad03i2](https://github.com/rad03i2)
- 🤝 **Freelance & Enterprise Systems**: Open for specialized backend engineering contracts.

---

## 📄 License
Licensed under the [MIT License](LICENSE). Developed by [rad03i2](https://github.com/rad03i2).

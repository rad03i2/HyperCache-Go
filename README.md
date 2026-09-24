# HyperCache-Go

A small, dependency-free, Redis RESP2-compatible in-memory cache server written in Go.

> **Status:** functional cache server intended for local development, learning, sidecars, and controlled internal workloads. It is not a full Redis replacement.

## English

### Why it exists
HyperCache-Go demonstrates a practical concurrent cache with a network protocol real Redis clients can speak. It favors a small auditable codebase over broad Redis compatibility.

### Features
- Sharded concurrent map using FNV-1a routing and per-shard locks.
- Binary-safe values with defensive copies on write/read.
- TTL expiration with active background cleanup and lazy expiration.
- RESP2 TCP server with `PING`, `SET`, `GET`, `DEL`, `EXISTS`, `TTL`, `DBSIZE`, `INFO`, and `QUIT`.
- `SET key value EX seconds` and `SET key value PX milliseconds`.
- Atomic hit/miss/set/delete counters.
- Graceful cache cleanup-loop shutdown through `Close()`.
- No runtime dependencies; Docker image runs as a non-root user.

### Requirements
Go 1.22+; optionally Docker and `redis-cli`.

### Install & run
```bash
git clone https://github.com/rad03i2/HyperCache-Go.git
cd HyperCache-Go
go run ./cmd/hypercache -addr=127.0.0.1:6379
```

Options:
```text
-addr     TCP listen address (default 127.0.0.1:6379)
-shards   cache shard count (default 64)
-cleanup  expired-key cleanup interval (default 5s)
```

### Usage
```bash
redis-cli -p 6379 PING
redis-cli -p 6379 SET greeting "hello"
redis-cli -p 6379 GET greeting
redis-cli -p 6379 SET session abc EX 60
redis-cli -p 6379 TTL session
redis-cli -p 6379 INFO
```

Docker:
```bash
docker build -t hypercache-go .
docker run --rm -p 6379:6379 hypercache-go
```

### Project structure
```text
cmd/hypercache/main.go   executable entrypoint
pkg/cache/               sharded cache, TTL and metrics
pkg/resp/                RESP2 parser/writer
pkg/server/              TCP command server
.github/workflows/ci.yml cross-platform CI
```

### Testing
```bash
go vet ./...
go test -race ./...
go build ./cmd/hypercache
```
CI runs these checks on Linux, Windows and macOS with Go 1.22 and 1.23.

### Preview guidance
This is a TCP service rather than a graphical app. For a project preview, capture a terminal with two panes: the HyperCache startup log on one side and `redis-cli` commands (`SET`, `GET`, `TTL`, `INFO`) on the other.

### Configuration
Configuration is intentionally CLI-only. No environment variables, credentials, or external services are required.

### Security & privacy
HyperCache does not provide authentication, TLS, ACLs, encryption at rest, persistence, or network isolation. The default bind address is loopback for safer local use. If you expose it beyond a trusted machine/network, put it behind appropriate network and transport controls. Cached values live only in process memory.

### Limitations
- Data is lost when the process stops.
- RESP2 support is deliberately partial; unsupported Redis commands return an error.
- No replication, clustering, transactions, pub/sub, Lua, persistence, authentication, TLS, LRU/LFU memory cap, or Redis modules.
- TTL precision exposed by `TTL` is seconds.
- Performance claims require workload-specific benchmarks; this project makes no fixed throughput claim.

### Optional roadmap
Bounded-memory eviction, graceful TCP server shutdown, RESP parser limits, benchmarks, and opt-in authenticated/TLS deployment modes are reasonable future additions.

### Contributing
See [CONTRIBUTING.md](CONTRIBUTING.md). Please include tests for behavior changes and keep the dependency footprint minimal.

### License
MIT — see [LICENSE](LICENSE).

### Author
**Radwan Abdulhadi Ahmed**  
**رضوان عبدالهادي أحمد**  
GitHub: **@rad03i2**

---

## العربية

### نظرة عامة
**HyperCache-Go** خادم تخزين مؤقت صغير يعمل في الذاكرة ومكتوب بلغة Go دون اعتماديات تشغيل خارجية. يستخدم RESP2 بحيث يمكن التعامل معه عبر `redis-cli` والعملاء المتوافقين مع Redis للأوامر المدعومة.

### لماذا يوجد المشروع؟
الهدف هو تقديم تطبيق عملي وواضح لمخزن مفاتيح/قيم متزامن مع بروتوكول شبكي حقيقي، مع إبقاء قاعدة الكود صغيرة وقابلة للمراجعة بدل الادعاء بأنه بديل كامل لـ Redis.

### المزايا
- تقسيم الذاكرة إلى Shards لتقليل التنافس على الأقفال.
- نسخ القيم عند الإدخال والإخراج لمنع تعديل الذاكرة الداخلية من الخارج.
- انتهاء صلاحية TTL مع تنظيف خلفي وفحص عند القراءة.
- الأوامر: `PING`, `SET`, `GET`, `DEL`, `EXISTS`, `TTL`, `DBSIZE`, `INFO`, `QUIT`.
- دعم `EX` بالثواني و`PX` بالميلي ثانية مع `SET`.
- عدادات Hits وMisses وSets وDeletes ذرّية.
- إيقاف حلقة التنظيف الخلفية عبر `Close()`.
- صورة Docker تعمل بمستخدم غير root.

### المتطلبات والتثبيت
يتطلب Go 1.22 أو أحدث. للتشغيل:
```bash
git clone https://github.com/rad03i2/HyperCache-Go.git
cd HyperCache-Go
go run ./cmd/hypercache -addr=127.0.0.1:6379
```

### مثال استخدام
```bash
redis-cli -p 6379 SET name "Radwan"
redis-cli -p 6379 GET name
redis-cli -p 6379 SET temp value EX 30
redis-cli -p 6379 TTL temp
redis-cli -p 6379 INFO
```

### الاختبارات
```bash
go vet ./...
go test -race ./...
go build ./cmd/hypercache
```
ويشغّل GitHub Actions هذه الفحوصات على Linux وWindows وmacOS.

### بنية المشروع
`cmd/hypercache` يحتوي نقطة التشغيل، و`pkg/cache` محرك الكاش، و`pkg/resp` محلل RESP2، و`pkg/server` خادم TCP.

### الإعداد
لا يحتاج المشروع ملف `.env` أو مفاتيح API أو أسرارًا. الإعدادات المتاحة حاليًا هي عنوان الاستماع وعدد الـShards وفاصل تنظيف المفاتيح المنتهية عبر خيارات سطر الأوامر.

### الخصوصية والأمان
لا يرسل البرنامج البيانات إلى خدمات خارجية. لكنه لا يوفر مصادقة أو TLS أو ACL أو تشفيرًا أو تخزينًا دائمًا. لذلك يرتبط افتراضيًا بـlocalhost، ولا ينبغي كشفه للإنترنت مباشرة.

### القيود
البيانات مؤقتة وتُفقد عند إيقاف البرنامج، وتوافق Redis جزئي فقط، ولا توجد replication أو clustering أو transactions أو pub/sub أو persistence أو حد ذاكرة/LRU. كما أن أي أرقام أداء يجب قياسها حسب الجهاز والحمل الفعلي.

### التطوير الاختياري
يمكن مستقبلًا إضافة حد للذاكرة وسياسة eviction، وإيقاف منظم لخادم TCP، وحدود أكثر صرامة لمحلل RESP، وbenchmarks موثقة.

### المساهمة والترخيص
راجع [CONTRIBUTING.md](CONTRIBUTING.md). المشروع مرخص وفق MIT؛ راجع [LICENSE](LICENSE).

### المؤلف
**Radwan Abdulhadi Ahmed**  
**رضوان عبدالهادي أحمد**  
GitHub: **@rad03i2**

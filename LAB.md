# Lab: Correlating Metrics to Logs in Grafana

A self-paced lab for learning Grafana **data links** and **correlations** — the features
that turn "I see a spike on a dashboard" into "here are the log lines behind that spike."

Everything runs locally in Docker. No company infrastructure is touched.

---

## Why this lab exists

Imagine you're investigating a real production dashboard, the diagnosis took an hour of manual work: read a dashboard, 
note the timestamps of the spikes, then hand-match those timestamps against `kubectl logs` output from two
different pods to discover which client was actually calling in.

Every one of those steps could have been a click, if metrics-to-logs correlation had been
configured. The point of this lab is to understand that feature well enough to implement yourself or to ask
platform admins for it with a specific, informed request rather than a vague one.

The lab reproduces the same shape of problem on a small service you own.

---

## What you will end up understanding

1. Why jumping from a metric to its logs is **not** automatic, and what has to line up.
2. The difference between a **Loki label** and a **parsed field**, and why cardinality
   decides which one you get.
3. Two mechanisms — panel **data links** and datasource **correlations** — and when each fits.
4. What that config looks like **as code**, which is what admins actually deploy.
5. Where correlation genuinely stops helping, so you don't oversell it.

---

## Prerequisites

- Docker Desktop running (the stack needs roughly 2 GB free).
- Ports free: `3000` (Grafana), `9090` (Prometheus), `3100` (Loki), `8080` (app),
  `5432`, `5672`, `15672`, `12345`.
- No Go toolchain needed. The app runs in a container for this lab.

---

## Part 1 — Boot the stack

```bash
make
```

That single command is the whole setup: it creates `.env` from `.env.sample`, builds
the image, starts every service, and waits while the database is migrated and seeded
before it returns. The `employees` table is created by the app's own migrations at
startup, and the `seed` job waits for `employee-service` to be healthy before it runs,
so seeding order is handled for you.

```bash
make ps                                # every service, with health
docker compose logs seed               # expect: DELETE 0 / INSERT 0 2
```

(`make` is a thin wrapper over `docker compose up -d --build`; use either.)

### What is now running

| Service | Port | Role |
|---|---|---|
| `employee-service` | 8080 | The app under observation. Emits Prometheus metrics + zap JSON logs. |
| `employee` | 5432 | Postgres. |
| `rabbitmq` | 5672 / 15672 | Required at boot — the app calls `log.Fatal` if it can't consume. |
| `prometheus` | 9090 | Scrapes `employee-service:8080/metrics` every 5s. |
| `loki` | 3100 | Log store. |
| `alloy` | 12345 | Log shipper: Docker socket → parse JSON → push to Loki. |
| `grafana` | 3000 | Anonymous admin access, no login needed. |
| `tempo` | 3200 / 4318 | Trace store. Receives OTLP spans and derives RED metrics + a service graph. |
| `verification-service` | 8090 | Downstream hop called by `employee-service` on every lookup, purely so traces span two processes. |
| `loadgen` | — | Generates a 200/400/404 mix continuously. |
| `seed` | — | One-shot job: seeds the `employees` table, then exits. Re-run with `make seed`. |

The original ELK stack (3 Elasticsearch nodes, Filebeat, Logstash, Kibana) is preserved
but **off by default** — it is heavy and Loki replaces it here. Start it with
`docker compose --profile elk up -d` if you ever want to compare the two approaches.

> **Where's the trace node graph?** This lab is about metrics↔logs correlation, but
> tracing rides along for free since `employee-service` is already instrumented. Don't
> look for a node graph in **Drilldown → Traces** — Grafana 12's Traces Drilldown app
> only shows a text "Service & Operation" tree under its **Service structure** tab, not
> a visual graph. The actual node-and-edges diagram lives in classic **Explore**: pick
> the **Tempo** datasource, switch the query type to **Service Graph**, and run it. You
> should see `user → employee-service → verification-service` with live request rates
> and durations on each edge.

### Checkpoint 1 — is the data actually flowing?

Run these before touching Grafana. If any fails, see [Troubleshooting](#troubleshooting).

```bash
# Metrics: expect counts for 200, 400 and 404
curl -s 'http://localhost:9090/api/v1/query?query=response_status' | jq -r \
  '.data.result[] | "\(.metric.status) -> \(.value[1])"'

# Logs: expect ["200","400","404"]
curl -s 'http://localhost:3100/loki/api/v1/label/status/values' | jq -c '.data'

# Grafana: expect both datasources with pinned uids
curl -s http://localhost:3000/api/datasources | jq -r '.[] | "\(.type) \(.uid)"'
```

Open <http://localhost:3000> → dashboard **Employee Service - Correlation Lab**
(folder *Observability Lab*). You should see traffic on three panels and live logs
on the fourth.

---

## Part 2 — The concept, before you click anything

### Why this isn't automatic

Grafana does not know that your Prometheus metrics and your Loki logs describe the same
requests. Nothing in either system says so. To bridge them you need a **join key**: a
piece of data present on *both* sides, with the *same name and the same values*.

Look at what this app emits.

**Metrics** (`pkg/middleware/metric.go`):

```
http_requests_total{path="/employees/{employee_id}"}
response_status{status="404"}
http_response_time_seconds{path="/employees/{employee_id}"}
```

**Logs** (`pkg/middleware/logger.go`, zap JSON):

```json
{"level":"info","status":404,"method":"GET",
 "endpoint":"/employees/00000000-0000-4000-8000-000000000000",
 "duration_ms":0.001222834}
```

Compare the two carefully:

| | Metrics | Logs | Usable as a join key? |
|---|---|---|---|
| status | `status="404"` | `"status":404` | **Yes** — same name, same values |
| method | *(absent)* | `"method":"GET"` | No — metrics don't carry it |
| route | `path="/employees/{employee_id}"` | `endpoint="/employees/00000000-..."` | **No** — different name *and* different values |

That last row is the important one. The metric carries the **route template**; the log
carries the **concrete URL**. They describe the same request and still cannot be matched
by equality. This is not sloppiness — it is deliberate, and the reason is cardinality.

### Labels vs parsed fields, and why cardinality decides

A Prometheus label or a Loki label is an **index key**. Every distinct combination of
label values creates a separate time series (Prometheus) or stream (Loki). Put a UUID in
a label and you get one series per UUID — millions of them. That is the classic way to
melt a metrics backend.

So the rule both systems follow:

- **Low cardinality → make it a label.** `status` has ~5 possible values ever. Cheap to
  index, fast to filter, safe as a join key.
- **High cardinality → keep it in the payload.** The concrete URL, the user id, the trace
  id. Stored, searchable, but *not* indexed.

Open `res/config/alloy.alloy` and find where this decision is made:

```alloy
stage.json {
  expressions = { level = "level", status = "status", method = "method",
                  endpoint = "endpoint", msg = "msg" }
}

stage.labels {
  values = { level = "", status = "", method = "" }   // endpoint deliberately absent
}
```

`stage.json` *parses* five fields. `stage.labels` *promotes only three* to Loki labels.
`endpoint` is parsed but never promoted, precisely because it contains a UUID.

Verify the consequence yourself:

```bash
# `status` IS a label - it appears in the label index
curl -s 'http://localhost:3100/loki/api/v1/labels' | jq -c '.data'

# `endpoint` is NOT there, yet you can still query on it via `| json`
curl -s -G 'http://localhost:3100/loki/api/v1/query_range' \
  --data-urlencode 'query={service="employee-service"} | json | endpoint=~".*not-a-valid.*"' \
  --data-urlencode 'limit=1' | jq -r '.data.result | length'
```

**The takeaway that transfers to work:** a correlation gets you to *the right pod, status
and second*. It does not get you to *this exact request* unless the metric happens to
carry an identifier — and usually it must not, for cardinality reasons. Knowing this
distinction is what separates a useful request to your admins from an impossible one.

---

## Part 3 — Exercise 1: a panel data link

Data links are the **older, panel-scoped** mechanism. One panel, one link. Good when the
jump is specific to a single chart.

The reliable way to build one is to construct the target query first, then templatize its
URL — do not hand-write Explore URLs from memory, the format changes between versions.

**Step 1 — build the query in Explore.**

1. Go to **Explore** (compass icon) → datasource **Loki**.
2. Switch to the code editor and enter:
   ```logql
   {service="employee-service", status="404"} | json
   ```
3. Run it. Confirm you get log lines.

**Step 2 — capture the URL.** Copy the whole thing from the browser address bar. It
encodes the datasource, the query and the time range. This is your template's skeleton,
and it is *known-correct for your Grafana version* because Grafana just generated it.

**Step 3 — templatize it.** In that URL, replace:

- the hardcoded `404` with `${__field.labels.status}`
- the `from`/`to` timestamps with `${__from}` and `${__to}`

Those are Grafana's built-in interpolation variables:

| Variable | Resolves to |
|---|---|
| `${__field.labels.status}` | the `status` label of the series you clicked |
| `${__from}` / `${__to}` | the panel's current time range, in epoch ms |
| `${__value.time}` | the timestamp of the exact point clicked |
| `${__series.name}` | the series display name |

**Step 4 — attach it to the panel.**

1. Dashboard → panel **Response status (rate/s)** → **Edit**.
2. Scroll the right-hand options to **Data links** → **Add link**.
3. Title: `View logs for this status`. URL: your templatized URL.
4. **Save dashboard.**

**Step 5 — use it.** Back on the dashboard, **click a point** on the `404` series (click,
not hover). A menu appears with your link. Follow it: Explore opens, filtered to
`status="404"`, over the same window.

You have just replaced the manual timestamp-matching that a real production incident
investigation would otherwise need.

> **A gotcha you will hit.** In the Explore results, expand a log line. You will see both
> `status` *and* `status_extracted`. When `| json` parses a field whose name collides with
> an existing label, Loki suffixes the parsed copy with `_extracted` rather than
> overwriting the label. Same for `level_extracted` and `method_extracted`. Filter on the
> **label** (`status="404"`, indexed and fast); use `_extracted` only if you need the
> parsed value specifically.

---

## Part 4 — Exercise 2: a correlation

Correlations are the **modern, datasource-scoped** mechanism. Define once, and it applies
to *every* panel using that datasource, across all dashboards. This is what you actually
want at company scale — nobody wants to hand-edit a link on 40 panels.

1. **Administration → Plugins and data → Correlations** → **Add correlation**.
   (Enabled here via `GF_FEATURE_TOGGLES_ENABLE=correlations`.)
2. **Source**: datasource `Prometheus`, results field **`status`**.
   This is the field whose clicked value gets passed along.
3. **Target**: datasource `Loki`, query:
   ```logql
   {service="employee-service", status="${status}"} | json
   ```
   `${status}` interpolates from the source field you just named.
4. Label: `View logs for this status`. Save.

Now go to **any** panel backed by Prometheus that exposes a `status` label and click a
point — the correlation appears without per-panel setup. Compare that with Exercise 1:
same outcome, but defined once instead of per panel.

**Why the target query is written this way.** Because Alloy promoted `status` to a real
label, `{status="${status}"}` is an *indexed stream selector* — Loki narrows to matching
streams before reading any data. Had `status` stayed a parsed field only, you would need:

```logql
{service="employee-service"} | json | status="${status}"
```

That is still correct, but it must read and parse every line in the service's streams
before filtering. On a busy production service that difference is large. **This is the
concrete argument for asking that a small set of low-cardinality fields be promoted to
labels at ingest time** — it is a shipper-config request, not a Grafana one, and it is
the kind of detail that makes an admin request actionable.

---

## Part 5 — Exercise 3: make it fail, and understand why

Try to correlate on the route instead of the status.

1. Add a correlation with source field **`path`** and target
   `{service="employee-service"} | json | endpoint="${path}"`.
2. Click a point on the **Requests by path** panel and follow it.

**It returns nothing.** `${path}` interpolates to `/employees/{employee_id}` — the literal
route template, braces and all — while every log line's `endpoint` holds a concrete URL
like `/employees/bc4cd1a1-...`. Equality can never match.

Make it work by matching on shape rather than value — replace the template's variable
segment with a regex:

```logql
{service="employee-service"} | json | endpoint=~"/employees/.*"
```

That returns lines, but note what you gave up: it now matches *every* employee request,
not the ones behind the point you clicked. **You cannot get from a route-template metric
back to one specific request.** The information was discarded at instrumentation time, on
purpose. Knowing this in advance means you won't promise your team something the data
model can't deliver.

---

## Part 6 — Exercise 4: correlate a real incident (500s)

So far only 200/400/404. Real correlation earns its keep on errors. Break the database:

```bash
docker compose stop employee
```

Watch the **Response status** panel. Within ~30s a `500` series appears (the handler's
`default:` branch turns DB errors into 500s). Now click it and follow your correlation.

You land on the error log lines — including the zap `stack` field, since
`pkg/middleware/logger.go` attaches `debug.Stack()` for any status ≥ 500. That is the
whole workflow: *see the spike → one click → read the stack trace.*

Restore:

```bash
docker compose start employee
```

> The app recovers on its own — `restart: unless-stopped` plus Postgres reconnection. If
> it doesn't, `docker compose restart employee-service`.

---

## Part 7 — Exercise 5: correlations as code

Your admins will not configure this by clicking. Grafana is managed as code, so hand them
a file.

First, get the **authoritative** payload for your Grafana version — don't trust anyone's
remembered schema, including this document's:

```bash
curl -s http://localhost:3000/api/datasources/correlations | jq
```

That is exactly what your hand-built correlation looks like to the API. Compare it with
the annotated starting point in
[`res/provisioning-examples/correlations.yaml`](res/provisioning-examples/correlations.yaml)
and reshape as needed.

To test it as real provisioning, mount that directory into Grafana by adding to the
`grafana` service in `docker-compose.yml`:

```yaml
      - ./res/provisioning-examples:/etc/grafana/provisioning/correlations:ro
```

then `docker compose up -d grafana`.

**Expect a surprise:** provisioned correlations are **read-only in the UI**. That is why
this lab has you build one by hand first — and it's worth knowing before you ask admins to
provision one, because it changes who can iterate on it afterwards. Delete your
hand-made duplicate to avoid two identical entries.

---

## Part 8 — Exercise 6: fix a real instrumentation bug

The dashboard's **p95 latency** panel reads ~0 forever. That is not a lab prop — it's a
genuine bug in this repo, and a good example of a metric that looks fine and is worthless.

Confirm it:

```bash
curl -s 'http://localhost:9090/api/v1/query?query=http_response_time_seconds_sum/http_response_time_seconds_count' \
  | jq -r '.data.result[].value[1]'
```

That prints roughly `0.000004` — **4 microseconds**. Meanwhile the logs report the truth
for the same requests:

```bash
curl -s -G 'http://localhost:3100/loki/api/v1/query_range' \
  --data-urlencode 'query={service="employee-service"} | json | line_format "{{.human_readable_duration}}"' \
  --data-urlencode 'limit=5' | jq -r '.data.result[].values[][1]'
```

Several hundred microseconds up to about a millisecond (e.g. `928.208µs`) — roughly
**200x** the histogram's claim.

The cause is in `pkg/middleware/metric.go`:

```go
err := before(w, r)                                              // handler runs HERE

route := mux.CurrentRoute(r)
path, _ := route.GetPathTemplate()

timer := prometheus.NewTimer(httpDurationMetric.WithLabelValues(path))  // starts AFTER
statusCode := v.StatusCode
responseStatusMetric.WithLabelValues(strconv.Itoa(statusCode)).Inc()
totalRequestsMetric.WithLabelValues(path).Inc()
timer.ObserveDuration()                                          // measures bookkeeping
```

The timer starts *after* the handler already ran, so it measures two counter increments
rather than the request. Every observation lands in the smallest bucket — check
`http_response_time_seconds_bucket` and you'll see `le="0.005"` already holds the full
count.

**Fix it** — start the timer before the handler and let `defer` close it:

```go
route := mux.CurrentRoute(r)
path, _ := route.GetPathTemplate()

timer := prometheus.NewTimer(httpDurationMetric.WithLabelValues(path))
defer timer.ObserveDuration()

err := before(w, r)
```

Then `docker compose up -d --build employee-service`, wait a minute, and watch the p95
panel show real values in the low milliseconds.

**Why this belongs in a correlation lab:** the logs were right the whole time. Correlation
is also how you *catch* a lying metric — you compare it against the ground truth sitting
next to it. A dashboard alone would never have told you.

---

## Part 9 — What to take to your admins

You now have specifics rather than "can we make dashboards clickable." Concretely:

**1. Is there a Loki (or equivalent) datasource for these clusters at all?**
Some clusters have no log-shipping stack at all — nothing to correlate *to*.
That is question zero; everything else depends on it.

**2. Which fields are promoted to labels at ingest?**
The correlation is only fast if the join key is an indexed label. Ask what the shipper
promotes today, and request a small, named set of low-cardinality fields — for an
internal gRPC service, something like `grpc_method` and `grpc_code`. Bring the
cardinality argument with you: you are asking for a handful of values, not UUIDs.

**3. Correlations, not per-panel data links.**
Datasource-scoped, so it works across every dashboard without per-panel edits.

**4. Provisioned, with the read-only trade-off understood.**
Provisioned correlations can't be edited in the UI. Decide deliberately who iterates.

**5. Be explicit about the limit.**
Correlation gets you to the right pod/status/second — not to one specific request, unless
a shared identifier (e.g. a trace id) exists on both sides. Say this up front so nobody
expects request-level tracing from it. If request-level *is* the goal, that's a different
ask: trace exemplars linking Prometheus to a tracing backend.

---

## Troubleshooting

**`employee-service` is unhealthy on boot.**
Check `docker compose logs employee-service`. If it's dialing `[::1]:5672`, RabbitMQ env
vars are wrong. `ardanlabs/conf` splits camelCase, so `RabbitMQ` becomes `RABBIT_MQ` —
the var is `API_RABBIT_MQ_HOST`, not `API_RABBITMQ_HOST`. Authoritative list:
`docker compose run --rm --no-deps employee-service ./main --help`.

**Loki has no `status` label.**
Alloy isn't parsing. `docker compose logs alloy`. Confirm the Docker socket is mounted and
that Docker Desktop allows it. Verify with
`curl -s localhost:3100/loki/api/v1/labels | jq`.

**Everything returns 404.**
The seed didn't run, or ran before migrations. Re-run `make seed` and
expect `INSERT 0 2`.

**The p95 panel is flat at zero.**
Expected — that's Exercise 6.

**Data link opens Explore but with no query.**
Your URL template is malformed. Rebuild it: make the query in Explore, copy the address
bar, then templatize. Don't hand-write it.

**Ports already in use.**
`docker compose down` here, or remap the host side of the `ports:` entries.

## Teardown

```bash
make stop     # stop, keep containers and data
make down     # remove containers, keep data
make reset    # delete volumes and start clean
```

---

## Reference

| File | Role |
|---|---|
| `docker-compose.yml` | Whole stack; ELK/test behind profiles, tracing (Tempo, verification-service) on by default |
| `res/config/tempo.yaml` | Tempo config — receivers, storage, metrics-generator processors |
| `res/config/alloy.alloy` | Log pipeline — **where labels vs fields is decided** |
| `res/config/loki.yaml` | Minimal single-binary Loki |
| `res/config/prometheus.yml` | Scrape config |
| `res/config/loadgen.sh` | Traffic generator (GET only, on purpose) |
| `res/provisioning/datasources/datasources.yaml` | Pinned datasource uids |
| `res/provisioning/dashboards/employee-service.json` | The lab dashboard |
| `res/provisioning-examples/correlations.yaml` | Correlation as code (not auto-loaded) |
| `pkg/middleware/metric.go` | Metrics — and Exercise 6's bug |
| `pkg/middleware/logger.go` | zap JSON logs, teed to file + stdout |

**A note on `loadgen.sh`:** it only issues `GET`s. The `PATCH` route publishes to RabbitMQ,
and `pkg/adapter/rabbitmq.go` calls `log.Fatal` if publishing fails — which would kill the
service mid-lab. Worth noticing as its own lesson in how an error-handling choice
constrains what you can safely test.

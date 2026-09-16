# Changelog

All notable changes to this repo are documented here, starting from when it
was pulled back out of storage (originally created 2021, ~4-5 years old at
time of revival). Format loosely follows [Keep a Changelog](https://keepachangelog.com/).

## [Unreleased]

### Fixed
- **Traces stopped appearing in Grafana ("Datasource was not found" in the
  Traces Drilldown app).** Root cause was three compounding issues, not a
  single Grafana version bump:
  1. Jaeger (the only trace backend in the stack) was gated behind the
     `--profile tracing` opt-in flag and never started by default.
  2. `cmd/employee/main.go` hardcoded the exporter endpoint to
     `http://localhost:14268/api/traces`. Since `employee-service` runs in
     its own container, `localhost` never reached the `jaeger` container -
     this must have worked previously only when the service ran directly on
     the host.
  3. Grafana had no trace datasource provisioned at all (only Prometheus and
     Loki), and Grafana 12's Traces Drilldown app requires a Tempo-compatible
     datasource specifically - the old Jaeger-native node graph in Explore
     isn't wired into that app.
- **Traces Drilldown showed "TraceQL metrics not configured" / "localblocks
  processor not found"** even after Tempo was wired up and traces were
  landing fine. The app's span-rate/breakdown panels need the `local-blocks`
  metrics-generator processor specifically - `service-graphs` and
  `span-metrics` alone aren't enough. Added it to `res/config/tempo.yaml`.

### Added
- **`verification-service`** (`cmd/verification`), a small standalone HTTP
  service that `employee-service` calls on every lookup. With only
  `employee-service` instrumented, Tempo's service graph only ever has one
  real node - the caller (loadgen/curl) doesn't propagate trace context, so
  it shows up as a "virtual" `user` node. `employee-service` now calls
  `verification-service` with an `otelhttp`-instrumented client
  (`pkg/handler/employee.go`), which injects the `traceparent` header, so the
  same trace id spans both processes. That gives a real three-node service
  graph (`user → employee-service → verification-service`) and a waterfall
  with more than one span to correlate across.

### Changed
- Replaced Jaeger with [Tempo](https://grafana.com/oss/tempo/) as the trace
  backend (`docker-compose.yml`, `res/config/tempo.yaml`). Tempo's
  metrics-generator derives RED metrics and a service graph from incoming
  spans and remote_writes them to Prometheus, which powers the "Service
  structure" tab and node graph in Traces Drilldown.
- Swapped the OTel exporter from the deprecated `otel/exporters/jaeger`
  package to `otlptracehttp` (`pkg/trace/provider.go`). The Jaeger exporter
  has been removed from recent OpenTelemetry Go SDK releases.
- Made the trace exporter endpoint and disabled flag configurable via
  `API_TRACE_OTLP_ENDPOINT` / `API_TRACE_DISABLED` (`pkg/config/config.go`)
  instead of hardcoded in `main.go`.
- Added a `Tempo` datasource to Grafana provisioning
  (`res/provisioning/datasources/datasources.yaml`) with trace-to-logs
  (Loki), trace-to-metrics and service-map correlation to Prometheus, and
  node graph enabled.
- Enabled Prometheus's remote-write receiver
  (`--web.enable-remote-write-receiver`) so Tempo's metrics-generator can
  push service-graph metrics to it.
- Bumped the OpenTelemetry Go SDK and related deps (`go.opentelemetry.io/otel`
  and friends v1.3.0 → v1.46.0) to pull in `otlptracehttp`.
- Tempo now starts by default with the rest of the lab stack; the
  `--profile tracing` opt-in profile (Jaeger) was removed.

use anyhow::Result;
use opentelemetry::KeyValue;
use opentelemetry::metrics::{Counter, Histogram, Meter};
use opentelemetry::trace::TracerProvider;
use opentelemetry_appender_tracing::layer::OpenTelemetryTracingBridge;
use opentelemetry_otlp::WithExportConfig;
use opentelemetry_sdk::{Resource, logs, metrics, trace};
use std::env;
use std::sync::OnceLock;
use tracing::info;
use tracing_opentelemetry::OpenTelemetryLayer;
use tracing_subscriber::{
    EnvFilter, Layer,
    fmt::{self, format::FmtSpan},
    layer::SubscriberExt,
    util::SubscriberInitExt,
};

pub struct ObservabilityConfig {
    pub log_level: String,
    pub json_logs: bool,
    pub otlp_endpoint: String,
    pub service_name: String,
}

impl ObservabilityConfig {
    pub fn from_env() -> Self {
        Self {
            log_level: env::var("LOG_LEVEL").unwrap_or_else(|_| "info".to_string()),
            json_logs: env::var("JSON_LOGS")
                .unwrap_or_else(|_| "false".to_string())
                .parse()
                .unwrap_or(false),
            otlp_endpoint: env::var("OTEL_EXPORTER_OTLP_ENDPOINT")
                .unwrap_or_else(|_| "http://otel-collector:4317".to_string()),
            service_name: env::var("OTEL_SERVICE_NAME")
                .unwrap_or_else(|_| "gateway-service".to_string()),
        }
    }
}

pub async fn init_observability(cfg: &ObservabilityConfig) -> Result<()> {
    let env_filter = EnvFilter::try_from_default_env()
        .or_else(|_| EnvFilter::try_new(&cfg.log_level))
        .unwrap_or_else(|_| EnvFilter::new("info"));

    let fmt_layer = if cfg.json_logs {
        fmt::layer()
            .json()
            .with_target(false)
            .with_span_events(FmtSpan::CLOSE)
            .with_thread_ids(true)
            .boxed()
    } else {
        fmt::layer()
            .with_target(true)
            .with_span_events(FmtSpan::CLOSE)
            .compact()
            .boxed()
    };

    // Initialize OpenTelemetry OTLP exporter
    let resource = Resource::new(vec![
        KeyValue::new("service.name", cfg.service_name.clone()),
        KeyValue::new("service.version", env!("CARGO_PKG_VERSION")),
    ]);

    // Initialize metrics
    let metrics_exporter = opentelemetry_otlp::new_exporter()
        .tonic()
        .with_endpoint(&cfg.otlp_endpoint)
        .build_metrics_exporter(Box::new(metrics::reader::DefaultTemporalitySelector::new()))?;

    let meter_provider = metrics::SdkMeterProvider::builder()
        .with_resource(resource.clone())
        .with_reader(
            metrics::PeriodicReader::builder(metrics_exporter, opentelemetry_sdk::runtime::Tokio)
                .with_interval(std::time::Duration::from_secs(5))
                .build(),
        )
        .build();

    opentelemetry::global::set_meter_provider(meter_provider);

    // Initialize tracing
    let trace_exporter = opentelemetry_otlp::new_exporter()
        .tonic()
        .with_endpoint(&cfg.otlp_endpoint)
        .build_span_exporter()?;

    let tracer_provider = trace::TracerProvider::builder()
        .with_batch_exporter(trace_exporter, opentelemetry_sdk::runtime::Tokio)
        .with_config(trace::Config::default().with_resource(resource.clone()))
        .build();

    let tracer = tracer_provider.tracer("aviation-gateway");
    let otel_layer = OpenTelemetryLayer::new(tracer);

    // Initialize log export
    let log_exporter = opentelemetry_otlp::new_exporter()
        .tonic()
        .with_endpoint(&cfg.otlp_endpoint)
        .build_log_exporter()?;

    let log_provider = logs::LoggerProvider::builder()
        .with_resource(resource.clone())
        .with_batch_exporter(log_exporter, opentelemetry_sdk::runtime::Tokio)
        .build();

    let otel_log_layer = OpenTelemetryTracingBridge::new(&log_provider);

    tracing_subscriber::registry()
        .with(env_filter)
        .with(fmt_layer)
        .with(otel_layer)
        .with(otel_log_layer)
        .init();

    // Initialize metrics after setting up the meter provider
    init_metrics();

    info!("Observability stack with OTLP initialized successfully");
    Ok(())
}

pub async fn shutdown_observability() {
    info!("Shutting down observability");
    opentelemetry::global::shutdown_tracer_provider();
}

static METER: OnceLock<Meter> = OnceLock::new();
static REQUEST_COUNTER: OnceLock<Counter<u64>> = OnceLock::new();
static REQUEST_DURATION: OnceLock<Histogram<f64>> = OnceLock::new();
static JWT_VERIFICATION_COUNTER: OnceLock<Counter<u64>> = OnceLock::new();
static JWT_VERIFICATION_DURATION: OnceLock<Histogram<f64>> = OnceLock::new();
static JWKS_FETCH_COUNTER: OnceLock<Counter<u64>> = OnceLock::new();
static JWKS_CACHE_MISSES: OnceLock<Counter<u64>> = OnceLock::new();
static JWKS_LOOKUPS: OnceLock<Counter<u64>> = OnceLock::new();

pub fn init_metrics() {
    let meter = opentelemetry::global::meter("aviation-gateway");

    let request_counter = meter
        .u64_counter("aviation_gateway_requests_total")
        .with_description("Total number of HTTP requests processed")
        .init();

    let request_duration = meter
        .f64_histogram("aviation_gateway_request_duration_ms")
        .with_description("HTTP request duration in milliseconds")
        .init();

    let jwt_verification_counter = meter
        .u64_counter("aviation_gateway_jwt_verifications_total")
        .with_description("Total number of JWT verifications")
        .init();

    let jwt_verification_duration = meter
        .f64_histogram("aviation_gateway_jwt_verification_duration_ms")
        .with_description("JWT verification duration in milliseconds")
        .init();

    let jwks_fetch_counter = meter
        .u64_counter("aviation_gateway_jwks_fetches_total")
        .with_description("Total number of JWKS fetches")
        .init();

    let jwks_cache_misses = meter
        .u64_counter("aviation_gateway_jwks_cache_misses_total")
        .with_description("Total number of JWKS cache misses")
        .init();

    let jwks_lookups = meter
        .u64_counter("aviation_gateway_jwks_lookups_total")
        .with_description("Total number of JWKS cache lookups (hits + misses)")
        .init();

    METER.set(meter).ok();
    REQUEST_COUNTER.set(request_counter).ok();
    REQUEST_DURATION.set(request_duration).ok();
    JWT_VERIFICATION_COUNTER.set(jwt_verification_counter).ok();
    JWT_VERIFICATION_DURATION
        .set(jwt_verification_duration)
        .ok();
    JWKS_FETCH_COUNTER.set(jwks_fetch_counter).ok();
    JWKS_CACHE_MISSES.set(jwks_cache_misses).ok();
    JWKS_LOOKUPS.set(jwks_lookups).ok();
}

// Metrics for gateway service
pub struct GatewayMetrics;

impl GatewayMetrics {
    pub fn request_started(method: &str, path: &str) {
        tracing::debug!("Request started: {} {}", method, path);
    }

    pub fn request_completed(method: &str, path: &str, status: u16, duration_ms: f64) {
        let labels = &[
            KeyValue::new("method", method.to_string()),
            KeyValue::new("path", path.to_string()),
            KeyValue::new("status", status.to_string()),
        ];

        if let Some(counter) = REQUEST_COUNTER.get() {
            counter.add(1, labels);
        }

        if let Some(histogram) = REQUEST_DURATION.get() {
            histogram.record(duration_ms, labels);
        }

        info!(
            method = method,
            path = path,
            status = status,
            duration_ms = duration_ms,
            "Request completed"
        );
    }

    pub fn router_request_started() {
        tracing::debug!("Router request started");
    }

    pub fn router_request_completed(status: u16, duration_ms: f64) {
        tracing::debug!(
            status = status,
            duration_ms = duration_ms,
            "Router request completed"
        );
    }

    pub fn jwt_verification_started() {
        tracing::debug!("JWT verification started");
    }

    pub fn jwt_verification_completed(success: bool, duration_ms: f64) {
        let labels = &[KeyValue::new("success", success.to_string())];

        if let Some(counter) = JWT_VERIFICATION_COUNTER.get() {
            counter.add(1, labels);
        }

        if let Some(histogram) = JWT_VERIFICATION_DURATION.get() {
            histogram.record(duration_ms, labels);
        }

        tracing::debug!(
            success = success,
            duration_ms = duration_ms,
            "JWT verification completed"
        );
    }

    pub fn jwks_fetch_started() {
        tracing::debug!("JWKS fetch started");
    }

    pub fn jwks_fetch_completed(success: bool, duration_ms: f64) {
        let labels = &[KeyValue::new("success", success.to_string())];

        if let Some(counter) = JWKS_FETCH_COUNTER.get() {
            counter.add(1, labels);
        }

        tracing::debug!(
            success = success,
            duration_ms = duration_ms,
            "JWKS fetch completed"
        );
    }

    pub fn jwks_cache_miss() {
        if let Some(counter) = JWKS_CACHE_MISSES.get() {
            counter.add(1, &[]);
        }

        tracing::debug!("JWKS cache miss");
    }

    pub fn jwks_lookup() {
        if let Some(counter) = JWKS_LOOKUPS.get() {
            counter.add(1, &[]);
        }
        tracing::debug!("JWKS cache lookup");
    }
}

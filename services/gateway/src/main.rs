mod config;
mod requests;
mod verify;
mod observability;

use crate::config::Config;
use crate::requests::forward_request;
use crate::verify::{JwksCache, verify_jwt};
use crate::observability::{init_observability, shutdown_observability, ObservabilityConfig, GatewayMetrics};
use axum::http::StatusCode;
use axum::response::Response;
use axum::{Router, extract::Request, routing::any};
use reqwest::Client;
use std::net::SocketAddr;
use std::sync::Arc;
use std::time::{Duration, Instant};
use tracing::{info, warn, error, debug, instrument, Span};

#[tokio::main]
async fn main() {
    let config = Config::from_env();
    let obs_config = ObservabilityConfig::from_env();

    if let Err(e) = init_observability(&obs_config).await {
        eprintln!("Failed to initialize observability: {}", e);
        std::process::exit(1);
    }

    info!("Starting aviation gateway service");
    config.validate();
    info!("Configuration validated successfully");
    info!(router_url = %config.router_url, port = config.port, "Gateway configuration loaded");
    let router_url = config.router_url.clone();

    let jwks_cache = Arc::new(JwksCache::new(
        config.jwks_url.clone(),
        Duration::from_secs(300),
    ));

    let client = Arc::new(
        Client::builder()
            .timeout(std::time::Duration::from_secs(30))
            .build()
            .expect("Failed to create HTTP client"),
    );

    let app = Router::new()
        .route(
            "/",
            any({
                let router_url = router_url.clone();
                let client = client.clone();
                let jwks_cache = jwks_cache.clone();
                move |req| {
                    route_handler(req, router_url, client, jwks_cache)
                }
            }),
        )
        .route(
            "/{*wildcard}",
            any({
                let router_url = router_url.clone();
                let client = client.clone();
                let jwks_cache = jwks_cache.clone();
                move |req| {
                    route_handler(req, router_url, client, jwks_cache)
                }
            }),
        );

    let addr = SocketAddr::from(([0, 0, 0, 0], config.port));
    info!(address = %addr, "Aviation gateway service started");

    let listener = tokio::net::TcpListener::bind(addr)
        .await
        .expect("Failed to bind to address");

    let result = axum::serve(listener, app).await;

    if let Err(e) = result {
        error!(error = %e, "Server error occurred");
    }

    info!("Shutting down gateway service");
    shutdown_observability().await;
}

#[instrument(
    name = "route_handler",
    skip(req, client, jwks_cache),
    fields(
        method = ?req.method(),
        uri = %req.uri(),
        user_agent = req.headers().get("user-agent").and_then(|v| v.to_str().ok()).unwrap_or("unknown")
    )
)]
async fn route_handler(
    req: Request,
    router_url: String,
    client: Arc<Client>,
    jwks_cache: Arc<JwksCache>,
) -> Response {
    let start_time = Instant::now();
    let method = req.method().to_string();
    let path = req.uri().path().to_string();
    let request_id = uuid::Uuid::new_v4().to_string();

    Span::current().record("request_id", &request_id);

    // Record request start
    GatewayMetrics::request_started(&method, &path);
    debug!(request_id = %request_id, "Processing incoming request");

    let auth_header = req
        .headers()
        .get("authorization")
        .and_then(|v| v.to_str().ok());

    let jwt = match auth_header {
        Some(header) if header.starts_with("Bearer ") => {
            debug!(request_id = %request_id, "Found Bearer token in Authorization header");
            header.trim_start_matches("Bearer ").trim()
        }
        _ => {
            let duration = start_time.elapsed().as_secs_f64() * 1000.0;
            warn!(request_id = %request_id, "Missing or invalid Authorization header");
            GatewayMetrics::request_completed(&method, &path, 401, duration);
            return Response::builder()
                .status(StatusCode::UNAUTHORIZED)
                .body(axum::body::Body::from(
                    "Missing or invalid Authorization header",
                ))
                .unwrap();
        }
    };

    let claims = match verify_jwt(jwt, &client, &jwks_cache).await {
        Ok(c) => {
            info!(
                request_id = %request_id,
                user_id = %c.sub,
                org_id = %c.org_id,
                roles = ?c.roles,
                "JWT verification successful"
            );
            c
        }
        Err(e) => {
            let duration = start_time.elapsed().as_secs_f64() * 1000.0;
            error!(
                request_id = %request_id,
                error = %e,
                "JWT verification failed"
            );
            GatewayMetrics::request_completed(&method, &path, 401, duration);
            return Response::builder()
                .status(StatusCode::UNAUTHORIZED)
                .body(axum::body::Body::from("Invalid or expired token"))
                .unwrap();
        }
    };

    let (mut parts, body) = req.into_parts();

    // Add user context headers for downstream services
    if let Ok(header_value) = claims.sub.parse() {
        parts.headers.insert("x-user-sub", header_value);
    }
    if let Ok(header_value) = claims.org_id.parse() {
        parts.headers.insert("x-org-id", header_value);
    }
    if let Ok(header_value) = claims.roles.join(",").parse() {
        parts.headers.insert("x-user-roles", header_value);
    }

    // Add request ID for distributed tracing
    if let Ok(header_value) = request_id.parse() {
        parts.headers.insert("x-request-id", header_value);
    }

    debug!(
        request_id = %request_id,
        "Added user context headers, forwarding request to router"
    );

    let req = Request::from_parts(parts, body);

    match forward_request(req, router_url, client).await {
        Ok(resp) => {
            let duration = start_time.elapsed().as_secs_f64() * 1000.0;
            let status = resp.status().as_u16();
            info!(
                request_id = %request_id,
                status = %status,
                duration_ms = %duration,
                "Request forwarded successfully"
            );
            GatewayMetrics::request_completed(&method, &path, status, duration);
            resp
        }
        Err(code) => {
            let duration = start_time.elapsed().as_secs_f64() * 1000.0;
            let status = code.as_u16();
            error!(
                request_id = %request_id,
                status = %status,
                duration_ms = %duration,
                "Error forwarding request to router"
            );
            GatewayMetrics::request_completed(&method, &path, status, duration);
            Response::builder()
                .status(code)
                .body(axum::body::Body::from("Error forwarding request"))
                .unwrap()
        }
    }
}

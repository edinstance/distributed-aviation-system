use crate::observability::GatewayMetrics;
use crate::requests::forward_request;
use crate::verify::{JwksCache, verify_jwt};
use axum::http::StatusCode;
use axum::response::{IntoResponse, Response};
use axum::{Json, extract::Request};
use reqwest::Client;
use serde_json::json;
use std::sync::Arc;
use std::time::Instant;
use tracing::{Span, debug, error, info, instrument, warn};

#[instrument(
    name = "route_handler",
    skip(req, client, jwks_cache),
    fields(
        method = ?req.method(),
        uri = %req.uri(),
        user_agent = req.headers().get("user-agent").and_then(|v| v.to_str().ok()).unwrap_or("unknown")
    )
)]
pub async fn route_handler(
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
            return (
                StatusCode::UNAUTHORIZED,
                Json(json!({
                    "error": "Missing or invalid Authorization header"
                })),
            )
                .into_response();
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
            return (
                StatusCode::UNAUTHORIZED,
                Json(json!({
                    "error": "Invalid or expired token"
                })),
            )
                .into_response();
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

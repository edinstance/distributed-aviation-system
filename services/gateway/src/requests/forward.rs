use axum::{body::to_bytes, extract::Request, http::StatusCode, response::Response};
use crate::observability::GatewayMetrics;
use reqwest::Client;
use std::{sync::Arc, time::Instant};
use tracing::{info, error, debug, instrument};

#[instrument(
    name = "forward_request",
    skip(req, client),
    fields(
        method = ?req.method(),
        uri = %req.uri(),
        router_url = %router_url,
        request_id = req.headers().get("x-request-id").and_then(|v| v.to_str().ok()).unwrap_or("unknown")
    )
)]
pub async fn forward_request(
    req: Request,
    router_url: String,
    client: Arc<Client>,
) -> Result<Response, StatusCode> {
    let (parts, body) = req.into_parts();

    let request_id = parts.headers
        .get("x-request-id")
        .and_then(|v| v.to_str().ok())
        .unwrap_or("unknown");

    let router_start = Instant::now();
    GatewayMetrics::router_request_started();

    debug!(
        request_id = %request_id,
        method = ?parts.method,
        uri = %parts.uri,
        "Preparing to forward request to router"
    );

    let body_bytes = to_bytes(body, 10 * 1024 * 1024)
        .await
        .map_err(|e| {
            error!(
                request_id = %request_id,
                error = %e,
                "Failed to read request body"
            );
            StatusCode::BAD_REQUEST
        })?;

    debug!(
        request_id = %request_id,
        body_size = body_bytes.len(),
        "Request body processed"
    );

    let mut headers = parts.headers.clone();
    headers.remove("content-length");
    headers.remove("transfer-encoding");
    headers.remove("host");

    debug!(
        request_id = %request_id,
        headers_count = headers.len(),
        "Request headers prepared"
    );

    let resp = client
        .request(parts.method, &router_url)
        .headers(headers)
        .body(body_bytes)
        .send()
        .await
        .map_err(|err| {
            let duration = router_start.elapsed().as_secs_f64() * 1000.0;
            if err.is_connect() {
                error!(
                    request_id = %request_id,
                    error = %err,
                    router_url = %router_url,
                    "Connection failed to router"
                );
                GatewayMetrics::router_request_completed(503, duration);
                StatusCode::SERVICE_UNAVAILABLE
            } else if err.is_timeout() {
                error!(
                    request_id = %request_id,
                    error = %err,
                    router_url = %router_url,
                    "Timeout contacting router"
                );
                GatewayMetrics::router_request_completed(504, duration);
                StatusCode::GATEWAY_TIMEOUT
            } else {
                error!(
                    request_id = %request_id,
                    error = %err,
                    router_url = %router_url,
                    "Error contacting router"
                );
                GatewayMetrics::router_request_completed(502, duration);
                StatusCode::BAD_GATEWAY
            }
        })?;

    let status = resp.status();
    let resp_headers = resp.headers().clone();

    debug!(
        request_id = %request_id,
        status = %status,
        response_headers_count = resp_headers.len(),
        "Received response from router"
    );

    let bytes = resp.bytes().await.map_err(|e| {
        let duration = router_start.elapsed().as_secs_f64() * 1000.0;
        error!(
            request_id = %request_id,
            error = %e,
            "Failed to read response body from router"
        );
        GatewayMetrics::router_request_completed(502, duration);
        StatusCode::BAD_GATEWAY
    })?;

    debug!(
        request_id = %request_id,
        response_size = bytes.len(),
        "Response body processed"
    );

    let mut builder = Response::builder().status(status);
    for (header_name, header_value) in resp_headers.iter() {
        if !header_name
            .as_str()
            .eq_ignore_ascii_case("transfer-encoding")
        {
            builder = builder.header(header_name, header_value);
        }
    }

    let duration = router_start.elapsed().as_secs_f64() * 1000.0;
    info!(
        request_id = %request_id,
        status = %status,
        response_size = bytes.len(),
        duration_ms = %duration,
        "Successfully forwarded request and received response"
    );

    GatewayMetrics::router_request_completed(status.as_u16(), duration);

    builder
        .body(axum::body::Body::from(bytes))
        .map_err(|e| {
            error!(
                request_id = %request_id,
                error = %e,
                "Failed to build response"
            );
            StatusCode::INTERNAL_SERVER_ERROR
        })
}

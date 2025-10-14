use axum::{body::to_bytes, extract::Request, http::StatusCode, response::Response};
use reqwest::Client;
use std::sync::Arc;

pub async fn forward_request(
    req: Request,
    router_url: String,
    client: Arc<Client>,
) -> Result<Response, StatusCode> {
    let (parts, body) = req.into_parts();
    let body_bytes = to_bytes(body, 10 * 1024 * 1024)
        .await
        .map_err(|_| StatusCode::BAD_REQUEST)?;

    let mut headers = parts.headers.clone();
    headers.remove("content-length");
    headers.remove("transfer-encoding");
    headers.remove("host");

    let resp = client
        .request(parts.method, &router_url)
        .headers(headers)
        .body(body_bytes)
        .send()
        .await
        .map_err(|err| {
            eprintln!("Error contacting router: {err}");
            if err.is_connect() {
                StatusCode::SERVICE_UNAVAILABLE
            } else if err.is_timeout() {
                StatusCode::GATEWAY_TIMEOUT
            } else {
                StatusCode::BAD_GATEWAY
            }
        })?;

    let status = resp.status();
    let resp_headers = resp.headers().clone();
    let bytes = resp.bytes().await.map_err(|_| StatusCode::BAD_GATEWAY)?;

    let mut builder = Response::builder().status(status);
    for (header_name, header_value) in resp_headers.iter() {
        if !header_name
            .as_str()
            .eq_ignore_ascii_case("transfer-encoding")
        {
            builder = builder.header(header_name, header_value);
        }
    }

    builder
        .body(axum::body::Body::from(bytes))
        .map_err(|_| StatusCode::INTERNAL_SERVER_ERROR)
}

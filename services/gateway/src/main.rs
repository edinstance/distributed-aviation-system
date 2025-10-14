mod config;
mod requests;
mod verify;

use crate::config::Config;
use crate::requests::forward_request;
use axum::response::Response;
use axum::{Router, extract::Request, routing::any};
use reqwest::Client;
use std::net::SocketAddr;
use std::sync::Arc;
use axum::http::StatusCode;
use crate::verify::verify_jwt;

#[tokio::main]
async fn main() {
    let config = Config::from_env();
    config.validate();
    println!("Config is valid");
    println!("Router URL: {}", config.router_url);
    let router_url = config.router_url.clone();

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
                move |req| route_handler(req, router_url.clone(), client.clone())
            }),
        )
        .route(
            "/{*wildcard}",
            any({
                let client = client.clone();
                move |req| route_handler(req, router_url, client)
            }),
        );

    let addr = SocketAddr::from(([0, 0, 0, 0], 1000));
    println!("🚀 Gateway running on {}", addr);
    axum::serve(tokio::net::TcpListener::bind(addr).await.unwrap(), app)
        .await
        .unwrap();
}

async fn route_handler(req: Request, router_url: String, client: Arc<Client>) -> Response {
    let auth_header = req.headers().get("authorization").and_then(|v| v.to_str().ok());
    let jwt = match auth_header {
        Some(header) if header.starts_with("Bearer ") => header.trim_start_matches("Bearer ").trim(),
        _ => {
            return Response::builder()
                .status(StatusCode::UNAUTHORIZED)
                .body(axum::body::Body::from("Missing or invalid Authorization header"))
                .unwrap();
        }
    };

    let claims = match verify_jwt(jwt, &client).await {
        Ok(c) => c,
        Err(e) => {
            eprintln!("JWT verification failed: {:?}", e);
            return Response::builder()
                .status(StatusCode::UNAUTHORIZED)
                .body(axum::body::Body::from("Invalid or expired token"))
                .unwrap();
        }
    };
    print!("Claims: {:#?}", claims);

    let (mut parts, body) = req.into_parts();
    if let Ok(header_value) = claims.sub.parse() {
        parts.headers.insert("x-user-sub", header_value);
    }
    if let Ok(header_value) = claims.org_id.parse() {
        parts.headers.insert("x-org-id", header_value);
    }
    if let Ok(header_value) = claims.roles.join(",").parse() {
        parts.headers.insert("x-user-roles", header_value);
    }

    let req = Request::from_parts(parts, body);

    match forward_request(req, router_url, client).await {
        Ok(resp) => resp,
        Err(code) => Response::builder()
            .status(code)
            .body(axum::body::Body::from("Error forwarding request"))
            .unwrap(),
    }
}

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
    println!("{:?}", req);
    match forward_request(req, router_url, client).await {
        Ok(resp) => resp,
        Err(code) => Response::builder()
            .status(code)
            .body(axum::body::Body::from("Error forwarding request"))
            .unwrap(),
    }
}

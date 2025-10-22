mod config;
mod handling;
mod observability;
mod requests;
mod verify;

use crate::config::Config;
use crate::handling::{route_handler, health_handler};
use crate::observability::{ObservabilityConfig, init_observability, shutdown_observability};
use crate::verify::JwksCache;
use axum::{Router, routing::any};
use reqwest::Client;
use std::net::SocketAddr;
use std::sync::Arc;
use std::time::Duration;
use tracing::{error, info};

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
            .timeout(Duration::from_secs(30))
            .build()
            .expect("Failed to create HTTP client"),
    );

    let app = Router::new()
        .route("/health", any(health_handler))
        .route(
            "/",
            any({
                let router_url = router_url.clone();
                let client = client.clone();
                let jwks_cache = jwks_cache.clone();
                move |req| {
                    route_handler(req, router_url.clone(), client.clone(), jwks_cache.clone())
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
                    route_handler(req, router_url.clone(), client.clone(), jwks_cache.clone())
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

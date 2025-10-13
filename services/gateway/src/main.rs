mod verify;
mod config;

use axum::{Router, extract::Request, routing::any};
use std::net::SocketAddr;

use crate::config::Config;

#[tokio::main]
async fn main() {
    let config = Config::from_env();
    config.validate();
    println!("Config is valid");

    let app = Router::new().route("/{*wildcard}", any(move |req| route_handler(req)));

    let addr = SocketAddr::from(([0, 0, 0, 0], 1000));
    println!("🚀 Gateway running on {}", addr);
    axum::serve(tokio::net::TcpListener::bind(addr).await.unwrap(), app)
        .await
        .unwrap();
}

async fn route_handler(req: Request) {
    println!("{:?}", req);
}

use std::env;
use dotenv::dotenv;

#[derive(Clone, Debug)]
pub struct Config {
    pub jwks_url: String,
    pub router_url: String,
    pub port: u16,
    pub log_level: String,
}

impl Config {
    pub fn from_env() -> Config {
        dotenv().ok();

        Self {
            jwks_url: env::var("JWKS_URL")
                .expect("Please set JWKS_URL (e.g. https://auth.example.com/.well-known/jwks.json)"),
            router_url: env::var("ROUTER_URL")
                .expect("Please set ROUTER_URL (e.g. https://apollo-router/)"),
            port: env::var("PORT")
                .ok()
                .and_then(|p| p.parse().ok())
                .unwrap_or(1000),
            log_level: env::var("LOG_LEVEL").unwrap_or_else(|_| "info".to_string()),
        }
    }
    pub fn validate(&self) {
        assert!(
            !self.jwks_url.is_empty(),
            "JWKS_URL must not be empty"
        );
        assert!(
            !self.router_url.is_empty(),
            "ROUTER_URL must not be empty"
        );
        assert!(
            self.port > 0,
            "PORT must be a positive integer"
        );
    }
}
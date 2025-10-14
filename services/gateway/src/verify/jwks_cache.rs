use crate::verify::Jwks;
use anyhow::Result;
use reqwest::Client;
use std::{
    sync::{Arc, RwLock},
    time::{Duration, Instant},
};

#[derive(Clone)]
pub struct JwksCache {
    inner: Arc<RwLock<Option<(Instant, Jwks)>>>,
    ttl: Duration,
    jwks_url: String,
}

impl JwksCache {
    pub fn new(jwks_url: String, ttl: Duration) -> Self {
        Self {
            inner: Arc::new(RwLock::new(None)),
            ttl,
            jwks_url,
        }
    }

    pub async fn get_jwks(&self, client: &Client) -> Result<Jwks> {
        // Fast path: check cache first
        if let Some((timestamp, jwks)) = self.inner.read().unwrap().clone() {
            if timestamp.elapsed() < self.ttl {
                return Ok(jwks);
            }
        }

        // Cache expired → refresh from network
        let jwks: Jwks = client
            .get(&self.jwks_url)
            .send()
            .await?
            .error_for_status()?
            .json()
            .await?;

        // Update cache
        *self.inner.write().unwrap() = Some((Instant::now(), jwks.clone()));

        Ok(jwks)
    }
}
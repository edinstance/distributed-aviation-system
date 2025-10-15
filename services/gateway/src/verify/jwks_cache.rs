use crate::verify::Jwks;
use crate::observability::GatewayMetrics;
use anyhow::Result;
use reqwest::Client;
use std::{
    sync::{Arc, RwLock},
    time::{Duration, Instant},
};
use tracing::{info, error, debug, instrument};

#[derive(Clone)]
pub struct JwksCache {
    inner: Arc<RwLock<Option<(Instant, Jwks)>>>,
    ttl: Duration,
    jwks_url: String,
}

impl JwksCache {
    pub fn new(jwks_url: String, ttl: Duration) -> Self {
        info!(
            jwks_url = %jwks_url,
            ttl_seconds = ttl.as_secs(),
            "Creating new JWKS cache"
        );
        Self {
            inner: Arc::new(RwLock::new(None)),
            ttl,
            jwks_url,
        }
    }

    #[instrument(
        name = "get_jwks",
        skip(self, client),
        fields(
            jwks_url = %self.jwks_url,
            ttl_seconds = self.ttl.as_secs()
        )
    )]
    pub async fn get_jwks(&self, client: &Client) -> Result<Jwks> {
        // Fast path: check cache first
        if let Some((timestamp, jwks)) = self.inner.read().unwrap().clone() {
            let age = timestamp.elapsed();
            if age < self.ttl {
                debug!(
                    age_seconds = age.as_secs(),
                    keys_count = jwks.keys.len(),
                    "Returning cached JWKS"
                );
                return Ok(jwks);
            } else {
                debug!(
                    age_seconds = age.as_secs(),
                    ttl_seconds = self.ttl.as_secs(),
                    "JWKS cache expired, refreshing from network"
                );
            }
        } else {
            debug!("No JWKS in cache, fetching from network");
            GatewayMetrics::jwks_cache_miss();
        }

        // Cache expired → refresh from network
        let fetch_start = Instant::now();
        GatewayMetrics::jwks_fetch_started();
        info!(jwks_url = %self.jwks_url, "Fetching JWKS from remote endpoint");

        let jwks: Jwks = client
            .get(&self.jwks_url)
            .send()
            .await
            .map_err(|e| {
                let duration = fetch_start.elapsed().as_secs_f64() * 1000.0;
                error!(
                    error = %e,
                    jwks_url = %self.jwks_url,
                    "Failed to fetch JWKS from remote endpoint"
                );
                GatewayMetrics::jwks_fetch_completed(false, duration);
                anyhow::anyhow!("JWKS fetch failed: {}", e)
            })?
            .error_for_status()
            .map_err(|e| {
                let duration = fetch_start.elapsed().as_secs_f64() * 1000.0;
                error!(
                    error = %e,
                    jwks_url = %self.jwks_url,
                    "JWKS endpoint returned error status"
                );
                GatewayMetrics::jwks_fetch_completed(false, duration);
                anyhow::anyhow!("JWKS endpoint error: {}", e)
            })?
            .json()
            .await
            .map_err(|e| {
                let duration = fetch_start.elapsed().as_secs_f64() * 1000.0;
                error!(
                    error = %e,
                    jwks_url = %self.jwks_url,
                    "Failed to parse JWKS JSON response"
                );
                GatewayMetrics::jwks_fetch_completed(false, duration);
                anyhow::anyhow!("JWKS JSON parse error: {}", e)
            })?;

        // Update cache
        *self.inner.write().unwrap() = Some((Instant::now(), jwks.clone()));

        let duration = fetch_start.elapsed().as_secs_f64() * 1000.0;
        info!(
            jwks_url = %self.jwks_url,
            keys_count = jwks.keys.len(),
            duration_ms = %duration,
            "Successfully cached fresh JWKS"
        );

        GatewayMetrics::jwks_fetch_completed(true, duration);
        Ok(jwks)
    }
}
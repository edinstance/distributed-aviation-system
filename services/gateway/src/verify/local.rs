use crate::verify::{Claims, JwksCache};
use crate::observability::GatewayMetrics;
use anyhow::Result;
use jsonwebtoken::{Algorithm, DecodingKey, Validation, decode, decode_header};
use reqwest::Client;
use std::{sync::Arc, time::Instant};
use tracing::{info, warn, error, debug, instrument};

#[instrument(
    name = "verify_jwt",
    skip(jwt, client, jwks_cache),
    fields(
        jwt_length = jwt.len()
    )
)]
pub async fn verify_jwt(
    jwt: &str,
    client: &Arc<Client>,
    jwks_cache: &Arc<JwksCache>,
) -> Result<Claims> {
    let start_time = Instant::now();
    GatewayMetrics::jwt_verification_started();
    debug!("Starting JWT verification process");

    // Decode jwt to read KID
    let header = decode_header(&jwt)
        .map_err(|e| {
            let duration = start_time.elapsed().as_secs_f64() * 1000.0;
            error!(error = %e, "Failed to decode JWT header");
            GatewayMetrics::jwt_verification_completed(false, duration);
            anyhow::anyhow!("Invalid JWT format: {}", e)
        })?;

    let kid = header
        .kid
        .ok_or_else(|| {
            let duration = start_time.elapsed().as_secs_f64() * 1000.0;
            error!("JWT header missing required 'kid' field");
            GatewayMetrics::jwt_verification_completed(false, duration);
            anyhow::anyhow!("Missing kid in JWT header")
        })?;

    debug!(kid = %kid, "Extracted key ID from JWT header");

    // Fetch JWKS from cache
    let jwks = jwks_cache.get_jwks(client).await
        .map_err(|e| {
            let duration = start_time.elapsed().as_secs_f64() * 1000.0;
            error!(error = %e, "Failed to fetch JWKS from cache");
            GatewayMetrics::jwt_verification_completed(false, duration);
            e
        })?;

    debug!(keys_count = jwks.keys.len(), "Retrieved JWKS from cache");

    // Find matching key
    let jwk = jwks
        .keys
        .iter()
        .find(|k| k.kid == kid)
        .ok_or_else(|| {
            let duration = start_time.elapsed().as_secs_f64() * 1000.0;
            error!(
                kid = %kid,
                available_kids = ?jwks.keys.iter().map(|k| &k.kid).collect::<Vec<_>>(),
                "No matching key found in JWKS"
            );
            GatewayMetrics::jwt_verification_completed(false, duration);
            anyhow::anyhow!("Matching kid not found in JWKS")
        })?;

    debug!(
        kid = %kid,
        "Found matching JWK for token verification"
    );

    // Build decoding key from base64url components
    let decoding_key = DecodingKey::from_rsa_components(&jwk.n, &jwk.e)
        .map_err(|e| {
            let duration = start_time.elapsed().as_secs_f64() * 1000.0;
            error!(
                error = %e,
                kid = %kid,
                "Failed to construct RSA decoding key from JWK components"
            );
            GatewayMetrics::jwt_verification_completed(false, duration);
            anyhow::anyhow!("Failed to build RSA decoder: {}", e)
        })?;

    debug!("Successfully constructed RSA decoding key");

    // Verify JWT
    let mut validation = Validation::new(Algorithm::RS256);
    validation.validate_exp = true;

    let token_data = decode::<Claims>(jwt, &decoding_key, &validation)
        .map_err(|e| {
            let duration = start_time.elapsed().as_secs_f64() * 1000.0;
            match e.kind() {
                jsonwebtoken::errors::ErrorKind::ExpiredSignature => {
                    warn!(kid = %kid, "JWT token has expired");
                    GatewayMetrics::jwt_verification_completed(false, duration);
                    anyhow::anyhow!("Token expired: {}", e)
                }
                jsonwebtoken::errors::ErrorKind::InvalidSignature => {
                    error!(kid = %kid, "JWT signature validation failed");
                    GatewayMetrics::jwt_verification_completed(false, duration);
                    anyhow::anyhow!("Invalid signature: {}", e)
                }
                _ => {
                    error!(error = %e, kid = %kid, "JWT validation failed");
                    GatewayMetrics::jwt_verification_completed(false, duration);
                    anyhow::anyhow!("JWT validation failed: {}", e)
                }
            }
        })?;

    let claims = token_data.claims;
    let duration = start_time.elapsed().as_secs_f64() * 1000.0;

    info!(
        kid = %kid,
        user_id = %claims.sub,
        org_id = %claims.org_id,
        roles = ?claims.roles,
        expires_at = claims.exp,
        duration_ms = %duration,
        "JWT verification successful"
    );

    GatewayMetrics::jwt_verification_completed(true, duration);
    Ok(claims)
}

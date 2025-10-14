use crate::verify::{Claims, JwksCache};
use anyhow::Result;
use jsonwebtoken::{Algorithm, DecodingKey, Validation, decode, decode_header};
use reqwest::Client;
use std::sync::Arc;

pub async fn verify_jwt(
    jwt: &str,
    client: &Arc<Client>,
    jwks_cache: &Arc<JwksCache>,
) -> Result<Claims> {
    // Decode jwt to read KID
    let header = decode_header(&jwt)?;
    let kid = header
        .kid
        .ok_or_else(|| anyhow::anyhow!("Missing kid in JWT header"))?;

    // Fetch JWKS from cache
    let jwks = jwks_cache.get_jwks(client).await?;

    // Find matching key
    let jwk = jwks
        .keys
        .iter()
        .find(|k| k.kid == kid)
        .ok_or_else(|| anyhow::anyhow!("Matching kid not found in JWKS"))?;

    // Build decoding key from base64url components
    let decoding_key = DecodingKey::from_rsa_components(&jwk.n, &jwk.e)
        .map_err(|_| anyhow::anyhow!("Failed to build RSA decoder"))?;

    // Verify JWT
    let mut validation = Validation::new(Algorithm::RS256);
    validation.validate_exp = true;

    let token_data = decode::<Claims>(jwt, &decoding_key, &validation)?;
    let claims = token_data.claims;
    Ok(claims)
}

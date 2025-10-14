use std::sync::Arc;
use crate::config::Config;
use crate::verify::{Claims, Jwks};
use anyhow::Result;
use jsonwebtoken::{Algorithm, DecodingKey, Validation, decode, decode_header};
use reqwest::Client;

pub async fn verify_jwt(jwt: &str, client: &Arc<Client>) -> Result<Claims> {
    let config = Config::from_env();
    let jwks_url = config.jwks_url;

    // Decode jwt to read KID
    let header = decode_header(&jwt)?;
    let kid = header
        .kid
        .ok_or_else(|| anyhow::anyhow!("Missing kid in JWT header"))?;

    // Fetch JWKS
    let jwks: Jwks = client.get(jwks_url).send().await?.json().await?;

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

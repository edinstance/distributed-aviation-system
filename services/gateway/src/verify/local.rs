use anyhow::Result;
use dotenv::dotenv;
use jsonwebtoken::{Algorithm, DecodingKey, Validation, decode, decode_header};
use reqwest::Client;
use std::env;

use crate::verify::{Claims, Jwks};

pub async fn verify_jwt(jwt: &str) -> Result<Claims> {
    dotenv().ok();

    let jwks_url = env::var("JWKS_URL")?;
    
    // Decode jwt to read KID
    let header = decode_header(&jwt)?;
    let kid = header
        .kid
        .ok_or_else(|| anyhow::anyhow!("Missing kid in JWT header"))?;

    // Fetch JWKS
    let client = Client::new();
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
    let validation = Validation::new(Algorithm::RS256);

    let token_data = decode::<Claims>(jwt, &decoding_key, &validation)?;
    let claims = token_data.claims;
    Ok(claims)
}

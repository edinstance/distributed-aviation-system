use std::env;
use anyhow::Result;
use dotenv::dotenv;
use jsonwebtoken::{Algorithm, DecodingKey, Validation, decode, decode_header};
use reqwest::Client;
use serde::{Deserialize, Serialize};

#[derive(Debug, Deserialize, Serialize, Clone)]
struct Claims {
    sub: String,
    org_id: String,
    roles: Vec<String>,
    exp: usize,
}

#[derive(Debug, Deserialize)]
struct Jwk {
    kid: String,
    n: String,
    e: String,
}

#[derive(Debug, Deserialize)]
struct Jwks {
    keys: Vec<Jwk>,
}

#[tokio::main]
async fn main() -> Result<()> {
    dotenv().ok();

    let jwks_url = env::var("JWKS_URL")?;
    
    let jwt = "eyJhbGciOiJSUzI1NiIsImtpZCI6InYxIiwidHlwIjoiSldUIn0.eyJ0b2tlbl90eXBlIjoiYWNjZXNzIiwiZXhwIjoxNzU5OTQ3ODU4LCJpYXQiOjE3NTk5NDQyNTgsImp0aSI6IjBhYmE5ODg3MWE2ODQ3MjI5YTc0NTMwNGQyMDMwMzk1Iiwic3ViIjoiZjVlOGQ0YjctYThlZS00OGJjLWJkOGUtZWNkMjMzM2M0ODJlIiwidXNlcm5hbWUiOiJhZG1pbiIsImVtYWlsIjoiYWRtaW5AbXljb21wYW55LmNvbSIsIm9yZ19pZCI6ImFmMTk5Y2RiLWE5YzktNGM4MS05ZDJjLWEwZGJjN2U3YTlkZSIsInJvbGVzIjpbXX0.MdNGHpkPZhqaoBZ1ge0QtNChnyjN0MSqYnBu6S1gPSgkUTDOJ1CNF1kT4_2k4VIPAT73Uc7AHgNuj9isWd4hK6peIPac4KJlH6zztRFFsWTIXuhLMqRTIMgtmpIuYOlGOTNIf-9gmfZHVorUyMtOcXdJIJoGYmhYMQC6f8HoM9lS4iFvudrPfOkqdTOb3u1rFXMhSABDBteCLCpZvP5F0l6wyBfpRrToUHBwUYIR2DF3wQjE8eHNcmRE1Va5p4k9d-cw0cIFBUB-Mneis2r_cf1wlOCM7m6uHkbR71ueHwfwCoYjaylS3jiYRbUaYLzzMmx83_lqXnZtaiK5GmF7dRU7bmay_d7-zdhmar5QK4PcAY676HUmHPlYAmXBJmiypGA7kpOdtnOG-aamkCKXm0vyKU3Rt0Lh0_S94V5TMqQXFej5Btkl-9wti5A4nlFIi87HkjXhqzWl8oslfkePBKl8uV6NuW7_sSa0A-A-oYcUpLvQK0BFTFNwlW_de-1F31FssxqyLqbUfmxQli9JRwPIFAO_o6TKq9KgzKb02mSWqHbupzNZrQi4w_Yw3QB5qCT-3MCgifbAz7jxGbpJy_EexqGNrWwqgOyU7C026qKrEx9eFgFrCTM0xk2XKtOG1hHAfpxopqFYP8f3lNRiXL5aokVGhzGad_Ek81IDIYI";

    // Decode header to read KID
    let header = decode_header(jwt)?;
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
    let token_data = decode::<Claims>(jwt, &decoding_key, &validation);

    match token_data {
        Ok(data) => println!("Verified! Claims: {:?}", data.claims),
        Err(err) => eprintln!("Verification failed: {:?}", err),
    }

    Ok(())
}

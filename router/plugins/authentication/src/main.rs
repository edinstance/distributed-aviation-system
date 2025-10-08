use anyhow::Result;
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
    let jwt = "eyJhbGciOiJSUzI1NiIsImtpZCI6InYxIiwidHlwIjoiSldUIn0.eyJ0b2tlbl90eXBlIjoiYWNjZXNzIiwiZXhwIjoxNzU5OTI3NzQwLCJpYXQiOjE3NTk5MjQxNDAsImp0aSI6ImI3ZTZlNDdmYjI2ZTQ4ZDg5OThlMTEzZDFjNjhjZmVhIiwic3ViIjoiZjVlOGQ0YjctYThlZS00OGJjLWJkOGUtZWNkMjMzM2M0ODJlIiwidXNlcm5hbWUiOiJhZG1pbiIsImVtYWlsIjoiYWRtaW5AbXljb21wYW55LmNvbSIsIm9yZ19pZCI6ImFmMTk5Y2RiLWE5YzktNGM4MS05ZDJjLWEwZGJjN2U3YTlkZSIsInJvbGVzIjpbXX0.uQ8LF0NrfUZSYbHPeseGaRaZ93lbOWMkMtzhKpQm_ic1icHwNDwzkyk7T_PHIN2a-8-Cc0cod2NzjQUQUtiHPi2mSqRpq9b0BGYLAl8hf0XjJF4BlkX78slI755v6ffucCS6K70w7RzwPtLsa1t_qWesOTMaVhfxWXh6VnDridDmv4CFebs4rL5Q5ygYFnnIT-vyUQ5LK9B0lS_aE1e2FYfLNbE95RLjjxP8eCiYkYFEqW1IoaANKv9v9mqjwlccpWy6bPsV8TDYBGl5bQGRcA6jOwavDIBIgQmA0PLpcYeHWtUCSuOyFrr0GzWe6MifhYqUuOVbFQLiVUoDOPUChyX9qA2CLpr0dsVeySDQrcOTzIcig3cW1E_RNt3wY02t8KGOoFHyHXet7Ub8le9olt60IvNDY8yF-PBXZaDt-82yFnvgD1p5vEaAibbtAXNn82S0b4S1Fxd7gQggGOg7cb8fVV2QW2w_wL_O9lNp6SkoSu13wIGRSIbypMpz9Wgu6SJomH9fBFWr6G1qI-cSaK-lY74C3gSd8bKO_spOh_hgZH_WyXzMRAtR7e4QjGNszyc3yZ-01I3ZqGCD3w7BUKLU_SCzmqUhaRoMuAVDL1wnnbVdNFkPPEkgYPCdxdJDwb6pfHm4sH39twJh_hxW6XaINVNVPv1NGUQw35_VlKg";
    let jwks_url = "http://localhost:8000/api/auth/jwks.json";

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

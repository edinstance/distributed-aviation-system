use std::fs;
use jsonwebtoken::{decode, DecodingKey, Validation, Algorithm};
use serde::{Deserialize, Serialize};

#[derive(Debug, Deserialize, Serialize, Clone)]
struct Claims {
    sub: String,
    org_id: String,
    roles: Vec<String>,
    exp: usize,
}

fn main() {
    let jwt = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbl90eXBlIjoiYWNjZXNzIiwiZXhwIjoxNzU5ODc4OTc1LCJpYXQiOjE3NTk4NzUzNzUsImp0aSI6ImFjZTUwZTVlNTQ3ZjQwMGNhYWUyMWViYzhkYjU4OGIzIiwic3ViIjoiM2MyMTY4MGUtYjhhNy00YzI1LTllN2QtNzFjMGRmY2RmZGE2IiwidXNlcm5hbWUiOiJhZG1pbiIsImVtYWlsIjoiYWRtaW5AbXljb21wYW55LmNvbSIsIm9yZ19pZCI6IjJiODczZjliLTA0NzEtNGEzZi1iOWUwLTM2ODU2MjE0NDYzMyIsInJvbGVzIjpbXX0.VzOKgqu0t4XpxvlJIIRgEkCLYqBuGRbc9cccrhSfzvkfP4iPIqBBKzr1OdgdqghcMG6Z0PGNTrpoPRJ6hFytB9ZJFtiMtKIDqWdJ6SdgMQOQmwDEqfZXSEnyv0-hOtZid-agRnAkSWgEF6oWOIQFj_pM2xSxMf0LnFbS741nTLp1A8k51WKFwicv1iI85AccQGvujEexff-iGk6Dz9UNWLNNp8ukFkKSf7qOsjS17YPIRwJdMNlR0bcCF2TAYmRDD3JMGOP2eXbVIDvWpgbYGv1O5VRQe5iighkfwc93zedV105Dn6QjzRF3kxGg6tNm9q1YGyxSXqh8PnFCHuHE6wYoR-v8IFPSmUFC2fpVLuDVSCOAOIfCw1hUIF8-okktCpGrCFaAx2q0ULwlWEql2H6pgsGwarfHGM4XGbZPot__tgCID-Beb960MzDU0S6YwmkLUJYdv_fE4tbMvXwzFuZ6VD1WkiS0EjXUwTnj1UbMv1bUXZx-QAzo4FhlZRgOPI7YJ6aWFCFOHVjDMRWwBNe9i0gBF6YwSBBAYOxC8hLPXGD0KbnYJGqz-vRFFR0TiBEQJWlsKUhArFwpcvqvDNF1IFnPSASM1fWw83BWH2wVghDOokdywHc1-2FNpCDUOAch5lZ0ojwR_B0YsnayeEiv5qx589LEeTOOT5vCyOI";

    let public_key_pem = fs::read("../../../services/authentication/keys/public.pem").expect("Cannot read public.pem");

    let key = DecodingKey::from_rsa_pem(&public_key_pem).expect("Invalid PEM");

    let validation = Validation::new(Algorithm::RS256);

    let token_data = decode::<Claims>(jwt, &key, &validation);

    match token_data {
        Ok(data) => println!("Decoded claims: {:?}", data.claims),
        Err(err) => println!("Error verifying JWT: {:?}", err),
    }
}
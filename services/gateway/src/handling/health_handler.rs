use axum::Json;
use serde_json::json;

pub async fn health_handler() -> impl axum::response::IntoResponse {
    Json(json!({ "status": "UP" }))
}